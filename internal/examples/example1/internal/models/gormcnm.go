package models

import "github.com/yyle88/gormcnm"

func (T *T用户) Columns() *T用户Columns {
	return &T用户Columns{
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
	// Embedding operation functions make it easy to use // 继承操作函数便于使用
	gormcnm.ColumnOperationClass
	// The column names and types of the model's columns // 模型各列的列名和类型
	ID   gormcnm.ColumnName[uint]
	U用户名 gormcnm.ColumnName[string]
	E邮箱  gormcnm.ColumnName[string]
	A年龄  gormcnm.ColumnName[int]
	D电话  gormcnm.ColumnName[string]
	J住所  gormcnm.ColumnName[string]
	S状态  gormcnm.ColumnName[string]
}
