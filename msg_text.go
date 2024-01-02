package feishu_bot_api

import "fmt"

var _ Message = (*textMessage)(nil)

type textMessage string

func (m textMessage) Apply(body *MessageBody) error {
	body.MsgType = "text"
	body.Content = &MessageBodyContent{
		Text: string(m),
	}

	return nil
}

// TextAtPerson @指定人
//
// 可填入用户的 Open ID 或 User ID，且必须是有效值（仅支持 @ 自定义机器人所在群的群成员），否则取名字展示，并不产生实际的 @ 效果
//
// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#756b882f
// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#272b1dee
func TextAtPerson(id, name string) string {
	return fmt.Sprintf(`<at user_id="%s">%s</at>`, id, name)
}

// TextAtEveryone @所有人
//
// 必须满足所在群开启 @ 所有人功能
func TextAtEveryone() string {
	return `<at user_id="all">everyone</at>`
}
