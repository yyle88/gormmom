package gormmom

import (
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type Param struct {
	path       string
	structName string
	columnsMap map[string]*schema.Field
}

func NewParam(path string, structName string, columnsMap map[string]*schema.Field) *Param {
	return &Param{
		path:       path,
		structName: structName,
		columnsMap: columnsMap,
	}
}

func NewParamV2[T any](path string) *Param {
	var object T
	return NewParamV3(path, object)
}

func NewParamV3(path string, object interface{}) *Param {
	structName := syntaxgo_reflect.GetTypeName(object)
	zaplog.LOG.Debug("new_param", zap.String("struct_name", structName))

	columnsMap := utils.GetSchemaFieldsMap(&object)
	zaplog.LOG.Debug("new_param", zap.Int("column_size", len(columnsMap)))

	return NewParam(path, structName, columnsMap)
}

func (param *Param) Validate() {
	if param.path == "" {
		panic("param.path is none")
	}
	if param.structName == "" {
		panic("param.struct_name is none")
	}
	if param.columnsMap == nil {
		panic("param.columns_map is none")
	}
}
