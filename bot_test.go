package feishu_bot_api

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/electricbubble/feishu-bot-api/v2/md"
)

func Test_bot_wait(t *testing.T) {
	const offsetDuration = 50 * time.Millisecond

	tests := []struct {
		count        int
		opts         *BotOptions
		wg           *sync.WaitGroup
		wantDuration time.Duration
	}{
		{
			count:        10,
			opts:         &BotOptions{},
			wantDuration: 1 * time.Second,
		},
		{
			count:        10,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 1 * time.Second,
		},
		{
			count:        11,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 2 * time.Second,
		},
		{
			count:        5,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 0,
		},
		{
			count:        15,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 2 * time.Second,
		},
		{
			count:        16,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 3 * time.Second,
		},
		{
			count:        1,
			opts:         &BotOptions{LimiterPerSecond: 1, LimiterPerMinute: 5},
			wantDuration: 0,
		},
		{
			count:        3,
			opts:         &BotOptions{LimiterPerSecond: 1, LimiterPerMinute: 5},
			wantDuration: 2 * time.Second,
		},
		{
			count:        5,
			opts:         &BotOptions{LimiterPerSecond: 1, LimiterPerMinute: 5},
			wantDuration: 4 * time.Second,
		},
		{
			count:        121,
			opts:         &BotOptions{LimiterPerSecond: 5, LimiterPerMinute: 100},
			wantDuration: 1*time.Minute + 4*time.Second,
		},
		{
			count:        6,
			opts:         &BotOptions{LimiterPerSecond: 1, LimiterPerMinute: 5},
			wantDuration: 1 * time.Minute,
		},
		{
			count:        15,
			opts:         &BotOptions{LimiterPerSecond: 1, LimiterPerMinute: 5},
			wantDuration: 2*time.Minute + 4*time.Second,
		},
		{
			count:        15,
			opts:         &BotOptions{LimiterPerSecond: 2, LimiterPerMinute: 10},
			wg:           new(sync.WaitGroup),
			wantDuration: 1*time.Minute + 2*time.Second,
		},
	}

	var wg sync.WaitGroup

	for _, tt := range tests {
		tt := tt

		wg.Add(1)
		go func() {
			defer wg.Done()

			b := NewBot("", tt.opts).(*bot)
			name := fmt.Sprintf("count-%d_s-%d_m-%d", tt.count, b.opts.LimiterPerSecond, b.opts.LimiterPerMinute)

			t.Run(name, func(t *testing.T) {

				{
					var (
						now  = time.Now()
						next time.Time
					)
					if tt.count <= b.opts.LimiterPerMinute {
						next = time.Unix(now.Unix(), 0).Add(time.Second)
					} else {
						next = time.Unix(now.Unix()-int64(now.Second()), 0).Add(time.Minute)
					}
					time.Sleep(next.Sub(now))
				}

				var (
					start    = time.Now()
					ms, mm   sync.Map
					bufDebug bytes.Buffer
				)

				for i := 1; i <= tt.count; i++ {
					if tt.wg != nil {
						tt.wg.Add(1)
						go func() {
							defer tt.wg.Done()

							if err := b.wait(); err != nil {
								t.Errorf("Received unexpected error:\n%+v", err)
							}
							now := time.Now()
							{
								counter := new(atomic.Int32)
								_v, _ := ms.LoadOrStore(now.Unix(), counter)
								counter = _v.(*atomic.Int32)
								counter.Add(1)
							}
							{
								counter := new(atomic.Int32)
								_v, _ := mm.LoadOrStore(now.Unix()-int64(now.Second()), counter)
								counter = _v.(*atomic.Int32)
								counter.Add(1)
							}
							bufDebug.WriteString(fmt.Sprintf("%s\t%d", now.Format("2006-01-02 15:04:05.000"), now.Unix()))
							// t.Log(now.Format("2006-01-02 15:04:05.000"), now.Unix())
						}()
					} else {
						if err := b.wait(); err != nil {
							t.Errorf("Received unexpected error:\n%+v", err)
						}
						now := time.Now()
						{
							counter := new(atomic.Int32)
							_v, _ := ms.LoadOrStore(now.Unix(), counter)
							counter = _v.(*atomic.Int32)
							counter.Add(1)
						}
						{
							counter := new(atomic.Int32)
							_v, _ := mm.LoadOrStore(now.Unix()-int64(now.Second()), counter)
							counter = _v.(*atomic.Int32)
							counter.Add(1)
						}
						bufDebug.WriteString(fmt.Sprintf("%s\t%d", now.Format("2006-01-02 15:04:05.000"), now.Unix()))
						// t.Log(now.Format("2006-01-02 15:04:05.000"), now.Unix())
					}
				}

				if tt.wg != nil {
					tt.wg.Wait()
				}

				d := time.Since(start)
				if d < tt.wantDuration-offsetDuration {
					t.Errorf("Actual min duration: %s, want: %s\nDEBUG:\n%s", d, tt.wantDuration-offsetDuration, bufDebug.String())
				}
				if d > tt.wantDuration+offsetDuration {
					t.Errorf("Actual max duration: %s, want: %s\nDEBUG:\n%s", d, tt.wantDuration+offsetDuration, bufDebug.String())
				}

				{
					keys := make([]int64, 0, int(tt.wantDuration.Seconds()))
					ms.Range(func(key, value any) bool {
						keys = append(keys, key.(int64))
						return true
					})
					slices.SortFunc(keys, func(a, b int64) int {
						return cmp.Compare(a, b)
					})

					for _, key := range keys {
						_v, _ := ms.Load(key)
						n := _v.(*atomic.Int32).Load()
						if int(n) > b.opts.LimiterPerSecond {
							t.Errorf("Actual per second: %d, want: %d\nDEBUG:\n%s", n, b.opts.LimiterPerSecond, bufDebug.String())
						}
					}
				}
				{
					keys := make([]int64, 0, int(tt.wantDuration.Minutes()))
					mm.Range(func(key, value any) bool {
						keys = append(keys, key.(int64))
						return true
					})
					slices.SortFunc(keys, func(a, b int64) int {
						return cmp.Compare(a, b)
					})

					for _, key := range keys {
						_v, _ := mm.Load(key)
						n := _v.(*atomic.Int32).Load()
						if int(n) > b.opts.LimiterPerMinute {
							t.Errorf("Actual per minute: %d, want: %d\nDEBUG:\n%s", n, b.opts.LimiterPerMinute, bufDebug.String())
						}
					}
				}

			})

		}()
	}

	wg.Wait()
}

