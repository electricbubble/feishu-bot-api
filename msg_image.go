package feishu_bot_api

var _ Message = (*imageMessage)(nil)

type imageMessage string

func (m imageMessage) Apply(body *MessageBody) error {
	body.MsgType = "image"
	body.Content = &MessageBodyContent{
		ImageKey: string(m),
	}

	return nil
}
