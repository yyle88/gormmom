// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yyle88/gormcngen
// Generated from: gormcnm_test.go:43 -> models.TestGen.func2
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import (
	"time"

	"github.com/yyle88/gormcnm"
	"gorm.io/gorm"
)

func (T *Example) Columns() *ExampleColumns {
	return &ExampleColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
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

type ExampleColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
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
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
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
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
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

func (T *Example3) Columns() *Example3Columns {
	return &Example3Columns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		ID:        gormcnm.Cnm(T.ID, "id"),
		CreatedAt: gormcnm.Cnm(T.CreatedAt, "created_at"),
		UpdatedAt: gormcnm.Cnm(T.UpdatedAt, "updated_at"),
		DeletedAt: gormcnm.Cnm(T.DeletedAt, "deleted_at"),
		U账号:       gormcnm.Cnm(T.U账号, "username"),
		N昵称:       gormcnm.Cnm(T.N昵称, "nickname"),
		V分数:       gormcnm.Cnm(T.V分数, "score"),
		Age:       gormcnm.Cnm(T.Age, "age"),
		Uuid:      gormcnm.Cnm(T.Uuid, "uuid"),
		Rank:      gormcnm.Cnm(T.Rank, "rank"),
	}
}

type Example3Columns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID        gormcnm.ColumnName[uint]
	CreatedAt gormcnm.ColumnName[time.Time]
	UpdatedAt gormcnm.ColumnName[time.Time]
	DeletedAt gormcnm.ColumnName[gorm.DeletedAt]
	U账号       gormcnm.ColumnName[string]
	N昵称       gormcnm.ColumnName[string]
	V分数       gormcnm.ColumnName[uint64]
	Age       gormcnm.ColumnName[int]
	Uuid      gormcnm.ColumnName[int]
	Rank      gormcnm.ColumnName[int]
}
