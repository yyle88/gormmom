package utils

import (
	"os"
	"sync"

	"github.com/yyle88/done"
	"github.com/yyle88/gormcngen"
	"gorm.io/gorm/schema"
)

// GetSchemaFieldsMap 把字段列表由 slice 转换为 map，以结构体中go的字段名为主键
func GetSchemaFieldsMap(dest interface{}) map[string]*schema.Field {
	sch := done.VCE(schema.Parse(dest, &sync.Map{}, &schema.NamingStrategy{
		SingularTable: false, //和默认值相同
		NoLowerCase:   false, //和默认值相同
	})).Nice()

	gormcngen.ShowSchemaMessage(sch)

	var mp = make(map[string]*schema.Field, len(sch.Fields))
	for _, field := range sch.Fields {
		mp[field.Name] = field ////go结构体成员名称
	}
	return mp
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
