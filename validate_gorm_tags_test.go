package gormmom_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath"
)

type Product1 struct {
	ID  uint    `gorm:"primaryKey"`
	P名称 string  `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	P价格 float64 `gorm:"column:price;type:decimal(10,2);not null;index"`
	S库存 int     `gorm:"type:int;default:0"`
	D描述 string  `gorm:"column:description;type:text"`
	S状态 string  `gorm:"column:status;type:varchar(20);default:'active';index"`
}

// TableName specifies the table name for Product1 model
// TableName 指定 Product1 模型的表名
func (*Product1) TableName() string {
	return "product1s"
}

func TestConfigs_ValidateGormTags(t *testing.T) {
	objects := []interface{}{
		&Product1{},
	}

	params := gormmom.ParseObjects(runpath.PARENT.Path(), objects)

	cfg := gormmom.NewConfigs(params, gormmom.NewOptions())
	t.Log(cfg)

	require.Error(t, cfg.ValidateGormTags())
}

type Product2 struct {
	ID  uint    `gorm:"primaryKey"`
	P名称 string  `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	P价格 float64 `gorm:"column:price;type:decimal(10,2);not null;index"`
	S库存 int     `gorm:"column:stock;type:int;default:0;index"`
	D描述 string  `gorm:"column:description;type:text"`
	S状态 string  `gorm:"column:status;type:varchar(20);default:'active';index"`
}

// TableName specifies the table name for Product2 model
// TableName 指定 Product2 模型的表名
func (*Product2) TableName() string {
	return "product2s"
}

func TestConfigs_ValidateGormTags_2(t *testing.T) {
	objects := []interface{}{
		&Product2{},
	}

	params := gormmom.ParseObjects(runpath.PARENT.Path(), objects)

	cfg := gormmom.NewConfigs(params, gormmom.NewOptions())
	t.Log(cfg)

	require.Error(t, cfg.ValidateGormTags())
}
