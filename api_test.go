package fsBotAPI

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func Test_newRequest(t *testing.T) {
	botKey := os.Getenv("BOT_KEY")
	t.Log(botKey)

	webhook := fmt.Sprintf(_fmtWebhook, botKey)

	msg := []byte(`{
    "msg_type": "text",
    "content": {
        "text": "新更新提醒第二行"
    }
}`)

	msg = []byte(`{
    "msg_type": "image",
    "content": {
        "image_key": "img_ecffc3b9-8f14-400f-a014-05eca1a4310g"
    }
}`)
	// img_7ea74629-9191-4176-998c-2e603c9c5e8g
	// img_ecffc3b9-8f14-400f-a014-05eca1a4310g

	msg = []byte(`{
    "msg_type": "share_chat",
    "content": {
        "share_chat_id": "oc_f5b1a7eb27ae2c7b6adc2a74faf339ff"
    }
}`)

	msg = []byte(`{
    "msg_type": "audio",
    "content": {
        "file_key": "75235e0c-4f92-430a-a99b-8446610223cg"
    }
}`)

	msg = []byte(`{
	"msg_type": "interactive",
	"card": {
		"config": {
			"wide_screen_mode": true,
			"enable_forward": true
		},
		"elements": [{
			"tag": "div",
			"text": {
				"i18n": {"zh_cn":"zh_cn **西湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。"},
				"tag": "lark_md"
			}
		}, {
			"actions": [{
				"tag": "button",
				"text": {
					"content": "更多景点介绍 :玫瑰:",
					"tag": "lark_md"
				},
				"url": "https://www.example.com",
				"type": "default",
				"value": {}
			}],
			"tag": "action"
		}],
		"header": {
			"title": {
				"i18n": {"zh_cn":"zh_cn今日旅游推荐"},
				"tag": "plain_text"
			}
		}
	}
}`)
	msg = []byte(`{
	"msg_type": "interactive",
	"card": {
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "blue",
    "title": {
      "i18n": {
        "en_us": "What's New on Feishu Open Platform",
        "ja_jp": "What's New on Feishu Open Platform",
        "zh_cn": "飞书开放平台近期重要更新"
      },
      "tag": "plain_text"
    }
  },
  "i18n_elements": {
    "en_us": [
      {
        "tag": "div",
        "text": {
          "content": "Dear Developers,",
          "tag": "lark_md"
        }
      },
      {
        "tag": "div",
        "text": {
          "content": "In order to keep you updated about what recently happened here on Feishu Open Platform, we have selected some of the important updates in the past week.",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "tag": "div",
        "text": {
          "content": "🌟 **Testing companies and users**\nCustom apps now support [test companies and users](https://open.feishu.cn/document/home/introduction-to-custom-app-development/testing-enterprise-and-personnel-functions?lang=en-US&from=630bot) function. After associating your app with the test company in the Developer Console, you can create a test environment and start debugging, which requires no administrator review. After debugging, you can return to the formal environment to apply for release, which greatly improves the efficiency of app development.",
          "tag": "lark_md"
        }
      },
      {
        "tag": "div",
        "text": {
          "content": "🌟 **Gadget and web app** \nIn daily work, you probably need to switch back and forth between app and chat window. Now with [toggleChat API](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/toggleChat?from=630bot&lang=en-US), Gadget V4.1.0 supports opening private chat or group chat in the sidebar of the app on PC, giving you a smoother experience.",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "tag": "div",
        "text": {
          "content": "**Get started: Developer tutorials**",
          "tag": "lark_md"
        }
      },
      {
        "extra": {
          "alt": {
            "content": "Image",
            "tag": "plain_text"
          },
          "img_key": "img_v2_7441ffc8-ff10-4ac2-869f-576eadd14e2g",
          "tag": "img"
        },
        "tag": "div",
        "text": {
          "content": "**🌟 Intro to user IDs** - help you understand the concept and design logic of different Feishu user IDs \n[**Check the tutorial here >>**](https://open.feishu.cn/document/home/user-identity-introduction/introduction?from=630bot&lang=en-US)",
          "tag": "lark_md"
        }
      },
      {
        "extra": {
          "alt": {
            "content": "Image",
            "tag": "plain_text"
          },
          "img_key": "img_v2_770b3ba0-26e4-4b73-a123-2bbb9850f3dg",
          "tag": "img"
        },
        "tag": "div",
        "text": {
          "content": "🌟 **Quick intro to message card** - brings more vitality into bots in terms of content presentation and user interaction\n[**Check the tutorial here >>**](https://open.feishu.cn/document/home/build-a-beautiful-message-card-in-5-minutes/what-is-a-message-card?from=630bot&lang=en-US)",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "actions": [
          {
            "multi_url": {
              "android_url": "https://open.feishu.cn/changelog?lang=en-US&from=630bot",
              "ios_url": "https://open.feishu.cn/changelog?lang=en-US&from=630bot",
              "pc_url": "https://open.feishu.cn/changelog?lang=en-US&from=630bot",
              "url": "https://open.feishu.cn/changelog?lang=en-US&from=630bot"
            },
            "tag": "button",
            "text": {
              "content": "Learn More",
              "tag": "lark_md"
            },
            "type": "primary"
          }
        ],
        "tag": "action"
      }
    ],
    "zh_cn": [
      {
        "tag": "div",
        "text": {
          "content": "亲爱的飞书开发者：为了能让你及时了解开放平台的新功能，修复的 bug 及开发文档的变更，现向你推送最近一周飞书开放平台的精选动态。",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "tag": "div",
        "text": {
          "content": "🌟 **测试企业与人员**\n自建应用支持[「测试企业与人员」](https://open.feishu.cn/document/home/introduction-to-custom-app-development/testing-enterprise-and-personnel-functions?from=630bot&lang=zh-CN)功能。 在开发者后台，将开发中的应用与测试企业关联后，即可在测试环境进行调试，无需管理员审核。调试完成后，再回到正式环境申请发布，大大提升研发效率。",
          "tag": "lark_md"
        }
      },
      {
        "tag": "div",
        "text": {
          "content": "🌟 **小程序与网页应用** \n在日常工作中，经常会遇到在应用和聊天窗口来回切换的情景，十分繁琐。小程序与网页应用 V4.1.0 新增 [toggleChat API](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/toggleChat?from=630bot&lang=zh-CN)，支持 PC 端在应用中以侧边栏形式打开用户或群组会话，让体验更流畅。",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "tag": "div",
        "text": {
          "content": "**优质飞书应用开发教程分享**",
          "tag": "lark_md"
        }
      },
      {
        "extra": {
          "alt": {
            "content": "",
            "tag": "plain_text"
          },
          "img_key": "img_v2_ed4810d8-6697-465b-9d74-42464477b31g",
          "tag": "img"
        },
        "tag": "div",
        "text": {
          "content": "🌟 **用户身份体系介绍**：真人讲解视频，带你快速了解飞书各个用户 ID 的概念和背后的设计逻辑。\n[点击查看>>](https://open.feishu.cn/document/home/user-identity-introduction/introduction?from=630bot&lang=zh-CN)",
          "tag": "lark_md"
        }
      },
      {
        "extra": {
          "alt": {
            "content": "",
            "tag": "plain_text"
          },
          "img_key": "img_v2_21bc4b9d-f14a-44e6-98f8-575c01e0975g",
          "tag": "img"
        },
        "tag": "div",
        "text": {
          "content": "🌟 **快速了解消息卡片**：让你的机器人在内容呈现、用户交互上更有生命力。\n[点击查看>>](https://open.feishu.cn/document/home/build-a-beautiful-message-card-in-5-minutes/what-is-a-message-card?from=630bot&lang=zh-CN)",
          "tag": "lark_md"
        }
      },
      {
        "tag": "hr"
      },
      {
        "actions": [
          {
            "tag": "button",
            "text": {
              "content": "🔎 查看更新详情",
              "tag": "plain_text"
            },
            "type": "primary",
            "url": "https://open.feishu.cn/changelog?from=630bot&lang=zh-CN"
          }
        ],
        "tag": "action"
      }
    ]
  }
}
}`)

	msg = []byte(`{
  "card": {
"card_link": { 
        "url": "https://www.baidu.com",
        "android_url": "https://developer.android.com/",
        "ios_url": "https://developer.apple.com/",
        "pc_url": "https://www.windows.com"
    },
    "header": {
      "title": {
        "content": " ",
        "tag": "plain_text"
      },
      "template": "orange"
    },
    "config": {
      "enable_forward": true,
      "update_multi": false
    },
    "elements": [
      {
        "tag": "div",
        "text": {
          "content": "abcd~!@#$%^\u0026*()",
          "tag": "lark_md"
        }
      }
    ]
  },
  "msg_type": "interactive"
}`)

	req, err := newRequest(http.MethodPost, webhook, msg)
	requireNil(t, err)

	bsResp, err := executeHTTP(req)
	requireNil(t, err)

	// {"Extra":null,"StatusCode":0,"StatusMessage":"success"}
	t.Log(string(bsResp))
}
