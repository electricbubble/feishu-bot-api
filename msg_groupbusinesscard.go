package feishu_bot_api

var _ Message = (*groupBusinessCard)(nil)

type groupBusinessCard string

func (m groupBusinessCard) Apply(body *MessageBody) error {
	body.MsgType = "share_chat"
	body.Content = &MessageBodyContent{
		ShareChatID: string(m),
	}

	return nil
}