func Test_bot_SendText(t *testing.T) {
	t.Run("full_webhook_has_secret_key", func(t *testing.T) {
		var (
			webhook   = os.Getenv("webhook")
			secretKey = os.Getenv("secret_key")
			b         = NewBot(webhook, &BotOptions{SecretKey: secretKey})
		)
		requireNoError(t, b.SendText("hi"))
	})

	t.Run("only_webhook_access_token_has_secret_key", func(t *testing.T) {
		var (
			webhook   = os.Getenv("webhook")
			secretKey = os.Getenv("secret_key")
			b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
		)
		requireNoError(t, b.SendText("hi"+TextAtEveryone()))
	})

	t.Run("only_webhook_access_token_no_secret_key_at_nobody", func(t *testing.T) {
		var (
			webhook = os.Getenv("webhook")
			b       = NewBot(webhook, nil)
		)
		requireNoError(t, b.SendText("hi"+TextAtPerson("nonexistent", "nobody")))
	})

	t.Run("only_webhook_access_token_no_secret_key_at_user_id", func(t *testing.T) {
		var (
			webhook = os.Getenv("webhook")
			b       = NewBot(webhook, nil)
			userID  = os.Getenv("user_id")
		)
		requireNoError(t, b.SendText("hi"+TextAtPerson(userID, "")))
	})

	t.Run("only_webhook_access_token_no_secret_key_at_open_id", func(t *testing.T) {
		var (
			webhook = os.Getenv("webhook")
			b       = NewBot(webhook, nil)
			openID  = os.Getenv("open_id")
		)
		requireNoError(t, b.SendText("hi"+TextAtPerson(openID, "")))
	})
}

func Test_bot_SendGroupBusinessCar(t *testing.T) {
	var (
		webhook = os.Getenv("webhook")
		b       = NewBot(webhook, nil)
		chatID  = os.Getenv("chat_id")
	)
	requireNoError(t, b.SendGroupBusinessCard(chatID))
}

