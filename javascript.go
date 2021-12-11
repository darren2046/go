package golanglibs

import (
	"github.com/robertkrimen/otto"
)

type javascriptVMStruct struct {
	vm *otto.Otto
}

func getJavascriptVM() *javascriptVMStruct {
	return &javascriptVMStruct{vm: otto.New()}
}

func (m *javascriptVMStruct) run(javascript string) *javascriptVMStruct {
	_, err := m.vm.Run(javascript)
	panicerr(err)
	return m
}

func (m *javascriptVMStruct) get(variableName string) string {
	value, err := m.vm.Get(variableName)
	panicerr(err)
	valueStr, err := value.ToString()
	panicerr(err)
	if valueStr == "undefined" {
		panicerr("变量" + variableName + "未定义")
	}
	return valueStr
}

func (m *javascriptVMStruct) set(variableName string, variableValue interface{}) {
	err := m.vm.Set(variableName, variableValue)
	panicerr(err)
}

func (m *javascriptVMStruct) isdefined(variableName string) bool {
	value, err := m.vm.Get(variableName)
	panicerr(err)
	valueStr, err := value.ToString()
	panicerr(err)
	return valueStr != "undefined"
}
