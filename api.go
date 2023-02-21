package fsBotAPI

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const _fmtWebhook = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

type Bot interface {
	PushText(content string) error
	PushPost(p Post, ps ...Post) error
	PushCard(bgColor CardTitleBgColor, cfg CardConfig, c Card, more ...Card) error
	PushImage(imageKey string) error
	PushShareChat(chatID string) error
}

// 签名校验
//
//	https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#348211be
func genSign(secret string, timestamp int64) (string, error) {
	sign := fmt.Sprintf("%d\n%s", timestamp, secret)

	var data []byte
	h := hmac.New(sha256.New, []byte(sign))
	if _, err := h.Write(data); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

type Language string

const (
	LangChinese  Language = "zh_cn"
	LangEnglish  Language = "en_us"
	LangJapanese Language = "ja_jp"
)

// var _ fmt.Stringer = (*apiResponse)(nil)

type apiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// var HTTPClient = http.DefaultClient
//
// func newRequest(method string, rawUrl string, rawBody []byte) (request *http.Request, err error) {
// 	debugLog(fmt.Sprintf("--> %s %s\n%s", method, rawUrl, rawBody))
//
// 	if request, err = http.NewRequest(method, rawUrl, bytes.NewBuffer(rawBody)); err != nil {
// 		return nil, err
// 	}
// 	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
//
// 	return
// }
//
// func executeHTTP(req *http.Request) (rawResp []byte, err error) {
// 	start := time.Now()
// 	var resp *http.Response
// 	if resp, err = HTTPClient.Do(req); err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		_, _ = io.Copy(ioutil.Discard, resp.Body)
// 		_ = resp.Body.Close()
// 	}()
//
// 	rawResp, err = ioutil.ReadAll(resp.Body)
// 	debugLog(fmt.Sprintf("<-- %s %s %d %s\n%s\n", req.Method, req.URL.String(), resp.StatusCode, time.Since(start), rawResp))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var reply = new(struct {
// 		Code int    `json:"code"`
// 		Msg  string `json:"msg"`
// 	})
// 	if err = json.Unmarshal(rawResp, reply); err != nil {
// 		return nil, fmt.Errorf("unknown response: %w\nraw response: %s", err, rawResp)
// 	}
// 	if reply.Code != 0 {
// 		return nil, fmt.Errorf("%s (code: %d)", reply.Msg, reply.Code)
// 	}
//
// 	return
// }
//
// var debugFlag = false
//
// func SetDebug(debug bool) {
// 	debugFlag = debug
// }
//
// func debugLog(msg string) {
// 	if !debugFlag {
// 		return
// 	}
// 	log.Println("[DEBUG-FeiShu-Bot-API] " + msg)
// }
