package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/YutaKakiki/go-todo-api/testutil/fixture"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestEmbed(t *testing.T) {
	partOfWant := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, partOfWant) {
		t.Errorf("want %s,but got %s", partOfWant, rawPubKey)
	}

	partOfWant = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, partOfWant) {
		t.Errorf("want %s,but got %s", partOfWant, rawPrivKey)
	}
}

func TestGenerateToken(t *testing.T) {
	ctx := context.Background()
	wantID := entity.UserID(20)
	// テストデータUserを作成
	u := fixture.User(&entity.User{ID: wantID})
	// redisに保存する操作はスタブ
	// Saveメソッドの期待する動きを設定
	moq := &StoreMock{}
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sut.GenerateToken(ctx, u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}

func TestGetToken(t *testing.T) {
	t.Parallel()

	// 期待するjwtを発行
	c := clock.FixedClocker{}
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/YutaKakiki/go-todo-api`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := parse(rawPrivKey)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}

	// Loadメソッドの挙動をモックに。
	uid := entity.UserID(20)
	ctx := context.Background()
	mock := &StoreMock{}
	mock.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return uid, nil
	}
	sut, err := NewJWTer(mock, c)
	if err != nil {
		t.Fatal(err)
	}
	// リクエストのモック
	req := httptest.NewRequest(
		http.MethodGet,
		"https://github.com/YutaKakiki",
		nil,
	)
	// 発行されたjwtをヘッダーに付加
	req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
	// リクエストからjwtを取得
	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no error,but got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken got = %v,want = %v", got, want)
	}
}

// 明日の時刻に固定
type FixedTomorrowClocker struct{}

func (c FixedTomorrowClocker) Now() time.Time {
	// 固定している時刻に1時間追加
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

func TestGetToken_NG(t *testing.T) {
	t.Parallel()

	// 固定のタイムスタンプでjwtを発行
	c := clock.FixedClocker{}
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/YutaKakiki/go-todo-api`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := parse(rawPrivKey)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}
	// Load()モックが返す値の構造体
	type moq struct {
		entity.UserID
		err error
	}
	// 定義ずみのエラーErrNotFoundが出る/有効期限切れエラー出るパターンのテーブル
	tests := map[string]struct {
		c   clock.Clocker //時間いじる
		moq moq
	}{
		"expired": {
			// 期限切れ
			c: FixedTomorrowClocker{},
		},
		"notFoundInStore": {
			// 期限はクリアしてる
			c: clock.FixedClocker{},
			// redis上に見つからないエラー
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			moq := &StoreMock{}
			// モックで振る舞いを設定
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return tt.moq.UserID, tt.moq.err
			}
			sut, err := NewJWTer(moq, tt.c)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			req := httptest.NewRequest(
				http.MethodGet,
				`https://github.com/YutaKakiki/go-todo-api`,
				nil,
			)
			// 発行されたjwtをヘッダーに付加
			req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
			got, err := sut.GetToken(ctx, req)
			// エラーが吐かれるべき
			if err == nil {
				t.Error("want no err , but got nil")
			}
			//jwtは取得できないはず
			if got != nil {
				t.Errorf("want nil , but got %v", got)
			}

		})

	}

}
