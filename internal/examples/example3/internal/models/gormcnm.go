package models

import "github.com/yyle88/gormcnm"

// 这个文件将被 gormcngen 自动生成和更新

func (T *T사용자) Columns() *T사용자Columns {
	return &T사용자Columns{
		ID:    gormcnm.Cnm(T.ID, "id"),
		U사용자명: gormcnm.Cnm(T.U사용자명, "u_acc0_a9c6_90c7_85ba"),
		E이메일:  gormcnm.Cnm(T.E이메일, "e_74c7_54ba_7cc7"),
		N나이:   gormcnm.Cnm(T.N나이, "n_98b0_74c7"),
		J전화:   gormcnm.Cnm(T.J전화, "j_04c8_54d6"),
		J주소:   gormcnm.Cnm(T.J주소, "j_fcc8_8cc1"),
		S상태:   gormcnm.Cnm(T.S상태, "s_c1c0_dcd0"),
	}
}

type T사용자Columns struct {
	// Embedding operation functions make it easy to use // 继承操作函数便于使用
	gormcnm.ColumnOperationClass
	// The column names and types of the model's columns // 模型各列的列名和类型
	ID    gormcnm.ColumnName[uint]
	U사용자명 gormcnm.ColumnName[string]
	E이메일  gormcnm.ColumnName[string]
	N나이   gormcnm.ColumnName[int]
	J전화   gormcnm.ColumnName[string]
	J주소   gormcnm.ColumnName[string]
	S상태   gormcnm.ColumnName[string]
}
