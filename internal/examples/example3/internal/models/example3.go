package models

import "gorm.io/gorm"

type Example3 struct {
	gorm.Model
	U账号  string `gorm:"column:username;uniqueIndex:udx_example3s_USERNAME; NOT NULL" mom:"udx:CNM;"`
	N昵称  string `gorm:"column:nickname; NOT NULL"`
	V分数  uint64 `gorm:"column:score;index:idx_example3s_score;" mom:"idx:cnm"`
	Age  int    `gorm:"index:idx_example3s_age;" mom:"idx:cnm"`
	Uuid int    `gorm:"uniqueIndex:udx_example3s_uuid;index:idx_example3s_uuid;" mom:"idx:cnm;udx:cnm"`
	Rank int    `gorm:"index" mom:"udx:cnm"` // udx not affect index
}

func (*Example3) TableName() string {
	return "example3s"
}
