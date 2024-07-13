package gormmom

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"slices"

	"github.com/yyle88/done"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type Config struct {
	tagName string
	dftRule gormmomrule.Rule //默认检查规则
	nameMap map[gormmomrule.Rule]func(string) string
	skipAbc bool //是否跳过简单字段，有的字段虽然没有配置名称或者规则，但是它满足简单字段，就也不做任何处理
}

func NewConfig() *Config {
	return &Config{
		tagName: "mom",
		dftRule: gormmomrule.S63, //默认检查规则，就是查看是不是63个合法字符（即字母数组下划线等）
		nameMap: make(map[gormmomrule.Rule]func(string) string),
		skipAbc: true,
	}
}

func (cfg *Config) SetTagName(tagName string) *Config {
	cfg.tagName = tagName
	return cfg
}

func (cfg *Config) SetDftRule(dftRule gormmomrule.Rule) *Config {
	cfg.dftRule = dftRule
	return cfg
}

func (cfg *Config) SetNameMap(nameMap map[gormmomrule.Rule]func(string) string) *Config {
	cfg.nameMap = nameMap
	return cfg
}

func (cfg *Config) SetSkipAbc(skipAbc bool) *Config {
	cfg.skipAbc = skipAbc
	return cfg
}

func (cfg *Config) GenReplaces(params []*Param) {
	for _, param := range params {
		cfg.GenWrite(param, cfg.GenSource(param))
	}
}

func (cfg *Config) GenReplace(param *Param) {
	cfg.GenWrite(param, cfg.GenSource(param))
}

func (cfg *Config) GenWrite(param *Param, newCode []byte) {
	done.Done(utils.WriteFile(param.path, done.VAE(formatgo.FormatBytes(newCode)).Nice()))
}

func (cfg *Config) GenSource(param *Param) []byte {
	param.Validate()

	srcData := done.VAE(os.ReadFile(param.path)).Nice()
	astFile := done.VCE(syntaxgo_ast.NewAstFromSource(srcData)).Nice()

	structType := syntaxgo_ast.SeekStructXName(astFile, param.structName)
	if structType == nil {
		const reason = "CAN NOT FIND STRUCT TYPE"
		zaplog.LOG.Panic(reason, zap.String("struct_name", param.structName))
		panic(reason)
	}
	done.Done(ast.Print(token.NewFileSet(), structType))

	type changeType struct {
		node ast.Node
		code string
	}
	var changeSteps []*changeType

	// 遍历结构体的字段
	for _, field := range structType.Fields.List {
		// 打印字段名称和类型
		for _, nameIdent := range field.Names {
			zaplog.LOG.Debug("--")
			zaplog.LOG.Debug("process", zap.String("struct_name:", nameIdent.Name))
			zaplog.LOG.Debug("--")

			schemaField, exist := param.columnsMap[nameIdent.Name]
			if !exist { //比如字段是 "V哈哈" 就没事 而假如是 "v哈哈" 或者 "哈哈" 就不行，因为非以大写字母开始的字段，就没有gorm的列名
				zaplog.LOG.Debug("NO SCHEMA_FIELD - MAYBE NAME IS UNEXPORTED", zap.String("struct_name", nameIdent.Name))
				continue
			}

			if field.Tag == nil {
				zaplog.LOG.Debug("NO TAG", zap.String("struct_name", nameIdent.Name))
				if cfg.skipAbc && gormmomrule.S63.Validate(schemaField.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", nameIdent.Name))
					continue
				}
				rule := cfg.dftRule
				if !rule.Validate(schemaField.DBName) {
					if len(field.Names) >= 2 { //比如 a,b int 这种两个字段在一起，但其中一个字段的列名不正确时，就没法自动解决啦（其实有办法但不想实现，因为代价较大而没有收益）
						const reason = "CAN NOT HANDLE THIS SITUATION"
						zaplog.LOG.Panic(reason, zap.String("struct_name", nameIdent.Name))
						panic(reason)
					}
					newTag := cfg.newFixTagCode(schemaField, "``", rule) //这里应该走创建标签的逻辑，但和修改标签的逻辑是相同的
					changeSteps = append(changeSteps, &changeType{
						node: syntaxgo_ast.NewNode(field.End(), field.End()), //在尾部插入新的标签，要紧贴字段而且在换行符前面
						code: newTag,
					})
				}
			} else if ruleName := cfg.extractRuleName(field.Tag.Value); ruleName != "" {
				rule := gormmomrule.Rule(ruleName)
				zaplog.LOG.Debug("process", zap.String("rule", string(rule)))
				newTag := cfg.newFixTagCode(schemaField, field.Tag.Value, rule)
				changeSteps = append(changeSteps, &changeType{
					node: field.Tag, //完整替换原来的标签
					code: newTag,
				})
			} else {
				if cfg.skipAbc && gormmomrule.S63.Validate(schemaField.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", nameIdent.Name))
					continue
				}
				rule := cfg.dftRule
				if !rule.Validate(schemaField.DBName) { //按照比较宽泛的规则也校验不过的时候就需要修正字段名
					newTag := cfg.newFixTagCode(schemaField, field.Tag.Value, rule)
					changeSteps = append(changeSteps, &changeType{
						node: field.Tag, //完整替换原来的标签
						code: newTag,
					})
				} else {
					zaplog.LOG.Debug("meet rule skip", zap.String("name", nameIdent.Name), zap.String("tag", field.Tag.Value))
				}
			}
		}
	}

	//需要翻转下从后往前替换，因为替换以后源码会变，假如从前往后替换坐标就对不上啦，而从后往前替换则不存在这个问题
	slices.Reverse(changeSteps)

	//接下来替换代码，把需要 新增 或者 替换 的标签都设置到代码里
	newCode := srcData
	for _, step := range changeSteps {
		newCode = syntaxgo_ast.ChangeNodeBytes(newCode, step.node, []byte(step.code))
	}
	return newCode
}

