package fsBotAPI

func newMsgImage(imageKey string) map[string]interface{} {
	return map[string]interface{}{
		"msg_type": "image",
		"content": map[string]string{
			"image_key": imageKey,
		},
	}
}
