package example3

import "time"

type Example2 struct {
	ID        int32  `gorm:"column:id; primaryKey;" json:"id"`
	V名称       string `gorm:"column:V_0D54_F079;type:text" mom:"naming:S63;"` //使用大写字母的规则
	V字段       string `gorm:"column:v_575b_b56b;" mom:"naming:s63;"`          //使用小写字母的规则
	V性别       string `gorm:"column:v_2760_2b52;" mom:"naming:s63;"`
	V特殊       string `gorm:"column:V_7972_8A6B;type:int32" mom:"naming:S63;"`
	V年龄       int    `gorm:"column:v_745e_849f;index:idx_example2s_V_745E_849F" json:"age" mom:"naming:s63;idx:CNM;"` //理论上不要直接给model添加json标签，因为那是view层的逻辑，但实际上假如非这样做也能处理
	Rank      int32
	V身高       int32     `gorm:"column:V_AB8E_D89A;" mom:"naming:S63;"`
	V体重       int32     `gorm:"column:v_534f_cd91;" mom:"naming:s63;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (*Example2) TableName() string {
	return "example2s"
}
