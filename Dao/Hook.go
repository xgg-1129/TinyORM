package Dao

import (
	mylog "TinyORM/Log"
	"reflect"
)

/*
	钩子本质=一个函数
	这个函数会在一定会在某些函数前执行，或者某些函数后执行，钩子钩住的就是这些函数
	钩子钩住的函数执行逻辑不变，钩子函数本身变动
	https://www.jianshu.com/p/60c9a09d1ac3
	https://www.cnblogs.com/codingSoul/p/6018582.html
*/


// 8种钩子函数，分别钩住 增、删、查、改、函数的 执行前 和 执行后
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)
// 用来调用钩子函数，当然其实也可以直接用的BeforeQuery()方式调用，不过用CallMethod的方法能够解耦
func (s *session) CallMethod(method string, value interface{}) {
	//获取具体的钩子函数
	fm := reflect.ValueOf(s.RefTable().Object).MethodByName(method)

	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	//构造参数，每个钩子函数第一个参数是结构体自己，然后是*session
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		//调用钩子函数，参数不用传递自己嘛？
		if v := fm.Call(param); len(v) > 0 {
			//var a error=nil
			//	fmt.Println(a.(error)) ,nil没有实现error的string方法，所以无法断言成功
			if err, ok := v[0].Interface().(error); ok {
				mylog.Error(err)
			}
		}
	}
	return
}