package main

import (
	"bytes"
	fsBotAPI "github.com/electricbubble/feishu-bot-api"
	"github.com/electricbubble/feishu-bot-api/md"
	"log"
	"os"
	"strings"
)

func loadTestBot(wh *string, secretKey *string) {
	*wh = os.Getenv("FS_BOT_WEBHOOK")
	*secretKey = os.Getenv("FS_BOT_SECRET_KEY")
}

func main() {
	webhook := "https://open.feishu.cn/open-apis/bot/v2/hook/045b07bd-xxxx-xxxx-xxxx-2b49c4af1a2b"
	// 群机器人的 webhook 地址, 可以选择直接使用 👆 地址
	// 也可以选择使用该地址尾部的 36位🆔
	webhook = "045b07bd-xxxx-xxxx-xxxx-2b49c4af1a2b"

	// 开启签名校验时使用
	secretKey := "你的密钥"

	loadTestBot(&webhook, &secretKey)

	// 开启签名校验
	bot := fsBotAPI.NewBot(webhook, fsBotAPI.WithSecretKey(secretKey))
	// 如果未开启可忽略密钥 👇
	// bot := fsBotAPI.NewBot(webhook)

	{ // 发送普通文本消息
		buf := bytes.NewBufferString("新更新提醒\n")
		buf.WriteString("🤓所有人👉" + fsBotAPI.StrMentionAll() + "\n")
		buf.WriteString("🤔你是谁👉" + fsBotAPI.StrMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef") + "\n")

		err := bot.PushText(buf.String())
		if err != nil {
			log.Fatalln(err)
		}
	}

	{ // 发送富文本消息
		// 可同时设置 3种语言环境
		err := bot.PushPost(
			fsBotAPI.WithPost(fsBotAPI.LangChinese, "🇨🇳我是一个标题",
				fsBotAPI.WithPostElementText("🇨🇳第一行: "),
				fsBotAPI.WithPostElementLink("超链接", "https://www.feishu.cn"),
				fsBotAPI.WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
				fsBotAPI.WithPostElementText("\n"),
				fsBotAPI.WithPostElementMentionAll(),
				fsBotAPI.WithPostElementText("+"),
				fsBotAPI.WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
				fsBotAPI.WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
			),
			fsBotAPI.WithPost(fsBotAPI.LangEnglish, "🇺🇸🇬🇧 title",
				fsBotAPI.WithPostElementText("🇺🇸🇬🇧 first line: "),
				fsBotAPI.WithPostElementLink("link", "https://www.feishu.cn"),
				fsBotAPI.WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
				fsBotAPI.WithPostElementText("\n"),
				fsBotAPI.WithPostElementMentionAll(),
				fsBotAPI.WithPostElementText("+"),
				fsBotAPI.WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
				fsBotAPI.WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
			),
			fsBotAPI.WithPost(fsBotAPI.LangJapanese, "🇯🇵 タイトル",
				fsBotAPI.WithPostElementText("🇯🇵 1行目: "),
				fsBotAPI.WithPostElementLink("リンク", "https://www.feishu.cn"),
				fsBotAPI.WithPostElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g"),
				fsBotAPI.WithPostElementText("\n"),
				fsBotAPI.WithPostElementMentionAll(),
				fsBotAPI.WithPostElementText("+"),
				fsBotAPI.WithPostElementMentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
				fsBotAPI.WithPostElementImage("img_ecffc3b9-8f14-400f-a014-05eca1a4310g"),
			),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}

	{ // 发送卡片消息, 可同时设置 3种语言环境
		mdZhCn := `**title**
~~DEL~~
🙈 看不见的人 👉` + md.MentionByOpenID("ou_c99c5f35d542efc7ee492afe11af19ef") + "\n" +
			md.ColorGreen("这是一个绿色文本") + "\n" +
			md.ColorRed("这是一个红色文本") + "\n" +
			md.ColorGrey("这是一个灰色文本 ")

		mdEnUs := `~~empty~~`

		// 卡片消息可以设置是否允许转发，默认允许转发
		// fsBotAPI.WithCardConfig(fsBotAPI.WithCardConfigEnableForward(false))
		// 也可以将整个卡片作为链接点击, 同时支持多端跳转 👇
		// fsBotAPI.WithCardConfig(fsBotAPI.WithCardConfigEnableForward(false), fsBotAPI.WithCardConfigCardLink(
		// 	"https://www.feishu.cn",
		// 	"https://zlink.toutiao.com/kG12?apk=1",
		// 	"https://zlink.toutiao.com/h2Sw",
		// 	"https://www.feishu.cn/download",
		// ))

		// 第二参数为卡片的基础配置, 默认配置可直接 nil
		err := bot.PushCard(fsBotAPI.BgColorOrange, nil,
			fsBotAPI.WithCard(fsBotAPI.LangChinese, "标题",
				fsBotAPI.WithCardElementPlainText("文本内容"),
				fsBotAPI.WithCardElementHorizontalRule(),
				fsBotAPI.WithCardElementPlainText(strings.Repeat("文本内容2", 20), 2),
				fsBotAPI.WithCardElementHorizontalRule(),
				fsBotAPI.WithCardElementMarkdown(mdZhCn),
				fsBotAPI.WithCardElementHorizontalRule(),
				fsBotAPI.WithCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
					fsBotAPI.WithCardElementImageTitle("    *图片标题*", true),
					fsBotAPI.WithCardElementImageHover("被发现了"),
				),
				fsBotAPI.WithCardElementNote(
					fsBotAPI.WithCardElementPlainText("**普通文本**"),
					fsBotAPI.WithCardElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
						fsBotAPI.WithCardElementImageTitle("    *图片标题*", true),
						fsBotAPI.WithCardElementImageHover("被发现了"),
					),
					fsBotAPI.WithCardElementMarkdown("*test*"),
				),
				fsBotAPI.WithCardElementFields(
					fsBotAPI.WithCardElementField(fsBotAPI.WithCardElementPlainText("列1\nv1"), true),
					fsBotAPI.WithCardElementField(fsBotAPI.WithCardElementMarkdown("**列2**\nv2"), true),
					fsBotAPI.WithCardElementField(fsBotAPI.WithCardElementMarkdown("~无效的信息~"), false),
				),
				fsBotAPI.WithCardElementActions(
					fsBotAPI.WithCardElementAction(fsBotAPI.WithCardElementPlainText("入门必读"), "https://www.feishu.cn/hc/zh-CN/articles/360024881814", fsBotAPI.WithCardElementActionButton(fsBotAPI.ButtonDefault)),
					fsBotAPI.WithCardElementAction(fsBotAPI.WithCardElementPlainText("快速习惯飞书️"), "https://www.feishu.cn/hc/zh-CN/categories-detail?category-id=7018450035717259265", fsBotAPI.WithCardElementActionButton(fsBotAPI.ButtonPrimary)),
					fsBotAPI.WithCardElementAction(
						fsBotAPI.WithCardElementMarkdown("**多端跳转下载**"), "", fsBotAPI.WithCardElementActionButton(fsBotAPI.ButtonDanger), fsBotAPI.WithCardElementActionMultiURL(
							"https://www.feishu.cn",
							"https://zlink.toutiao.com/kG12?apk=1",
							"https://zlink.toutiao.com/h2Sw",
							"https://www.feishu.cn/download",
						),
					),
				),

				fsBotAPI.WithCardElementMarkdown("*TEST*", fsBotAPI.WithCardExtraElementImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g",
					fsBotAPI.WithCardElementImageTitle("    *图片标题*", true),
					fsBotAPI.WithCardElementImageHover("被发现了"),
				)),
			),

			fsBotAPI.WithCard(fsBotAPI.LangEnglish, "title",
				fsBotAPI.WithCardElementMarkdown(mdEnUs),
			),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}

	{ // 发送图片消息
		err := bot.PushImage("img_7ea74629-9191-4176-998c-2e603c9c5e8g")
		if err != nil {
			log.Fatalln(err)
		}
	}

	{ // 发送群名片
		err := bot.PushShareChat("oc_f5b1a7eb27ae2c7b6adc2a74faf339ff")
		if err != nil {
			log.Fatalln(err)
		}
	}
}
