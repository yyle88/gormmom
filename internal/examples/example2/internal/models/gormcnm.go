// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yyle88/gormcngen
// Generated from: gormcnm_test.go:37 -> models.TestGenGormMomAndCnm.func2
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yyle88/gormcnm"

// 这个文件将被 gormcngen 自动生成和更新

func (T *Tユーザー) Columns() *TユーザーColumns {
	return &TユーザーColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		ID:     gormcnm.Cnm(T.ID, "id"),
		Uユーザー名: gormcnm.Cnm(T.Uユーザー名, "u_e630_fc30_b630_fc30_0d54"),
		Eメール:   gormcnm.Cnm(T.Eメール, "e_e130_fc30_eb30"),
		A年齢:    gormcnm.Cnm(T.A年齢, "a_745e_629f"),
		D電話:    gormcnm.Cnm(T.D電話, "d_fb96_718a"),
		J住所:    gormcnm.Cnm(T.J住所, "j_4f4f_4062"),
		Sステータス: gormcnm.Cnm(T.Sステータス, "s_b930_c630_fc30_bf30_b930"),
	}
}

type TユーザーColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID     gormcnm.ColumnName[uint]
	Uユーザー名 gormcnm.ColumnName[string]
	Eメール   gormcnm.ColumnName[string]
	A年齢    gormcnm.ColumnName[int]
	D電話    gormcnm.ColumnName[string]
	J住所    gormcnm.ColumnName[string]
	Sステータス gormcnm.ColumnName[string]
}
