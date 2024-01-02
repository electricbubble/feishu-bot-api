package feishu_bot_api

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_cardMessageHeaderI18nTexts_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		ts      cardHeaderI18nTexts
		want    []byte
		wantErr bool
	}{
		{
			name:    "case_1",
			ts:      nil,
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "case_2",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "zh_cn", text: "这是主标题！"},
			},
			want:    []byte(`{"zh_cn":"这是主标题！"}`),
			wantErr: false,
		},
		{
			name: "case_3",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "zh_cn", text: "这是主标题！"},
				cardHeaderI18nText{language: "en_us", text: "It is the title!"},
			},
			want:    []byte(`{"zh_cn":"这是主标题！","en_us":"It is the title!"}`),
			wantErr: false,
		},
		{
			name: "case_4",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "zh_cn", text: `这是"主标题"！`},
				cardHeaderI18nText{language: "en_us", text: "It is the title!"},
			},
			want:    []byte(`{"zh_cn":"这是\"主标题\"！","en_us":"It is the title!"}`),
			wantErr: false,
		},
		{
			name: "case_5",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "zh_cn", text: `\这是"主标题"！`},
				cardHeaderI18nText{language: "en_us", text: "It is the title!"},
			},
			want:    []byte(`{"zh_cn":"\\这是\"主标题\"！","en_us":"It is the title!"}`),
			wantErr: false,
		},
		{
			name: "case_6",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "zh_cn", text: `这是主标题！`},
				cardHeaderI18nText{language: "en_us", text: "It is the title!"},
				cardHeaderI18nText{language: "zh_cn", text: "这是多余的标题"},
			},
			want:    []byte(`{"zh_cn":"这是主标题！","en_us":"It is the title!"}`),
			wantErr: false,
		},
		{
			name: "case_7",
			ts: cardHeaderI18nTexts{
				cardHeaderI18nText{language: "en_us", text: "It is the title!"},
				cardHeaderI18nText{language: "zh_cn", text: `这是主标题！`},
				cardHeaderI18nText{language: "zh_cn", text: "这是多余的标题"},
			},
			want:    []byte(`{"en_us":"It is the title!","zh_cn":"这是主标题！"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ts.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON()\n got = %s\nwant = %s", got, tt.want)
			}
		})
	}
}

func Test_cardMessageHeaderI18nTextTags_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tts     cardHeaderI18nTextTags
		want    []byte
		wantErr bool
	}{
		{
			name:    "case_1",
			tts:     nil,
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "case_2",
			tts: cardHeaderI18nTextTags{
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容"},
					Color:    "carmine",
				},
			},
			want:    []byte(`{"zh_cn":[{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容"},"color":"carmine"}]}`),
			wantErr: false,
		},
		{
			name: "case_3",
			tts: cardHeaderI18nTextTags{
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容2"},
					Color:    "carmine",
				},
			},
			want:    []byte(`{"zh_cn":[{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容"},"color":"carmine"},{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容2"},"color":"carmine"}]}`),
			wantErr: false,
		},
		{
			name: "case_4",
			tts: cardHeaderI18nTextTags{
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容2"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "en_us",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "Tag content"},
					Color:    "carmine",
				},
			},
			want:    []byte(`{"zh_cn":[{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容"},"color":"carmine"},{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容2"},"color":"carmine"}],"en_us":[{"tag":"text_tag","text":{"tag":"plain_text","content":"Tag content"},"color":"carmine"}]}`),
			wantErr: false,
		},
		{
			name: "case_5",
			tts: cardHeaderI18nTextTags{
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "en_us",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "Tag content"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "zh_cn",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "标签内容2"},
					Color:    "carmine",
				},
				cardHeaderI18nTextTag{
					language: "en_us",
					Tag:      "text_tag",
					Text:     cardHeaderComponentPlainText{Tag: "plain_text", Content: "Tag content 2"},
					Color:    "",
				},
			},
			want:    []byte(`{"zh_cn":[{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容"},"color":"carmine"},{"tag":"text_tag","text":{"tag":"plain_text","content":"标签内容2"},"color":"carmine"}],"en_us":[{"tag":"text_tag","text":{"tag":"plain_text","content":"Tag content"},"color":"carmine"},{"tag":"text_tag","text":{"tag":"plain_text","content":"Tag content 2"}}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tts.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON()\n got = %s\nwant = %s", got, tt.want)
			}
		})
	}
}

func Test_cardMessage_buildHeader_output(t *testing.T) {
	m := cardMessage{
		globalConf: &CardGlobalConfig{
			headerIcon:     &cardHeaderIcon{ImgKey: "img_v2_881d0c9c-8717-49a7-b075-1cca6562443g"},
			headerTemplate: "red",
			config:         nil,
			link:           nil,
		},
		builders: []*CardBuilder{
			{
				language:              "zh_cn",
				headerTitleContent:    "中文环境下的主标题",
				headerSubtitleContent: "中文环境下的副标题",
				headerI18nTextTags: &(cardHeaderI18nTextTags{
					cardHeaderI18nTextTag{
						language: "zh_cn",
						Tag:      "text_tag",
						Text: cardHeaderComponentPlainText{
							Tag:     "plain_text",
							Content: "标题标签",
						},
						Color: "carmine",
					},
				}),
			},
			{
				language:              "en_us",
				headerTitleContent:    "英语环境下的主标题",
				headerSubtitleContent: "英语环境下的副标题",
				headerI18nTextTags: &(cardHeaderI18nTextTags{
					cardHeaderI18nTextTag{
						language: "en_us",
						Tag:      "text_tag",
						Text: cardHeaderComponentPlainText{
							Tag:     "plain_text",
							Content: "tagDemo",
						},
						Color: "carmine",
					},
				}),
			},
			nil,
			{
				language:              "ja_jp",
				headerTitleContent:    "日语环境下的主标题",
				headerSubtitleContent: "日语环境下的副标题",
				headerI18nTextTags:    nil,
			},
		},
	}

	bs, err := json.Marshal(m.buildHeader())
	requireNoError(t, err)

	t.Logf("\n%s", bs)

}
