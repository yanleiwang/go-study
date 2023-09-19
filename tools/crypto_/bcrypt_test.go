package crypto_

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestBcrypto(t *testing.T) {
	/*
		bcrypt 是只能加密， 比较， 无法解密的
		他把随机生成的盐值， 存储在加密串中， 所以无需在数据库存储盐值。
	*/

	pwd := "123456"
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	require.NoError(t, err)

	// 比较
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(pwd))
	require.NoError(t, err)
}