func Test_bot_SendImage(t *testing.T) {
	var (
		webhook = os.Getenv("webhook")
		b       = NewBot(webhook, nil)
		imgKey  = os.Getenv("image_key")
	)
	requireNoError(t, b.SendImage(imgKey))
}

func Test_bot_SendRichText(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
		userID    = os.Getenv("user_id")
	)
	err := b.SendRichText(
		NewRichText(LanguageChinese, "🇨🇳 标题").
			Text("🇨🇳 文本", false).
			Hyperlink("超链接", "https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#f62e72d5").
			Image("img_7ea74629-9191-4176-998c-2e603c9c5e8g").
			Text("\n", false).
			AtEveryone().
			Text("+", false).
			At(userID, "").
			Image("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		NewRichText(LanguageEnglish, "🇬🇧🇺🇸 title").
			Text("🇨🇳 text", false).
			Hyperlink("hyper link", "https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#f62e72d5").
			Image("img_7ea74629-9191-4176-998c-2e603c9c5e8g").
			Text("\n", false).
			AtEveryone().
			Text("+", false).
			At(userID, "").
			Image("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		NewRichText(LanguageJapanese, "🇯🇵 タイトル").
			Text("🇯🇵 テキスト", false).
			Hyperlink("ハイパーリンク", "https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#f62e72d5").
			Image("img_7ea74629-9191-4176-998c-2e603c9c5e8g").
			Text("\n", false).
			AtEveryone().
			Text("+", false).
			At(userID, "").
			Image("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		NewRichText("zh_hk", "🇨🇳🇭🇰 标题").
			Text("🇨🇳🇭🇰 文本", false).
			Hyperlink("超链接", "https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#f62e72d5").
			Image("img_7ea74629-9191-4176-998c-2e603c9c5e8g").
			Text("\n", false).
			AtEveryone().
			Text("+", false).
			At(userID, "").
			Image("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
	)
	requireNoError(t, err)
}

func Test_bot_SendCard(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)
	err := b.SendCard(
		NewCardGlobalConfig().
			HeaderIcon("img_ecffc3b9-8f14-400f-a014-05eca1a4310g").
			HeaderTemplate(CardHeaderTemplateGreen).
			CardLink(
				"https://www.feishu.cn",
				"https://www.windows.com",
				"https://developer.apple.com",
				"https://developer.android.com",
			),
		NewCard(LanguageChinese, "🇨🇳 标题").
			HeaderSubtitle("🇨🇳 副标题").
			HeaderTextTags([]CardHeaderTextTag{
				{Content: "标题标签", Color: CardHeaderTextTagColorCarmine},
			}),
		NewCard(LanguageEnglish, "🇬🇧🇺🇸 title").
			HeaderSubtitle("🇬🇧🇺🇸 subtitle").
			HeaderTextTags([]CardHeaderTextTag{
				{Content: "tagDemo", Color: CardHeaderTextTagColorCarmine},
			}),
	)
	requireNoError(t, err)
}

func Test_bot_SendCard_CardElementDiv(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	// 文本
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text#e22d4592
	t.Run("text", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementDiv().LarkMarkdown("**text**/~~text~~"),
					NewCardElementDiv().PlainText("测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本", 2),
				}),
		)
		requireNoError(t, err)
	})

	// 双列文本
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field#e22d4592
	t.Run("fields", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "你有一个休假申请待审批").
				Elements([]CardElement{
					NewCardElementDiv().Fields([]CardElementDivFieldText{
						{
							IsShort: true,
							Mode:    CardElementDivTextModeLarkMarkdown,
							Content: "**申请人**\n王晓磊",
						},
						{
							IsShort: true,
							Mode:    CardElementDivTextModeLarkMarkdown,
							Content: "**休假类型：**\n年假",
						},
						{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: ""},
						{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "**时间：**\n2020-4-8 至 2020-4-10（共3天）"},
						{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: ""},
						{
							IsShort: true,
							Mode:    CardElementDivTextModeLarkMarkdown,
							Content: "**备注**\n因家中有急事，需往返老家，故请假",
						},
					}),
				}),
		)
		requireNoError(t, err)
	})

	// 附加图片元素
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/image
	t.Run("extra_image", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementDiv().
						PlainText("image element 1", 0).
						ExtraImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g", true, "hover图片后的tips文案"),
					NewCardElementHorizontalRule(),
					NewCardElementDiv().
						PlainText("image element 2", 0).
						ExtraImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g", true, "hover图片后的tips文案hover图片后的tips文案hover图片后的tips文案hover图片后的tips文案hover图片后的tips文案hover图片后的tips文案hover图片后的tips文案"),
					NewCardElementHorizontalRule(),
				}),
		)
		requireNoError(t, err)
	})
}

