package Dialect

import "reflect"

type field struct {
	Name    string
	SqlType string //这里是数据库中支持的类型
	Tag     string
}
type Schema struct {
	Object     interface{}
	Name       string
	Fields     []*field
	FieldsName []string
	FieldsMap  map[string]*field
}

func ParseObect(object interface{},d Dialect)*Schema {
	modelType:=reflect.Indirect(reflect.ValueOf(object)).Type()
	res:=new(Schema)
	res.Object =object
	res.Name=modelType.Name()
	res.FieldsMap =make(map[string]*field)
	for i:=0;i<modelType.NumField();i++{
		temp:=modelType.Field(i)

		filed:=&field{
			Name:    temp.Name,
			SqlType: d.DataTypeOf(reflect.Indirect(reflect.New(temp.Type))),
		}
		if str ,ok := temp.Tag.Lookup("orm");ok{
			filed.Tag = str
		}
		res.Fields =append(res.Fields,filed)
		res.FieldsMap[filed.Name]=filed
		res.FieldsName =append(res.FieldsName,filed.Name)
	}
	return res
}
func (s *Schema) RecordValues(object interface{}) []interface{}{
	value:=reflect.Indirect(reflect.ValueOf(object))
	var res []interface{}
	for _,fd := range s.Fields{
		res=append(res,value.FieldByName(fd.Name).Interface())
	}
	return res
}
