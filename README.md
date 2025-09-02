[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/gormmom/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/gormmom/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/gormmom)](https://pkg.go.dev/github.com/yyle88/gormmom)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/gormmom/main.svg)](https://coveralls.io/github/yyle88/gormmom?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.22%2C%201.23-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yyle88/gormmom.svg)](https://github.com/yyle88/gormmom/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/gormmom)](https://goreportcard.com/report/github.com/yyle88/gormmom)

# 🌍 GORMMOM - Native Language Programming Revolution with GORM

**gormmom** is the **native language programming engine** that breaks down language barriers in database development. As the **smart tag generation engine** of the GORM ecosystem, it empowers teams worldwide to write database models in native languages while automatically generating database-compatible GORM tags and column names.

> 🎯 **Language Liberation**: Code in Chinese, Arabic, Japanese, and various languages - gormmom bridges the gap between human expression and database requirements.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## 🌟 The Problem & Solution

### ⚡ Standard Approach
```go
// ❌ Common approach: Developers constrained to English naming
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

### ✅ GORMMOM Solution
```go
// ✅ GORMMOM: Program in native language!
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
    return "users"  // Database-compatible table name
}
```

## 🌍 Multi-Language Examples

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

## 🏗️ GORMMOM in the Ecosystem

```
┌─────────────────────────────────────────────────────────────────────┐
│                    GORM Type-Safe Ecosystem                         │
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

**GORMMOM** is the **native language bridge** that enables worldwide teams to participate in the type-safe GORM ecosystem.

---

## 🚀 Installation

```bash
go get github.com/yyle88/gormmom
```

---

## 💡 Core Features

### 🌍 Language Support
- **Unicode compatible**: Chinese, Japanese, Korean, and other languages
- **Smart mapping**: Native names to database-safe column names
- **Deterministic output**: Consistent generation across runs
- **Go export rules**: Uppercase prefixes (e.g., `U用户名`, `Uユーザー名`)

### 🤖 Smart Generation
- **Auto GORM tags**: `column`, `index`, `uniqueIndex` and more
- **Naming patterns**: `snake_case`, `camelCase`, custom formats
- **Index management**: Descriptive names with auto-generation
- **Tag preservation**: Maintains existing tags while adding new ones

### ⚠️ Requirements
- **Struct naming**: Start with uppercase `type T用户` (Go export rules)
- **TableName() method**: Required to avoid `regexp does not match` errors
- **Two-pass generation**: First run updates tags, second run generates methods

---

## 🛠️ Configuration and Usage

### Basic Configuration

```go
// Configuration is done via Options in test file:
options := gormmom.NewOptions()
// Various configuration options available

// With gormcngen:
options := gormcngen.NewOptions().
    WithColumnClassExportable(true).     // Export T用户Columns
    WithColumnsMethodRecvName("T").      // Method name
    WithColumnsCheckFieldType(true)      // Type checking
```

### Usage

Gormmom works through test files to generate database tags and type-safe methods. See the [examples](internal/examples/) directories with complete working examples.


---

## 🌟 Benefits

| Standard Approach | GORMMOM Approach |
|---------------------|------------------|
| ❌ **Foreign language fields** | ✅ **Native language fields** |
| ❌ **Lost business context** | ✅ **Business semantics preserved** |
| ❌ **Translation overhead** | ✅ **Direct domain modeling** |
| ❌ **Context disconnect** | ✅ **Natural alignment** |
| ❌ **Team communication gaps** | ✅ **Unified understanding** |

---

## 🎯 Getting Started

### 1. Define Native Language Model
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

### 2. Generate Tags and Methods
Create a test file to run the generation process. See [examples](internal/examples/) directories with complete setup.

```bash
# Step 1: Generate GORM tags for native language fields
go test -v -run TestGen/GenGormMom

# Step 2: Generate type-safe column methods
go test -v -run TestGen/GenGormCnm
```

### 3. Use Type-Safe Queries
```go
var user T用户
cls := user.Columns()

err := db.Where(cls.U用户名.Eq("张三")).First(&user).Error
```

### 4. Complete Example
```go
// The struct now has database-compatible tags
// Use it with standard GORM operations
db.AutoMigrate(&T用户{})

// Create records
user := T用户{U用户名: "张三", E邮箱: "zhang@example.com", A年龄: 25, D电话: "13800138000", J住所: "北京市海淀区", S状态: "活跃"}
db.Create(&user)

// Operations with type-safe columns
var users []T用户
cls := (&T用户{}).Columns()
db.Where(cls.A年龄.Gt(18)).Find(&users)
```

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-08-28 08:33:43.829511 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share your use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize by reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo for new releases and features
- 🌟 **Success stories?** Share how this package improved your workflow
- 💬 **General feedback?** All suggestions and comments are welcome

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage interface).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement your changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation for user-facing changes and use meaningful commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a pull request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Happy Coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yyle88/gormmom.svg?variant=adaptive)](https://starchart.cc/yyle88/gormmom)

---

## 🔗 Related Projects

- 🏗️ **[gormcnm](https://github.com/yyle88/gormcnm)** - Type-safe column foundation
- 🤖 **[gormcngen](https://github.com/yyle88/gormcngen)** - Smart code generation
- 🏢 **[gormrepo](https://github.com/yyle88/gormrepo)** - Enterprise storage pattern  
- 🌍 **[gormmom](https://github.com/yyle88/gormmom)** - Native language programming (this package)
