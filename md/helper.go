package md

import (
	"fmt"
)

// 仅支持部分
// 语法详情: https://open.feishu.cn/document/ukTMukTMukTM/uADOwUjLwgDM14CM4ATN

func Italics(s string) string {
	return fmt.Sprintf("*%s*", s)
}

func Bold(s string) string {
	return fmt.Sprintf("**%s**", s)
}

func Strikethrough(s string) string {
	return fmt.Sprintf("~~%s~~", s)
}

func Link(url string) string {
	return fmt.Sprintf("<a>%s</a>", url)
}

func TextLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func Image(hoverText, imageKey string) string {
	return "!" + TextLink(hoverText, imageKey)
}

func HorizontalRule() string {
	return ` ---`
}

func MentionByOpenID(id string) string {
	return fmt.Sprintf(`<at id=%s></at>`, id)
}

func ColorGreen(s string) string {
	return fmt.Sprintf(`<font color='green'>%s</font>`, s)
}

func ColorRed(s string) string {
	return fmt.Sprintf(`<font color='red'>%s</font>`, s)
}

func ColorGrey(s string) string {
	return fmt.Sprintf(`<font color='grey'>%s</font>`, s)
}
