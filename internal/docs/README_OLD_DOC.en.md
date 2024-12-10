# gormmom
**Empowering Native Language Programming**

---

Huawei recently introduced the Cangjie programming language. However, it still lacks support for programming directly in Chinese, despite the growing demand for such a feature.

---

In Go programming, it's actually possible to write code using your native language.

For example:
```go
type A学生信息 struct {
    V姓名 string
    V性别 bool
    V年龄 int
}

func (a *A客户信息) Get姓名() string { return a.V姓名 }
```

This approach works seamlessly.

---

However, when working with GORM, challenges arise because column names need to be explicitly defined using GORM tags.

This tool offers a simple solution to help you easily generate GORM tags, making it possible to continue programming in your native language while working with Go.

---

Due to my limited proficiency in English and the difficulty in finding straightforward, descriptive words, I named this package **"gormmom."**

While mastering English is essential for programmers, I believe that even individuals with advanced English skills—whether they’ve passed CET-4/6 (China's English proficiency tests) or have over ten years of programming experience—can still struggle to express themselves fluently in code. In such cases, it may be time to consider an alternative approach.

---

For example, during this project, I encountered difficulty when trying to extract content from tags. The recommended term was `extract`, but I initially only thought of alternatives such as `get`, `parse`, or `obtain`.

This experience reinforced my belief in the necessity of programming in your native language. Without it, code ends up as a patchwork of imprecise terms, which can lead to confusion and inefficiency.

---

Another issue I frequently face is that in English, opposite meanings are often expressed with words of varying lengths. For instance:  
`up/down`, `left/right`, `open/close`, `start/stop`, `public/private`.

Programmers are often forced to choose from a limited set of words that only loosely fit the context, such as:
- `get/set`
- `start/close`
- `create/update/select/delete`.

For non-native English speakers, this inconsistency can create confusion and hinder clarity.

---

Currently, this tool does not include a translation feature and cannot automatically convert native language code into English.

It focuses on simplifying the process of generating GORM tags.

In the future, we may add translation capabilities for multiple languages or even introduce pinyin-based encoding. However, there are no immediate plans for these features, as the current functionality is sufficient for most use cases.

---

Huawei missed a significant opportunity by not fully leveraging the potential of the "Cangjie" name.

Chinese is a language with concise semantics, where characters don’t require separators, thus avoiding issues like camelCase or underscores. This makes it particularly well-suited for a programming language.

Perhaps, in the future, someone will develop a fully functional Chinese programming language (and perhaps such efforts are already underway).

---

**gormmom** helps bridge the gap, enabling developers to write meaningful, efficient code in their native language while utilizing GORM in Go.

--- 
