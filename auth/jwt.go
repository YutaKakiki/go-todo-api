package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	// redisのことだろうけど、依存しないためにここでインターフェースに依存するようにしている
	Store   Store
	Clocker clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

// Storeは、コンストラクタインジェクションによってDI
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privatekey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJTWer: private key :%w", err)
	}
	publickey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJTWer: private key :%w", err)
	}
	j.PrivateKey = privatekey
	j.PublicKey = publickey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}

// カスタムクレームの値
const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

// JWTを発行する
func (j *JWTer) GenerateToken(ctx context.Context, u *entity.User) ([]byte, error) {
	uuid := uuid.New().String()
	// ヘッダーはデフォルトで作成される
	// ペイロードを作成
	tok, err := jwt.NewBuilder().
		JwtID(uuid).
		Issuer("github.com/YutaKakiki/go-todo-api").
		Subject("access_toketn").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role). //カスタムクレーム：payloadの中に埋め込む情報
		Claim(UserNameKey, u.Name).
		Build() //トークンを作成
	if err != nil {
		return nil, fmt.Errorf("GenerateToken: failed to build token: %w", err)
	}
	// トークンをredisに保存
	// トークンのUUID : ユーザーID
	err = j.Store.Save(ctx, tok.JwtID(), u.ID)
	if err != nil {
		return nil, err
	}
	// 秘密鍵を使用し、指定のアルゴリズムで署名を行う
	sined, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return sined, nil
}

// 公開鍵で署名を検証する
func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	// 独自に検証を実装するのでfalse
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}
	// jwtの有効期限を検証
	err = jwt.Validate(token, jwt.WithClock(j.Clocker))
	if err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token :%w", err)
	}
	// redisからデータが削除されていた場合
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired :%w", token.JwtID(), err)
	}
	return token, nil
}

// JWTから取得したデータをリクエストスコープのコンテクストに書き込んでクローン
// ここで検証も行なっている
func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	// リクエストからjwtを取得（リクエストのヘッダから取得）
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	// jwt のID（uuid）をキーとしてredis上にあるデータを検索し、ユーザーIDを取得
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	// リクエストスコープからコンテクストを取得
	ctx := r.Context()
	// JWTから取得したUserIDとRoleKeyをコンテキストに追加
	ctx = SetUserID(ctx, uid)
	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

// コンテキストに含めるキーの型を定義
type userIDKey struct{}

// ユーザーIDをcontextにセット
func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	// any型を型アサーションで型変換
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

type roleKey struct{}

// roleをcontextにセット
func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey) //RokeKeyは、jwtのペイロードに設定したカスタムクレーム
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}

	return context.WithValue(ctx, roleKey{}, get)

}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

// 管理者権限の有無を検証
func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
