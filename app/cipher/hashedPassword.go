package cipher

import (
	"crypto/sha256"
	"encoding/hex"
)

// 文字列のハッシュ化
func HashStr(str string) string {

	//byteのスライスに変換
	b := []byte(str)
	//sha256を使用してパスワードをハッシュ化する
	sha256 := sha256.Sum256(b)
	hashed := hex.EncodeToString(sha256[:])

	return hashed
}