func Test_bot_SendCard_CardElementMarkdown(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	// Markdown 组件
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#7e4a981e
	t.Run("markdown", func(t *testing.T) {
		err := b.SendCard(
			NewCardGlobalConfig().HeaderTemplate(CardHeaderTemplateBlue),
			NewCard(LanguageChinese, "这是卡片标题栏").
				Elements([]CardElement{
					NewCardElementMarkdown("普通文本\n标准emoji😁😢🌞💼🏆❌✅\n*斜体*\n**粗体**\n~~删除线~~\n[文字链接](www.example.com)\n[差异化跳转]($urlVal)\n<at id=all></at>").
						Href(
							"https://www.feishu.cn",
							"https://www.windows.com",
							"https://developer.apple.com",
							"https://developer.android.com",
						),
					NewCardElementHorizontalRule(),
					NewCardElementMarkdown("上面是一行分割线\n![hover_text](img_v2_16d4ea4f-6cd5-48fa-97fd-25c8d4e79b0g)\n上面是一个图片标签"),
				}),
		)
		requireNoError(t, err)
	})

	// text 的 lark_md 模式
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#aaf4b702
	t.Run("lark_md", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementDiv().
						PlainText("text-lark_md", 1).
						Fields([]CardElementDivFieldText{
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "<a>https://open.feishu.cn</a>"},
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "ready\nnew line"},
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "*Italic*"},
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "**Bold**"},
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "~~delete line~~"},
							{IsShort: false, Mode: CardElementDivTextModeLarkMarkdown, Content: "<at id=all></at>"},
						}),
					NewCardElementHorizontalRule(),
					NewCardElementDiv().Fields([]CardElementDivFieldText{
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.AtPerson("abcd", "nobody")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.AtEveryone()},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.Hyperlink("https://open.feishu.cn")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.TextLink("开放平台", "https://open.feishu.cn")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.HorizontalRule()},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.FeiShuEmoji("DONE")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.GreenText("绿色文本")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.RedText("红色文本")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.GreyText("灰色文本")},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: "默认的白底黑字样式"},
						{Mode: CardElementDivTextModeLarkMarkdown, Content: md.TextTag(md.TextTagColorRed, "红色标签")},
					}),
				}),
		)
		requireNoError(t, err)
	})
}

func Test_bot_SendCard_CardElementImage(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#e22d4592
	err := b.SendCard(nil,
		NewCard(LanguageChinese, "").
			Elements([]CardElement{
				NewCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g", "Hover图片后的tips提示，不需要可以传空").
					// TitleWithPlainText("Block-img").
					TitleWithLarkMarkdown("*Block*-img").
					Mode(CardElementImageModeFitHorizontal).
					CompactWidth(false),
			}),
	)
	requireNoError(t, err)
}

func Test_bot_SendCard_CardElementNote(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/notes-module#3827dadd
	err := b.SendCard(nil,
		NewCard(LanguageChinese, "").
			Elements([]CardElement{
				NewCardElementNote().
					AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", true, "这是备注图片1").
					AddElementWithPlainText("备注信息1").
					AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", false, "这是备注图片").
					AddElementWithPlainText("备注信息2").
					AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", false, "alt_3").
					AddElementWithLarkMarkdown("*备注*~~信息~~" + md.RedText("3")),
			}),
	)
	requireNoError(t, err)
}

