package feishu_bot_api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var _ Message = (*richTextMessage)(nil)

type richTextMessage []*RichTextBuilder

func (m richTextMessage) Apply(body *MessageBody) error {
	raw, err := m.marshal()
	if err != nil {
		return fmt.Errorf("rich text message: %w", err)
	}

	body.MsgType = "post"
	body.Content = &MessageBodyContent{
		Post: &raw,
	}
	return nil
}

func (m richTextMessage) marshal() (json.RawMessage, error) {
	var buf bytes.Buffer
	buf.Grow(len(m) * 48)

	buf.WriteString("{")

	nnm := make(richTextMessage, 0, len(m))
	for i := range m {
		if m[i] == nil {
			continue
		}
		nnm = append(nnm, m[i])
	}

	commas := len(nnm)
	for _, rtb := range nnm {
		commas--

		if rtb == nil {
			continue
		}

		bs, err := json.Marshal(rtb.body)
		if err != nil {
			return nil, fmt.Errorf("marshal(%s): %w", rtb.language, err)
		}
		s := fmt.Sprintf(`"%s":%s`, _quoteEscaper.Replace(string(rtb.language)), bs)

		if commas > 0 {
			s += ","
		}

		buf.WriteString(s)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

// --------------------------------------------------------------------------------

type (
	RichTextBuilder struct {
		language Language
		body     richTextBody
	}

	richTextBody struct {
		Title   string              `json:"title,omitempty"`
		Content []richTextParagraph `json:"content,omitempty"`
	}
	richTextParagraph []richTextLabel

	richTextLabel struct {
		Tag  string `json:"tag"`
		Text string `json:"text,omitempty"`

		// 仅 文本标签(text) 使用；表示是否 unescape 解码。默认值为 false，未用到 unescape 时可以不填
		UnEscape *bool `json:"un_escape,omitempty"`

		// 仅 超链接标签(a) 使用；链接地址，需要确保链接地址的合法性，否则消息会发送失败
		Href string `json:"href,omitempty"`

		// 仅 @ 标签(at) 使用
		UserID string `json:"user_id,omitempty"`
		// 仅 @ 标签(at) 使用
		UserName string `json:"user_name,omitempty"`

		// 仅 图片标签(img) 使用；图片的唯一标识。可通过 上传图片 接口获取 image_key
		// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/image/create
		ImageKey string `json:"image_key,omitempty"`
	}
)

func NewRichText(language Language, title string) *RichTextBuilder {
	rtb := &RichTextBuilder{
		language: language,
		body: richTextBody{
			Title:   title,
			Content: make([]richTextParagraph, 0, 1),
		},
	}
	return rtb
}

// Text 文本标签
//
// unEscape 表示是否 unescape 解码，未用到 unescape 时传入 false
func (rtb *RichTextBuilder) Text(text string, unEscape bool) *RichTextBuilder {
	lbl := richTextLabel{
		Tag:      "text",
		Text:     text,
		UnEscape: &unEscape,
	}
	paragraph := rtb.lastParagraph()
	paragraph = append(paragraph, lbl)
	rtb.updateLastParagraph(paragraph)
	return rtb
}

// Hyperlink 超链接标签
func (rtb *RichTextBuilder) Hyperlink(text, href string) *RichTextBuilder {
	lbl := richTextLabel{
		Tag:  "a",
		Text: text,
		Href: href,
	}
	paragraph := rtb.lastParagraph()
	paragraph = append(paragraph, lbl)
	rtb.updateLastParagraph(paragraph)
	return rtb
}

// At @ 标签
//
// id: 用户的 Open ID 或 User ID
//
//	@ 单个用户时，id 字段必须是有效值（仅支持 @ 自定义机器人所在群的群成员）
//	@ 所有人时，填 all (也可以使用 RichTextBuilder.AtEveryone)
func (rtb *RichTextBuilder) At(id, name string) *RichTextBuilder {
	lbl := richTextLabel{
		Tag:      "at",
		UserID:   id,
		UserName: name,
	}
	paragraph := rtb.lastParagraph()
	paragraph = append(paragraph, lbl)
	rtb.updateLastParagraph(paragraph)
	return rtb
}

func (rtb *RichTextBuilder) AtEveryone() *RichTextBuilder {
	lbl := richTextLabel{
		Tag:    "at",
		UserID: "all",
	}
	paragraph := rtb.lastParagraph()
	paragraph = append(paragraph, lbl)
	rtb.updateLastParagraph(paragraph)
	return rtb
}

// Image 图片标签
//
// 图片的唯一标识。可通过 上传图片 接口获取
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/image/create
func (rtb *RichTextBuilder) Image(imgKey string) *RichTextBuilder {
	lbl := richTextLabel{
		Tag:      "img",
		ImageKey: imgKey,
	}
	paragraph := rtb.lastParagraph()
	paragraph = append(paragraph, lbl)
	rtb.updateLastParagraph(paragraph)
	return rtb
}

func (rtb *RichTextBuilder) lastParagraph() richTextParagraph {
	if len(rtb.body.Content) == 0 {
		rtb.body.Content = append(rtb.body.Content, make(richTextParagraph, 0, 4))
	}

	return rtb.body.Content[len(rtb.body.Content)-1]
}

func (rtb *RichTextBuilder) updateLastParagraph(paragraph richTextParagraph) {
	rtb.body.Content[len(rtb.body.Content)-1] = paragraph
}
