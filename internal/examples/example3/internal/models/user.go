package models

// 文档中的韩语示例模型（原始版本）
/*
type T사용자 struct {
    ID      uint   `gorm:"primaryKey"`
    U사용자명 string `gorm:"uniqueIndex"`
    E이메일  string `gorm:"index"`
    N나이    int    `gorm:""`
    J전화    string `gorm:""`
    J주소    string `gorm:""`
    S상태    string `gorm:"index"`
}
*/

// 经过 gormmom 生成器处理后的韩语模型
type T사용자 struct {
	ID    uint   `gorm:"primaryKey"`
	U사용자명 string `gorm:"column:u_acc0_a9c6_90c7_85ba;uniqueIndex:udx_users_u_acc0_a9c6_90c7_85ba" mom:"mcp:s63;udx:cnm;"`
	E이메일  string `gorm:"column:e_74c7_54ba_7cc7;index:idx_users_e_74c7_54ba_7cc7" mom:"mcp:s63;idx:cnm;"`
	N나이   int    `gorm:"column:n_98b0_74c7;" mom:"mcp:s63;"`
	J전화   string `gorm:"column:j_04c8_54d6;" mom:"mcp:s63;"`
	J주소   string `gorm:"column:j_fcc8_8cc1;" mom:"mcp:s63;"`
	S상태   string `gorm:"column:s_c1c0_dcd0;index:idx_users_s_c1c0_dcd0" mom:"mcp:s63;idx:cnm;"`
}

func (*T사용자) TableName() string {
	return "users"
}
