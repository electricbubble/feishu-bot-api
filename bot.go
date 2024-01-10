package feishu_bot_api

import (
	"errors"
	"fmt"
	"time"

	"github.com/electricbubble/xhttpclient"
	"golang.org/x/time/rate"
)

type bot struct {
	webhookAccessToken           string
	opts                         *BotOptions
	limiterSecond, limiterMinute *rate.Limiter
	cli                          *xhttpclient.XClient
}

type apiRequest struct {
	MessageBody
	Timestamp int64  `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
}

type apiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (b *bot) SendText(content string) error {
	return b.SendMessage(textMessage(content))
}

func (b *bot) SendRichText(rt *RichTextBuilder, multiLanguage ...*RichTextBuilder) error {
	return b.SendMessage(richTextMessage(append([]*RichTextBuilder{rt}, multiLanguage...)))
}

func (b *bot) SendGroupBusinessCard(chatID string) error {
	return b.SendMessage(groupBusinessCard(chatID))
}

func (b *bot) SendImage(imgKey string) error {
	return b.SendMessage(imageMessage(imgKey))
}

func (b *bot) SendCard(globalConf *CardGlobalConfig, card *CardBuilder, multiLanguage ...*CardBuilder) error {
	return b.SendMessage(cardMessage{
		globalConf: globalConf,
		builders:   append([]*CardBuilder{card}, multiLanguage...),
	})
}
func (b *bot) SendCardViaTemplate(id string, variables any) error {
	return b.SendMessage(cardMessageViaTemplate{id: id, variables: variables})
}

func (b *bot) SendMessage(msg Message) (err error) {
	if b.opts.limiterEnabled() {
		if err := b.wait(); err != nil {
			return err
		}
	}

	var req apiRequest
	if s := b.opts.SecretKey; s != "" {
		req.Timestamp = time.Now().Unix()
		if req.Sign, err = genSignature(req.Timestamp, s); err != nil {
			return fmt.Errorf("gen signature: %w", err)
		}
	}

	if err := msg.Apply(&req.MessageBody); err != nil {
		return fmt.Errorf("apply: %w", err)
	}

	if f := b.opts.HookAfterMessageApply; f != nil {
		if err := f(&req.MessageBody); err != nil {
			return fmt.Errorf("hook(AfterMessageApply): %w", err)
		}
	}

	var resp apiResponse
	_, respBody, err := b.cli.Do(&resp, nil,
		xhttpclient.
			NewPost().
			Path("/open-apis/bot/v2/hook", b.webhookAccessToken).
			Body(req),
	)
	if err != nil {
		return fmt.Errorf("unexpected: %w (resp body: %s)", err, respBody)
	}

	if resp.Code != 0 {
		return fmt.Errorf("api error: %s", respBody)
	}

	return
}

func (b *bot) wait() error {

REDO:

	now := time.Now()

	{
		ts := time.Unix(now.Unix(), 0)

		if b.limiterSecond.TokensAt(ts) <= 0 {
			time.Sleep(ts.Add(time.Second).Sub(now))
			goto REDO
		}

		rs := b.limiterSecond.ReserveN(ts, 1)
		if !rs.OK() {
			return errors.New("limiter(second): not allowed to act")
		}

		switch d := rs.DelayFrom(ts); d {
		case rate.InfDuration:
			return errors.New("limiter(second): cannot grant the token")
		case 0:
		default:
			time.Sleep(d)
		}
	}

	{
		tm := time.Unix(now.Unix()-int64(now.Second()), 0)

		if b.limiterMinute.TokensAt(tm) <= 0 {
			time.Sleep(tm.Add(time.Minute).Sub(now))
			goto REDO
		}

		rm := b.limiterMinute.ReserveN(tm, 1)
		if !rm.OK() {
			return errors.New("limiter(minute): not allowed to act")
		}
		switch d := rm.DelayFrom(tm); d {
		case rate.InfDuration:
			return errors.New("limiter(minute): cannot grant the token")
		case 0:
		default:
			time.Sleep(d)
		}
	}

	return nil
}