func Test_bot_SendCard_CardElementAction(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button#1a79927c

	t.Run("button_url", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementAction().Actions([]CardElementActionComponent{
						NewCardElementActionButton(CardElementDivTextModePlainText, "主按钮").
							Type(CardElementActionButtonTypePrimary).
							URL("https://open.feishu.cn/document"),
					}),
				}),
		)
		requireNoError(t, err)
	})

	t.Run("button_multi_url", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementAction().Actions([]CardElementActionComponent{
						NewCardElementActionButton(CardElementDivTextModePlainText, "主按钮").
							Type(CardElementActionButtonTypePrimary).
							MultiURL(
								"https://www.baidu.com",
								"https://www.windows.com",
								"lark://msgcard/unsupported_action",
								"https://developer.android.com",
							),
					}),
				}),
		)
		requireNoError(t, err)
	})

	// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/overflow#1a79927c
	t.Run("overflow", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementAction().Actions([]CardElementActionComponent{
						NewCardElementActionButton(CardElementDivTextModePlainText, "主按钮").
							Type(CardElementActionButtonTypePrimary).
							URL("https://open.feishu.cn/document"),
						NewCardElementActionOverflow().
							AddOptionWithMultiURL(
								"Option-1",
								"https://www.baidu.com",
								"https://www.windows.com",
								"https://developer.apple.com",
								"https://developer.android.com",
							).
							AddOptionWithURL("baidu", "https://www.baidu.com").
							AddOptionWithURL("开发文档", "https://open.feishu.cn/document/home/index").Confirm("Confirmation", "Content"),
					}),
				}),
		)
		requireNoError(t, err)
	})

	// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button?lang=zh-CN#e22d4592
	t.Run("example", func(t *testing.T) {
		defaultURL := "https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button?lang=zh-CN#e22d4592"
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").
				Elements([]CardElement{
					NewCardElementDiv().
						LarkMarkdown("Button element").
						// ExtraAction(
						// 	NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "Secondary confirmation").
						// 		Type(CardElementActionButtonTypeDefault).
						// 		URL(defaultURL).
						// 		Confirm("Confirmation", "Content"),
						// ).
						ExtraAction(
							NewCardElementActionOverflow().
								AddOptionWithMultiURL(
									"Option-1",
									"https://www.baidu.com",
									"https://www.windows.com",
									"https://developer.apple.com",
									"https://developer.android.com",
								).
								AddOptionWithURL("baidu", "https://www.baidu.com").
								AddOptionWithURL("开发文档", "https://open.feishu.cn/document/home/index").Confirm("Confirmation", "Content"),
						),
					NewCardElementHorizontalRule(),
					NewCardElementAction().Actions([]CardElementActionComponent{
						NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "default style").Type(CardElementActionButtonTypeDefault).URL(defaultURL),
						NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "primary style").Type(CardElementActionButtonTypePrimary).URL(defaultURL),
						NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "danger style").Type(CardElementActionButtonTypeDanger).URL(defaultURL),
						NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "target url").Type(CardElementActionButtonTypeDefault).URL("https://www.baidu.com"),
						NewCardElementActionButton(CardElementDivTextModeLarkMarkdown, "multi url").
							Type(CardElementActionButtonTypePrimary).
							MultiURL(
								"https://www.baidu.com",
								"https://www.windows.com",
								"https://developer.apple.com",
								"https://developer.android.com",
							),
					}),
					NewCardElementNote().AddElementWithPlainText("hello World"),
				}),
		)
		requireNoError(t, err)
	})
}

