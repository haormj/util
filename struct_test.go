package util

import (
	"encoding/json"
	"testing"
)

type Hello struct {
	Name   *string
	Age    int
	Weight int64
	World  World
}

type World struct {
	Sex   *string
	Habit string
}

func TestGetStructField(t *testing.T) {
	name := "hao"
	sex := "male"
	hello := &Hello{
		Name: &name,
		Age:  17,
		World: World{
			Sex: &sex,
		},
	}
	t.Log(GetStructField(hello, "Name").String())
	t.Log(GetStructField(hello, "Age").Int())
	t.Log(GetStructField(hello, "World", "Sex").String())
	t.Log(GetStructField(hello, "Ran").String())
}

func TestSetStructField(t *testing.T) {
	hello := &Hello{}
	t.Log(SetStructField(hello, "hao", "Name"))
	t.Log(SetStructField(hello, 18, "Age"))
	t.Log(SetStructField(hello, "male", "World", "Sex"))
	t.Log(SetStructField(hello, "eat", "World", "Habit"))
	t.Log(SetStructField(hello, 20, "Weight"))
	helloBytes, err := json.Marshal(hello)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(helloBytes))
}
