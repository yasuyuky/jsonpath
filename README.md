jsonpath for go
===============

[![Build Status](https://travis-ci.org/yasuyuky/jsonpath.png?branch=master)](https://travis-ci.org/yasuyuky/jsonpath)
[![Coverage Status](https://coveralls.io/repos/yasuyuky/jsonpath/badge.png?branch=master)](https://coveralls.io/r/yasuyuky/jsonpath?branch=master)

A simple and elastic access interface to json object for golang.

Install
=======

    go get github.com/yasuyuky/jsonpath

Usage
=====

Basic usage (decode and get)
--------------------------

```go
import (
	"github.com/yasuyuky/jsonpath"
)

// fist you should decode string
data, err := jsonpath.DecodeString(json_string)

// or io.Reader
data, err := jsonpath.DecodeReader(json_reader)

// then you can get element using jsonpath.Get
// 1st arg is decoded data
// 2nd arg is path([]interface{}) for element
//   it contain string(for object), int(for array)
// 3rd arg is default value
v, err := jsonpath.Get(data, []interface{}{"foo", 5}, nil)
```

Get with type
-------------

```go
// you can also use GetString/GetNumber/GetBool
// these function can get element with type assertion
s, err := jsonpath.GetString(data, []interface{}{"bar", "baz"}, "")
```

Read directly
-------------
```go
// or you can also read directly from io.Reader
v, err := jsonpath.Read(json_reader, []interface{}{"baz", 5}, nil)
```

Slice for array
---------------

```go
// jsonpath.Slice{start, stop} for get range 'start <= x < stop'
// index starts from 0
a, err := jsonpath.Get(data, []interface{}{"baz", jsonpath.Slice{1, 4}}, nil)
```

Using filter function
---------------------
you can use `func(string/int, interface{}) bool` for complex filtering.

see jsonpath_test.go for more detailed usage

License
=======

2-clause BSD license
