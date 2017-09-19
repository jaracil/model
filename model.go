package model

import (
	"errors"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

var dataInstance *js.Object
var vmInstance *js.Object

func chkInit() {
	if dataInstance == nil {
		dataInstance = js.Global.Get("Object").New()
		dataInstance.Set("data", js.Global.Get("Object").New())
		dataInstance.Set("methods", js.Global.Get("Object").New())
		dataInstance.Set("computed", js.Global.Get("Object").New())
		dataInstance.Set("watch", js.Global.Get("Object").New())
	}
}

func Object(path string) *js.Object {
	chkInit()
	obj := dataInstance.Get("data")
	chunks := strings.Split(path, ".")
	for _, chunk := range chunks {
		tmp := obj.Get(chunk)
		if tmp == js.Undefined {
			tmp = js.Global.Get("Object").New()
			obj.Set(chunk, tmp)
		}
		obj = tmp
	}
	return obj
}

func AddMethod(key string, val interface{}) error {
	chkInit()
	if dataInstance.Get("methods").Get(key) != js.Undefined {
		return errors.New("Duplicated key in AddMethod")
	}
	dataInstance.Get("methods").Set(key, val)
	return nil
}

func AddComputed(key string, val interface{}) error {
	chkInit()
	if dataInstance.Get("computed").Get(key) != js.Undefined {
		return errors.New("Duplicated key in AddComputed")
	}
	dataInstance.Get("computed").Set(key, val)
	return nil
}

func AddWatch(key string, cbf interface{}) error {
	chkInit()
	last := dataInstance.Get("watch").Get(key)
	f := func(val, prev interface{}) {
		if !(last == js.Undefined) {
			last.Invoke(val, prev)
		}
		cbf.(func(interface{}, interface{}))(val, prev)
	}
	dataInstance.Get("watch").Set(key, f)
	return nil
}

func Ready(mount string) {
	if mount != "" {
		dataInstance.Set("el", mount)
	}
	vmInstance = js.Global.Get("Vue").New(dataInstance)
	js.Global.Set("vm", vmInstance)
}

func InternalDataModel() *js.Object {
	return dataInstance
}

func Vm() *js.Object {
	return vmInstance
}

func Set(obj *js.Object, key interface{}, value interface{}) {
	js.Global.Get("Vue").Call("set", obj, key, value)
}

func Delete(obj *js.Object, key interface{}) {
	js.Global.Get("Vue").Call("delete", obj, key)
}
