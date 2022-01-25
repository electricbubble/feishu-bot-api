package fsBotAPI

import (
	"os"
	"strings"
	"testing"
)

func requireNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_bot_PushText(t *testing.T) {
	key := os.Getenv("BOT_KEY")
	t.Log(key)

	var (
		bot = NewBot(key)
		err error
	)

	content := "test content"
	// err := bot.PushText(content)
	// requireNil(t, err)

	content = `first line
second line`
	// err = bot.PushText(content)
	// requireNil(t, err)

	content = `first line
second line
` + StrMentionAll()
	// err = bot.PushText(content)
	// requireNil(t, err)

	content = `first line
second line
` + StrMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef") + `
这个人不在这个群聊里`
	// err = bot.PushText(content)
	// requireNil(t, err)

	content = `first line
second line
` + StrMentionByOpenID("ou_invalid", "不存在的人")
	err = bot.PushText(content)
	requireNil(t, err)
}

func Test_bot_PushPost(t *testing.T) {
	key := os.Getenv("BOT_KEY")
	t.Log(key)

	err := NewBot(key).PushPost(
		WithPost(LangChinese, "🇨🇳我是一个标题",
			WithPostElementText("🇨🇳第一行："),
			WithPostElementLink("超链接", "https://www.feishu.cn"),
			WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
			WithPostElementText("\n"),
			WithPostElementMentionAll(),
			WithPostElementText("+"),
			WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
			WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		),
		WithPost(LangEnglish, "🇺🇸🇬🇧 title",
			WithPostElementText("🇺🇸🇬🇧 first line："),
			WithPostElementLink("link", "https://www.feishu.cn"),
			WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
			WithPostElementText("\n"),
			WithPostElementMentionAll(),
			WithPostElementText("+"),
			WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
			WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		),
		WithPost(LangJapanese, "🇯🇵 見出し",
			WithPostElementText("🇯🇵 最初の行："),
			WithPostElementLink("リンク", "https://www.feishu.cn"),
			WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
			WithPostElementText("\n"),
			WithPostElementMentionAll(),
			WithPostElementText("+"),
			WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
			WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
		),
	)
	requireNil(t, err)
}

func Test_bot_PushImage(t *testing.T) {
	key := os.Getenv("BOT_KEY")
	t.Log(key)

	err := NewBot(key).PushImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g")
	requireNil(t, err)

}

func Test_bot_PushShareChat(t *testing.T) {
	key := os.Getenv("BOT_KEY")
	t.Log(key)

	err := NewBot(key).PushShareChat("oc_f5b1a7eb27ae2c7b6adc2a74faf339ff")
	requireNil(t, err)
}

func Test_bot_PushCard(t *testing.T) {
	key := os.Getenv("BOT_KEY")
	t.Log(key)

	bot := NewBot(key)

	mdZhCn := `**title**
~~DEL~~
`

	err := bot.PushCard(BgColorOrange, WithCardConfig(WithCardConfigCardLink(
		"https://www.feishu.cn",
		"https://zlink.toutiao.com/kG12?apk=1",
		"https://zlink.toutiao.com/h2Sw",
		"https://www.feishu.cn/download",
	)),
		WithCard(LangChinese, "标题",
			WithCardElementPlainText("文本内容"),
			WithCardElementHorizontalRule(),
			WithCardElementPlainText(strings.Repeat("文本内容2", 20), 2),
			WithCardElementHorizontalRule(),
			WithCardElementMarkdown(mdZhCn),
			WithCardElementHorizontalRule(),
			WithCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
				WithCardElementImageTitle("    *图片标题*", true),
				WithCardElementImageHover("被发现了"),
			),
			WithCardElementNote(
				WithCardElementPlainText("**普通文本**"),
				WithCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
					WithCardElementImageTitle("    *图片标题*", true),
					WithCardElementImageHover("被发现了"),
				),
				WithCardElementMarkdown("*test*"),
			),
			WithCardElementFields(
				WithCardElementField(WithCardElementPlainText("列1\nv1"), true),
				WithCardElementField(WithCardElementMarkdown("**列2**\nv2"), true),
				WithCardElementField(WithCardElementMarkdown("~无效的信息~"), false),
			),
			WithCardElementActions(
				WithCardElementAction(WithCardElementPlainText("入门必读"), "https://www.feishu.cn/hc/zh-CN/articles/360024881814", WithCardElementActionButton(ButtonDefault)),
				WithCardElementAction(WithCardElementPlainText("快速习惯飞书️"), "https://www.feishu.cn/hc/zh-CN/categories-detail?category-id=7018450035717259265", WithCardElementActionButton(ButtonPrimary)),
				WithCardElementAction(
					WithCardElementMarkdown("**多端跳转下载**"), "", WithCardElementActionButton(ButtonDanger), WithCardElementActionMultiURL(
						"https://www.feishu.cn",
						"https://zlink.toutiao.com/kG12?apk=1",
						"https://zlink.toutiao.com/h2Sw",
						"https://www.feishu.cn/download",
					),
				),
			),

			WithCardElementMarkdown("*TEST*", WithCardExtraElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
				WithCardElementImageTitle("    *图片标题*", true),
				WithCardElementImageHover("被发现了"),
			)),
		),

		WithCard(LangEnglish, "title",
			WithCardElementMarkdown("~~empty~~"),
		),
	)

	requireNil(t, err)
}
