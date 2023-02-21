package fsBotAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/electricbubble/xhttpclient"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BotOption func(*bot)

func WithSecretKey(key string) BotOption {
	return func(b *bot) {
		b.secretKey = strings.TrimSpace(key)
	}
}

func WithDebugOutput() BotOption {
	return WithBodyCodec(&debugOutputBodyCodec{})
}

func WithBodyCodec(codec xhttpclient.BodyCodec) BotOption {
	return func(b *bot) {
		b.cli.WithBodyCodec(codec)
	}
}

var (
	_bcPoolDebugOutput = sync.Pool{
		New: func() any {
			return new(debugOutputBodyCodec)
		},
	}

	_ xhttpclient.BodyCodec             = (*debugOutputBodyCodec)(nil)
	_ xhttpclient.BodyHeaderContentType = (*debugOutputBodyCodec)(nil)
	_ xhttpclient.BodyCodecOnSend       = (*debugOutputBodyCodec)(nil)
	_ xhttpclient.BodyCodecOnReceive    = (*debugOutputBodyCodec)(nil)
)

type debugOutputBodyCodec struct {
	start time.Time
}

func (d *debugOutputBodyCodec) Get() xhttpclient.BodyCodec {
	return _bcPoolDebugOutput.Get().(*debugOutputBodyCodec)
}

func (d *debugOutputBodyCodec) Put(codec xhttpclient.BodyCodec) {
	_bcPoolDebugOutput.Put(codec)
}

func (d *debugOutputBodyCodec) Encode(body any) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}
	return buf, nil
}

func (d *debugOutputBodyCodec) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func (d *debugOutputBodyCodec) ContentType() string {
	return xhttpclient.ContentTypeValueJSON
}

func (d *debugOutputBodyCodec) OnSend(req *http.Request) {
	br, err := req.GetBody()
	if err != nil {
		log.Printf("[DEBUG-FeiShu-Bot-API] on send: get body: %s", err)
		return
	}
	defer func() {
		if err := br.Close(); err != nil {
			log.Printf("[DEBUG-FeiShu-Bot-API] on send: close body: %s", err)
		}
	}()

	bsBody, err := io.ReadAll(br)
	if err != nil {
		log.Printf("[DEBUG-FeiShu-Bot-API] on send: read body: %s", err)
		return
	}

	log.Printf("[DEBUG-FeiShu-Bot-API] %s",
		fmt.Sprintf("--> %s %s\n%s", req.Method, req.URL.String(), bsBody),
	)

	d.start = time.Now()
}

func (d *debugOutputBodyCodec) OnReceive(req *http.Request, resp *http.Response) {
	bsBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[DEBUG-FeiShu-Bot-API] on receive: read body: %s", err)
		return
	}
	defer func() {
		resp.Body = io.NopCloser(bytes.NewReader(bsBody))
	}()
	if err = resp.Body.Close(); err != nil {
		log.Printf("[DEBUG-FeiShu-Bot-API] on receive: close body: %s", err)
		return
	}

	log.Printf("[DEBUG-FeiShu-Bot-API] %s",
		fmt.Sprintf("<-- %s %s %d %s\n%s\n", req.Method, req.URL.String(), resp.StatusCode, time.Since(d.start), bsBody),
	)
}
