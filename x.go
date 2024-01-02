package feishu_bot_api

import "strings"

var _quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
