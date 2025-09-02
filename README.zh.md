[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/gormmom/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/gormmom/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/gormmom)](https://pkg.go.dev/github.com/yyle88/gormmom)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/gormmom/main.svg)](https://coveralls.io/github/yyle88/gormmom?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.22%2C%201.23-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yyle88/gormmom.svg)](https://github.com/yyle88/gormmom/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/gormmom)](https://goreportcard.com/report/github.com/yyle88/gormmom)

# 🌍 GORMMOM - GORM 原生语言编程革命

**gormmom** 是 **原生语言编程引擎**，打破数据库开发中的语言壁垒。作为 GORM 生态系统的 **智能标签生成引擎**，它赋能全球团队使用原生语言编写数据库模型，同时自动生成数据库兼容的 GORM 标签和列名。

> 🎯 **语言解放**: 用中文、阿拉伯语、日语和各种语言编程 - gormmom 架起人类表达与数据库需求之间的桥梁。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 🌟 原生语言编程革命

### ⚡ 全球开发挑战

**常见限制** - 英语编程：
```go
// ❌ 常见方法：开发者被限制在英语命名
type User struct {
    ID          uint      `gorm:"primaryKey"`
    Username    string    `gorm:"column:username;uniqueIndex"`
    Email       string    `gorm:"column:email;index"`
    Age         int       `gorm:"column:age"`
    PhoneNumber string    `gorm:"column:phone_number"`
    Address     string    `gorm:"column:address"`
    Status      string    `gorm:"column:status;index"`
}
```

**GORMMOM 解决方案** - 真正的原生语言编程：
```go
// ✅ GORMMOM: 用原生语言编程！
type T用户 struct {
    ID    uint   `gorm:"primaryKey"`
    U用户名 string `gorm:"uniqueIndex"`
    E邮箱  string `gorm:"index"`
    A年龄  int    `gorm:""`
    D电话  string `gorm:""`
    J住所  string `gorm:""`
    S状态  string `gorm:"index"`
}

func (*T用户) TableName() string {
    return "users"  // 数据库兼容的表名
}
```

## 🌍 多语言示例

### 繁體中文
```go
type T用戶 struct {
    ID    uint      `gorm:"primaryKey"`
    U用戶名 string    `gorm:"uniqueIndex"`
    E郵箱  string    `gorm:"index"`
    A年齡  int       `gorm:""`
    D電話  string    `gorm:""`
    J住所  string    `gorm:""`
    S狀態  string    `gorm:"index"`
}

func (*T用戶) TableName() string {
    return "users"
}
```

### 日本語
```go
type Tユーザー struct {
    ID         uint      `gorm:"primaryKey"`
    Uユーザー名  string    `gorm:"uniqueIndex"`
    Eメール     string    `gorm:"index"`
    A年齢       int       `gorm:""`
    D電話       string    `gorm:""`
    J住所       string    `gorm:""`
    Sステータス  string    `gorm:"index"`
}

func (*Tユーザー) TableName() string {
    return "users"
}
```

### 한국어
```go
type T사용자 struct {
    ID      uint      `gorm:"primaryKey"`
    U사용자명 string    `gorm:"uniqueIndex"`
    E이메일  string    `gorm:"index"`
    A나이    int       `gorm:""`
    J전화    string    `gorm:""`
    J주소    string    `gorm:""`
    S상태    string    `gorm:"index"`
}

func (*T사용자) TableName() string {
    return "users"
}
```

---

## 🏗️ 生态系统中的 GORMMOM

```
┌─────────────────────────────────────────────────────────────────────┐
│                    GORM 类型安全生态系统                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │
│  │  gormzhcn   │    │  gormmom    │    │  gormrepo   │              │
│  │ Chinese API │───▶│ Native Lang │───▶│  Package    │─────┐        │
│  │  Localize   │    │  Smart Tags │    │  Pattern    │     │        │
│  └─────────────┘    └─────────────┘    └─────────────┘     │        │
│         │                   │                              │        │
│         │                   ▼                              ▼        │
│         │            ┌─────────────┐              ┌─────────────┐   │
│         │            │ gormcngen   │              │Application  │   │
│         │            │Code Generate│─────────────▶│Custom Code  │   │
│         │            │AST Operation│              │             │   │
│         │            └─────────────┘              └─────────────┘   │
│         │                   │                              ▲        │
│         │                   ▼                              │        │
│         └────────────▶┌─────────────┐◄─────────────────────┘        │
│                       │   GORMCNM   │                               │
│                       │ FOUNDATION  │                               │
│                       │ Type-Safe   │                               │
│                       │ Core Logic  │                               │
│                       └─────────────┘                               │
│                              │                                      │
│                              ▼                                      │
│                       ┌─────────────┐                               │
│                       │    GORM     │                               │
│                       │  Database   │                               │
│                       └─────────────┘                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**GORMMOM** 是 **原生语言桥梁**，让全球团队参与类型安全的 GORM 生态系统。

---

## 🚀 安装

```bash
go get github.com/yyle88/gormmom
```

---

## 💡 核心功能

### 🌍 语言支持
- **Unicode 兼容**：中文、日文、韩文和其他语言
- **智能映射**：原生名称转换为数据库安全列名
- **确定性输出**：跨运行的一致生成
- **Go 导出规则**：大写前缀（如 `U用户名`、`Uユーザー名`）

### 🤖 智能生成
- **自动 GORM 标签**：`column`、`index`、`uniqueIndex` 等
- **命名模式**：`snake_case`、`camelCase`、自定义格式
- **索引管理**：自动生成描述性名称
- **标签保留**：保留现有标签同时添加新标签

### ⚠️ 要求
- **结构体命名**：以大写开头 `type T用户`（Go 导出规则）
- **TableName() 方法**：必需，避免 `regexp does not match` 错误
- **两遍生成**：第一次运行更新标签，第二次运行生成方法

---

## 🛠️ 配置和使用

### 基本配置

```go
// 配置通过测试文件中的选项完成：
options := gormmom.NewOptions()
// 各种配置选项可用

// 配合 gormcngen：
options := gormcngen.NewOptions().
    WithColumnClassExportable(true).     // 导出 T用户Columns
    WithColumnsMethodRecvName("T").      // 方法名
    WithColumnsCheckFieldType(true)      // 类型检查
```

### 使用方法

Gormmom 通过测试文件生成数据库标签和类型安全方法。查看 [examples](internal/examples/) DIR 获取完整的工作示例。

---

## 🌟 优势

| 标准方法 | GORMMOM 方法 |
|---------------------|------------------|
| ❌ **外语字段** | ✅ **母语字段** |
| ❌ **丢失业务上下文** | ✅ **业务语义保留** |
| ❌ **翻译开销** | ✅ **直接领域建模** |
| ❌ **上下文断开** | ✅ **自然对齐** |
| ❌ **团队沟通差距** | ✅ **统一理解** |

---

## 🎯 入门指南

### 1. 定义原生语言模型
```go
type T用户 struct {
    ID    uint      `gorm:"primaryKey"`
    U用户名 string    `gorm:"uniqueIndex"`
    E邮箱  string    `gorm:"index"`
    A年龄  int       `gorm:""`
    D电话  string    `gorm:""`
    J住所  string    `gorm:""`
    S状态  string    `gorm:"index"`
}

func (*T用户) TableName() string {
    return "users"
}
```

### 2. 生成标签和方法
创建测试文件运行生成过程。查看 [examples](internal/examples/) DIR 获取完整设置。

```bash
# 步骤 1：生成原生语言字段的 GORM 标签
go test -v -run TestGen/GenGormMom

# 步骤 2：生成类型安全列方法
go test -v -run TestGen/GenGormCnm
```

### 3. 使用类型安全查询
```go
var user T用户
cls := user.Columns()

err := db.Where(cls.U用户名.Eq("张三")).First(&user).Error
```

### 4. 完整示例
```go
// 结构体现在具有数据库兼容标签
// 与标准 GORM 操作一起使用
db.AutoMigrate(&T用户{})

// 创建记录
user := T用户{U用户名: "张三", E邮箱: "zhang@example.com", A年龄: 25, D电话: "13800138000", J住所: "北京市海淀区", S状态: "活跃"}
db.Create(&user)

// 使用类型安全列操作
var users []T用户
cls := (&T用户{}).Columns()
db.Where(cls.A年龄.Gt(18)).Find(&users)
```

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-08-28 08:33:43.829511 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **意见反馈？** 欢迎所有建议和宝贵意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Pull Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Pull Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**使用这个包快乐编程！** 🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/yyle88/gormmom.svg?variant=adaptive)](https://starchart.cc/yyle88/gormmom)

---

## 🔗 相关项目

- 🏗️ **[gormcnm](https://github.com/yyle88/gormcnm)** - 类型安全列基础
- 🤖 **[gormcngen](https://github.com/yyle88/gormcngen)** - 智能代码生成
- 🏢 **[gormrepo](https://github.com/yyle88/gormrepo)** - 企业存储模式  
- 🌍 **[gormmom](https://github.com/yyle88/gormmom)** - 原生语言编程（本包）