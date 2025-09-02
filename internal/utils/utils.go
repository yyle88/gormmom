package utils

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

// NewSchemaFieldsMap 把字段列表由 slice 转换为 map，以结构体中go的字段名为主键
func NewSchemaFieldsMap(gormSchema *schema.Schema) *linkedhashmap.Map[string, *schema.Field] {
	gormcngen.ShowSchemaEnglish(gormSchema)
	gormcngen.ShowSchemaChinese(gormSchema)

	res := linkedhashmap.New[string, *schema.Field]()
	for _, field := range gormSchema.Fields {
		res.Put(field.Name, field) //键是Go结构体成员名称
	}
	return res
}

// ListGoFiles 获取指定目录下的所有 .go 文件路径（不递归子目录）
func ListGoFiles(root string) []string {
	var paths []string
	for _, one := range rese.A1(os.ReadDir(root)) {
		// 检查是否是文件和扩展名为 .go
		if !one.IsDir() && filepath.Ext(one.Name()) == ".go" {
			paths = append(paths, filepath.Join(root, one.Name()))
		}
	}
	return paths
}

func ParseSchema(object interface{}) *schema.Schema {
	return rese.P1(schema.Parse(object, &sync.Map{}, &schema.NamingStrategy{
		SingularTable: false, //和默认值相同
		NoLowerCase:   false, //和默认值相同
	}))
}

func ParseTags[T any](sourceCode []byte, structObject *T) *linkedhashmap.Map[string, string] {
	astBundle := rese.P1(syntaxgo_ast.NewAstBundleV1(sourceCode))
	astFile, fileSet := astBundle.GetBundle()
	zaplog.LOG.Debug("ast-get-package-name", zap.String("package_name", astFile.Name.Name))

	structName := syntaxgo_reflect.GetTypeNameV4(structObject)
	zaplog.LOG.Debug("reflect-get-struct-name", zap.String("name", structName))

	structType := resb.P1(syntaxgo_search.FindStructTypeByName(astFile, structName))
	must.Done(ast.Print(fileSet, structType))

	var results = linkedhashmap.New[string, string]()
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				if field.Tag != nil {
					results.Put(name.Name, field.Tag.Value)
				} else {
					results.Put(name.Name, "")
				}
			}
		}
	}
	return results
}

func TrimQuotes(s string) string {
	return strings.Trim(s, `"`)
}

func TrimBackticks(s string) string {
	return strings.Trim(s, "`")
}

func ParseTagsTrimBackticks[T any](sourceCode []byte, structObject *T) *linkedhashmap.Map[string, string] {
	results := linkedhashmap.New[string, string]()
	ParseTags(sourceCode, structObject).Each(func(key string, value string) {
		results.Put(key, TrimBackticks(value))
	})
	return results
}

// NewCommonRegexp 创建一个正则表达式，检查列名的长度和字符
// 当列名前带个前导空格 比如 `gorm:"column: name;"` 时，在gorm中也是可以用的，但该规则里不允许这种情况，避免出现其它问题
func NewCommonRegexp(maxLen int) *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,` + strconv.Itoa(maxLen) + `}$`)
}

// MustMatchRegexp 检查字符串是否匹配正则表达式，如果不匹配则抛出 panic
func MustMatchRegexp(regexpRegexp *regexp.Regexp, value string) {
	if !regexpRegexp.MatchString(value) {
		zaplog.LOG.Panic("regexp does not match", zap.String("regexp", regexpRegexp.String()), zap.String("value", value))
	}
}

// ValidateTableName validates table name for database compatibility
// Checks if table name contains ASCII characters suitable for index generation
// Provides helpful message with solution when validation fails
//
// ValidateTableName 验证表名的数据库兼容性
// 检查表名是否包含适合索引生成的 ASCII 字符
// 在验证失败时提供带解决方案的有用信息
func ValidateTableName(tableName string, structName string) {
	// Check if table name contains non-ASCII characters
	hasNonASCII := false
	for _, c := range tableName {
		if c > 127 {
			hasNonASCII = true
			break
		}
	}

	if hasNonASCII {
		zaplog.LOG.Panic(
			"Table name contains non-ASCII characters which is not compatible. Please add a TableName() method to the struct with an ASCII table name.",
			zap.String("struct_name", structName),
			zap.String("current_table_name", tableName),
			zap.String("solution", fmt.Sprintf("Add this method to the struct: func (%s) TableName() string { return \"ascii_table_name\" }", structName)),
			zap.String("example", fmt.Sprintf("func (%s) TableName() string { return \"users\" }", structName)),
		)
	}

	// Check length constraint for table names
	if len(tableName) > 63 {
		zaplog.LOG.Panic(
			"Table name exceeds maximum length (63 characters). Please add a TableName() method to the struct with a compact name.",
			zap.String("struct_name", structName),
			zap.String("current_table_name", tableName),
			zap.Int("current_length", len(tableName)),
			zap.Int("max_length", 63),
			zap.String("solution", fmt.Sprintf("Add this method to the struct: func (%s) TableName() string { return \"compact_name\" }", structName)),
		)
	}
}
