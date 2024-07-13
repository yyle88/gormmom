package gormmom

import (
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
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

// NewParam 创建参数信息
func NewParam(path string, structName string, columnsMap map[string]*schema.Field) *Param {
	return &Param{
		path:       path,
		structName: structName,
		columnsMap: columnsMap,
	}
}

// NewParamV2 使用泛型创建参数信息。T 只能传类型名称而非带指针的类型名
func NewParamV2[T any](path string) *Param {
	var object T //这个时候T要传类型名，而不是指针类型，否则就不能初始化类型
	return NewParamV3(path, &object)
}

// NewParamV3 使用对象创建参数信息 object 传对象或者对象指针都行
func NewParamV3(path string, object interface{}) *Param {
	structName := syntaxgo_reflect.GetTypeNameV3(object)
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

func CreateParams(root string, models []interface{}) []*Param {
	var params = make([]*Param, 0, len(models))
	var paths = utils.GetLsGoFilePaths(root)
	var idxSet = make(map[int]bool, len(models)) //记住已经处理的数据
	for _, path := range paths {
		astFile, err := syntaxgo_ast.NewAstXFilepath(path)
		if err != nil {
			zaplog.LOG.Warn("something is wrong then warn", zap.String("path", path), zap.Error(err))
			continue
		}
		for objIdx, object := range models {
			if idxSet[objIdx] {
				//说明这种情况已经处理过
				//当然也有可能是已经错误
				continue
			}
			structName := syntaxgo_reflect.GetTypeNameV3(object)
			if structName == "" {
				zaplog.LOG.Warn("object doesn't have struct name", zap.Int("idx", objIdx))
				idxSet[objIdx] = true
				continue
			}
			structType := syntaxgo_ast.SeekStructXName(astFile, structName)
			if structType == nil {
				//这种情况下没有错误，而是说明这个文件里没有定义这个模型
				//但是在其他文件里可能有因此这里不是错误
				continue
			}
			params = append(params, NewParamV3(path, object))
			idxSet[objIdx] = true
		}
	}
	return params
}
