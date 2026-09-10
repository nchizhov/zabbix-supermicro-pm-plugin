package ipmi

import (
	"fmt"
	"reflect"
)

func ReflectHandler(data *Data, params []string, handler string) (any, error) {
	v := reflect.ValueOf(data)

	method := v.MethodByName(handler)

	if !method.IsValid() {
		return nil, fmt.Errorf("uknown handler %s", handler)
	}

	args := []reflect.Value{}
	if len(params) != 0 {
		for _, param := range params {
			args = append(args, reflect.ValueOf(param))
		}
	}

	var resVal any
	var errVal error
	var results []reflect.Value
	results = method.Call(args)
	if results[0].IsValid() {
		resVal = results[0].Interface()
	}
	if !results[1].IsNil() {
		errVal = results[1].Interface().(error)
	}
	return resVal, errVal
}
