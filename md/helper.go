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

// MentionByID 通过 Open ID 或 User ID @指定人
//
// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#e1cdee9f
func MentionByID(id string) string {
	return fmt.Sprintf(`<at id=%s></at>`, id)
}

// MentionByEmail 通过邮箱 @指定人
//
// https://open.feishu.cn/document/ugTN1YjL4UTN24CO1UjN/uUzN1YjL1cTN24SN3UjN?from=mcb#acc98e1b
func MentionByEmail(email string) string {
	return fmt.Sprintf(`<at email=%s></at>`, email)
}

// MentionAll @所有人
//
// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#e1cdee9f
func MentionAll() string {
	return `<at id=all></at>`
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
