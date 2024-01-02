package feishu_bot_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
)

var _ Message = (*cardMessage)(nil)

type cardMessage struct {
	globalConf *CardGlobalConfig

	builders []*CardBuilder
}

func (m cardMessage) Apply(body *MessageBody) error {
	card := &MessageBodyCard{
		Header:       nil,
		Elements:     nil,
		I18nElements: nil,
		Config:       nil,
		CardLink:     nil,
	}

	rawHeader, err := json.Marshal(m.buildHeader())
	if err != nil {
		return fmt.Errorf("card message: marshal header: %w", err)
	}
	card.Header = (*json.RawMessage)(&rawHeader)

	rawI18nElements, err := m.marshalI18nElements()
	if err != nil {
		return fmt.Errorf("card message: marshal i18n_elements: %w", err)
	}
	card.I18nElements = &rawI18nElements

	if m.globalConf != nil && m.globalConf.config != nil {
		raw, err := json.Marshal(m.globalConf.config)
		if err != nil {
			return fmt.Errorf("card message: marshal config: %w", err)
		}
		card.Config = (*json.RawMessage)(&raw)
	}

	if m.globalConf != nil && m.globalConf.link != nil {
		raw, err := json.Marshal(m.globalConf.link)
		if err != nil {
			return fmt.Errorf("card message: marshal card_link: %w", err)
		}
		card.CardLink = (*json.RawMessage)(&raw)
	}

	rawCard, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("card message: marshal: %w", err)
	}

	body.MsgType = "interactive"
	body.Card = (*json.RawMessage)(&rawCard)
	return nil
}

func (m cardMessage) buildHeader() cardHeader {
	nnBuilders := make([]*CardBuilder, 0, len(m.builders))
	for i := range m.builders {
		if m.builders[i] == nil {
			continue
		}
		nnBuilders = append(nnBuilders, m.builders[i])
	}

	ret := cardHeader{
		Title: cardHeaderTitle{
			Tag:  "plain_text",
			I18n: make(cardHeaderI18nTexts, len(nnBuilders)),
		},
		Subtitle:        nil,
		Icon:            nil,
		Template:        "",
		I18nTextTagList: nil,
	}

	for i := range nnBuilders {
		ret.Title.I18n[i] = cardHeaderI18nText{
			language: nnBuilders[i].language,
			text:     nnBuilders[i].headerTitleContent,
		}
	}

	for i := range nnBuilders {
		b := nnBuilders[i]
		if b.headerSubtitleContent != "" {
			ret.Subtitle = &cardHeaderTitle{
				Tag:  "plain_text",
				I18n: make(cardHeaderI18nTexts, len(nnBuilders)),
			}
			for i := range ret.Subtitle.I18n {
				ret.Subtitle.I18n[i] = cardHeaderI18nText{
					language: nnBuilders[i].language,
					text:     nnBuilders[i].headerSubtitleContent,
				}
			}
			break
		}
	}

	if m.globalConf != nil && m.globalConf.headerIcon != nil {
		ret.Icon = m.globalConf.headerIcon
	}

	if m.globalConf != nil && m.globalConf.headerTemplate != "" {
		ret.Template = m.globalConf.headerTemplate
	}

	i18nTextTagList := make(cardHeaderI18nTextTags, 0, len(nnBuilders))
	for i := range nnBuilders {
		b := nnBuilders[i]
		if b.headerI18nTextTags != nil {
			i18nTextTagList = append(i18nTextTagList, *b.headerI18nTextTags...)
		}
	}
	if len(i18nTextTagList) > 0 {
		ret.I18nTextTagList = &i18nTextTagList
	}

	return ret
}

