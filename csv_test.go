package csv

import (
	"encoding/json"
	"testing"
)

func TestAll(t *testing.T) {
	var err error
	var data = `id,name,animal
#int,#str,#animal
123,S2776,"{""Name"":""Platypus"",""Order"":""Monotremata""}"
456,A1889,"{""Name"":""Quoll"",""Order"":""Dasyuromorphia""}"`

	type Base struct {
		Id int
	}
	type Animal struct {
		Name  string
		Order string
	}
	type Server struct {
		Base   `csv:"extends"`
		Name   string
		Animal Animal
	}

	var r1 []Server
	err = ReadFile("csv_test.data", &r1)
	if err != nil {
		t.Errorf("TestAll error: %v", err)
	}

	var r2 []*Server
	err = ReadString(data, &r2)
	if err != nil {
		t.Errorf("TestAll error: %v", err)
	}

	b1, _ := json.Marshal(r1)
	b2, _ := json.Marshal(r2)
	if string(b1) != string(b2) {
		t.Errorf("TestAll error!")
	}
	if string(b1) != `[{"Id":123,"Name":"S2776","Animal":{"Name":"Platypus","Order":"Monotremata"}},{"Id":456,"Name":"A1889","Animal":{"Name":"Quoll","Order":"Dasyuromorphia"}}]` {
		t.Errorf("TestAll error!")
	}
}
