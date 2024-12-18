package example1enus

import "time"

type Example struct {
	ID        int32  `gorm:"column:id; primaryKey;" json:"id"`
	V名称       string `gorm:"column:v_0d54_f079;type:text" mom:"naming:s63;"`
	V字段       string `gorm:"column:v_575b_b56b;" mom:"naming:s63;"`
	V性别       string `gorm:"column:v_2760_2b52;" mom:"naming:s63;"`
	V特殊       string `gorm:"column:v_7972_8a6b;type:int32" mom:"naming:s63;"`
	V年龄       int    `gorm:"column:v_745e_849f;" json:"age" mom:"naming:s63;"` //理论上不要直接给model添加json标签，因为那是view层的逻辑，但实际上假如非这样做也能处理
	Rank      int32
	V身高       int32     `gorm:"column:v_ab8e_d89a;" mom:"naming:s63;"`
	V体重       int32     `gorm:"column:v_534f_cd91;" mom:"naming:s63;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
