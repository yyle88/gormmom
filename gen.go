package gormmom

import (
	"fmt"
	"go/ast"
	"os"
	"slices"

	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_astnode"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type Config struct {
	ruleTagName string
	defaultRule gormmomrule.MomRULE //默认检查规则
	nameGenMap  map[gormmomrule.MomRULE]gormmomrule.CnmMakeIFace
	skipAbc123  bool //是否跳过简单字段，有的字段虽然没有配置名称或者规则，但是它满足简单字段，就也不做任何处理
	genIdxName  bool
	idxNameMap  map[gormidxname.IdxNAME]gormidxname.IdxNameIFace
}

func NewConfig() *Config {
	return &Config{
		ruleTagName: "mom",
		defaultRule: gormmomrule.DEFAULT, //默认检查规则，就是查看是不是63个合法字符（即字母数组下划线等）
		nameGenMap:  gormmomrule.GetPresetCnmMakeMap(),
		skipAbc123:  true,
		genIdxName:  true,
		idxNameMap:  gormidxname.GetPresetNameImpMap(),
	}
}

func (cfg *Config) SetRuleTagName(ruleTagName string) *Config {
	cfg.ruleTagName = ruleTagName
	return cfg
}

func (cfg *Config) SetDefaultRule(defaultRule gormmomrule.MomRULE) *Config {
	cfg.defaultRule = defaultRule
	return cfg
}

func (cfg *Config) SetNameGenImp(momRULE gormmomrule.MomRULE, cnmMakeImp gormmomrule.CnmMakeIFace) *Config {
	cfg.nameGenMap[momRULE] = cnmMakeImp
	return cfg
}

func (cfg *Config) SetSkipSimple(skipSimple bool) *Config {
	cfg.skipAbc123 = skipSimple
	return cfg
}

func (cfg *Config) SetGenIdxName(genIdxName bool) *Config {
	cfg.genIdxName = genIdxName
	return cfg
}

func (cfg *Config) SetIdxNameImp(rule gormidxname.IdxNAME, idxNameImp gormidxname.IdxNameIFace) *Config {
	cfg.idxNameMap[rule] = idxNameImp
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
	param.CheckParam()

	srcData := done.VAE(os.ReadFile(param.path)).Nice()
	astBundle := done.VCE(syntaxgo_ast.NewAstBundleV1(srcData)).Nice()

	astFile, fileSet := astBundle.GetBundle()

	structContent, ok := syntaxgo_search.FindStructTypeByName(astFile, param.structName)
	if !ok {
		const reason = "CAN NOT FIND STRUCT TYPE"
		zaplog.LOG.Panic(reason, zap.String("struct_name", param.structName))
		panic(reason)
	}
	done.Done(ast.Print(fileSet, structContent))

	var srcChanges []*changeType

	// 遍历结构体的字段
	for _, field := range structContent.Fields.List {
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
				if cfg.skipAbc123 && cfg.nameGenMap[gormmomrule.DEFAULT].CheckName(schemaField.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", nameIdent.Name))
					continue
				}
				momRULE := cfg.defaultRule
				if !cfg.nameGenMap[momRULE].CheckName(schemaField.DBName) {
					if len(field.Names) >= 2 { //比如 a,b int 这种两个字段在一起，但其中一个字段的列名不正确时，就没法自动解决啦（其实有办法但不想实现，因为代价较大而没有收益）
						const reason = "CAN NOT HANDLE THIS SITUATION"
						zaplog.LOG.Panic(reason, zap.String("struct_name", nameIdent.Name))
						panic(reason) //这种情况下当有错时，就不处理这种情况，就需要程序员先把两个字段定义到两行里
					}
					changeTag := cfg.newFixTagCode(schemaField, "``", momRULE) //这里应该走创建标签的逻辑，但和修改标签的逻辑是相同的
					srcChanges = append(srcChanges, &changeType{
						vFieldName: nameIdent.Name,
						oldTagNode: syntaxgo_astnode.NewNode(field.End(), field.End()), //在尾部插入新的标签，要紧贴字段而且在换行符前面
						newTagCode: changeTag,
					})
				}
			} else {
				if ruleName := cfg.extractRuleField(field.Tag.Value); ruleName != "" {
					momRULE := gormmomrule.MomRULE(ruleName)
					zaplog.LOG.Debug("process", zap.String("rule", string(momRULE)))
					changeTag := cfg.newFixTagCode(schemaField, field.Tag.Value, momRULE)
					srcChanges = append(srcChanges, &changeType{
						vFieldName: nameIdent.Name,
						oldTagNode: field.Tag, //完整替换原来的标签
						newTagCode: changeTag,
					})
				} else {
					if cfg.skipAbc123 && cfg.nameGenMap[gormmomrule.DEFAULT].CheckName(schemaField.DBName) {
						zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", nameIdent.Name))
						continue
					}
					momRULE := cfg.defaultRule
					if !cfg.nameGenMap[momRULE].CheckName(schemaField.DBName) { //按照比较宽泛的规则也校验不过的时候就需要修正字段名
						changeTag := cfg.newFixTagCode(schemaField, field.Tag.Value, momRULE)
						srcChanges = append(srcChanges, &changeType{
							vFieldName: nameIdent.Name,
							oldTagNode: field.Tag, //完整替换原来的标签
							newTagCode: changeTag,
						})
					} else {
						zaplog.LOG.Debug("meet rule skip", zap.String("name", nameIdent.Name), zap.String("tag", field.Tag.Value))
					}
				}
			}
		}
	}

	zaplog.LOG.Debug("change_column_names")
	for _, rep := range srcChanges {
		zaplog.LOG.Debug("check_column:", zap.String("name", rep.vFieldName), zap.String("code", rep.newTagCode))
	}

	if cfg.genIdxName {
		//这里增加个新逻辑，就是单列索引的索引名称不正确，需要也校正索引名，因此这个函数会补充标签内容
		cfg.rewriteIndexNames(param, srcChanges)

		zaplog.LOG.Debug("change_index_names")
		for _, rep := range srcChanges {
			zaplog.LOG.Debug("check_index:", zap.String("name", rep.vFieldName), zap.String("code", rep.newTagCode))
		}
	}

	//需要翻转下从后往前替换，因为替换以后源码会变，假如从前往后替换坐标就对不上啦，而从后往前替换则不存在这个问题
	slices.Reverse(srcChanges)

	//接下来替换代码，把需要 新增 或者 替换 的标签都设置到代码里
	newCode := srcData
	for _, step := range srcChanges {
		newCode = syntaxgo_astnode.ChangeNodeCode(newCode, step.oldTagNode, []byte(step.newTagCode))
	}
	return newCode
}

