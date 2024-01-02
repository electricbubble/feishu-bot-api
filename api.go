package feishu_bot_api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/electricbubble/xhttpclient"
	"golang.org/x/time/rate"
)

type Bot interface {
	// SendText 发送文本消息
	//
	// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#756b882f
	//
	// @指定人: TextAtPerson
	// @所有人: TextAtEveryone
	SendText(content string) error

	// SendRichText 发送富文本消息
	//
	// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#f62e72d5
	SendRichText(rt *RichTextBuilder, multiLanguage ...*RichTextBuilder) error

	// SendGroupBusinessCard 发送群名片
	//
	// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#897b5321
	//
	// 群 ID 获取方式: https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-id-description
	SendGroupBusinessCard(chatID string) error

	// SendImage 发送图片
	//
	// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#132a114c
	//
	// image_key 获取方式: https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/image/create
	SendImage(imgKey string) error

	// SendCard 发送消息卡片
	//
	// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#4996824a
	SendCard(globalConf *CardGlobalConfig, card *CardBuilder, multiLanguage ...*CardBuilder) error

	// SendCardViaTemplate 使用卡片 ID 发送消息
	//
	// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/send-message-card/send-message-using-card-id
	SendCardViaTemplate(id string, variables any) error

	SendMessage(msg Message) error
}

type Message interface {
	Apply(body *MessageBody) error
}

type MessageBody struct {
	MsgType string              `json:"msg_type"`
	Content *MessageBodyContent `json:"content,omitempty"`
	Card    *json.RawMessage    `json:"card,omitempty"`
}

type MessageBodyContent struct {
	Text        string           `json:"text,omitempty"`
	Post        *json.RawMessage `json:"post,omitempty"`
	ShareChatID string           `json:"share_chat_id,omitempty"`
	ImageKey    string           `json:"image_key,omitempty"`
}

type MessageBodyCard struct {
	Header       *json.RawMessage `json:"header,omitempty"`
	Elements     *json.RawMessage `json:"elements,omitempty"`
	I18nElements *json.RawMessage `json:"i18n_elements"`
	Config       *json.RawMessage `json:"config,omitempty"`
	CardLink     *json.RawMessage `json:"card_link,omitempty"`
}

type MessageBodyCardTemplate struct {
	Type string                      `json:"type"`
	Data MessageBodyCardTemplateData `json:"data"`
}

type MessageBodyCardTemplateData struct {
	TemplateID       string           `json:"template_id"`
	TemplateVariable *json.RawMessage `json:"template_variable,omitempty"`
}

func NewBot(webhook string, opts *BotOptions) Bot {
	if opts == nil {
		opts = &BotOptions{}
	}
	opts.init()

	b := &bot{opts: opts}

	if s := strings.TrimSpace(webhook); strings.Contains(s, "/open-apis/bot") {
		b.webhookAccessToken = path.Base(s)
	} else {
		b.webhookAccessToken = s
	}

	if opts.limiterEnabled() {
		b.limiterSecond = rate.NewLimiter(rate.Limit(opts.LimiterPerSecond), opts.LimiterPerSecond)
		b.limiterMinute = rate.NewLimiter(rate.Every(time.Minute/time.Duration(opts.LimiterPerMinute)), opts.LimiterPerMinute)
	}

	b.cli = xhttpclient.NewClient().BaseURL(opts.BaseURL)

	return b
}

type BotOptions struct {
	BaseURL                            string
	LimiterPerSecond, LimiterPerMinute int
	SecretKey                          string
}

func NewBotOptions() *BotOptions { return &BotOptions{} }

func (opts *BotOptions) init() {
	if strings.TrimSpace(opts.BaseURL) == "" {
		opts.BaseURL = "https://open.feishu.cn"
	}

	if opts.limiterEnabled() {
		if opts.LimiterPerSecond == 0 {
			opts.LimiterPerSecond = 5
		}

		if opts.LimiterPerMinute == 0 {
			opts.LimiterPerMinute = 100
		}
	}

	opts.SecretKey = strings.TrimSpace(opts.SecretKey)
}

func (opts *BotOptions) limiterEnabled() bool {
	if opts.LimiterPerSecond <= -1 || opts.LimiterPerMinute <= -1 {
		return false
	}

	return true
}

func (opts *BotOptions) SetBaseURL(s string) *BotOptions {
	opts.BaseURL = strings.TrimSpace(s)
	return opts
}

func (opts *BotOptions) SetLimiterPerSecond(n int) *BotOptions {
	opts.LimiterPerSecond = n
	return opts
}

func (opts *BotOptions) SetLimiterPerMinute(n int) *BotOptions {
	opts.LimiterPerMinute = n
	return opts
}

func (opts *BotOptions) SetSecretKey(s string) *BotOptions {
	opts.SecretKey = s
	return opts
}

// --------------------------------------------------------------------------------

// 签名校验
//
// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#3c6592d6
func genSignature(timestamp int64, secretKey string) (string, error) {
	s := strconv.FormatInt(timestamp, 10) + "\n" + secretKey

	var data []byte
	h := hmac.New(sha256.New, []byte(s))
	if _, err := h.Write(data); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// --------------------------------------------------------------------------------

type Language string

const (
	LanguageChinese  Language = "zh_cn"
	LanguageEnglish  Language = "en_us"
	LanguageJapanese Language = "ja_jp"
)
