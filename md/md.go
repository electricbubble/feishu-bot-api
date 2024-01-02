package md

import "fmt"

// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags?lang=zh-CN#abc9b025

func LineBreak() string {
	return "\n"
}

func Italic(s string) string {
	return fmt.Sprintf("*%s*", s)
}

func Bold(s string) string {
	return fmt.Sprintf("**%s**", s)
}

func Strikethrough(s string) string {
	return fmt.Sprintf("~~%s~~", s)
}

// AtPerson @指定人
//
// 自定义机器人仅支持使用 open_id、user_id @指定人
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags?lang=zh-CN#abc9b025
func AtPerson(id, name string) string {
	return fmt.Sprintf("<at id=%s>%s</at>", id, name)
}

func AtEveryone() string {
	return "<at id=all></at>"
}

func Hyperlink(s string) string {
	return fmt.Sprintf("<a href='%s'></a>", s)
}

func TextLink(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

// Image
//   - 仅支持 Markdown 组件
//   - 不支持在 text 元素的 lark_md 类型中使用
func Image(imgKey, hoverText string) string {
	return fmt.Sprintf("![%s](%s)", hoverText, imgKey)
}

func HorizontalRule() string {
	return "\n ---\n"
}

// FeiShuEmoji 飞书表情
//
// https://open.feishu.cn/document/server-docs/im-v1/message-reaction/emojis-introduce
func FeiShuEmoji(emojiKey string) string {
	return fmt.Sprintf(":%s:", emojiKey)
}

func GreenText(s string) string {
	return fmt.Sprintf("<font color='green'>%s</font>", s)
}

func RedText(s string) string {
	return fmt.Sprintf("<font color='red'>%s</font>", s)
}

func GreyText(s string) string {
	return fmt.Sprintf("<font color='grey'>%s</font>", s)
}

type TextTagColor string

const (
	TextTagColorNeutral   TextTagColor = "neutral"
	TextTagColorBlue      TextTagColor = "blue"
	TextTagColorTurquoise TextTagColor = "turquoise"
	TextTagColorLime      TextTagColor = "lime"
	TextTagColorOrange    TextTagColor = "orange"
	TextTagColorViolet    TextTagColor = "violet"
	TextTagColorIndigo    TextTagColor = "indigo"
	TextTagColorWathet    TextTagColor = "wathet"
	TextTagColorGreen     TextTagColor = "green"
	TextTagColorYellow    TextTagColor = "yellow"
	TextTagColorRed       TextTagColor = "red"
	TextTagColorPurple    TextTagColor = "purple"
	TextTagColorCarmine   TextTagColor = "carmine"
)

// TextTag 标签
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags?lang=zh-CN#abc9b025
func TextTag(color TextTagColor, s string) string {
	return fmt.Sprintf("<text_tag color='%s'>%s</text_tag>", color, s)
}
