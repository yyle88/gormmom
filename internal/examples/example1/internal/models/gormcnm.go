// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yyle88/gormcngen
// Generated from: gormcnm_test.go:41 -> models.TestGen.func2
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yyle88/gormcnm"

func (T *T用户) Columns() *T用户Columns {
	return &T用户Columns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		ID:   gormcnm.Cnm(T.ID, "id"),
		U用户名: gormcnm.Cnm(T.U用户名, "u_2875_3762_0d54"),
		E邮箱:  gormcnm.Cnm(T.E邮箱, "e_ae90_b17b"),
		A年龄:  gormcnm.Cnm(T.A年龄, "a_745e_849f"),
		D电话:  gormcnm.Cnm(T.D电话, "d_3575_dd8b"),
		J住所:  gormcnm.Cnm(T.J住所, "j_4f4f_4062"),
		S状态:  gormcnm.Cnm(T.S状态, "s_b672_0160"),
	}
}

type T用户Columns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID   gormcnm.ColumnName[uint]
	U用户名 gormcnm.ColumnName[string]
	E邮箱  gormcnm.ColumnName[string]
	A年龄  gormcnm.ColumnName[int]
	D电话  gormcnm.ColumnName[string]
	J住所  gormcnm.ColumnName[string]
	S状态  gormcnm.ColumnName[string]
}
