package fsBotAPI

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type i18nCard struct {
	lang     string        // 卡片所属的语言环境
	title    string        // 卡片所属语言环境的标题
	elements []interface{} // 卡片所属语言环境的所有元素
}

type Card func() i18nCard

// WithCard 卡片消息, 可指定语言环境
//  支持元素如下:
//  普通文本: WithCardElementPlainText
//  MarkDown: WithCardElementMarkdown
//  可并排字段: WithCardElementFields
//  按钮: WithCardElementActions
//  分割线: WithCardElementHorizontalRule
//  图片: WithCardElementImage
//  备注: WithCardElementNote
func WithCard(lang Language, title string, elem CardElement, elements ...CardElement) Card {
	elements = append([]CardElement{elem}, elements...)
	es := make([]interface{}, 0, len(elements))
	for _, fn := range elements {
		if fn == nil {
			continue
		}
		es = append(es, fn(false))
	}
	return func() i18nCard {
		return i18nCard{
			lang:     string(lang),
			title:    title,
			elements: es,
		}
	}
}

func GenMsgCard(bgColor CardTitleBgColor, cfg CardConfig, c Card, more ...Card) map[string]interface{} {
	more = append([]Card{c}, more...)
	cards := make([]i18nCard, 0, len(more))
	for _, fn := range more {
		cards = append(cards, fn())
	}

	i18nTitle := make(map[string]string, 3)
	i18nElements := make(map[string]interface{}, 3)
	for _, c := range cards {
		i18nTitle[c.lang] = c.title
		i18nElements[c.lang] = c.elements
	}

	sub := map[string]interface{}{
		"header":        buildCardHeader(bgColor, i18nTitle),
		"i18n_elements": i18nElements,
	}

	if cfg != nil {
		_cfg := cfg()
		if _cfg.mCfg != nil {
			sub["config"] = _cfg.mCfg
		}
		if _cfg.mCardLink != nil {
			sub["card_link"] = _cfg.mCardLink
		}
	}

	return map[string]interface{}{
		"msg_type": "interactive",
		"card":     sub,
	}
}

func buildCardHeader(bgColor CardTitleBgColor, i18nTitle map[string]string) *msgCardHeader {
	header := defaultMsgCardHeader()
	header.Template = string(bgColor)
	header.Title.I18n = i18nTitle

	return header
}

// CardTitleBgColor 标题背景色
//  最佳实践：https://open.feishu.cn/document/ukTMukTMukTM/ukTNwUjL5UDM14SO1ATN#8239feff
//  - 绿色（Green）代表完成/成功
//  - 橙色（Orange）代表警告/警示
//  - 红色（Red）代表错误/异常
//  - 灰色（Grey）代表失效
type CardTitleBgColor string

const (
	BgColorDefault   CardTitleBgColor = ""
	BgColorBlue      CardTitleBgColor = "blue"
	BgColorWathet    CardTitleBgColor = "wathet"
	BgColorTurquoise CardTitleBgColor = "turquoise"
	BgColorGreen     CardTitleBgColor = "green"
	BgColorYellow    CardTitleBgColor = "yellow"
	BgColorOrange    CardTitleBgColor = "orange"
	BgColorRed       CardTitleBgColor = "red"
	BgColorCarmine   CardTitleBgColor = "carmine"
	BgColorViolet    CardTitleBgColor = "violet"
	BgColorPurple    CardTitleBgColor = "purple"
	BgColorIndigo    CardTitleBgColor = "indigo"
	BgColorGrey      CardTitleBgColor = "grey"
)

type msgCardHeader struct {
	Title    msgCardTitle `json:"title"`
	Template string       `json:"template,omitempty"` // 控制标题背景的颜色
}

type msgCardTitle struct {
	// Content string            `json:"content,omitempty"`
	I18n map[string]string `json:"i18n,omitempty"`
	Tag  string            `json:"tag"`
}

func defaultMsgCardHeader() *msgCardHeader {
	return &msgCardHeader{
		Title: msgCardTitle{
			// Content: " ",
			I18n: nil,
			Tag:  "plain_text", // 仅支持"plain_text"
		},
		Template: "",
	}
}

type cardConfig struct {
	mCfg      map[string]interface{}
	mCardLink map[string]string
}

type CardConfig func() cardConfig

type CardConfigOption func(*cardConfig)

// WithCardConfigEnableForward 设置是否允许卡片被转发, 默认允许转发
func WithCardConfigEnableForward(b bool) CardConfigOption {
	return func(cfg *cardConfig) {
		if cfg.mCfg == nil {
			cfg.mCfg = make(map[string]interface{}, 2)
		}
		cfg.mCfg["enable_forward"] = b
	}
}