type changeType struct {
	vFieldName string   //结构体的字段名
	oldTagNode ast.Node //标签的起止位置-就是在src源码中的位置，便于后面的替换代码
	newTagCode string   //新标签的新内容-就是标签的完整全部内容
}

func (cfg *Config) extractRuleField(tagCode string) string {
	return cfg.extractSomeField(tagCode, cfg.ruleTagName, "rule")
}

func (cfg *Config) extractSomeField(tagCode string, key1 string, key2 string) string {
	tagValue := syntaxgo_tag.ExtractTagValue(tagCode, key1)
	if tagValue == "" {
		return ""
	}
	tagField := syntaxgo_tag.ExtractTagField(tagValue, key2, syntaxgo_tag.EXCLUDE_WHITESPACE_PREFIX)
	if tagField == "" {
		return ""
	}
	return tagField
}

func (cfg *Config) newFixTagCode(schemaField *schema.Field, tag string, momRULE gormmomrule.MomRULE) string {
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("tag", tag))
	//在 gorm 里修改 column 内容
	newTag := cfg.newFixGormTag(schemaField, tag, momRULE)
	//在 规则 里修改 rule 内容
	newTag = cfg.newFixRuleTag(newTag, momRULE)
	//这是替换后的结果，即替换整个标签内容，获得新的完整标签内容
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("new_tag", newTag))
	return newTag
}

func (cfg *Config) newFixGormTag(schemaField *schema.Field, tag string, momRULE gormmomrule.MomRULE) string {
	var columnName = cfg.nameGenMap[momRULE].GenNewCnm(schemaField.Name)
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
	return cfg.newFixTagField(tag, "gorm", "column", columnName, TOP)
}

func (cfg *Config) newFixRuleTag(tag string, momRULE gormmomrule.MomRULE) string {
	zaplog.LOG.Debug("new_fix_rule_tag", zap.String("rule_name", string(momRULE)))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, cfg.ruleTagName)
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`%s:"rule:%s;"`, cfg.ruleTagName, string(momRULE))
		if rch := tag[len(tag)-2]; rch != ' ' && rch != '`' {
			part = " " + part //说明前面还有别的标签
		}
		p := len(tag) - 1 //插在最后一个"`"的前面，即最后一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return cfg.newFixTagField(tag, cfg.ruleTagName, "rule", string(momRULE), TOP)
}

type enumInsertLocation string

const (
	TOP enumInsertLocation = "TOP"
	END enumInsertLocation = "END"
)

func (cfg *Config) newFixTagField(tag string, tagName string, tagFieldName string, tagFieldValue string, insertLocation enumInsertLocation) string {
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("tag", tag))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, tagName)
	if sdx < 0 || edx < 0 {
		zaplog.LOG.Panic("IMPOSSIBLE") //能进到这个函数里的都是已经找到标签的
	}
	zaplog.LOG.Debug("new_fix_tag_field", zap.String("tag_value", tagValue), zap.Int("sdx", sdx), zap.Int("edx", edx))

	tagField, s2x, e2x := syntaxgo_tag.ExtractTagFieldIndex(tagValue, tagFieldName, syntaxgo_tag.INCLUDE_WHITESPACE_PREFIX)
	if s2x < 0 || e2x < 0 { //表示没找到 rule 自定义的内容
		part := fmt.Sprintf(tagFieldName+":%s;", tagFieldValue)
		if insertLocation == TOP {
			p := sdx //插在gorm:的后面
			return tag[:p] + part + tag[p:]
		} else if insertLocation == END {
			p := edx
			if p > 0 {
				if c := tag[p-1]; c == '"' || c == ';' || c == ' ' {
					//当是第一个或者前一个已经带分号时，就不需要加分号
				} else {
					part = ";" + part //否则就需要在前面添加个分号
				}
			}
			return tag[:p] + part + tag[p:]
		} else {
			panic(erero.New("WRONG"))
		}
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
