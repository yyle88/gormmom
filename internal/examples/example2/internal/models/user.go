package models

// 文档中的日语示例模型（原始版本）
/*
type Tユーザー struct {
    ID         uint   `gorm:"primaryKey"`
    Uユーザー名  string `gorm:"uniqueIndex"`
    Eメール     string `gorm:"index"`
    A年齢       int    `gorm:""`
    D電話       string `gorm:""`
    J住所       string `gorm:""`
    Sステータス  string `gorm:"index"`
}
*/

// 经过 gormmom 生成器处理后的日语模型
type Tユーザー struct {
	ID     uint   `gorm:"primaryKey"`
	Uユーザー名 string `gorm:"column:u_e630_fc30_b630_fc30_0d54;uniqueIndex:udx_users_u_e630_fc30_b630_fc30_0d54" mom:"mcp:s63;udx:cnm;"`
	Eメール   string `gorm:"column:e_e130_fc30_eb30;index:idx_users_e_e130_fc30_eb30" mom:"mcp:s63;idx:cnm;"`
	A年齢    int    `gorm:"column:a_745e_629f;" mom:"mcp:s63;"`
	D電話    string `gorm:"column:d_fb96_718a;" mom:"mcp:s63;"`
	J住所    string `gorm:"column:j_4f4f_4062;" mom:"mcp:s63;"`
	Sステータス string `gorm:"column:s_b930_c630_fc30_bf30_b930;index:idx_users_s_b930_c630_fc30_bf30_b930" mom:"mcp:s63;idx:cnm;"`
}

func (*Tユーザー) TableName() string {
	return "users"
}