func Test_bot_SendCard_CardElementColumnSet(t *testing.T) {
	var (
		webhook   = os.Getenv("webhook")
		secretKey = os.Getenv("secret_key")
		b         = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
	)

	t.Run("columns_elements", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "多列布局").Elements([]CardElement{
				NewCardElementMarkdown("*斜体*\n**粗体**\n~~删除线~~" + md.LineBreak()),
				NewCardElementColumnSet().
					BackgroundStyle(CardElementColumnSetBackgroundStyleGrey).
					FlexMode(CardElementColumnSetFlexModeFlow).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{
								NewCardElementMarkdown("**Markdown** *组件*"),
								NewCardElementDiv().LarkMarkdown("**再加一个**/~~text~~"),
							}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementDiv().LarkMarkdown("**text**/~~text~~")}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementDiv().PlainText("测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本", 2)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementDiv().PlainText("测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本测试文本", 0)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{
								NewCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g", "Hover图片后的tips提示，不需要可以传空").
									// TitleWithPlainText("Block-img").
									TitleWithLarkMarkdown("*Block*-img").
									Mode(CardElementImageModeFitHorizontal).
									CompactWidth(false),
							}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementHorizontalRule()}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{
								NewCardElementNote().
									AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", true, "这是备注图片1").
									AddElementWithPlainText("备注信息1").
									AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", false, "这是备注图片").
									AddElementWithPlainText("备注信息2").
									AddElementWithImage("img_v2_041b28e3-5680-48c2-9af2-497ace79333g", false, "alt_3").
									AddElementWithLarkMarkdown("*备注*~~信息~~" + md.RedText("3")),
							}),
					}),
				NewCardElementColumnSet().
					BackgroundStyle(CardElementColumnSetBackgroundStyleGrey).
					FlexMode(CardElementColumnSetFlexModeFlow).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().Width(CardElementColumnSetColumnWidthWeighted).Weight(1).
							Elements([]CardElement{
								NewCardElementDiv().
									LarkMarkdown("ISV产品接入及企业自主开发，更好地对接现有系统，满足不同组织的需求。").
									ExtraAction(
										NewCardElementActionOverflow().
											AddOptionWithURL("打开飞书应用目录", "https://app.feishu.cn").
											AddOptionWithURL("打开飞书开发文档", "https://open.feishu.cn").
											AddOptionWithURL("打开飞书官网", "https://www.feishu.cn"),
									),
							}),
					}),
			}),
		)
		requireNoError(t, err)
	})

	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set?lang=zh-CN#5a267d0d
	t.Run("example_1", func(t *testing.T) {
		err := b.SendCard(nil,
			NewCard(LanguageChinese, "").Elements([]CardElement{
				NewCardElementMarkdown("**个人审批效率总览**" + md.LineBreak()),
				NewCardElementColumnSet().
					FlexMode(CardElementColumnSetFlexModeBisect).
					BackgroundStyle(CardElementColumnSetBackgroundStyleGrey).
					HorizontalSpacing(CardElementColumnSetHorizontalSpacingDefault).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementMarkdown("已审批单量\n**29单**\n" + md.GreenText("领先团队59%")).TextAlign(CardElementMarkdownTextAlignCenter)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementMarkdown("平均审批耗时\n**0.9小时**\n" + md.GreenText("领先团队100%")).TextAlign(CardElementMarkdownTextAlignCenter)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							Elements([]CardElement{NewCardElementMarkdown("待批率\n**25%**\n" + md.RedText("落后团队10%")).TextAlign(CardElementMarkdownTextAlignCenter)}),
					}),
				NewCardElementMarkdown("**团队审批效率参考：**"),
				NewCardElementColumnSet().
					// FlexMode(CardElementColumnSetFlexModeNone).
					BackgroundStyle(CardElementColumnSetBackgroundStyleGrey).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
							Elements([]CardElement{NewCardElementMarkdown("**审批人**").TextAlign(CardElementMarkdownTextAlignCenter)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
							Elements([]CardElement{NewCardElementMarkdown("**审批时长**").TextAlign(CardElementMarkdownTextAlignCenter)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
							Elements([]CardElement{NewCardElementMarkdown("**对比上周变化**").TextAlign(CardElementMarkdownTextAlignCenter)}),
					}),
				NewCardElementColumnSet().Columns([]*CardElementColumnSetColumn{
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown("王大明").TextAlign(CardElementMarkdownTextAlignCenter)}),
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown("小于1小时").TextAlign(CardElementMarkdownTextAlignCenter)}),
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown(md.GreenText("⬇️12%")).TextAlign(CardElementMarkdownTextAlignCenter)}),
				}),
				NewCardElementColumnSet().Columns([]*CardElementColumnSetColumn{
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown("张军").TextAlign(CardElementMarkdownTextAlignCenter)}),
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown("2小时").TextAlign(CardElementMarkdownTextAlignCenter)}),
					NewCardElementColumnSetColumn().
						Width(CardElementColumnSetColumnWidthWeighted).
						Weight(1).
						VerticalAlign(CardElementColumnSetColumnVerticalAlignTop).
						Elements([]CardElement{NewCardElementMarkdown(md.RedText("⬆️5%")).TextAlign(CardElementMarkdownTextAlignCenter)}),
				}),
			}),
		)
		requireNoError(t, err)
	})

	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set?lang=zh-CN#8efc1337
	t.Run("example_2", func(t *testing.T) {
		err := b.SendCard(
			NewCardGlobalConfig().HeaderTemplate(CardHeaderTemplateGreen),
			NewCard(LanguageChinese, "🏨 酒店申请已通过，请选择房型").Elements([]CardElement{
				NewCardElementMarkdown("**入住酒店：**[杭州xxxx酒店](https://open.feishu.cn/)\n<font color='grey'>📍 浙江省杭州市西湖区</font>"),
				NewCardElementHorizontalRule(),
				NewCardElementColumnSet().
					FlexMode(CardElementColumnSetFlexModeNone).
					BackgroundStyle(CardElementColumnSetBackgroundStyleDefault).
					HorizontalSpacing(CardElementColumnSetHorizontalSpacingDefault).
					ActionMultiURL(
						"https://open.feishu.cn",
						"https://www.windows.com",
						"https://developer.apple.com",
						"https://developer.android.com",
					).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							VerticalAlign(CardElementColumnSetColumnVerticalAlignCenter).
							Elements([]CardElement{NewCardElementImage("img_v2_120b03c8-27e3-456f-89c0-90ede1aa59ag", "").Mode(CardElementImageModeFitHorizontal)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(3).
							Elements([]CardElement{NewCardElementMarkdown("**高级双床房**\n<font color='grey'>双早|40-47㎡|有窗户|双床</font>\n<font color='red'>¥699</font> 起").TextAlign(CardElementMarkdownTextAlignLeft)}),
					}),
				NewCardElementHorizontalRule(),
				NewCardElementColumnSet().
					FlexMode(CardElementColumnSetFlexModeNone).
					BackgroundStyle(CardElementColumnSetBackgroundStyleDefault).
					HorizontalSpacing(CardElementColumnSetHorizontalSpacingDefault).
					ActionMultiURL(
						"https://open.feishu.cn",
						"https://www.windows.com",
						"https://developer.apple.com",
						"https://developer.android.com",
					).
					Columns([]*CardElementColumnSetColumn{
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(1).
							VerticalAlign(CardElementColumnSetColumnVerticalAlignCenter).
							Elements([]CardElement{NewCardElementImage("img_v2_120b03c8-27e3-456f-89c0-90ede1aa59ag", "").Mode(CardElementImageModeFitHorizontal)}),
						NewCardElementColumnSetColumn().
							Width(CardElementColumnSetColumnWidthWeighted).
							Weight(3).
							Elements([]CardElement{NewCardElementMarkdown("**精品大床房**\n<font color='grey'>双早|40-47㎡|有窗户|大床</font>\n<font color='red'>¥666</font> 起").TextAlign(CardElementMarkdownTextAlignLeft)}),
					}),
			}),
		)
		requireNoError(t, err)
	})
}

