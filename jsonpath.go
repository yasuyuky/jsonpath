package jsonpath

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Slice struct{ Start, Stop int }

type pathTypeMismatch struct{}
type notFound struct{}
type unknownPathType struct{}

func recursiveGet(data interface{}, path []interface{}) interface{} {

	if len(path) == 0 {
		switch data.(type) {
		case string:
			return data.(string)
		case float64:
			return data.(float64)
		case bool:
			return data.(bool)
		case nil:
			return nil
		case []interface{}:
			return data
		case map[string]interface{}:
			return data
		}
	}

	switch path[0].(type) {

	case string:
		switch data.(type) {
		case map[string]interface{}:
			for k, v := range data.(map[string]interface{}) {
				if k == path[0].(string) {
					return recursiveGet(v, path[1:])
				}
			}
			return notFound{}
		default:
			return pathTypeMismatch{}
		}

	case int:
		switch data.(type) {
		case []interface{}:
			for i, v := range data.([]interface{}) {
				if i == path[0].(int) {
					return recursiveGet(v, path[1:])
				}
			}
			return notFound{}
		default:
			return pathTypeMismatch{}
		}

	case func(int, interface{}) bool:
		switch data.(type) {
		case []interface{}:
			ret := make([]interface{}, 0)
			for i, v := range data.([]interface{}) {
				if path[0].(func(int, interface{}) bool)(i, v) {
					val := recursiveGet(v, path[1:])
					switch val.(type) {
					case pathTypeMismatch, notFound:
						continue
					case unknownPathType:
						return val
					}
					ret = append(ret, val)
				}
			}
			return ret
		}

	case func(string, interface{}) bool:
		switch data.(type) {
		case map[string]interface{}:
			ret := make([]interface{}, 0)
			for k, v := range data.(map[string]interface{}) {
				if path[0].(func(string, interface{}) bool)(k, v) {
					val := recursiveGet(v, path[1:])
					switch val.(type) {
					case pathTypeMismatch, notFound:
						continue
					case unknownPathType:
						return val
					}
					ret = append(ret, val)
				}
			}
			return ret
		}

	case Slice:
		switch data.(type) {
		case []interface{}:
			ret := make([]interface{}, 0)
			start := path[0].(Slice).Start
			stop := path[0].(Slice).Stop
			for i, v := range data.([]interface{}) {
				if start <= i && i < stop {
					ret = append(ret, recursiveGet(v, path[1:]))
				}
			}
			return ret
		}
	}

	return unknownPathType{}

}

func Get(decoded interface{}, path []interface{}, defaultValue interface{}) (interface{}, error) {
	val := recursiveGet(decoded, path)
	switch val.(type) {
	case notFound:
		return defaultValue, fmt.Errorf("not found")
	case unknownPathType:
		return defaultValue, fmt.Errorf("unknown path type")
	case pathTypeMismatch:
		return defaultValue, fmt.Errorf("mismatched path type")
	}
	return val, nil
}

func GetString(decoded interface{}, path []interface{}, defaultValue string) (string, error) {
	value, err := Get(decoded, path, nil)
	if err != nil {
		return defaultValue, err
	}
	switch value.(type) {
	case string:
		return value.(string), nil
	}
	return defaultValue, fmt.Errorf("unexpected type")
}

func GetNumber(decoded interface{}, path []interface{}, defaultValue float64) (float64, error) {
	value, err := Get(decoded, path, nil)
	if err != nil {
		return defaultValue, err
	}
	switch value.(type) {
	case float64:
		return value.(float64), nil
	}
	return defaultValue, fmt.Errorf("unexpected type")
}

func GetBool(decoded interface{}, path []interface{}, defaultValue bool) (bool, error) {
	value, err := Get(decoded, path, nil)
	if err != nil {
		return defaultValue, err
	}
	switch value.(type) {
	case bool:
		return value.(bool), nil
	}
	return defaultValue, fmt.Errorf("unexpected type")
}

func DecodeString(s string) (interface{}, error) {
	return DecodeReader(strings.NewReader(s))
}

func DecodeReader(r io.Reader) (interface{}, error) {
	var data interface{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&data)
	return data, err
}

func Read(r io.Reader, path []interface{}, defaultValue interface{}) (interface{}, error) {
	data, err := DecodeReader(r)
	if err != nil {
		return nil, err
	}
	return Get(data, path, defaultValue)
}

func ReadString(r io.Reader, path []interface{}, defaultValue string) (string, error) {
	data, err := DecodeReader(r)
	if err != nil {
		return defaultValue, err
	}
	return GetString(data, path, defaultValue)
}

func ReadNumber(r io.Reader, path []interface{}, defaultValue float64) (float64, error) {
	data, err := DecodeReader(r)
	if err != nil {
		return defaultValue, err
	}
	return GetNumber(data, path, defaultValue)
}

func ReadBool(r io.Reader, path []interface{}, defaultValue bool) (bool, error) {
	data, err := DecodeReader(r)
	if err != nil {
		return defaultValue, err
	}
	return GetBool(data, path, defaultValue)
}
