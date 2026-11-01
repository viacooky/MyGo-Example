package app

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

type pageDocumentData struct {
	Markdown markdown.Options
}

const markdownSample = `# MyGo Example

一个展示 **Markdown** 原生渲染的页面，支持标题、列表、表格、引用与代码块。

## 特性

- 标题层级与段落
- **粗体**、*斜体* 与 ~~删除线~~
- 有序与无序列表
- 任务列表
- 表格、引用与分隔线
- 链接与代码块

## 任务

- [x] 渲染标题
- [ ] 支持图片
- [ ] 导出 PDF

## 表格

| 语言 | 用途 | 年份 |
| --- | --- | --- |
| Go | 后端 | 2009 |
| Python | 脚本 | 1991 |
| Rust | 系统 | 2010 |

> 引用：简洁是可靠的前提。

## 代码

~~~go
package main

import "fmt"

func main() {
	fmt.Println("Hello, MyGo")
}
~~~

## 链接

访问 [MyGo 仓库](https://github.com/egoist/mygo) 了解更多。
`

func (app *app) pageDocument(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Document").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
			markdown.View(c, markdownSample, &app.pageDocumentData.Markdown)
		})
	})
}
