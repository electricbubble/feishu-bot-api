package fsBotAPI

import "fmt"

// TextMentionAll @所有人
func TextMentionAll() string {
	return `<at user_id="all"></at>`
}

// TextMentionByOpenID @单个用户
//
// 如果 Open ID 无效，则取 name 展示
//
// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#e1cdee9f
func TextMentionByOpenID(id string, name ...string) string {
	var s string
	if len(name) != 0 {
		s = name[0]
	}
	return fmt.Sprintf(`<at user_id="%s">%s</at>`, id, s)
}

func newMsgText(text string) map[string]interface{} {
	msgType := "text"
	return map[string]interface{}{
		"msg_type": msgType,
		"content": map[string]string{
			msgType: text,
		},
	}
}
