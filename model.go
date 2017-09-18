package model

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var dataInstance *js.Object

func chkInit() {
	if dataInstance == nil {
		dataInstance = js.Global.Get("Object").New()
		dataInstance.Set("data", js.Global.Get("Object").New())
		dataInstance.Set("methods", js.Global.Get("Object").New())
		dataInstance.Set("computed", js.Global.Get("Object").New())
		dataInstance.Set("watch", js.Global.Get("Object").New())
	}
}

type Watcher func(val, prev interface{})

func AddRoot(key string, val interface{}) error {
	chkInit()
	if dataInstance.Get(key) != js.Undefined {
		return errors.New("Duplicated key in AddRoot")
	}
	dataInstance.Set(key, val)
	return nil
}

func AddData(key string, val interface{}) error {
	chkInit()
	if dataInstance.Get("data").Get(key) != js.Undefined {
		return errors.New("Duplicated key in AddData")
	}
	dataInstance.Get("data").Set(key, val)
	return nil
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

func Ready() {
	js.Global.Set("vm", js.Global.Get("Vue").New(dataInstance))
}

func Object() *js.Object {
	return js.Global.Get("Object").New()
}
