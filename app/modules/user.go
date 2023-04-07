package modules

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// ユーザー情報
type User struct {
	Id         int64
	Name       string
	Email      string
	Password   string
	Created_at time.Time
	Updated_at time.Time
}

// ユーザーが登録画面で設定するパラメーター
type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

func Encrypt(char string) string {
	encryptText := fmt.Sprintf("%x", sha256.Sum256([]byte(char)))
	return encryptText
}
