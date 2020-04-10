package util

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewPeriod(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	type args struct {
		st int64
		et int64
	}
	tests := []struct {
		name    string
		args    args
		want    Period
		wantErr error
	}{
		{
			name: "point",
			args: args{
				st: 0,
				et: 0,
			},
			want:    newPeriod(0, 0),
			wantErr: nil,
		},
		{
			name: "negative st",
			args: args{
				st: -1,
				et: 0,
			},
			want:    newPeriod(0, 0),
			wantErr: ErrPeriodStNegative,
		},
		{
			name: "negative et",
			args: args{
				st: 0,
				et: -1,
			},
			want:    newPeriod(0, -1),
			wantErr: nil,
		},
		{
			name: "st > et",
			args: args{
				st: 10,
				et: 5,
			},
			want:    newPeriod(0, 0),
			wantErr: ErrPeriodInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriod(tt.args.st, tt.args.et)
			if err != tt.wantErr {
				t.Errorf("NewPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPeriods(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	type args struct {
		se []int64
	}
	tests := []struct {
		name    string
		args    args
		want    []Period
		wantErr error
	}{
		{
			name: "empty period set",
			args: args{
				se: []int64{},
			},
			want:    []Period{},
			wantErr: nil,
		},
		{
			name: "normal periods",
			args: args{
				se: []int64{0, 1, 2, 3},
			},
			want:    []Period{newPeriod(0, 1), newPeriod(2, 3)},
			wantErr: nil,
		},
		{
			name: "illegal arg",
			args: args{
				se: []int64{0, 1, 2},
			},
			want:    nil,
			wantErr: ErrPeriodIllegalArg,
		},
		{
			name: "invalid period",
			args: args{
				se: []int64{0, 1, 2, 3, 4, -1, 5, 3},
			},
			want:    nil,
			wantErr: ErrPeriodInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriods(tt.args.se...)
			if err != tt.wantErr {
				t.Errorf("NewPeriods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeriods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodContains(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	type args struct {
		t int64
		p Period
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				t: -1,
				p: newPeriod(30, 40),
			},
			want: false,
		},
		{
			args: args{
				t: 0,
				p: newPeriod(1, -1),
			},
			want: false,
		},
		{
			args: args{
				t: 35,
				p: newPeriod(30, -1),
			},
			want: true,
		},
		{
			args: args{
				t: 35,
				p: newPeriod(30, 40),
			},
			want: true,
		},
		{
			args: args{
				t: 35,
				p: newPeriod(30, 32),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodContains(tt.args.t, tt.args.p); got != tt.want {
				t.Errorf("PeriodContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodUnion(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		p1 Period
		p2 Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		// p1 p2 all infinity
		{
			args: args{
				p1: newPeriod(20, -1),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(20, -1),
		},
		{
			args: args{
				p1: newPeriod(10, -1),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(10, -1),
		},
		// p1 infinity
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(10, 20),
			},
			want: newPeriods(30, -1, 10, 20),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(20, 30),
			},
			want: newPeriods(20, -1),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(25, 35),
			},
			want: newPeriods(25, -1),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(10, 35),
			},
			want: newPeriods(10, -1),
		},
		// p2 infinity
		{
			args: args{
				p1: newPeriod(10, 20),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(30, -1, 10, 20),
		},
		{
			args: args{
				p1: newPeriod(20, 30),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(20, -1),
		},
		{
			args: args{
				p1: newPeriod(25, 35),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(25, -1),
		},
		{
			args: args{
				p1: newPeriod(10, 35),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(10, -1),
		},
		// other
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(10, 20),
			},
			want: newPeriods(30, 40, 10, 20),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(20, 30),
			},
			want: newPeriods(20, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(25, 35),
			},
			want: newPeriods(25, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(35, 40),
			},
			want: newPeriods(30, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(10, 60),
			},
			want: newPeriods(10, 60),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(30, 60),
			},
			want: newPeriods(30, 60),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(35, 60),
			},
			want: newPeriods(30, 60),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(40, 60),
			},
			want: newPeriods(30, 60),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(45, 60),
			},
			want: newPeriods(30, 40, 45, 60),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodUnion(tt.args.p1, tt.args.p2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodIntersect(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		p1 Period
		p2 Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		// p1,p2 negative
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(10, -1),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				p1: newPeriod(10, -1),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(30, -1),
		},
		// p1 negative
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(10, 20),
			},
			want: newPeriods(),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(20, 30),
			},
			want: newPeriods(30, 30),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(25, 35),
			},
			want: newPeriods(30, 35),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(30, 35),
			},
			want: newPeriods(30, 35),
		},
		{
			args: args{
				p1: newPeriod(30, -1),
				p2: newPeriod(40, 50),
			},
			want: newPeriods(40, 50),
		},
		// p2 negative
		{
			args: args{
				p1: newPeriod(10, 20),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(),
		},
		{
			args: args{
				p1: newPeriod(20, 30),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(30, 30),
		},
		{
			args: args{
				p1: newPeriod(25, 35),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(30, 35),
		},
		{
			args: args{
				p1: newPeriod(30, 35),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(30, 35),
		},
		{
			args: args{
				p1: newPeriod(40, 50),
				p2: newPeriod(30, -1),
			},
			want: newPeriods(40, 50),
		},
		// other
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(10, 20),
			},
			want: newPeriods(),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(20, 30),
			},
			want: newPeriods(30, 30),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(25, 35),
			},
			want: newPeriods(30, 35),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(25, 40),
			},
			want: newPeriods(30, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(25, 45),
			},
			want: newPeriods(30, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(35, 45),
			},
			want: newPeriods(35, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(40, 45),
			},
			want: newPeriods(40, 40),
		},
		{
			args: args{
				p1: newPeriod(30, 40),
				p2: newPeriod(41, 45),
			},
			want: newPeriods(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodIntersect(tt.args.p1, tt.args.p2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodsContains(t *testing.T) {
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		t  int64
		ps []Period
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				t:  0,
				ps: newPeriods(),
			},
			want: false,
		},
		{
			args: args{
				t:  30,
				ps: newPeriods(0, -1),
			},
			want: true,
		},
		{
			args: args{
				t:  30,
				ps: newPeriods(10, 20, 20, 30),
			},
			want: true,
		},
		{
			args: args{
				t:  30,
				ps: newPeriods(0, 5, 5, 10, 10, 20, 25, 28),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodsContains(tt.args.t, tt.args.ps); got != tt.want {
				t.Errorf("PeriodsContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddPeriodToResultSet(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		p  Period
		rs []Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		// p -> rs => result
		// [10, 20] -> [] => [10, 20]
		// [10, -1] -> [] => [10, -1]
		// [10, 20] -> [30, 40] => [30, 40] [10, 20]
		// [20, 30] -> [30, 40] => [20, 40]
		// [25, 35] -> [30, 40] => [25, 40]
		// [40, 50] -> [30, 40] => [30, 50]
		// [50, 60] -> [30, 40] =>  [30, 40] [50, 60]
		// [10, -1] -> [30, 40] [50, 60] => [10, -1]
		{
			args: args{
				p:  newPeriod(10, 20),
				rs: newPeriods(),
			},
			want: newPeriods(10, 20),
		},
		{
			args: args{
				p:  newPeriod(10, -1),
				rs: newPeriods(),
			},
			want: newPeriods(10, -1),
		},
		{
			args: args{
				p:  newPeriod(10, 20),
				rs: newPeriods(30, 40),
			},
			want: newPeriods(30, 40, 10, 20),
		},
		{
			args: args{
				p:  newPeriod(20, 30),
				rs: newPeriods(30, 40),
			},
			want: newPeriods(20, 40),
		},
		{
			args: args{
				p:  newPeriod(25, 35),
				rs: newPeriods(30, 40),
			},
			want: newPeriods(25, 40),
		},
		{
			args: args{
				p:  newPeriod(40, 50),
				rs: newPeriods(30, 40),
			},
			want: newPeriods(30, 50),
		},
		{
			args: args{
				p:  newPeriod(50, 60),
				rs: newPeriods(30, 40),
			},
			want: newPeriods(30, 40, 50, 60),
		},
		{
			args: args{
				p:  newPeriod(10, -1),
				rs: newPeriods(30, 40, 50, 60),
			},
			want: newPeriods(10, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddPeriodToResultSet(tt.args.p, tt.args.rs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddPeriodToResultSet() args %v, got %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestPeriodsUnion(t *testing.T) {
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		a []Period
		b []Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		// normal
		// [10, 20] [30, 40] => [10, 20] [30, 40]
		// [20, 30] [30, 40] => [20, 40]
		// [25, 35] [30, 40] => [25, 40]
		// [30, 40] [30, 40] => [30, 40]
		// [35, 45] [30, 40] => [30, 45]
		// [40, 50] [30, 40] => [30, 50]
		// [50, 60] [30, 40] => [50, 60] [30, 40]
		{
			args: args{
				a: newPeriods(10, 20, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(10, 20, 30, 40),
		},
		{
			args: args{
				a: newPeriods(20, 30, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(20, 40),
		},
		{
			args: args{
				a: newPeriods(25, 35, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(25, 40),
		},
		{
			args: args{
				a: newPeriods(30, 40, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, 40),
		},
		{
			args: args{
				a: newPeriods(35, 45, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, 45),
		},
		{
			args: args{
				a: newPeriods(40, 50, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, 50),
		},
		{
			args: args{
				a: newPeriods(50, 60, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(50, 60, 30, 40),
		},

		// rigth positive infinity
		// [10, 20] [30, -1] => [10, 20] [30, -1]
		// [20, 30] [30, -1] => [20, -1]
		// [25, 35] [30, -1] => [25, -1]
		// [30, 40] [30, -1] => [30, -1]
		// [35, 45] [30, -1] => [30, -1]
		// [40, 50] [30, -1] => [30, -1]
		// [50, 60] [30, -1] => [30, -1]
		{
			args: args{
				a: newPeriods(10, 20, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(10, 20, 30, -1),
		},
		{
			args: args{
				a: newPeriods(20, 30, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(20, -1),
		},
		{
			args: args{
				a: newPeriods(25, 35, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(25, -1),
		},
		{
			args: args{
				a: newPeriods(30, 40, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(35, 45, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(40, 50, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(50, 60, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		// left positive infinity
		// [10, -1] [30, 40] => [10, -1]
		// [20, -1] [30, 40] => [20, -1]
		// [25, -1] [30, 40] => [25, -1]
		// [30, -1] [30, 40] => [30, -1]
		// [35, -1] [30, 40] => [30, -1]
		// [40, -1] [30, 40] => [30, -1]
		// [50, -1] [30, 40] => [50, -1] [30, 40]
		{
			args: args{
				a: newPeriods(10, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(10, -1),
		},
		{
			args: args{
				a: newPeriods(20, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(20, -1),
		},
		{
			args: args{
				a: newPeriods(25, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(25, -1),
		},
		{
			args: args{
				a: newPeriods(30, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(35, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(40, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(50, -1, 30, 40),
				b: newPeriods(),
			},
			want: newPeriods(50, -1, 30, 40),
		},
		// all infinity
		// [10, -1] [30, -1] => [10, -1]
		// [30, -1] [30, -1] => [30, -1]
		// [40, -1] [30, -1] => [30, -1]
		{
			args: args{
				a: newPeriods(10, -1, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(10, -1),
		},
		{
			args: args{
				a: newPeriods(30, -1, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		{
			args: args{
				a: newPeriods(40, -1, 30, -1),
				b: newPeriods(),
			},
			want: newPeriods(30, -1),
		},
		// other
		// [10, 20] [20, 30] [30, 40] [40, 50] [50, 60] => [10, 60]
		// [10, -1] => [10, -1]
		// [10, 20] [10, -1] [0, 5] => [10, -1] [0, 5]
		// [10, 20] [30, 40] [50, 60] => [10, 20] [30, 40] [50, 60]
		// [30, 40] [10, 60] => [10, 60]
		{
			args: args{
				a: newPeriods(10, 20, 20, 30, 30, 40, 40, 50, 50, 60),
				b: newPeriods(),
			},
			want: newPeriods(10, 60),
		},
		{
			args: args{
				a: newPeriods(10, -1),
				b: newPeriods(),
			},
			want: newPeriods(10, -1),
		},
		{
			args: args{
				a: newPeriods(10, 20, 10, -1, 0, 5),
				b: newPeriods(),
			},
			want: newPeriods(10, -1, 0, 5),
		},
		{
			args: args{
				a: newPeriods(10, 20, 30, 40, 50, 60),
				b: newPeriods(),
			},
			want: newPeriods(10, 20, 30, 40, 50, 60),
		},
		{
			args: args{
				a: newPeriods(30, 40, 10, 60),
				b: newPeriods(),
			},
			want: newPeriods(10, 60),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodsUnion(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodsUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodsIntersect(t *testing.T) {
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		a []Period
		b []Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		{
			args: args{
				a: newPeriods(0, 10, 30, -1),
				b: newPeriods(5, 9, 20, -1),
			},
			want: newPeriods(5, 9, 30, -1),
		},
		{
			args: args{
				a: newPeriods(30, 40, 50, -1),
				b: newPeriods(10, 20, 30, 30, 32, 34, 38, 40, 60, -1),
			},
			want: newPeriods(30, 30, 32, 34, 38, 40, 60, -1),
		},
		{
			args: args{
				a: newPeriods(30, 40, 50, -1),
				b: newPeriods(),
			},
			want: newPeriods(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodsIntersect(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodsIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodsSort(t *testing.T) {
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		a []Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		{
			args: args{
				a: newPeriods(3, 4, 1, 2, 8, 9, 20, 30, 12, 14),
			},
			want: newPeriods(1, 2, 3, 4, 8, 9, 12, 14, 20, 30),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := make([]Period, len(tt.args.a))
			copy(tmp, tt.args.a)
			sort.Sort(Periods(tmp))
			got := []Period(tmp)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestPeriodsSort() args: %v, got: %v, want %v", tt.args.a, got, tt.want)
			}
		})
	}

}

func TestPeriodPartition(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	type args struct {
		p        Period
		interval int64
	}
	tests := []struct {
		name string
		args args
		want map[int64]Period
	}{
		{
			args: args{
				p:        newPeriod(12, -1),
				interval: 10,
			},
			want: map[int64]Period{
				10: newPeriod(12, -1),
			},
		},
		{
			args: args{
				p:        newPeriod(12, 18),
				interval: 10,
			},
			want: map[int64]Period{
				10: newPeriod(12, 18),
			},
		},
		{
			args: args{
				p:        newPeriod(12, 27),
				interval: 10,
			},
			want: map[int64]Period{
				10: newPeriod(12, 19),
				20: newPeriod(20, 27),
			},
		},
		{
			args: args{
				p:        newPeriod(12, 39),
				interval: 10,
			},
			want: map[int64]Period{
				10: newPeriod(12, 19),
				20: newPeriod(20, 29),
				30: newPeriod(30, 39),
			},
		},
		{
			args: args{
				p:        newPeriod(10, 20),
				interval: 10,
			},
			want: map[int64]Period{
				10: newPeriod(10, 19),
				20: newPeriod(20, 20),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodPartition(tt.args.p, tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodPartition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodsPartition(t *testing.T) {
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		ps       []Period
		interval int64
	}
	tests := []struct {
		name string
		args args
		want map[int64][]Period
	}{
		{
			args: args{
				ps:       newPeriods(12, 18, 23, 47),
				interval: 10,
			},
			want: map[int64][]Period{
				10: newPeriods(12, 18),
				20: newPeriods(23, 29),
				30: newPeriods(30, 39),
				40: newPeriods(40, 47),
			},
		},
		{
			args: args{
				ps:       newPeriods(12, 18, 23, 47, 15, 27, 12, 15),
				interval: 10,
			},
			want: map[int64][]Period{
				10: newPeriods(12, 18, 15, 19, 12, 15),
				20: newPeriods(23, 29, 20, 27),
				30: newPeriods(30, 39),
				40: newPeriods(40, 47),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodsPartition(tt.args.ps, tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodsPartition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodComplement(t *testing.T) {
	var newPeriod = func(st, et int64) Period {
		p, _ := NewPeriod(st, et)
		return p
	}
	var newPeriods = func(se ...int64) []Period {
		ps, _ := NewPeriods(se...)
		return ps
	}
	type args struct {
		a Period
		b Period
	}
	tests := []struct {
		name string
		args args
		want []Period
	}{
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(10, -1),
			},
			want: newPeriods(10, 29),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(30, -1),
			},
			want: newPeriods(),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(40, -1),
			},
			want: newPeriods(),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(10, 20),
			},
			want: newPeriods(10, 20),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(20, 30),
			},
			want: newPeriods(20, 29),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(25, 35),
			},
			want: newPeriods(25, 29),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(30, 40),
			},
			want: newPeriods(),
		},
		{
			args: args{
				a: newPeriod(30, -1),
				b: newPeriod(40, 50),
			},
			want: newPeriods(),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(10, -1),
			},
			want: newPeriods(10, 29, 41, -1),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(30, -1),
			},
			want: newPeriods(41, -1),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(35, -1),
			},
			want: newPeriods(41, -1),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(40, -1),
			},
			want: newPeriods(41, -1),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(45, -1),
			},
			want: newPeriods(45, -1),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(10, 20),
			},
			want: newPeriods(10, 20),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(20, 30),
			},
			want: newPeriods(20, 29),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(35, 45),
			},
			want: newPeriods(41, 45),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(40, 50),
			},
			want: newPeriods(41, 50),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(50, 60),
			},
			want: newPeriods(50, 60),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(10, 60),
			},
			want: newPeriods(10, 29, 41, 60),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(30, 60),
			},
			want: newPeriods(41, 60),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(35, 60),
			},
			want: newPeriods(41, 60),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(40, 60),
			},
			want: newPeriods(41, 60),
		},
		{
			args: args{
				a: newPeriod(30, 40),
				b: newPeriod(50, 60),
			},
			want: newPeriods(50, 60),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PeriodComplement(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeriodComplement() = %v, want %v", got, tt.want)
			}
		})
	}
}
