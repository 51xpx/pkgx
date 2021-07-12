package convert

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/51xpx/pkgx/pkg/pagex"
)

// AssignPageData 转换特殊字段，输出结构化列表信息
func AssignPageData(pageInfo pagex.PaginateResult, data interface{}) error {
	// @todo 检查为什么_格式的key不能读取的问题
	for _, item := range pageInfo.List {
		for k, v := range item {
			nk := strings.ReplaceAll(k, "_", "")
			item[nk] = v
		}
	}

	if err := Decode(pageInfo, &data); err != nil {
		return err
	}
	return nil
}

// AssignData 转换特殊字段，输出结构化列表信息
func AssignData(input interface{}, output interface{}) error {
	// @todo 先转map保证解析不会异常
	data := input
	k := reflect.TypeOf(input).Kind()
	if k == reflect.Struct || k == reflect.Ptr {
		data = structs.Map(input)
	}

	if err := Decode(data, &output); err != nil {
		return err
	}
	return nil
}