// WithCardConfigEnableUpdateMulti 设置是否为共享卡片, 默认不共享
//  true: 是共享卡片，也即更新卡片的内容对所有收到这张卡片的人员可见。
//  false: 是独享卡片，仅操作用户可见卡片的更新内容。
func WithCardConfigEnableUpdateMulti(b bool) CardConfigOption {
	return func(cfg *cardConfig) {
		if cfg.mCfg == nil {
			cfg.mCfg = make(map[string]interface{}, 2)
		}
		cfg.mCfg["update_multi"] = b
	}
}

// WithCardConfigCardLink 设置卡片的多端跳转链接
func WithCardConfigCardLink(url, android, ios, pc string) CardConfigOption {
	return func(cfg *cardConfig) {
		if cfg.mCardLink == nil {
			cfg.mCardLink = make(map[string]string, 4)
		}
		cfg.mCardLink = map[string]string{
			"url":         url,
			"android_url": android,
			"ios_url":     ios,
			"pc_url":      pc,
		}

	}
}

// WithCardConfig 卡片消息的属性配置
//  - 是否允许卡片消息被转发, 默认值: true WithCardConfigEnableForward
//  - 是否为共享卡片, 默认值: false WithCardConfigEnableUpdateMulti
//  - 设置卡片跳转链接 WithCardConfigCardLink
func WithCardConfig(opt CardConfigOption, opts ...CardConfigOption) CardConfig {
	opts = append([]CardConfigOption{opt}, opts...)
	var ret cardConfig
	for _, fn := range opts {
		fn(&ret)
	}
	return func() cardConfig {
		return ret
	}
}

type CardElement func(isEmbedded bool) interface{}

// WithCardElementPlainText 普通文本内容
//  lines: 内容显示行数
func WithCardElementPlainText(text string, lines ...int) CardElement {
	return func(isEmbedded bool) interface{} {
		sub := map[string]interface{}{
			"tag":     "plain_text",
			"content": text,
		}
		if len(lines) != 0 && lines[0] > 0 {
			sub["lines"] = lines[0]
		}
		if isEmbedded {
			return sub
		}
		elem := map[string]interface{}{
			"tag":  "div",
			"text": sub,
		}
		return elem
	}
}

type CardExtraElement func() (key string, v interface{})

func WithCardExtraElementImage(imgKey string, opts ...CardElemImageOption) CardExtraElement {
	return func() (key string, v interface{}) {
		key, v = "extra", WithCardElementImage(imgKey, opts...)(true)
		return
	}
}

// WithCardElementMarkdown MarkDown 语法展示文本内容
//  语法仅支持部分, 语法详情: https://open.feishu.cn/document/ukTMukTMukTM/uADOwUjLwgDM14CM4ATN
func WithCardElementMarkdown(md string, extra ...CardExtraElement) CardElement {
	return func(isEmbedded bool) interface{} {
		sub := map[string]interface{}{
			"tag":     "lark_md",
			"content": md,
		}
		if isEmbedded {
			return sub
		}

		elem := map[string]interface{}{
			"tag":  "div",
			"text": sub,
		}
		for _, fn := range extra {
			k, v := fn()
			elem[k] = v
		}
		return elem
	}
}

type CardElementField func() interface{}

func WithCardElementField(elem CardElement, isShort bool) CardElementField {
	return func() interface{} {
		return map[string]interface{}{
			"text":     elem(true),
			"is_short": isShort,
		}
	}
}

// WithCardElementFields 能并排布局的字段元素
//  支持元素:
//  - WithCardElementPlainText
//  - WithCardElementMarkdown
func WithCardElementFields(f CardElementField, fields ...CardElementField) CardElement {
	fields = append([]CardElementField{f}, fields...)
	fs := make([]interface{}, 0, len(fields))
	for _, fn := range fields {
		if fn == nil {
			continue
		}
		fs = append(fs, fn())
	}

	return func(bool) interface{} {
		return map[string]interface{}{
			"tag":    "div",
			"fields": fs,
		}
	}
}

type ElementButton string

const (
	ButtonDefault ElementButton = "default"
	ButtonPrimary ElementButton = "primary"
	ButtonDanger  ElementButton = "danger"
)

type CardElementActionOption func() (key string, v interface{})

func WithCardElementActionButton(btn ElementButton) CardElementActionOption {
	return func() (key string, v interface{}) {
		key, v = "type", string(btn)
		return
	}
}

func WithCardElementActionMultiURL(url, android, ios, pc string) CardElementActionOption {
	return func() (key string, v interface{}) {
		key, v = "multi_url", map[string]string{
			"url":         url,
			"android_url": android,
			"ios_url":     ios,
			"pc_url":      pc,
		}
		return
	}
}

type CardElementAction func() interface{}

