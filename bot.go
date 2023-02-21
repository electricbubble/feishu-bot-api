package fsBotAPI

import (
	"fmt"
	"github.com/electricbubble/xhttpclient"
	"strings"
	"time"
)

type bot struct {
	cli       *xhttpclient.XClient
	webhook   string
	secretKey string
}

func NewBot(wh string, opts ...BotOption) Bot {
	wh = strings.TrimSpace(wh)
	b := &bot{cli: xhttpclient.NewClient()}
	if !strings.Contains(wh, "open.feishu.cn") {
		b.webhook = fmt.Sprintf(_fmtWebhook, wh)
	} else {
		b.webhook = wh
	}
	for _, fn := range opts {
		fn(b)
	}
	return b
}

func (b *bot) PushText(content string) error {
	return b.pushMsg(newMsgText(content))
}

func (b *bot) PushPost(p Post, ps ...Post) error {
	return b.pushMsg(newMsgPost(p, ps...))
}

func (b *bot) PushCard(bgColor CardTitleBgColor, cfg CardConfig, c Card, more ...Card) error {
	return b.pushMsg(GenMsgCard(bgColor, cfg, c, more...))
}

func (b *bot) PushImage(imageKey string) error {
	return b.pushMsg(newMsgImage(imageKey))
}

func (b *bot) PushShareChat(chatID string) error {
	return b.pushMsg(newMsgShareChat(chatID))
}

func (b *bot) pushMsg(msg map[string]interface{}) (err error) {
	if b.secretKey != "" {
		ts := time.Now().Unix()
		signed, err := genSign(b.secretKey, ts)
		if err != nil {
			return err
		}
		msg["timestamp"] = ts
		msg["sign"] = signed
	}

	var reply apiResponse
	_, respBody, err := b.cli.Do(&reply, nil, xhttpclient.
		NewPost().
		Path(b.webhook).
		Body(msg),
	)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	if reply.Code != 0 {
		return fmt.Errorf("unknown error: %s", respBody)
	}

	return
}
