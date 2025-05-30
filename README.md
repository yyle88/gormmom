[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/gormmom/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/gormmom/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/gormmom)](https://pkg.go.dev/github.com/yyle88/gormmom)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/gormmom/master.svg)](https://coveralls.io/github/yyle88/gormmom?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.22%2C%201.23-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yyle88/gormmom.svg)](https://github.com/yyle88/gormmom/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/gormmom)](https://goreportcard.com/report/github.com/yyle88/gormmom)

# gormmom

**Empowering Native Language Programming, Simplifying GORM Tag Generation**

---

`gormmom` is a tool designed to automatically generate GORM tags, aimed at helping developers program in their native language while simplifying the process of defining GORM tags. The tool processes struct fields and automatically generates GORM-compliant tags, ensuring field names follow specific naming conventions.

---

## CHINESE README

[中文说明](README.zh.md)

---

## Features

- **Automatic GORM Tag Generation**: Automatically generates GORM tags for struct fields, such as `column`, `index`, `unique`, etc.
- **Native Language Programming Support**: Allows developers to define struct fields in their native language (e.g., Chinese), reducing the difficulty of understanding business.

---

## Installation

```bash
go get github.com/yyle88/gormmom
```

---

## Usage Example

### Original Code (Native Language Programming)
```go
type Example struct {
    V证号 string `gorm:"primaryKey"`
    V姓名 string `gorm:"index"`
    V年龄 int    `gorm:"unique"`
    V性别 bool   `gorm:"column:sex;uniqueIndex" mom:"mcp:S63"`
}
```

### Generated Code (Automatic GORM Tag Generation)
```go
type Example struct {
    V证号 string `gorm:"column:v_c18b_f753;primaryKey" mom:"mcp:s63;"`
    V姓名 string `gorm:"column:v_d359_0d54;index:idx_example_v_d359_0d54" mom:"mcp:s63;idx:cnm;"`
    V年龄 int    `gorm:"column:v_745e_849f;unique" mom:"mcp:s63;"`
    V性别 bool   `gorm:"column:V_2760_2B52;uniqueIndex:udx_example_V_2760_2B52" mom:"mcp:S63;udx:cnm;"`
}
```

### Gorm Select Usage(Select with gorm repo)
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

## Configuration Options

- **mcp**: Configures the naming convention for database column names.
- **idx**: Configures the naming convention for single-column indexes.
- **udx**: Configures the naming convention for single-column unique indexes.

---

## Design Ideas

[README OLD DOC](internal/docs/README_OLD_DOC.en.md)

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repo on GitHub (using the webpage interface).
2. Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. Navigate to the cloned project (`cd repo-name`)
4. Create a feature branch (`git checkout -b feature/xxx`).
5. Stage changes (`git add .`)
6. Commit changes (`git commit -m "Add feature xxx"`).
7. Push to the branch (`git push origin feature/xxx`).
8. Open a pull request on GitHub (on the GitHub webpage).

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

If you find this package valuable, give me some stars on GitHub! Thank you!!!

**Thank you for your support!**

**Happy Coding with `gormmom`!** 🎉

Give me stars. Thank you!!!

## GitHub Stars

[![starring](https://starchart.cc/yyle88/gormmom.svg?variant=adaptive)](https://starchart.cc/yyle88/gormmom)
