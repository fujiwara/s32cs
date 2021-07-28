package s32cs_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fujiwara/s32cs"
)

var sdfs = []struct {
	Name   string
	Src    []byte
	Expect []byte
}{
	{
		Name:   "add",
		Src:    []byte(`{"id":"123","type":"add","fields":{"foo":"bar","bar":["A","B"]}}`),
		Expect: []byte(`{"id":"123","type":"add","fields":{"foo":"bar","bar":["A","B"]}}`),
	},
	{
		Name:   "add2",
		Src:    []byte(`{"id":"xxx123","type":"add","fields":{"foo":"bar","bar":["A","B"]}}`),
		Expect: []byte(`{"id":"xxx123","type":"add","fields":{"foo":"bar","bar":["A","B"]}}`),
	},
	{
		Name:   "delete",
		Src:    []byte(`{"id":"123","type":"delete"}`),
		Expect: []byte(`{"id":"123","type":"delete"}`),
	},
	{
		Name:   "broken",
		Src:    []byte(`{"id":"123","fields":{"foo":"bar","bar":["A","B"]`),
		Expect: nil,
	},
	{
		Name:   "no type",
		Src:    []byte(`{"id":"123","fields":{"foo":"bar","bar":["A","B"]}`),
		Expect: nil,
	},
	{
		Name:   "no id",
		Src:    []byte(`{type:"add","fields":{"foo":"bar","bar":["A","B"]}`),
		Expect: nil,
	},
	{
		Name:   "invalid type",
		Src:    []byte(`{"id":"123","type":"xxx","fields":{"foo":"bar","bar":["A","B"]}`),
		Expect: nil,
	},
	{
		Name:   "invalid char removed",
		Src:    []byte(`{"id":"123","type":"add","fields":{"foo":"bar\u001d\u0000bar","bar":["AA","B\uFFFFB"]}}`),
		Expect: []byte(`{"id":"123","type":"add","fields":{"foo":"barbar","bar":["AA","BB"]}}`),
	},
}

func TestSDFUnmarshal(t *testing.T) {
	for _, c := range sdfs {
		var record, expect s32cs.SDFRecord
		err := json.Unmarshal(c.Src, &record)
		if err == nil {
			err = record.Validate()
		}
		if c.Expect == nil {
			if err == nil {
				t.Errorf("%s must be failed but no error returned", c.Name)
			}
			continue
		}
		if err != nil {
			t.Error("testcase", c.Name, "errored with", err)
			continue
		}
		json.Unmarshal(c.Expect, &expect)
		if !reflect.DeepEqual(&record, &expect) {
			t.Errorf("%s unexpected marshal expected %#v got %#v", c.Name, &expect, &record)
		}
	}
}
