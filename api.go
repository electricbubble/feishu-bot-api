package fsBotAPI

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const _fmtWebhook = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

type Bot interface {
	PushText(content string) error
	PushPost(p Post, ps ...Post) error
	PushCard(bgColor CardTitleBgColor, cfg CardConfig, c Card, more ...Card) error
	PushImage(imageKey string) error
	PushShareChat(chatID string) error
}

// 签名校验
//
//	https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#348211be
func genSign(secret string, timestamp int64) (string, error) {
	sign := fmt.Sprintf("%d\n%s", timestamp, secret)

	var data []byte
	h := hmac.New(sha256.New, []byte(sign))
	if _, err := h.Write(data); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

type Language string

const (
	LangChinese  Language = "zh_cn"
	LangEnglish  Language = "en_us"
	LangJapanese Language = "ja_jp"
)

type apiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
