package example3

import (
	"time"

	"github.com/yyle88/gormcnm"
)

func (T *Example1) Columns() *Example1Columns {
	return &Example1Columns{
		ID:        gormcnm.Cnm(T.ID, "id"),
		V名称:       gormcnm.Cnm(T.V名称, "v_0d54_f079"),
		V字段:       gormcnm.Cnm(T.V字段, "v_575b_b56b"),
		V性别:       gormcnm.Cnm(T.V性别, "v_2760_2b52"),
		V特殊:       gormcnm.Cnm(T.V特殊, "v_7972_8a6b"),
		V年龄:       gormcnm.Cnm(T.V年龄, "v_745e_849f"),
		Rank:      gormcnm.Cnm(T.Rank, "rank"),
		V身高:       gormcnm.Cnm(T.V身高, "v_ab8e_d89a"),
		V体重:       gormcnm.Cnm(T.V体重, "v_534f_cd91"),
		CreatedAt: gormcnm.Cnm(T.CreatedAt, "created_at"),
		UpdatedAt: gormcnm.Cnm(T.UpdatedAt, "updated_at"),
	}
}

type Example1Columns struct {
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

func (T *Example2) Columns() *Example2Columns {
	return &Example2Columns{
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

type Example2Columns struct {
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