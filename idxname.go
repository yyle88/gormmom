package gormmom

import (
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func (cfg *Config) CheckIndexes(params []*Param) {
	for _, param := range params {
		schIndexes := param.sch.ParseIndexes()
		zaplog.LOG.Debug("check_indexes", zap.String("object_class", param.sch.Name), zap.String("table_name", param.sch.Table), zap.Int("index_count", len(schIndexes)))
		for _, node := range schIndexes {
			zaplog.LOG.Debug("foreach_index")
			zaplog.LOG.Debug("check_a_index", zap.String("index_name", node.Name), zap.Int("field_size", len(node.Fields)))
			if len(node.Fields) == 1 { //只检查单列索引，因为复合索引就得手写名称，因此没有问题
				cfg.checkSingleIndex(node.Name, node.Fields[0].Name)
			}
		}
	}
}

func (cfg *Config) checkSingleIndex(indexName string, fieldName string) {
	zaplog.LOG.Debug("check_single_index", zap.String("index_name", indexName), zap.String("field_name", fieldName))
}