func WithCardElementAction(elem CardElement, url string, opts ...CardElementActionOption) CardElementAction {
	ret := map[string]interface{}{
		"tag":  "button",
		"text": elem(true),
		"url":  url,
	}
	for _, fn := range opts {
		k, v := fn()
		ret[k] = v
	}

	return func() interface{} {
		return ret
	}
}

// WithCardElementActions 按钮, 可指定但固定跳转, 或多端跳转
func WithCardElementActions(act CardElementAction, actions ...CardElementAction) CardElement {
	actions = append([]CardElementAction{act}, actions...)
	as := make([]interface{}, 0, len(actions))
	for _, fn := range actions {
		as = append(as, fn())
	}

	return func(bool) interface{} {
		return map[string]interface{}{
			"tag":     "action",
			"actions": as,
		}
	}
}

// WithCardElementHorizontalRule 分割线
func WithCardElementHorizontalRule() CardElement {
	return func(bool) interface{} {
		return map[string]interface{}{
			"tag": "hr",
		}
	}
}

type CardElemImageOption func() (key string, v interface{})

// WithCardElementImageHover hover 图片时弹出的Tips文案
//  仅支持普通文本格式
func WithCardElementImageHover(text string) CardElemImageOption {
	return func() (key string, v interface{}) {
		key, v = "alt", map[string]interface{}{
			"tag":     "plain_text",
			"content": text,
		}
		return
	}
}

// WithCardElementImageTitle 图片的标题
//  默认普通文本格式
//  `md` 传入 `true`, 可支持 Markdown
func WithCardElementImageTitle(text string, md ...bool) CardElemImageOption {
	var isMD bool
	if len(md) != 0 && md[0] {
		isMD = md[0]
	}
	return func() (key string, v interface{}) {
		tag := "plain_text"
		if isMD {
			tag, text = "lark_md", trimPrefixSpace(text)
		}
		key, v = "title", map[string]interface{}{
			"tag":     tag,
			"content": text,
		}
		return
	}
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func trimPrefixSpace(s string) string {
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			return strings.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}
	return s[start:]
}

type ImageMode string

const (
	ImageModeCropCenter    ImageMode = "crop_center"
	ImageModeFitHorizontal ImageMode = "fit_horizontal"
)

// WithCardElementImageMode 图片显示模式
//  默认 居中裁剪模式
//  ImageModeCropCenter：居中裁剪模式，对长图会限高，并居中裁剪后展示
//  ImageModeFitHorizontal：平铺模式，宽度撑满卡片完整展示上传的图片。该属性会覆盖custom_width 属性
func WithCardElementImageMode(mode ImageMode) CardElemImageOption {
	return func() (key string, v interface{}) {
		key, v = "mode", string(mode)
		return
	}
}

// WithCardElementImageCustomWidth 自定义图片的最大展示宽度
//  默认展示宽度撑满卡片的通栏图片
//  可在 278px~580px 范围内指定最大展示宽度
//  在飞书4.0以上版本生效
func WithCardElementImageCustomWidth(w int) CardElemImageOption {
	min, max := 278, 580
	if w < min {
		w = min
	}
	if w > max {
		w = max
	}
	return func() (key string, v interface{}) {
		key, v = "custom_width", w
		return
	}
}

// WithCardElementImageCompactWidth 是否展示为紧凑型的图片
//  默认为 false
//  若配置为 true，则展示最大宽度为278px的紧凑型图片
func WithCardElementImageCompactWidth(b bool) CardElemImageOption {
	return func() (key string, v interface{}) {
		key, v = "compact_width", b
		return
	}
}

// WithCardElementImagePreview 点击后是否放大图片
//  缺省为true
//  在配置 card_link 后可设置为false，使用户点击卡片上的图片也能响应card_link链接跳转
func WithCardElementImagePreview(b bool) CardElemImageOption {
	return func() (key string, v interface{}) {
		key, v = "preview", b
		return
	}
}

func WithCardElementImage(imgKey string, opts ...CardElemImageOption) CardElement {
	// hover 默认为空，不展示
	opts = append([]CardElemImageOption{WithCardElementImageHover("")}, opts...)

	return func(bool) interface{} {
		elem := map[string]interface{}{
			"tag":     "img",
			"img_key": imgKey,
		}

		for _, fn := range opts {
			if k, v := fn(); v != nil {
				elem[k] = v
			}
		}

		return elem
	}
}

// WithCardElementNote 卡片的备注信息
//  支持元素:
//  - WithCardElementPlainText
//  - WithCardElementMarkdown
//  - WithCardElementImage
func WithCardElementNote(elem CardElement, elements ...CardElement) CardElement {
	elements = append([]CardElement{elem}, elements...)
	es := make([]interface{}, 0, len(elements))
	for _, fn := range elements {
		es = append(es, fn(true))
	}
	return func(bool) interface{} {
		return map[string]interface{}{
			"tag":      "note",
			"elements": es,
		}
	}
}
