package fsBotAPI

type i18nPost struct {
	lang     string        // 富文本的语言环境
	title    string        // 富文本标题
	elements []interface{} // 段落的所有元素
}

type Post func() i18nPost

// WithPost 富文本消息, 可指定语言环境
//  支持元素如下:
//  普通文本: WithPostElementText
//  文字超链接: WithPostElementLink
//  图片: WithPostElementImage
//  @所有人: WithPostElementMentionAll
//  @指定用户(OpenID): WithPostElementMentionByOpenID
func WithPost(lang Language, title string, elements ...PostElement) Post {
	return func() i18nPost {
		es := make([]interface{}, 0, len(elements))
		p := make([]interface{}, 0, len(elements))
		for _, fn := range elements {
			elem := fn()
			if elem.isImage {
				// 图片元素必须是独立的一个段落
				es = append(es, p, []interface{}{elem.elem})
				p = make([]interface{}, 0, len(elements))
			} else {
				p = append(p, elem.elem)
			}
		}
		es = append(es, p)
		return i18nPost{
			lang:     string(lang),
			title:    title,
			elements: es,
		}
	}
}

type postElement struct {
	elem    interface{}
	isImage bool // 图片元素必须是独立的一个段落
}

type PostElement func() postElement

// WithPostElementText 富文本消息的文字元素
//  isUnescape 表示是不是 unescape 解码，默认为 false ，不用可以不填
func WithPostElementText(text string, isUnescape ...bool) PostElement {
	return func() postElement {
		elem := map[string]interface{}{
			"tag":  "text",
			"text": text,
		}
		if len(isUnescape) != 0 && isUnescape[0] {
			elem["un_escape"] = isUnescape[0]
		}
		return postElement{
			elem:    elem,
			isImage: false,
		}
	}
}

// WithPostElementLink 富文本消息的文字超链接元素
func WithPostElementLink(text, href string) PostElement {
	return func() postElement {
		return postElement{
			elem: map[string]interface{}{
				"tag":  "a",
				"text": text,
				"href": href,
			},
			isImage: false,
		}
	}
}

// WithPostElementImage 富文本消息的图片元素
func WithPostElementImage(imgKey string) PostElement {
	return func() postElement {
		return postElement{
			elem: map[string]interface{}{
				"tag":       "img",
				"image_key": imgKey,
			},
			isImage: true,
		}
	}
}

// WithPostElementMentionAll 富文本消息的 @所有人
func WithPostElementMentionAll() PostElement {
	return func() postElement {
		return postElement{
			elem: map[string]interface{}{
				"tag":     "at",
				"user_id": "all",
			},
			isImage: false,
		}
	}
}

// WithPostElementMentionByOpenID 富文本消息的 @用户
//  Open ID 必须是有效值，否则仅显示 `@` 符号（实际效果不同于 PushText 时会显示 name）
func WithPostElementMentionByOpenID(id string, name ...string) PostElement {
	return func() postElement {
		elem := map[string]interface{}{
			"tag":     "at",
			"user_id": id,
		}
		if len(name) != 0 {
			elem["user_name"] = name[0]
		}
		return postElement{
			elem:    elem,
			isImage: false,
		}
	}
}

func newMsgPost(p Post, more ...Post) map[string]interface{} {
	more = append([]Post{p}, more...)
	i18nPosts := make([]i18nPost, 0, len(more))
	for _, fn := range more {
		i18nPosts = append(i18nPosts, fn())
	}
	post := make(map[string]interface{}, 3)
	for _, p := range i18nPosts {
		post[p.lang] = map[string]interface{}{
			"title":   p.title,
			"content": p.elements,
		}
	}

	msgType := "post"
	return map[string]interface{}{
		"msg_type": msgType,
		"content": map[string]interface{}{
			msgType: post,
		},
	}
}
