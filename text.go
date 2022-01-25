package fsBotAPI

import "fmt"

// StrMentionAll @所有人
func StrMentionAll() string {
	return `<at user_id="all"></at>`
}

// StrMentionByOpenID @单个用户
//  如果 Open ID 无效，则取 name 展示
func StrMentionByOpenID(id string, name ...string) string {
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
