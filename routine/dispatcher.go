package routine

import (
	"log"
	"runtime/debug"
)

// Dispatcher define
// inspired by http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	maxWorkers int
	maxJobs    int
	JobQueue   chan Job
}

// NewDispatcher create dispatcher
func NewDispatcher(maxWorkers, maxJobs int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		maxWorkers: maxWorkers,
		maxJobs:    maxJobs,
		JobQueue:   make(chan Job, maxJobs),
	}
}

// Add job to queue
func (d *Dispatcher) Add(job Job) bool {
	select {
	case d.JobQueue <- job:
		return true
	default:
		return false
	}
}

// Run dispatcher
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// Close dispatcher
func (d *Dispatcher) Close() {

}

// CurrentJobs current job queue
func (d *Dispatcher) CurrentJobs() int {
	return len(d.JobQueue)
}

// CurrentWorkers current workers
func (d *Dispatcher) CurrentWorkers() int {
	return d.maxWorkers
}

// MaxWorkers max workers
func (d *Dispatcher) MaxWorkers() int {
	return d.maxWorkers
}

// MaxJobs max jobs
func (d *Dispatcher) MaxJobs() int {
	return d.maxJobs
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// try to obtain a worker job channel that is available.
			// this will block until a worker is idle
			jobChannel := <-d.WorkerPool

			// dispatch the job to the worker job channel
			jobChannel <- job
		}
	}
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// NewWorker create worker
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		defer func() {
			// if Do function is panic, worker will restart
			if err := recover(); err != nil {
				log.Printf("[ERR] execute Do panic, %v\n", err)
				log.Printf("stack: %s\n", debug.Stack())
				go w.Start()
			}
		}()
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				job.Do()
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// Job something need to do
type Job interface {
	Do()
}
