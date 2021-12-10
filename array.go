package golanglibs

import "reflect"

type arrayStruct struct {
	array interface{}
}

func Array(array interface{}) *arrayStruct {
	return &arrayStruct{array: array}
}

func (a *arrayStruct) Has(item interface{}) bool {
	// 获取值的列表
	arr := reflect.ValueOf(a.array)

	// 手工判断值的类型
	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panicerr("Invalid data type of param \"array\": Not an Array")
	}

	// 遍历值的列表
	for i := 0; i < arr.Len(); i++ {
		// 取出值列表的元素并转换为interface
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
