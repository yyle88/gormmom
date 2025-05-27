package example2

import (
	"time"

	"github.com/yyle88/gormcnm"
)

func (T *Example) Columns() *exampleColumns {
	return &exampleColumns{
		ID:        gormcnm.Cnm(T.ID, "id"),
		V名称:       gormcnm.Cnm(T.V名称, "V_0D54_F079"),
		V字段:       gormcnm.Cnm(T.V字段, "v_575b_b56b"),
		V性别:       gormcnm.Cnm(T.V性别, "v_2760_2b52"),
		V特殊:       gormcnm.Cnm(T.V特殊, "V_7972_8A6B"),
		V年龄:       gormcnm.Cnm(T.V年龄, "v_745e_849f"),
		Rank:      gormcnm.Cnm(T.Rank, "rank"),
		V身高:       gormcnm.Cnm(T.V身高, "V_AB8E_D89A"),
		V体重:       gormcnm.Cnm(T.V体重, "v_534f_cd91"),
		CreatedAt: gormcnm.Cnm(T.CreatedAt, "created_at"),
		UpdatedAt: gormcnm.Cnm(T.UpdatedAt, "updated_at"),
	}
}

type exampleColumns struct {
	// Embedding operation functions make it easy to use // 继承操作函数便于使用
	gormcnm.ColumnOperationClass
	// The column names and types of the model's columns // 模型各列的列名和类型
	ID        gormcnm.ColumnName[int32]
	V名称       gormcnm.ColumnName[string]
	V字段       gormcnm.ColumnName[string]
	V性别       gormcnm.ColumnName[string]
	V特殊       gormcnm.ColumnName[string]
	V年龄       gormcnm.ColumnName[int]
	Rank      gormcnm.ColumnName[int32]
	V身高       gormcnm.ColumnName[int32]
	V体重       gormcnm.ColumnName[int32]
	CreatedAt gormcnm.ColumnName[time.Time]
	UpdatedAt gormcnm.ColumnName[time.Time]
}