func (m cardMessage) marshalI18nElements() (json.RawMessage, error) {
	nnBuilders := make([]*CardBuilder, 0, len(m.builders))
	for i := range m.builders {
		if m.builders[i] == nil {
			continue
		}
		nnBuilders = append(nnBuilders, m.builders[i])
	}

	if len(nnBuilders) == 0 {
		return []byte("null"), nil
	}

	var (
		languages = make([]Language, 0, len(nnBuilders))
		esValues  = make([][]any, 0, len(nnBuilders))
	)
	for i := range nnBuilders {
		b := nnBuilders[i]
		if slices.Contains(languages, b.language) {
			continue
		}
		languages = append(languages, b.language)
		esValues = append(esValues, b.elements)
	}

	var buf bytes.Buffer
	buf.Grow(len(languages) * 48)

	buf.WriteString("{")

	commas := len(languages)
	for i := range languages {
		commas--

		language := string(languages[i])
		es := esValues[i]

		bs, err := json.Marshal(es)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", language, err)
		}

		buf.WriteString(fmt.Sprintf(`"%s":%s`, _quoteEscaper.Replace(language), bs))

		if commas > 0 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

// --------------------------------------------------------------------------------

type (
	CardGlobalConfig struct {
		headerIcon     *cardHeaderIcon
		headerTemplate CardHeaderTemplate
		config         *cardConfig
		link           *cardLink
	}

	CardBuilder struct {
		language Language

		headerTitleContent    string
		headerSubtitleContent string
		headerI18nTextTags    *cardHeaderI18nTextTags

		elements []any
	}

	CardHeaderTextTag struct {
		Content string
		Color   CardHeaderTextTagColor
	}
)

func NewCardGlobalConfig() *CardGlobalConfig {
	return &CardGlobalConfig{}
}

func NewCard(language Language, title string) *CardBuilder {
	return &CardBuilder{
		language:           language,
		headerTitleContent: title,

		elements: make([]any, 0, 4),
	}
}

// HeaderIcon 标题的前缀图标
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#3827dadd
func (gConf *CardGlobalConfig) HeaderIcon(imgKey string) *CardGlobalConfig {
	gConf.headerIcon = &cardHeaderIcon{ImgKey: imgKey}
	return gConf
}

// HeaderTemplate 标题主题颜色
func (gConf *CardGlobalConfig) HeaderTemplate(template CardHeaderTemplate) *CardGlobalConfig {
	gConf.headerTemplate = template
	return gConf
}

// ConfigEnableForward 是否允许转发卡片
//   - true：允许
//   - false：不允许
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-structure/card-configuration#3827dadd
func (gConf *CardGlobalConfig) ConfigEnableForward(b bool) *CardGlobalConfig {
	if gConf.config == nil {
		gConf.config = new(cardConfig)
	}
	gConf.config.EnableForward = b
	return gConf
}

// ConfigUpdateMulti 是否为共享卡片
//   - true：是共享卡片，更新卡片的内容对所有收到这张卡片的人员可见
//   - false：非共享卡片，即独享卡片，仅操作用户可见卡片的更新内容
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-structure/card-configuration#3827dadd
func (gConf *CardGlobalConfig) ConfigUpdateMulti(b bool) *CardGlobalConfig {
	if gConf.config == nil {
		gConf.config = new(cardConfig)
	}
	gConf.config.UpdateMulti = b
	return gConf
}

// CardLink 消息卡片跳转链接
//
// 用于指定卡片整体的点击跳转链接，可以配置默认链接，也可以分别为 PC 端、Android 端、iOS 端配置不同的跳转链接
//   - 如果未配置 pc、ios、android，则默认跳转至 defaultURL
//   - 如果配置了 pc、ios、android，则优先生效各端指定的跳转链接
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7bfe6950
func (gConf *CardGlobalConfig) CardLink(defaultURL, pc, ios, android string) *CardGlobalConfig {
	gConf.link = &cardLink{
		URL:     defaultURL,
		PC:      pc,
		IOS:     ios,
		Android: android,
	}
	return gConf
}

// ----------------------------------------

// HeaderSubtitle 卡片的副标题信息
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#3827dadd
func (cb *CardBuilder) HeaderSubtitle(subtitle string) *CardBuilder {
	cb.headerSubtitleContent = subtitle
	return cb
}

// HeaderTextTags 标题的标签属性
//
// 最多可配置 3 个标签内容，如果配置的标签数量超过 3 个，则取前 3 个标签进行展示。
// 标签展示顺序与数组顺序一致
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#3827dadd
func (cb *CardBuilder) HeaderTextTags(tags []CardHeaderTextTag) *CardBuilder {
	tts := make(cardHeaderI18nTextTags, len(tags))
	for i := range tags {
		tts[i] = cardHeaderI18nTextTag{
			language: cb.language,
			Tag:      "text_tag",
			Text: cardHeaderComponentPlainText{
				Tag:     "plain_text",
				Content: tags[i].Content,
			},
			Color: tags[i].Color,
		}
	}
	cb.headerI18nTextTags = &tts
	return cb
}

// Elements 卡片的正文内容
func (cb *CardBuilder) Elements(elements []CardElement) *CardBuilder {
	for i := range elements {
		if elements[i] == nil {
			continue
		}
		cb.elements = append(cb.elements, elements[i].Entity())
	}
	return cb
}

// --------------------------------------------------------------------------------

type (
	// cardConfig
	//
	// 参数说明: https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-structure/card-configuration#3827dadd
	cardConfig struct {
		// 是否允许转发卡片
		//  true：允许
		//  false：不允许
		EnableForward bool `json:"enable_forward"`

		// 是否为共享卡片
		//  true：是共享卡片，更新卡片的内容对所有收到这张卡片的人员可见
		//  false：非共享卡片，即独享卡片，仅操作用户可见卡片的更新内容
		UpdateMulti bool `json:"update_multi"`
	}

	// cardLink
	//
	// 参数说明: https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7bfe6950
	cardLink differentJumpLinks

	// cardHeader
	//
	// 参数说明: https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#3827dadd
	cardHeader struct {
		// 卡片的主标题信息
		Title cardHeaderTitle `json:"title"`

		// 卡片的副标题信息
		Subtitle *cardHeaderTitle `json:"subtitle,omitempty"`

		// 标题的前缀图标。一个卡片仅可配置一个标题图标
		Icon *cardHeaderIcon `json:"icon,omitempty"`

		// 标题主题颜色
		Template CardHeaderTemplate `json:"template,omitempty"`

		// 标题标签的国际化属性。
		//
		// 最多可配置 3 个标签内容，如果配置的标签数量超过 3 个，则取前 3 个标签进行展示。
		// 标签展示顺序与数组顺序一致
		I18nTextTagList *cardHeaderI18nTextTags `json:"i18n_text_tag_list,omitempty"`
	}
	cardHeaderTitle struct {
		// 文本标识。固定取值：plain_text
		Tag string `json:"tag"`

		// 国际化文本内容
		I18n cardHeaderI18nTexts `json:"i18n"`
	}
	cardHeaderIcon struct {
		ImgKey string `json:"img_key,omitempty"`
	}
)

// ----------------------------------------

type differentJumpLinks struct {
	// 默认的链接地址
	URL string `json:"url"`

	// PC 端的链接地址
	PC string `json:"pc_url,omitempty"`

	// iOS 端的链接地址
	IOS string `json:"ios_url,omitempty"`

	// Android 端的链接地址
	Android string `json:"android_url,omitempty"`
}

// ----------------------------------------

// CardHeaderTemplate 标题样式
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#a19af820
//
// 样式建议: https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#88aad907
type CardHeaderTemplate string

const (
	CardHeaderTemplateBlue      CardHeaderTemplate = "blue"
	CardHeaderTemplateWathet    CardHeaderTemplate = "wathet"
	CardHeaderTemplateTurquoise CardHeaderTemplate = "turquoise"
	CardHeaderTemplateGreen     CardHeaderTemplate = "green"
	CardHeaderTemplateYellow    CardHeaderTemplate = "yellow"
	CardHeaderTemplateOrange    CardHeaderTemplate = "orange"
	CardHeaderTemplateRed       CardHeaderTemplate = "red"
	CardHeaderTemplateCarmine   CardHeaderTemplate = "carmine"
	CardHeaderTemplateViolet    CardHeaderTemplate = "violet"
	CardHeaderTemplatePurple    CardHeaderTemplate = "purple"
	CardHeaderTemplateIndigo    CardHeaderTemplate = "indigo"
	CardHeaderTemplateGrey      CardHeaderTemplate = "grey"
	CardHeaderTemplateDefault   CardHeaderTemplate = "default"
)

// ----------------------------------------

var _ json.Marshaler = (*cardHeaderI18nTextTags)(nil)

type cardHeaderI18nTextTags []cardHeaderI18nTextTag

func (tts cardHeaderI18nTextTags) MarshalJSON() ([]byte, error) {
	if tts == nil {
		return []byte("null"), nil
	}

	var (
		languages = make([]Language, 0, len(tts))
		ttsValues = make([]cardHeaderI18nTextTags, 0, len(tts))
	)
	for _, tt := range tts {
		idx := slices.Index(languages, tt.language)
		if idx == -1 {
			languages = append(languages, tt.language)
			ttsValues = append(ttsValues, cardHeaderI18nTextTags{
				cardHeaderI18nTextTag{Tag: tt.Tag, Text: tt.Text, Color: tt.Color},
			})
			continue
		}
		tts := ttsValues[idx]
		tts = append(tts, cardHeaderI18nTextTag{Tag: tt.Tag, Text: tt.Text, Color: tt.Color})
		ttsValues[idx] = tts
	}

	var buf bytes.Buffer
	buf.Grow(len(languages) * 48)

	buf.WriteString("{")

	commas := len(languages)
	for i := range languages {
		commas--

		language := string(languages[i])

		buf.WriteString(fmt.Sprintf(`"%s":`, _quoteEscaper.Replace(language)))

		{
			buf.WriteString("[")
			commas := len(ttsValues[i])
			for _, tt := range ttsValues[i] {
				commas--

				bs, err := json.Marshal(tt)
				if err != nil {
					return nil, fmt.Errorf("marshal(%s): %w", language, err)
				}
				buf.Write(bs)

				if commas > 0 {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]")
		}

		if commas > 0 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

type cardHeaderI18nTextTag struct {
	language Language

	// 标题标签的标识。固定取值：text_tag
	Tag string `json:"tag"`

	// 标题标签的内容。基于文本组件的 plain_text 模式定义内容
	Text cardHeaderComponentPlainText `json:"text"`

	// 标题标签的颜色，默认为蓝色（blue）
	Color CardHeaderTextTagColor `json:"color,omitempty"`
}

type cardHeaderComponentPlainText struct {
	// 固定取值：plain_text
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// ----------------------------------------

// CardHeaderTextTagColor 标签样式
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-header#616d726a
type CardHeaderTextTagColor string

const (
	CardHeaderTextTagColorNeutral   CardHeaderTextTagColor = "neutral"
	CardHeaderTextTagColorBlue      CardHeaderTextTagColor = "blue"
	CardHeaderTextTagColorTurquoise CardHeaderTextTagColor = "turquoise"
	CardHeaderTextTagColorLime      CardHeaderTextTagColor = "lime"
	CardHeaderTextTagColorOrange    CardHeaderTextTagColor = "orange"
	CardHeaderTextTagColorViolet    CardHeaderTextTagColor = "violet"
	CardHeaderTextTagColorIndigo    CardHeaderTextTagColor = "indigo"
	CardHeaderTextTagColorWathet    CardHeaderTextTagColor = "wathet"
	CardHeaderTextTagColorGreen     CardHeaderTextTagColor = "green"
	CardHeaderTextTagColorYellow    CardHeaderTextTagColor = "yellow"
	CardHeaderTextTagColorRed       CardHeaderTextTagColor = "red"
	CardHeaderTextTagColorPurple    CardHeaderTextTagColor = "purple"
	CardHeaderTextTagColorCarmine   CardHeaderTextTagColor = "carmine"
)

// --------------------------------------------------------------------------------

var _ json.Marshaler = (*cardHeaderI18nTexts)(nil)

type cardHeaderI18nTexts []cardHeaderI18nText

func (ts cardHeaderI18nTexts) MarshalJSON() ([]byte, error) {
	if ts == nil {
		return []byte("null"), nil
	}

	var (
		languages = make([]Language, 0, len(ts))
		values    = make(cardHeaderI18nTexts, 0, len(ts))
	)
	for _, t := range ts {
		if slices.Contains(languages, t.language) {
			continue
		}
		languages = append(languages, t.language)
		values = append(values, cardHeaderI18nText{text: t.text})
	}

	var buf bytes.Buffer
	buf.Grow(len(languages) * 24)

	buf.WriteString("{")

	commas := len(languages)
	for i := range languages {
		commas--

		language := string(languages[i])
		t := values[i]

		buf.WriteString(fmt.Sprintf(
			`"%s":"%s"`,
			_quoteEscaper.Replace(language), _quoteEscaper.Replace(t.text),
		))

		if commas > 0 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

type cardHeaderI18nText struct {
	language Language
	text     string
}

// --------------------------------------------------------------------------------

type CardElement interface {
	Entity() any
}

var (
	_ CardElement = (*CardElementColumnSet)(nil)
	_ CardElement = (*CardElementDiv)(nil)
	_ CardElement = (*CardElementMarkdown)(nil)
	_ CardElement = (*CardElementAction)(nil)
	_ CardElement = (*CardElementHorizontalRule)(nil)
	_ CardElement = (*CardElementImage)(nil)
	_ CardElement = (*CardElementNote)(nil)
)

// ----------------------------------------

// CardElementDiv
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#6bdb3f37
type CardElementDiv struct {
	div cardElementDiv
}

func (e *CardElementDiv) Entity() any {
	return e.div
}

type CardElementDivTextMode string

const (
	CardElementDivTextModePlainText    CardElementDivTextMode = "plain_text"
	CardElementDivTextModeLarkMarkdown CardElementDivTextMode = "lark_md"
)

type CardElementDivFieldText struct {
	// 是否并排布局
	//  - true：并排
	//  - false：不并排
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field#3827dadd
	IsShort bool

	// Mode
	//  - CardElementDivTextModePlainText
	//  - CardElementDivTextModeLarkMarkdown
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text#3827dadd
	Mode CardElementDivTextMode

	Content string

	// 内容显示行数
	//
	// 该字段仅支持 plain_text 模式，不支持 lark_md 模式
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text#3827dadd
	Lines int
}

func NewCardElementDiv() *CardElementDiv {
	return &CardElementDiv{div: cardElementDiv{Tag: "div"}}
}

// PlainText 单个文本内容(普通文本内容)
//
// lines: 内容显示行数
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text#3827dadd
func (e *CardElementDiv) PlainText(content string, lines int) *CardElementDiv {
	e.div.Text = &cardElementDivText{
		Tag:     string(CardElementDivTextModePlainText),
		Content: content,
		Lines:   lines,
	}
	return e
}

// LarkMarkdown 单个文本内容(支持部分 Markdown 语法的文本内容)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags
func (e *CardElementDiv) LarkMarkdown(content string) *CardElementDiv {
	e.div.Text = &cardElementDivText{
		Tag:     string(CardElementDivTextModeLarkMarkdown),
		Content: content,
	}
	return e
}

// Fields 双列文本
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field
func (e *CardElementDiv) Fields(fields []CardElementDivFieldText) *CardElementDiv {
	val := make(cardElementDivFields, len(fields))
	for i := range fields {
		f := fields[i]
		val[i] = cardElementDivField{
			IsShort: f.IsShort,
			Text: cardElementDivText{
				Tag:     string(f.Mode),
				Content: f.Content,
				Lines:   f.Lines,
			},
		}
	}
	e.div.Fields = &val
	return e
}

// ExtraImage 在文本右侧附加图片元素
//
// preview: 点击后是否放大图片。在配置 card_link 后可设置为false，使用户点击卡片上的图片也能响应card_link链接跳转
// altContent: 图片hover说明
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/image
func (e *CardElementDiv) ExtraImage(imgKey string, preview bool, altContent string) *CardElementDiv {
	e.div.Extra = cardElementExtraImage{
		Tag:    "img",
		ImgKey: imgKey,
		Alt: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: altContent,
		},
		Preview: &preview,
	}
	return e
}

// ExtraAction 在文本右侧附加交互组件
//   - NewCardElementActionButton
//   - NewCardElementActionOverflow
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#6bdb3f37
func (e *CardElementDiv) ExtraAction(component CardElementActionComponent) *CardElementDiv {
	if component == nil {
		return e
	}
	e.div.Extra = component.ActionEntity()
	return e
}

func (e *CardElementDiv) _() *CardElementDiv {
	return e
}

// ----------------------------------------

// CardElementMarkdown
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags
type CardElementMarkdown struct {
	md cardElementMarkdown
}

func (e *CardElementMarkdown) Entity() any {
	return e.md
}

type CardElementMarkdownTextAlign string

const (
	CardElementMarkdownTextAlignLeft   CardElementMarkdownTextAlign = "left"
	CardElementMarkdownTextAlignCenter CardElementMarkdownTextAlign = "center"
	CardElementMarkdownTextAlignRight  CardElementMarkdownTextAlign = "right"
)

func NewCardElementMarkdown(content string) *CardElementMarkdown {
	return &CardElementMarkdown{md: cardElementMarkdown{Tag: "markdown", Content: content}}
}

func (e *CardElementMarkdown) Content(content string) *CardElementMarkdown {
	e.md.Content = content
	return e
}

// TextAlign 文本内容的对齐方式
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#3827dadd
func (e *CardElementMarkdown) TextAlign(textAlign CardElementMarkdownTextAlign) *CardElementMarkdown {
	e.md.TextAlign = string(textAlign)
	return e
}

// Href 差异化跳转。仅在 PC 端、移动端需要跳转不同链接时使用
//
// [差异化跳转]($urlVal)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#3827dadd
func (e *CardElementMarkdown) Href(defaultURL, pc, ios, android string) *CardElementMarkdown {
	e.md.Href = &cardElementMarkdownHref{
		URLVal: differentJumpLinks{
			URL:     defaultURL,
			PC:      pc,
			IOS:     ios,
			Android: android,
		},
	}
	return e
}

// ----------------------------------------

// CardElementHorizontalRule 分割线
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/divider-line-module
type CardElementHorizontalRule struct {
	hr cardElementHorizontalRule
}

func (e *CardElementHorizontalRule) Entity() any {
	return e.hr
}

func NewCardElementHorizontalRule() CardElement {
	return &CardElementHorizontalRule{hr: cardElementHorizontalRule{Tag: "hr"}}
}

// ----------------------------------------

// CardElementImage
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
type CardElementImage struct {
	img cardElementImage
}

func (e *CardElementImage) Entity() any {
	return e.img
}

func NewCardElementImage(imgKey, altContent string) *CardElementImage {
	return &CardElementImage{img: cardElementImage{
		Tag:    "img",
		ImgKey: imgKey,
		Alt: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: altContent,
		},
		Title:        nil,
		CustomWidth:  nil,
		CompactWidth: nil,
		Mode:         "",
		Preview:      nil,
	}}
}

// TitleWithPlainText 图片标题(plain_text)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
func (e *CardElementImage) TitleWithPlainText(title string) *CardElementImage {
	e.img.Title = &cardElementDivText{
		Tag:     string(CardElementDivTextModePlainText),
		Content: title,
	}
	return e
}

// TitleWithLarkMarkdown 图片标题(lark_md)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
func (e *CardElementImage) TitleWithLarkMarkdown(title string) *CardElementImage {
	e.img.Title = &cardElementDivText{
		Tag:     string(CardElementDivTextModeLarkMarkdown),
		Content: title,
	}
	return e
}

// CustomWidth 自定义图片的最大展示宽度，支持在 278px ~ 580px 范围内指定最大展示宽度
//
// 默认情况下图片宽度与图片组件所占区域的宽度一致
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
func (e *CardElementImage) CustomWidth(px int) *CardElementImage {
	e.img.CustomWidth = &px
	return e
}

// CompactWidth 是否展示为紧凑型的图片
//
// 默认值为 false。如果配置为 true，则展示最大宽度为 278px 的紧凑型图片
func (e *CardElementImage) CompactWidth(b bool) *CardElementImage {
	e.img.CompactWidth = &b
	return e
}

type CardElementImageMode string

const (
	CardElementImageModeCropCenter    CardElementImageMode = "crop_center"
	CardElementImageModeFitHorizontal CardElementImageMode = "fit_horizontal"
	CardElementImageModeStretch       CardElementImageMode = "stretch"
	CardElementImageModeLarge         CardElementImageMode = "large"
	CardElementImageModeMedium        CardElementImageMode = "medium"
	CardElementImageModeSmall         CardElementImageMode = "small"
	CardElementImageModeTiny          CardElementImageMode = "tiny"
)

// Mode 图片显示模式
//   - crop_center：居中裁剪模式，对长图会限高，并居中裁剪后展示
//   - fit_horizontal：平铺模式，宽度撑满卡片完整展示上传的图片
//   - stretch：自适应。图片宽度撑满卡片宽度，当图片 高:宽 小于 16:9 时，完整展示原图。当图片 高:宽 大于 16:9 时，顶部对齐裁剪图片，并在图片底部展示 长图 脚标
//   - large：大图，尺寸为 160 × 160，适用于多图混排
//   - medium：中图，尺寸为 80 × 80，适用于图文混排的封面图
//   - small：小图，尺寸为 40 × 40，适用于人员头像
//   - tiny：超小图，尺寸为 16 × 16，适用于图标、备注
//
// 注意：设置该参数后，会覆盖 custom_width 参数。更多信息参见 消息卡片设计规范 https://open.feishu.cn/document/tools-and-resources/design-specification/message-card-design-specifications
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
func (e *CardElementImage) Mode(mode CardElementImageMode) *CardElementImage {
	e.img.Mode = string(mode)
	return e
}

// Preview 点击后是否放大图片
//
// 默认值为 true，即点击后放大图片
//
// 如果你为卡片配置了 消息卡片跳转链接，
// 可将该参数设置为 false，后续用户点击卡片上的图片也能响应 card_link 链接跳转
func (e *CardElementImage) Preview(b bool) *CardElementImage {
	e.img.Preview = &b
	return e
}

// ----------------------------------------

// CardElementNote 备注
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/notes-module
type CardElementNote struct {
	note cardElementNote
}

func (e *CardElementNote) Entity() any {
	return e.note
}

func NewCardElementNote() *CardElementNote {
	return &CardElementNote{note: cardElementNote{
		Tag:      "note",
		Elements: make([]any, 0, 2),
	}}
}

// AddElementWithPlainText 添加文本(plain_text)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/notes-module#3827dadd
func (e *CardElementNote) AddElementWithPlainText(s string) *CardElementNote {
	e.note.Elements = append(e.note.Elements,
		cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: s,
		},
	)
	return e
}

// AddElementWithLarkMarkdown 添加文本(lark_md)
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/notes-module#3827dadd
func (e *CardElementNote) AddElementWithLarkMarkdown(s string) *CardElementNote {
	e.note.Elements = append(e.note.Elements,
		cardElementDivText{
			Tag:     string(CardElementDivTextModeLarkMarkdown),
			Content: s,
		},
	)
	return e
}

// AddElementWithImage 添加图片
//
// preview: 点击后是否放大图片。在配置 card_link 后可设置为false，使用户点击卡片上的图片也能响应card_link链接跳转
// altContent: 图片hover说明
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/notes-module#3827dadd
func (e *CardElementNote) AddElementWithImage(imgKey string, preview bool, altContent string) *CardElementNote {
	e.note.Elements = append(e.note.Elements,
		cardElementExtraImage{
			Tag:    "img",
			ImgKey: imgKey,
			Alt: cardElementDivText{
				Tag:     string(CardElementDivTextModePlainText),
				Content: altContent,
			},
			Preview: &preview,
		},
	)
	return e
}

// ----------------------------------------

// CardElementAction 交互模块（action）
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#60ddc64e
type CardElementAction struct {
	action cardElementAction
}

func (e *CardElementAction) Entity() any {
	return e.action
}

func NewCardElementAction() *CardElementAction {
	return &CardElementAction{action: cardElementAction{
		Tag:     "action",
		Actions: make([]any, 0, 3),
		Layout:  "",
	}}
}

func (e *CardElementAction) Actions(actions []CardElementActionComponent) *CardElementAction {
	for i := range actions {
		if actions[i] == nil {
			continue
		}
		e.action.Actions = append(e.action.Actions, actions[i].ActionEntity())
	}
	return e
}

type CardElementActionLayout string

const (
	CardElementActionLayoutBisected   CardElementActionLayout = "bisected"
	CardElementActionLayoutTrisection CardElementActionLayout = "trisection"
	CardElementActionLayoutFlow       CardElementActionLayout = "flow"
)

// Layout 设置窄屏自适应布局方式
//   - bisected：二等分布局，每行两列交互元素
//   - trisection：三等分布局，每行三列交互元素
//   - flow：流式布局，元素会按自身大小横向排列并在空间不够的时候折行
func (e *CardElementAction) Layout(layout CardElementActionLayout) *CardElementAction {
	e.action.Layout = string(layout)
	return e
}

// ----------------------------------------

// CardElementColumnSet 多列布局
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set
type CardElementColumnSet struct {
	cs cardElementColumnSet
}

func (e *CardElementColumnSet) Entity() any {
	return e.cs
}

func NewCardElementColumnSet() *CardElementColumnSet {
	return &CardElementColumnSet{cs: cardElementColumnSet{
		Tag:               "column_set",
		FlexMode:          string(CardElementColumnSetFlexModeNone),
		BackgroundStyle:   "",
		HorizontalSpacing: "",
		Columns:           make([]cardElementColumnSetColumn, 0, 2),
		Action:            nil,
	}}
}

type CardElementColumnSetFlexMode string

const (
	CardElementColumnSetFlexModeNone    CardElementColumnSetFlexMode = "none"
	CardElementColumnSetFlexModeStretch CardElementColumnSetFlexMode = "stretch"
	CardElementColumnSetFlexModeFlow    CardElementColumnSetFlexMode = "flow"
	CardElementColumnSetFlexModeBisect  CardElementColumnSetFlexMode = "bisect"
	CardElementColumnSetFlexModeTrisect CardElementColumnSetFlexMode = "trisect"
)

// FlexMode 移动端和 PC 端的窄屏幕下，各列的自适应方式
//   - none：不做布局上的自适应，在窄屏幕下按比例压缩列宽度
//   - stretch：列布局变为行布局，且每列（行）宽度强制拉伸为 100%，所有列自适应为上下堆叠排布
//   - flow：列流式排布（自动换行），当一行展示不下一列时，自动换至下一行展示
//   - bisect：两列等分布局
//   - trisect：三列等分布局
//
// 默认值：none
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
func (e *CardElementColumnSet) FlexMode(mode CardElementColumnSetFlexMode) *CardElementColumnSet {
	e.cs.FlexMode = string(mode)
	return e
}

type CardElementColumnSetBackgroundStyle string

const (
	CardElementColumnSetBackgroundStyleDefault CardElementColumnSetBackgroundStyle = "default"
	CardElementColumnSetBackgroundStyleGrey    CardElementColumnSetBackgroundStyle = "grey"
)

// BackgroundStyle 多列布局的背景色样式
//   - default：默认的白底样式，dark mode 下为黑底
//   - grey：灰底样式
//
// 当存在多列布局的嵌套时，上层多列布局的颜色覆盖下层多列布局的颜色
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
func (e *CardElementColumnSet) BackgroundStyle(bs CardElementColumnSetBackgroundStyle) *CardElementColumnSet {
	e.cs.BackgroundStyle = string(bs)
	return e
}

type CardElementColumnSetHorizontalSpacing string

const (
	CardElementColumnSetHorizontalSpacingDefault CardElementColumnSetHorizontalSpacing = "default"
	CardElementColumnSetHorizontalSpacingSmall   CardElementColumnSetHorizontalSpacing = "small"
)

// HorizontalSpacing 多列布局内，各列之间的水平分栏间距
//   - default：默认间距
//   - small：窄间距
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
func (e *CardElementColumnSet) HorizontalSpacing(hs CardElementColumnSetHorizontalSpacing) *CardElementColumnSet {
	e.cs.HorizontalSpacing = string(hs)
	return e
}

// ActionMultiURL 设置点击布局容器时的交互配置。当前仅支持跳转交互。如果布局容器内有交互组件，则优先响应交互组件定义的交互
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
func (e *CardElementColumnSet) ActionMultiURL(defaultURL, pc, ios, android string) *CardElementColumnSet {
	e.cs.Action = &cardElementColumnSetAction{
		MultiURL: differentJumpLinks{
			URL:     defaultURL,
			PC:      pc,
			IOS:     ios,
			Android: android,
		},
	}
	return e
}

// Columns 多列布局容器内，各个列容器的配置信息
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
func (e *CardElementColumnSet) Columns(columns []*CardElementColumnSetColumn) *CardElementColumnSet {
	for i := range columns {
		if columns[i] == nil || columns[i].csc.Tag == "" {
			continue
		}
		e.cs.Columns = append(e.cs.Columns, columns[i].csc)
	}
	return e
}

type CardElementColumnSetColumn struct {
	csc cardElementColumnSetColumn
}

func NewCardElementColumnSetColumn() *CardElementColumnSetColumn {
	return &CardElementColumnSetColumn{csc: cardElementColumnSetColumn{
		Tag:           "column",
		Width:         "",
		Weight:        nil,
		VerticalAlign: "",
		Elements:      nil,
	}}
}

type CardElementColumnSetColumnWidth string

const (
	CardElementColumnSetColumnWidthAuto     CardElementColumnSetColumnWidth = "auto"
	CardElementColumnSetColumnWidthWeighted CardElementColumnSetColumnWidth = "weighted"
)

// Width 列宽度属性
//   - auto：列宽度与列内元素宽度一致
//   - weighted：列宽度按 weight 参数定义的权重分布
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#1cf15c63
func (e *CardElementColumnSetColumn) Width(width CardElementColumnSetColumnWidth) *CardElementColumnSetColumn {
	e.csc.Width = string(width)
	return e
}

// Weight 当 width 取值 weighted 时生效，表示当前列的宽度占比。取值范围：1 ~ 5
//
// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#1cf15c63
func (e *CardElementColumnSetColumn) Weight(n int) *CardElementColumnSetColumn {
	e.csc.Weight = &n
	return e
}

type CardElementColumnSetColumnVerticalAlign string

const (
	CardElementColumnSetColumnVerticalAlignTop    CardElementColumnSetColumnVerticalAlign = "top"
	CardElementColumnSetColumnVerticalAlignCenter CardElementColumnSetColumnVerticalAlign = "center"
	CardElementColumnSetColumnVerticalAlignBottom CardElementColumnSetColumnVerticalAlign = "bottom"
)

// VerticalAlign 列内成员垂直对齐方式
//   - top：顶对齐
//   - center：居中对齐
//   - bottom：底部对齐
func (e *CardElementColumnSetColumn) VerticalAlign(va CardElementColumnSetColumnVerticalAlign) *CardElementColumnSetColumn {
	e.csc.VerticalAlign = string(va)
	return e
}

func (e *CardElementColumnSetColumn) Elements(elements []CardElement) *CardElementColumnSetColumn {
	for i := range elements {
		if elements[i] == nil {
			continue
		}

		e.csc.Elements = append(e.csc.Elements, elements[i].Entity())
	}
	return e
}

// ----------------------------------------

type CardElementActionComponent interface {
	ActionEntity() any
}

var (
	_ CardElementActionComponent = (*CardElementActionButton)(nil)
	_ CardElementActionComponent = (*CardElementActionOverflow)(nil)
)

// CardElementActionButton
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button
type CardElementActionButton struct {
	button cardElementActionButton
}

func (e *CardElementActionButton) ActionEntity() any {
	return e.button
}

func NewCardElementActionButton(mode CardElementDivTextMode, content string) *CardElementActionButton {
	return &CardElementActionButton{button: cardElementActionButton{
		Tag: "button",
		Text: cardElementDivText{
			Tag:     string(mode),
			Content: content,
		},
		URL:      "",
		MultiURL: nil,
		Type:     "",
	}}
}

// URL 点击按钮后的跳转链接。不可与 MultiURL 同时设置
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button#3827dadd
func (e *CardElementActionButton) URL(s string) *CardElementActionButton {
	e.button.URL = s
	return e
}

// MultiURL 基于 url 元素配置多端跳转链接，不可与 URL 同时设置
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button#3827dadd
func (e *CardElementActionButton) MultiURL(defaultURL, pc, ios, android string) *CardElementActionButton {
	e.button.MultiURL = &differentJumpLinks{
		URL:     defaultURL,
		PC:      pc,
		IOS:     ios,
		Android: android,
	}
	return e
}

type CardElementActionButtonType string

const (
	CardElementActionButtonTypeDefault CardElementActionButtonType = "default"
	CardElementActionButtonTypePrimary CardElementActionButtonType = "primary"
	CardElementActionButtonTypeDanger  CardElementActionButtonType = "danger"
)

// Type 配置按钮样式
//   - default：默认样式
//   - primary：强调样式
//   - danger：警示样式
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button#3827dadd
func (e *CardElementActionButton) Type(typ CardElementActionButtonType) *CardElementActionButton {
	e.button.Type = string(typ)
	return e
}

// Confirm 设置二次确认弹框
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button?lang=zh-CN#3827dadd
func (e *CardElementActionButton) Confirm(title, text string) *CardElementActionButton {
	e.button.Confirm = &cardElementActionConfirm{
		Title: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: title,
		},
		Text: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: text,
		},
	}
	return e
}

// ----------------------------------------

// CardElementActionOverflow 折叠按钮组（overflow）
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/overflow
type CardElementActionOverflow struct {
	overflow cardElementActionOverflow
}

func (e *CardElementActionOverflow) ActionEntity() any {
	return e.overflow
}

func NewCardElementActionOverflow() *CardElementActionOverflow {
	return &CardElementActionOverflow{overflow: cardElementActionOverflow{
		Tag:     "overflow",
		Options: make([]cardElementActionOption, 0, 2),
		Confirm: nil,
	}}
}

// AddOptionWithURL 添加跳转链接的选项
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#9fa21514
func (e *CardElementActionOverflow) AddOptionWithURL(text, defaultURL string) *CardElementActionOverflow {
	e.overflow.Options = append(e.overflow.Options, cardElementActionOption{
		Text: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: text,
		},
		URL:      defaultURL,
		MultiURL: nil,
	})
	return e
}

// AddOptionWithMultiURL 添加多端跳转链接的选项
//
// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#9fa21514
func (e *CardElementActionOverflow) AddOptionWithMultiURL(text, defaultURL, pc, ios, android string) *CardElementActionOverflow {
	e.overflow.Options = append(e.overflow.Options, cardElementActionOption{
		Text: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: text,
		},
		URL: "",
		MultiURL: &differentJumpLinks{
			URL:     defaultURL,
			PC:      pc,
			IOS:     ios,
			Android: android,
		},
	})
	return e
}

