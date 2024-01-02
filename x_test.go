package feishu_bot_api

import (
	"testing"
)

func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Received unexpected error:\n%+v", err)
	}
}
