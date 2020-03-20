package main

import (
	"context"
	"log"

	"github.com/haormj/util/invoker"
	"github.com/haormj/util/invoker/function"
	"github.com/haormj/util/invoker/receiver"
)

type Test struct {
}

func (t *Test) SayHello(ctx context.Context, input *string, output *string) error {
	log.Println("receive from client: ", *input)
	*output = "hi"
	return nil
}

func SayWorld(ctx context.Context, input *string, output *string) error {
	log.Println("receive from client: ", *input)
	*output = "hi world"
	return nil
}

func main() {
	//TestReceiver()
	TestFunction()
}

func TestReceiver() {
	log.Println("--- Test receiver invoker begin ---")
	inv := receiver.NewInvoker(&Test{})
	if err := inv.Init(); err != nil {
		log.Fatalln(err)
	}
	log.Println("receiver name: ", inv.Name())
	mi := invoker.NewMessage()
	mi.SetFuncName("SayHello")
	input := "hello"
	output := ""
	log.Println("input: ", input)
	mi.SetParameters([]interface{}{context.Background(), &input, &output})
	mo, err := inv.Invoke(context.Background(), mi)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("output: ", output)
	log.Println("err: ", mo.Parameters()[0])
	log.Println("--- Test receiver invoker end ---")
}

func TestFunction() {
	log.Println("--- Test function invoker begin ---")
	inv := function.NewInvoker(SayWorld)
	if err := inv.Init(); err != nil {
		log.Fatalln(err)
	}
	log.Println("receiver name: ", inv.Name())
	mi := invoker.NewMessage()
	mi.SetFuncName("SayWorld")
	input := "hello world"
	output := ""
	log.Println("input: ", input)
	mi.SetParameters([]interface{}{context.Background(), &input, &output})
	mo, err := inv.Invoke(context.Background(), mi)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("output: ", output)
	log.Println("err: ", mo.Parameters()[0])
	log.Println("--- Test receiver invoker end ---")

}
