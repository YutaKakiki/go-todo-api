package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/YutaKakiki/go-todo-api/entity"
)

func User(u *entity.User) *entity.User {
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),
		Name:     strconv.Itoa(rand.Int())[:5],
		Password: "password",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}

	// ユーザーフィールドの値が設定されていれば、そちらの情報を優先させて更新する

	if u == nil {
		return result
	}
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.Created.IsZero() {
		result.Created = u.Created
	}
	if !u.Modified.IsZero() {
		result.Modified = u.Modified
	}
	return result
}
