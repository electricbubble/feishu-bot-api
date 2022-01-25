package fsBotAPI

func newMsgShareChat(chatID string) map[string]interface{} {
	return map[string]interface{}{
		"msg_type": "share_chat",
		"content": map[string]string{
			"share_chat_id": chatID,
		},
	}
}
