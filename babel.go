package babel

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/dop251/goja"
)

func Transform(src io.Reader, opts map[string]interface{}) (io.Reader, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	res, err := TransformString(string(data), opts)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(res), nil
}

func TransformString(src string, opts map[string]interface{}) (string, error) {
	if opts == nil {
		opts = map[string]interface{}{}
	}
	vm := goja.New()
	transform, err := loadBabel(vm)
	if err != nil {
		return "", err
	}
	v, err := transform(src, opts)
	if err != nil {
		return "", err
	}
	return v.ToObject(vm).Get("code").Export().(string), nil
}

func loadBabel(vm *goja.Runtime) (func(string, map[string]interface{}) (goja.Value, error), error) {
	babelsrc, err := Asset("babel.js")
	if err != nil {
		return nil, err
	}
	_, err = vm.RunScript("babel.js", string(babelsrc))
	if err != nil {
		return nil, fmt.Errorf("unable to load babel.js: %s", err)
	}
	var transform goja.Callable
	babel := vm.Get("Babel")
	if err := vm.ExportTo(babel.ToObject(vm).Get("transform"), &transform); err != nil {
		return nil, fmt.Errorf("unable to export transform fn: %s", err)
	}
	return func(src string, opts map[string]interface{}) (goja.Value, error) {
		return transform(babel, vm.ToValue(src), vm.ToValue(opts))
	}, nil
}
