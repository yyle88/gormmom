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
    V性别 bool   `gorm:"column:sex;uniqueIndex" mom:"naming:S63"`
}
```

### Generated Code (Automatic GORM Tag Generation)
```go
type Example struct {
    V证号 string `gorm:"column:v_c18b_f753;primaryKey" mom:"naming:s63;"`
    V姓名 string `gorm:"column:v_d359_0d54;index:idx_example_v_d359_0d54" mom:"naming:s63;idx:cnm;"`
    V年龄 int    `gorm:"column:v_745e_849f;unique" mom:"naming:s63;"`
    V性别 bool   `gorm:"column:V_2760_2B52;uniqueIndex:udx_example_V_2760_2B52" mom:"naming:S63;udx:cnm;"`
}
```

---

## Configuration Options

- **naming**: Configures the naming convention for database column names.
- **idx**: Configures the naming convention for single-column indexes.
- **udx**: Configures the naming convention for single-column unique indexes.

---

## Design Ideas

[README OLD DOC](internal/docs/README_OLD_DOC.en.md)

---

## License

`gormmom` is open-source and released under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

## Support

Welcome to contribute to this project by submitting pull requests or reporting issues.

If you find this package helpful, give it a star on GitHub!

**Thank you for your support!**

**Happy Coding with `gormmom`!** 🎉

Give me stars. Thank you!!!

## See stars
[![see stars](https://starchart.cc/yyle88/gormmom.svg?variant=adaptive)](https://starchart.cc/yyle88/gormmom)
