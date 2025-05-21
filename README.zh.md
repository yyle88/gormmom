# gormmom

**赋能母语编程，简化 GORM 标签生成**

---

`gormmom` 是一个用于自动生成 GORM 标签的工具，旨在帮助开发者在编写 Go 代码时，使用母语进行编程，同时简化 GORM 标签的定义。该工具通过对结构体字段进行处理，自动生成符合 GORM 规范的标签，并确保字段名符合特定命名规则。

---

## 英文文档

[ENGLISH README](README.md)

---

## 特性

- **自动生成 GORM 标签**：根据字段定义自动生成 GORM 的标签，例如 `column`、`index`、`unique` 等。
- **母语编程支持**：允许开发者使用母语（如中文）来定义结构体字段，能够降低业务的理解难度。

---

## 安装

```bash
go get github.com/yyle88/gormmom
```

---

## 使用示例

### 原始代码（母语编程）
```go
type Example struct {
    V证号 string `gorm:"primaryKey"`
    V姓名 string `gorm:"index"`
    V年龄 int    `gorm:"unique"`
    V性别 bool   `gorm:"column:sex;uniqueIndex" mom:"mcp:S63"`
}
```

### 输出代码（自动生成 GORM 标签）
```go
type Example struct {
    V证号 string `gorm:"column:v_c18b_f753;primaryKey" mom:"mcp:s63;"`
    V姓名 string `gorm:"column:v_d359_0d54;index:idx_example_v_d359_0d54" mom:"mcp:s63;idx:cnm;"`
    V年龄 int    `gorm:"column:v_745e_849f;unique" mom:"mcp:s63;"`
    V性别 bool   `gorm:"column:V_2760_2B52;uniqueIndex:udx_example_V_2760_2B52" mom:"mcp:S63;udx:cnm;"`
}
```

### 查询语句 (使用 repo 查询数据)
```go
example, err := repo.Repo(caseDB).First(func(db *gorm.DB, cls *exampleColumns) *gorm.DB {
    return db.Where(cls.V名称.Eq(name))
})
```

```go
results, err := repo.Repo(caseDB).Find(func(db *gorm.DB, cls *exampleColumns) *gorm.DB {
    return db.Where(cls.V名称.In(names))
})
```

---

## 配置选项

- **mcp**：配置数据表列名的命名规则。
- **idx**：配置单键普通索引的命名规则。
- **udx**：配置单键唯一索引的命名规则。

---

## 设计思路

[旧版说明](internal/docs/README_OLD_DOC.zh.md)

---

## 许可

项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE)。

## 贡献与支持

欢迎通过提交 pull request 或报告问题来贡献此项目。

如果你觉得这个包对你有帮助，请在 GitHub 上给个 ⭐，感谢支持！！！

**感谢你的支持！**

**祝编程愉快！** 🎉

Give me stars. Thank you!!!
