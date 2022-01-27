package fsBotAPI

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWithCardTitle(t *testing.T) {
	// fnCT := WithCardHeader(BgColorDefault).
	// 	WithTitle(LangChinese, "中文 标题").
	// 	WithTitle(LangEnglish, "英文 title").
	// 	WithTitle(LangJapanese, "日语 見出し")
	// // msgCT := new(msgCardTitle)
	// // fnCT(msgCT)
	//
	// msgCT := fnCT.titleI18n()
	//
	// for lang, s := range msgCT {
	// 	t.Log(lang, s)
	// }
	//
	// t.Log()
	//
	// WithCard(
	// 	WithCardHeader(BgColorCarmine).
	// 		WithTitle(LangChinese, "中文 标题").
	// 		WithTitle(LangEnglish, "英文 title").
	// 		WithTitle(LangJapanese, "日语 見出し"),
	// 	WithCardConfig().WithEnableForward(false).WithEnableUpdateMulti(false),
	// 	nil,
	// )
	//
	// t.Log()
	//
	// WithCard(nil, nil, nil)
	//
	// t.Log()

	mdZhCn := `普通文本
标准emoji 😁😢🌞💼🏆❌✅
*斜体*
**粗体**
~~删除线~~
[文字链接](https://www.feishu.cn)
[差异化跳转（桌面端和移动配置不同跳转链接）]($urlVal)
<at id=all></at>
<at id=ou_c99c5f35d542efc7ee492afe11af19ef></at>
 ---
上面是一行分割线
![光标hover图片上的tips文案可不填](img_7ea74629-9191-4176-998c-2e603c9c5e8g)
上面是一个图片标签
`

	// msgCard := WithCard(WithCardHeader(BgColorOrange).WithTitle(LangChinese, "标题"), nil,
	// 	WithCardElement(LangChinese, WithCardElementMarkdown(mdZhCn)),
	// 	WithCardElement(LangEnglish, WithCardElementPlainText("content")),
	// 	WithCardElement(LangJapanese, WithCardElementPlainText("japanese")),
	// )
	//
	// t.Log(strings.Repeat("-", 50))
	//
	// bs, err := json.MarshalIndent(msgCard(), "", "  ")
	// requireNil(t, err)
	//
	// fmt.Println(string(bs))
	//
	// t.Log(strings.Repeat("-", 50))

	card := GenMsgCard(BgColorOrange, nil,
		WithCard(LangChinese, "标题", WithCardElementMarkdown(mdZhCn)),
		WithCard(LangEnglish, "TITLE", WithCardElementPlainText("content")),
		WithCard(LangJapanese, "タイトル", WithCardElementPlainText("japanese")),
	)

	bs, err := json.MarshalIndent(card, "", "  ")
	requireNil(t, err)

	fmt.Println(string(bs))

	// cardCfg := WithCardConfig()
	//
	// bs, err := json.MarshalIndent(cardCfg, "", "  ")
	// requireNil(t, err)
	//
	// fmt.Println(string(bs))

	// bs, err := json.MarshalIndent(msgCT, "", "  ")
	// requireNil(t, err)
	//
	// fmt.Println(string(bs))
}
