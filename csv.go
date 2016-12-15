package csv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ReadFile(path string, v interface{}) (err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return Decode(f, v)
}

func ReadByte(data []byte, v interface{}) (err error) {
	in := bytes.NewReader(data)
	return Decode(in, v)
}

func ReadString(str string, v interface{}) (err error) {
	in := strings.NewReader(str)
	return Decode(in, v)
}

//func ReadStdin(){}

// decode
func Decode(in io.Reader, v interface{}) (err error) {
	m, err := csv2map(in)
	if err != nil {
		return
	}
	return map2list(m, v)
}

// io => []map[string]string
func csv2map(in io.Reader) (ret []map[string]string, err error) {
	r := csv.NewReader(in)
	var records [][]string
	records, err = r.ReadAll()
	if err != nil {
		err = fmt.Errorf("csv2map error: read fail: %v", err)
		return
	}
	name := records[0]
	for i := 1; i < len(records); i++ {
		m := make(map[string]string, len(records))
		for j, k := range name {
			value := records[i][j]
			if value == "" {
				continue
			}
			if strings.HasPrefix(value, "#") {
				continue
			}
			if k == "" {
				continue
			}
			m[k] = value
		}
		if len(m) == 0 {
			continue
		}
		ret = append(ret, m)
	}
	return
}

// []map[string]string => []interface{}
func map2list(m []map[string]string, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("csv: Decode error: must ptr")
	}

	for {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			continue
		}
		break
	}
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("must slice, type=%v", rv.Kind())
	}
	ret := reflect.MakeSlice(rv.Type(), 0, 0)
	for _, i := range m {
		d := reflect.New(rv.Type().Elem())
		if err := map2struct(i, d); err != nil {
			return err
		}
		ret = reflect.Append(ret, d.Elem())
	}
	rv.Set(ret)
	return
}

// map[string]string => struct
func map2struct(m map[string]string, v reflect.Value) (err error) {
	v = v.Elem()
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.New(v.Type().Elem()))
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("not struct: %v", v.Kind())
	}
	return value(m, v)
}

func getMapValue(m map[string]string, key string) (ret string) {
	k := strings.ToLower(key)
	ret, _ = m[k]
	return
}

func value(m map[string]string, v reflect.Value) (err error) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		t := v.Type().Field(i).Tag.Get("csv")
		k := v.Type().Field(i).Name
		if !f.CanSet() {
			continue
		}
		d := getMapValue(m, k)
		switch f.Kind() {
		default:
			err = fmt.Errorf("not support")
		case reflect.Struct:
			if t == "extends" {
				err = value(m, f)
			} else {
				err = valueStruct(d, f)
			}
		case reflect.String:
			err = valueString(d, f)
		case reflect.Int:
			err = valueInt(d, f)
		case reflect.Int32:
			err = valueIntx(d, f, 32)
		case reflect.Int64:
			err = valueIntx(d, f, 64)
			//TODO: support bool
			//case reflect.Bool:
		}
		if err != nil {
			return
		}
	}
	return
}

func valueStruct(s string, v reflect.Value) (err error) {
	b := []byte(s)
	err = json.Unmarshal(b, v.Addr().Interface())
	return
}
func valueString(s string, v reflect.Value) (err error) {
	v.Set(reflect.ValueOf(s))
	return
}
func valueInt(s string, v reflect.Value) (err error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(n))
	return
}
func valueIntx(s string, v reflect.Value, x int) (err error) {
	n, err := strconv.ParseInt(s, 10, x)
	if err != nil {
		return
	}
	v.Set(reflect.ValueOf(n))
	return
}
