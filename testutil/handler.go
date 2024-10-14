package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// JSONの期待値・実際の差分を
func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any

	// JSONを構造体にしたものをjw変数に格納
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", got, err)
	}

	if diff := cmp.Diff(jg, jw); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

func AsertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	// 実際のレスポンスを読み取る
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	// ステータスコードを検証
	if got.StatusCode != status {
		t.Fatalf("want status %d,but got %d,body:%q", status, got.StatusCode, gb)
	}
	if len(gb) == 0 && len(body) == 0 {
		// ボディがなければ、JSON比較もできない
		return
	}
	AssertJSON(t, body, gb)
}

// ファイルから入力値、期待値を取得
// のちのゴールデンテストで使用
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
