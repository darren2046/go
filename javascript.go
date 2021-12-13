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

func (m *javascriptVMStruct) Run(javascript string) *javascriptVMStruct {
	_, err := m.vm.Run(javascript)
	Panicerr(err)
	return m
}

func (m *javascriptVMStruct) Get(variableName string) string {
	value, err := m.vm.Get(variableName)
	Panicerr(err)
	valueStr, err := value.ToString()
	Panicerr(err)
	if valueStr == "undefined" {
		Panicerr("变量" + variableName + "未定义")
	}
	return valueStr
}

func (m *javascriptVMStruct) Set(variableName string, variableValue interface{}) {
	err := m.vm.Set(variableName, variableValue)
	Panicerr(err)
}

func (m *javascriptVMStruct) Isdefined(variableName string) bool {
	value, err := m.vm.Get(variableName)
	Panicerr(err)
	valueStr, err := value.ToString()
	Panicerr(err)
	return valueStr != "undefined"
}
