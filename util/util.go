package util

import "reflect"

//ElementInStrSlice 判断string类型element是否在string类型的slice里
func ElementInStrSlice(s string, e []string) bool {
	for _, v := range e {
		if s == v {
			return true
		}
	}
	return false
}

//RemoveRepeatElement 移除一个string类型slice里的重复元素
func RemoveRepeatElement(s []string) []string {
	result := make([]string, 0)

	for _, v := range s {
		if !ElementInStrSlice(v, result) {
			result = append(result, v)
		}
	}

	return result
}

//IsFn 判断给出的数据是否是一个函数
func IsFn(fn interface{}) bool {
	return reflect.TypeOf(fn).Kind() == reflect.Func
}
