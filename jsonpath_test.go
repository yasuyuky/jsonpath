package jsonpath

import (
	"strings"
	"testing"
)

var (
	test_json = `
{"foo": null,
 "baz": ["1", 2, null, ["10",20], true, {"a": "b"}]
}
`
	wrong_json = `
{"foo": null,
 "baz": ["1", 2, null, ["10", 20]]
`
)

func TestGet(t *testing.T) {
	data, err := DecodeString(test_json)
	if err != nil {
		t.Errorf("json decode err")
	}
	v, err := Get(data, []interface{}{"baz", 5}, nil)
	if err != nil {
		t.Errorf("json read err")
	}
	if v == nil {
		t.Errorf("path not found")
	}
}

func TestRead(t *testing.T) {
	v, err := Read(strings.NewReader(test_json), []interface{}{"baz", 5}, nil)
	if err != nil {
		t.Errorf("json read err")
	}
	if v == nil {
		t.Errorf("path not found")
	}
	v, err = Read(strings.NewReader(test_json), []interface{}{"baz", 2}, 0)
	if v != nil {
		t.Errorf("value must be nil")
	}
	v, err = Read(strings.NewReader(test_json), []interface{}{"bar", 3, 1}, nil)
	if err == nil {
		t.Errorf("path must be not found")
	}
}

func TestReadString(t *testing.T) {
	v, err := ReadString(strings.NewReader(test_json), []interface{}{"baz", 0}, "")
	if v != "1" {
		t.Errorf("path{\"baz\", 0} must be \"1\"")
	}
	v, err = ReadString(strings.NewReader(test_json), []interface{}{"baz", 1}, "")
	if err == nil {
		t.Errorf("ReadString for path{\"baz\", 1} must be error")
	}
	v, err = ReadString(strings.NewReader(test_json), []interface{}{1}, "")
	if err == nil {
		t.Errorf("path must be mismatched")
	}
}

func TestReadNumber(t *testing.T) {
	v, err := ReadNumber(strings.NewReader(test_json), []interface{}{"baz", 1}, 10)
	if v != 2.0 {
		t.Errorf("path{\"baz\", 1} must be 2.0")
	}
	v, err = ReadNumber(strings.NewReader(test_json), []interface{}{"baz", 0}, 10)
	if err == nil {
		t.Errorf("ReadNumber for path{\"baz\", 1} must be error")
	}
	v, err = ReadNumber(strings.NewReader(test_json), []interface{}{"baz", 10}, 0)
	if err == nil {
		t.Errorf("path must be not found")
	}
}

func TestReadBool(t *testing.T) {
	v, err := ReadBool(strings.NewReader(test_json), []interface{}{"baz", 4}, false)
	if v != true {
		t.Errorf("path{\"baz\", 4} must be true")
	}
	v, err = ReadBool(strings.NewReader(test_json), []interface{}{"baz", 3}, false)
	if err == nil {
		t.Errorf("ReadBool for path{\"baz\", 1} must be error")
	}
	v, err = ReadBool(strings.NewReader(test_json), []interface{}{"baz", "baz"}, false)
	if err == nil {
		t.Errorf("path must be mismatched")
	}
}

func TestFilter(t *testing.T) {
	select_all_elements := func(int, interface{}) bool { return true }
	v, err := Read(strings.NewReader(test_json), []interface{}{"baz", select_all_elements, 0}, nil)
	if v == nil {
		t.Errorf("path not found")
	}

	v, err = Read(strings.NewReader(test_json), []interface{}{"baz", select_all_elements, 0.1}, nil)
	if err == nil {
		t.Errorf("must be error")
	}

	key_contains_a := func(k string, v interface{}) bool { return strings.Contains(k, "a") }
	v, err = Read(strings.NewReader(test_json), []interface{}{key_contains_a, 0}, nil)
	if v == nil {
		t.Errorf("path not found")
	}

	v, err = Read(strings.NewReader(test_json), []interface{}{key_contains_a, "a"}, nil)
	if err != nil {
		t.Errorf("must not be error")
	}

	v, err = Read(strings.NewReader(test_json), []interface{}{key_contains_a, 0.1}, nil)
	if err == nil {
		t.Errorf("must be error")
	}
}

func TestSlice(t *testing.T) {
	v, err := Read(strings.NewReader(test_json), []interface{}{"baz", Slice{1, 4}, 1}, nil)
	if err != nil {
		t.Errorf("json read err")
	}
	if v == nil {
		t.Errorf("path not found")
	}
}

func TestWrongPath(t *testing.T) {
	d, err := Read(strings.NewReader(test_json), []interface{}{"baz", 0.1}, nil)
	if err == nil {
		t.Errorf("it must be error")
	}
	if d != nil {
		t.Errorf("it must be null")
	}
}

func TestWrongJson(t *testing.T) {
	d, err := Read(strings.NewReader(wrong_json), []interface{}{"baz", Slice{1, 4}, 1}, nil)
	if err == nil {
		t.Errorf("it must be error")
	}
	if d != nil {
		t.Errorf("it must be null")
	}

	d, err = ReadString(strings.NewReader(wrong_json), []interface{}{"baz", Slice{1, 4}, 1}, "")
	if err == nil {
		t.Errorf("it must be error")
	}
	if d != "" {
		t.Errorf("it must be default value")
	}

	d, err = ReadNumber(strings.NewReader(wrong_json), []interface{}{"baz", Slice{1, 4}, 1}, 0)
	if err == nil {
		t.Errorf("it must be error")
	}
	if d != 0. {
		t.Errorf("it must be default value")
	}

	d, err = ReadBool(strings.NewReader(wrong_json), []interface{}{"baz", Slice{1, 4}, 1}, false)
	if err == nil {
		t.Errorf("it must be error")
	}
	if d != false {
		t.Errorf("it must be default value")
	}

}