func Test_bot_SendCardViaTemplate(t *testing.T) {
	var (
		webhook    = os.Getenv("webhook")
		secretKey  = os.Getenv("secret_key")
		b          = NewBot(webhook, NewBotOptions().SetSecretKey(secretKey))
		templateID = os.Getenv("template_id")
	)

	type (
		tplVarsGroupTable struct {
			Person   string `json:"person"`
			Time     string `json:"time"`
			WeekRate string `json:"week_rate"`
		}

		tplVars struct {
			TotalCount   string              `json:"total_count"`
			TotalPercent string              `json:"total_percent"`
			Hours        string              `json:"hours"`
			HoursPercent string              `json:"hours_percent"`
			Pending      string              `json:"pending"`
			PendingRate  string              `json:"pending_rate"`
			GroupTable   []tplVarsGroupTable `json:"group_table"`
		}
	)

	variables := tplVars{
		TotalCount:   "29",
		TotalPercent: "<font color='green'>领先团队59%</font>",
		Hours:        "0.9",
		HoursPercent: "<font color='green'>领先团队100%</font>",
		Pending:      "25%",
		PendingRate:  "<font color='red'>落后团队10%</font>",
		GroupTable: []tplVarsGroupTable{
			{Person: "王大明", Time: "小于1小时", WeekRate: "<font color='green'>↓12%</font>"},
			{Person: "张军", Time: "2小时", WeekRate: "<font color='red'>↑5%</font>"},
			{Person: "李小方", Time: "3小时", WeekRate: "<font color='green'>↓25%</font>"},
		},
	}

	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set?lang=zh-CN#5a267d0d
	err := b.SendCardViaTemplate(templateID, variables)
	requireNoError(t, err)
}
