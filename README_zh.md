# pdfx

一个类似 Git 风格的 PDF 命令行工具，用于合并、拆分和提取 PDF 页面。

## 功能

- **merge**：将目录中的所有 PDF 文件合并为一个文件
- **split**：在指定页码处将 PDF 拆分为两个文件
- **extract**：提取 PDF 的指定页面并合并为一个新文件

## 构建

```bash
go mod tidy
go build -o pdfx.exe .
```

## 用法

```bash
# 合并当前目录下的所有 PDF
./pdfx.exe merge

# 拆分 PDF（从第 N 页开始分成第二份）
./pdfx.exe split -n 输入文件 -f 页码

# 提取指定页面
./pdfx.exe extract -n 输入文件 -p 1,2,3
```

## 命令详解

### merge

合并指定目录中的所有 PDF 文件。

```bash
./pdfx.exe merge [flags]
```

| 参数 | 说明 | 是否必须 | 默认值 |
|------|------|----------|--------|
| `--dir, -d` | 包含 PDF 文件的目录 | 否 | 当前目录 |
| `--output, -o` | 输出文件路径 | 否 | `merged.pdf` |

**示例：**

```bash
# 合并当前目录下的所有 PDF，输出为 merged.pdf
./pdfx.exe merge

# 合并指定目录
./pdfx.exe merge -d ./docs

# 合并并指定输出文件名
./pdfx.exe merge -d ./docs -o 合并结果.pdf
```

### split

在指定页码处将 PDF 拆分为两个文件。

```bash
./pdfx.exe split [flags]
```

| 参数 | 说明 | 是否必须 | 默认值 |
|------|------|----------|--------|
| `--name, -n` | 输入 PDF 文件路径 | 是 | - |
| `--from, -f` | 第二份文件的起始页码（从 1 开始） | 是 | - |

**输出文件命名规则：** `原文件名_起始页-结束页.pdf`

**示例：**

```bash
# 从第 16 页开始拆分
./pdfx.exe split -n 线性代数-基础600题 -f 16

# 输出：
# 线性代数-基础600题_1-15.pdf
# 线性代数-基础600题_16-56.pdf
```

### extract

提取 PDF 的指定页面并合并为一个新文件。

```bash
./pdfx.exe extract [flags]
```

| 参数 | 说明 | 是否必须 | 默认值 |
|------|------|----------|--------|
| `--name, -n` | 输入 PDF 文件路径 | 是 | - |
| `--pages, -p` | 要提取的页码，逗号分隔，如 1,2,3,4 | 是 | - |
| `--output, -o` | 输出文件路径 | 否 | `<输入名>_extracted.pdf` |

**示例：**

```bash
# 提取第 1、2、3 页
./pdfx.exe extract -n 文档.pdf -p 1,2,3

# 指定输出文件名
./pdfx.exe extract -n 文档.pdf -p 1,2,3 -o 摘要.pdf
```

## 注意事项

- 输入文件路径无需手动添加 `.pdf` 后缀，工具会自动补全
- 输出文件路径也无需手动添加 `.pdf` 后缀，工具会自动补全
- `split` 命令的 `-f` 参数表示第二份文件的**起始页码**，即第一份文件包含 `1` 到 `from-1` 页
