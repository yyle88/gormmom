package example2enus

import "time"

type Example struct {
	ID        int32  `gorm:"column:id; primaryKey;" json:"id"`
	V名称       string `gorm:"column:V_0D54_F079;type:text;uniqueIndex:udx_examples_V_0D54_F079" mom:"naming:S63;udx:cnm;"` //使用大写字母的规则
	V字段       string `gorm:"column:v_575b_b56b;" mom:"naming:s63;"`                                                       //使用小写字母的规则
	V性别       string `gorm:"column:v_2760_2b52;" mom:"naming:s63;"`
	V特殊       string `gorm:"column:V_7972_8A6B;type:int32;index:idx_examples_V_7972_8A6B" mom:"naming:S63;idx:cnm;"`
	V年龄       int    `gorm:"column:v_745e_849f;" json:"age" mom:"naming:s63;"` //理论上不要直接给model添加json标签，因为那是view层的逻辑，但实际上假如非这样做也能处理
	Rank      int32
	V身高       int32     `gorm:"column:V_AB8E_D89A;" mom:"naming:S63;"`
	V体重       int32     `gorm:"column:v_534f_cd91;" mom:"naming:s63;"`
	CreatedAt time.Time `gorm:"autoCreateTime;index;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;uniqueIndex;"`
}