func (cfg *Config) extractRuleName(tag string) string {
	tagValue := syntaxgo_tag.ExtractTagValue(tag, cfg.tagName)
	if tagValue == "" {
		return ""
	}
	tagField := syntaxgo_tag.ExtractTagField(tagValue, "rule", true)
	if tagField == "" {
		return ""
	}
	return tagField
}

func (cfg *Config) newFixTagCode(schemaField *schema.Field, tag string, rule gormmomrule.Rule) string {
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("tag", tag))
	//在 gorm 里修改 column 内容
	newTag := cfg.newFixGormTag(schemaField, tag, rule)
	//在 规则 里修改 rule 内容
	newTag = cfg.newFixRuleTag(newTag, rule)
	//这是替换后的结果，即替换整个标签内容，获得新的完整标签内容
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("new_tag", newTag))
	return newTag
}

func (cfg *Config) newFixGormTag(schemaField *schema.Field, tag string, rule gormmomrule.Rule) string {
	var columnName = gormmomrule.MakeName(rule, schemaField.Name, cfg.nameMap)
	zaplog.LOG.Debug("new_fix_gorm_tag", zap.String("name", schemaField.Name), zap.String("column_name", columnName))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, "gorm")
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`gorm:"column:%s;"`, columnName)
		if tag[1] != ' ' && tag[1] != '`' {
			part += " " //说明后面还有别的标签
		}
		p := 1 //插在第一个"`"的后面，即第一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return cfg.newFixTagField(tag, "gorm", "column", columnName)
}

func (cfg *Config) newFixRuleTag(tag string, rule gormmomrule.Rule) string {
	var ruleName = string(rule)
	zaplog.LOG.Debug("new_fix_rule_tag", zap.String("rule_name", ruleName))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, cfg.tagName)
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`%s:"rule:%s;"`, cfg.tagName, ruleName)
		if rch := tag[len(tag)-2]; rch != ' ' && rch != '`' {
			part = " " + part //说明前面还有别的标签
		}
		p := len(tag) - 1 //插在最后一个"`"的前面，即最后一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return cfg.newFixTagField(tag, cfg.tagName, "rule", ruleName)
}

func (cfg *Config) newFixTagField(tag string, tagName string, tagFieldName string, tagFieldValue string) string {
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("tag", tag))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, tagName)
	if sdx < 0 || edx < 0 {
		zaplog.LOG.Panic("IMPOSSIBLE") //能进到这个函数里的都是已经找到标签的
	}
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("tag_value", tagValue), zap.Int("sdx", sdx), zap.Int("edx", edx))

	tagField, s2x, e2x := syntaxgo_tag.ExtractTagFieldIndex(tagValue, tagFieldName, false)
	if s2x < 0 || e2x < 0 { //表示没找到 rule 自定义的内容
		part := fmt.Sprintf(tagFieldName+":%s;", tagFieldValue)
		p := sdx //插在gorm:的后面
		return tag[:p] + part + tag[p:]
	}
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("tag_field", tagField), zap.Int("s2x", s2x), zap.Int("e2x", e2x))

	spx := sdx + s2x //把起点坐标补上前面的
	epx := sdx + e2x
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("old_value", tag[spx:epx]))

	part := tagFieldValue
	if tag[epx] != ';' {
		part += ";" //当没有分号的时候就补个分号，没有也行但是有的话更安全些
	}
	return tag[:spx] + part + tag[epx:]
}
