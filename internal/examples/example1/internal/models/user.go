package models

/*
	type T用户 struct {
		ID    uint   `gorm:"primaryKey"`
		U用户名 string `gorm:"uniqueIndex"`
		E邮箱  string `gorm:"index"`
		A年龄  int    `gorm:""`
		D电话  string `gorm:""`
		J住所  string `gorm:""`
		S状态  string `gorm:"index"`
	}
*/
type T用户 struct {
	ID   uint   `gorm:"primaryKey"`
	U用户名 string `gorm:"column:u_2875_3762_0d54;uniqueIndex:udx_users_u_2875_3762_0d54" mom:"mcp:s63;udx:cnm;"`
	E邮箱  string `gorm:"column:e_ae90_b17b;index:idx_users_e_ae90_b17b" mom:"mcp:s63;idx:cnm;"`
	A年龄  int    `gorm:"column:a_745e_849f;" mom:"mcp:s63;"`
	D电话  string `gorm:"column:d_3575_dd8b;" mom:"mcp:s63;"`
	J住所  string `gorm:"column:j_4f4f_4062;" mom:"mcp:s63;"`
	S状态  string `gorm:"column:s_b672_0160;index:idx_users_s_b672_0160" mom:"mcp:s63;idx:cnm;"`
}

func (*T用户) TableName() string {
	return "users"
}
