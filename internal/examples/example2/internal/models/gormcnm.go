package models

import "github.com/yyle88/gormcnm"

// 这个文件将被 gormcngen 自动生成和更新

func (T *Tユーザー) Columns() *TユーザーColumns {
	return &TユーザーColumns{
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
	// Embedding operation functions make it easy to use // 继承操作函数便于使用
	gormcnm.ColumnOperationClass
	// The column names and types of the model's columns // 模型各列的列名和类型
	ID     gormcnm.ColumnName[uint]
	Uユーザー名 gormcnm.ColumnName[string]
	Eメール   gormcnm.ColumnName[string]
	A年齢    gormcnm.ColumnName[int]
	D電話    gormcnm.ColumnName[string]
	J住所    gormcnm.ColumnName[string]
	Sステータス gormcnm.ColumnName[string]
}