// Confirm 设置二次确认弹框
//
// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/overflow#3827dadd
func (e *CardElementActionOverflow) Confirm(title, text string) *CardElementActionOverflow {
	e.overflow.Confirm = &cardElementActionConfirm{
		Title: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: title,
		},
		Text: cardElementDivText{
			Tag:     string(CardElementDivTextModePlainText),
			Content: text,
		},
	}
	return e
}

// ----------------------------------------

type (
	// cardElementDiv 内容模块（div）
	//
	// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#6bdb3f37
	cardElementDiv struct {
		// 内容模块的标识。固定取值：div
		Tag string `json:"tag"`

		// 单个文本内容
		//
		// 参数配置详情: https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text
		Text *cardElementDivText `json:"text,omitempty"`

		// 双列文本
		//
		// 参数配置详情: https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field
		Fields *cardElementDivFields `json:"fields,omitempty"`

		// 附加元素，添加后展示在文本右侧。支持附加的元素：
		//  - 图片（image）
		//  - 按钮（button）
		//  - 列表选择器（selectMenu）
		//  - 折叠按钮组（overflow）
		//  - 日期选择器（datePicker）
		Extra any `json:"extra,omitempty"`
	}

	// cardElementDivText
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text
	cardElementDivText struct {
		// 文本元素的标签
		//  - plain_text：普通文本内容
		//  - lark_md：支持部分 Markdown 语法的文本内容。关于 Markdown 语法的详细介绍，可参见 Markdown https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags
		Tag string `json:"tag"`

		Content string `json:"content"`

		// 内容显示行数
		//
		// 该字段仅支持 text 的 plain_text 模式，不支持 lark_md 模式
		Lines int `json:"lines,omitempty"`
	}

	// cardElementDivFields 双列文本
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field
	cardElementDivFields []cardElementDivField
	// cardElementDivField
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/field#3827dadd
	cardElementDivField struct {
		// 是否并排布局
		//  - true：并排
		//  - false：不并排
		IsShort bool `json:"is_short"`

		Text cardElementDivText `json:"text"`
	}

	// cardElementExtraImage 内容元素的一种，可用于内容块的extra字段和备注块的elements字段
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/image
	cardElementExtraImage struct {
		// 元素标签。固定取值：img
		Tag string `json:"tag"`

		// 图片资源，获取方式：https://open.feishu.cn/document/server-docs/im-v1/image/create?appId=cli_a2850f3fbb38900d
		ImgKey string `json:"img_key"`

		// 图片hover说明
		Alt cardElementDivText `json:"alt"`

		// 点击后是否放大图片，缺省为true。
		//
		// 在配置 card_link 后可设置为false，使用户点击卡片上的图片也能响应card_link链接跳转
		Preview *bool `json:"preview,omitempty"`
	}

	// cardElementMarkdown
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags
	cardElementMarkdown struct {
		// Markdown 组件的标识。固定取值：markdown
		Tag string `json:"tag"`

		// 使用已支持的 Markdown 语法构造 Markdown 内容。
		//
		// 语法详情：https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#abc9b025
		Content string `json:"content"`

		// 文本内容的对齐方式
		//  - left：左对齐
		//  - center：居中对齐
		//  - right：右对齐
		TextAlign string `json:"text_align,omitempty"`

		// 差异化跳转。仅在 PC 端、移动端需要跳转不同链接时使用
		Href *cardElementMarkdownHref `json:"href,omitempty"`
	}

	cardElementMarkdownHref struct {
		URLVal differentJumpLinks `json:"urlVal"`
	}

	// cardElementHorizontalRule
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/divider-line-module
	cardElementHorizontalRule struct {
		// 分割线模块标识，固定取值：hr
		Tag string `json:"tag"`
	}

	// cardElementImage
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/image-module#3827dadd
	cardElementImage struct {
		// 图片组件的标签，固定取值：img
		Tag string `json:"tag"`

		ImgKey string `json:"img_key"`

		// 悬浮（hover）图片时弹出的说明文案
		//
		// 使用文本组件的数据结构展示文案，详情参见文本组件
		//
		// 当文本组件的 content 参数取值为空时，不展示图片文案内容
		Alt cardElementDivText `json:"alt"`

		// 图片标题
		Title *cardElementDivText `json:"title,omitempty"`

		// 自定义图片的最大展示宽度，支持在 278px ~ 580px 范围内指定最大展示宽度
		//
		// 默认情况下图片宽度与图片组件所占区域的宽度一致
		CustomWidth *int `json:"custom_width,omitempty"`

		// 是否展示为紧凑型的图片。
		//
		// 默认值为 false。如果配置为 true，则展示最大宽度为 278px 的紧凑型图片。
		CompactWidth *bool `json:"compact_width,omitempty"`

		// 图片显示模式
		//  - crop_center：居中裁剪模式，对长图会限高，并居中裁剪后展示
		//  - fit_horizontal：平铺模式，宽度撑满卡片完整展示上传的图片
		//  - stretch：自适应。图片宽度撑满卡片宽度，当图片 高:宽 小于 16:9 时，完整展示原图。当图片 高:宽 大于 16:9 时，顶部对齐裁剪图片，并在图片底部展示 长图 脚标
		//  - large：大图，尺寸为 160 × 160，适用于多图混排
		//  - medium：中图，尺寸为 80 × 80，适用于图文混排的封面图
		//  - small：小图，尺寸为 40 × 40，适用于人员头像
		//  - tiny：超小图，尺寸为 16 × 16，适用于图标、备注
		//
		// 注意：设置该参数后，会覆盖 custom_width 参数。更多信息参见 消息卡片设计规范 https://open.feishu.cn/document/tools-and-resources/design-specification/message-card-design-specifications
		Mode string `json:"mode,omitempty"`

		// 点击后是否放大图片
		//
		// 默认值为 true，即点击后放大图片
		//
		// 如果你为卡片配置了 消息卡片跳转链接，
		// 可将该参数设置为 false，后续用户点击卡片上的图片也能响应 card_link 链接跳转
		Preview *bool `json:"preview,omitempty"`
	}

	cardElementNote struct {
		// 备注组件的标识。固定取值：note
		Tag string `json:"tag"`

		// 备注信息。支持添加的元素：
		//  - 文本组件的数据结构，构成备注信息的文本内容。数据结构参见文本组件 https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text
		//  - image 元素的数据结构，构成备注信息的小尺寸图片。数据结构参见 image 元素 https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#a974e363
		Elements []any `json:"elements"`
	}

	cardElementAction struct {
		// 交互模块的标识。固定取值：action
		Tag string `json:"tag"`

		// 添加可交互的组件。支持添加的组件：
		//  - 按钮（button）
		//  - 列表选择器（selectMenu）
		//  - 折叠按钮组（overflow）
		//  - 日期选择器（datePicker）
		Actions []any `json:"actions"`

		// 设置窄屏自适应布局方式
		//  - bisected：二等分布局，每行两列交互元素
		//  - trisection：三等分布局，每行三列交互元素
		//  - flow：流式布局，元素会按自身大小横向排列并在空间不够的时候折行
		Layout string `json:"layout,omitempty"`
	}

	// cardElementActionButton
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/button
	cardElementActionButton struct {
		// 按钮组件的标识。固定取值：button
		Tag string `json:"tag"`

		// 按钮中的文本。基于文本组件的数据结构配置文本内容，详情参见文本组件 https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/embedded-non-interactive-elements/text
		Text cardElementDivText `json:"text"`

		// 点击按钮后的跳转链接。该字段与 multi_url 字段不可同时设置
		URL string `json:"url,omitempty"`

		// 基于 url 元素配置多端跳转链接，详情参见url 元素。该字段与 url 字段不可同时设置 https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#09a320b3
		MultiURL *differentJumpLinks `json:"multi_url,omitempty"`

		// 配置按钮样式
		//  - default：默认样式
		//  - primary：强调样式
		//  - danger：警示样式
		Type string `json:"type,omitempty"`

		// value 该字段用于交互组件的回传交互方式,当用户点击交互组件后，会将 value 的值返回给接收回调数据的服务器。后续你可以通过服务器接收的 value 值进行业务处理
		//
		// 自定义机器人发送的消息卡片，只支持通过按钮、文字链方式跳转 URL，不支持点击后回调信息到服务端的回传交互
		// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#4996824a
		// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interaction-module

		// 设置二次确认弹框
		//
		// confirm 元素的配置方式可参见 confirm https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7f700aa9
		Confirm *cardElementActionConfirm `json:"confirm,omitempty"`
	}

	// cardElementActionConfirm
	//
	// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7f700aa9
	cardElementActionConfirm struct {
		// 弹窗标题。由文本组件构成（仅支持文本组件的 plain_text 模式）
		//
		// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7f700aa9
		Title cardElementDivText `json:"title"`

		// 弹窗内容。由文本组件构成（仅支持文本组件的 plain_text 模式）
		//
		// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7f700aa9
		Text cardElementDivText `json:"text"`
	}

	// cardElementActionOverflow
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/add-card-interaction/interactive-components/overflow
	cardElementActionOverflow struct {
		// 折叠按钮组的标签。固定取值：overflow
		Tag string `json:"tag"`

		// 折叠按钮组当中的选项按钮。按钮基于 option 元素进行配置
		//
		// 详情参见 option 元素 https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#9fa21514
		Options []cardElementActionOption `json:"options"`

		// 设置二次确认弹框。confirm 元素的配置方式可参见 confirm https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#7f700aa9
		Confirm *cardElementActionConfirm `json:"confirm,omitempty"`
	}

	// cardElementActionOption
	//
	// https://open.feishu.cn/document/ukTMukTMukTM/uYzM3QjL2MzN04iNzcDN/component-list/common-components-and-elements#9fa21514
	cardElementActionOption struct {
		// 选项显示的内容
		Text cardElementDivText `json:"text"`

		// 选项的跳转链接，仅支持在折叠按钮组（overflow）中设置
		//
		// url 和 multi_url 字段必须且仅能填写其中一个
		URL string `json:"url,omitempty"`

		// 选项的跳转链接,仅支持在折叠按钮组（overflow）中设置。支持按操作系统设置不同的链接，参数配置详情参见 链接元素（url）
		//
		// url 和 multi_url 字段必须且仅能填写其中一个
		MultiURL *differentJumpLinks `json:"multi_url,omitempty"`
	}

	// cardElementColumnSet
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set
	cardElementColumnSet struct {
		// 多列布局容器的标识，固定取值：column_set
		Tag string `json:"tag"`

		// 移动端和 PC 端的窄屏幕下，各列的自适应方式
		//  - none：不做布局上的自适应，在窄屏幕下按比例压缩列宽度
		//  - stretch：列布局变为行布局，且每列（行）宽度强制拉伸为 100%，所有列自适应为上下堆叠排布
		//  - flow：列流式排布（自动换行），当一行展示不下一列时，自动换至下一行展示
		//  - bisect：两列等分布局
		//  - trisect：三列等分布局
		//
		// 默认值：none
		FlexMode string `json:"flex_mode"`

		// 多列布局的背景色样式
		//  - default：默认的白底样式，dark mode 下为黑底
		//  - grey：灰底样式
		//
		// 当存在多列布局的嵌套时，上层多列布局的颜色覆盖下层多列布局的颜色
		BackgroundStyle string `json:"background_style,omitempty"`

		// 多列布局内，各列之间的水平分栏间距
		//  - default：默认间距
		//  - small：窄间距
		HorizontalSpacing string `json:"horizontal_spacing,omitempty"`

		// 多列布局容器内，各个列容器的配置信息
		//
		// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
		Columns []cardElementColumnSetColumn `json:"columns,omitempty"`

		// 设置点击布局容器时的交互配置。当前仅支持跳转交互。如果布局容器内有交互组件，则优先响应交互组件定义的交互
		//
		// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#3827dadd
		Action *cardElementColumnSetAction `json:"action,omitempty"`
	}
	cardElementColumnSetAction struct {
		MultiURL differentJumpLinks `json:"multi_url"`
	}

	// cardElementColumnSetColumn
	//
	// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#1cf15c63
	cardElementColumnSetColumn struct {
		// 列容器标识，固定取值：column
		Tag string `json:"tag"`

		// 列宽度属性
		//  - auto：列宽度与列内元素宽度一致
		//  - weighted：列宽度按 weight 参数定义的权重分布
		Width string `json:"width,omitempty"`

		// 当 width 取值 weighted 时生效，表示当前列的宽度占比。取值范围：1 ~ 5
		Weight *int `json:"weight,omitempty"`

		// 列内成员垂直对齐方式
		//  - top：顶对齐
		//  - center：居中对齐
		//  - bottom：底部对齐
		VerticalAlign string `json:"vertical_align,omitempty"`

		// 需要在列内展示的卡片元素
		//
		// https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/column-set#1cf15c63
		Elements []any `json:"elements,omitempty"`
	}
)
