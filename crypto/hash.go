package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/goinggo/tracelog"
)

type (
	SecureEntity interface {
		TokenBytes() ([]byte, error)
	}
)

const key = "~*&^*lnkldfnljdf&^*)%*^%*^%ksjdlj3984mn38JL:8k3km}[71&$@@%^*())km2pumsdjgmHJGSL:JHFdkjfj(&%#$22)=="

//Generates a SHA256 Token from a String.
func Hash(value string) string {

	data := []byte(value)
	hash := sha256.New()
	hash.Write(data)
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

//Generates a Signed SHA256 Hash and then encodes it to Base64
func SignedEncodedHash(message string, keyString string) string {
	return base64.StdEncoding.EncodeToString(SignedHash(message, keyString))
}

//Generates a Signed SHA256 Hash
func SignedHash(message string, keyToSign string) []byte {
	if keyToSign == "" {
		tracelog.TRACE("go-common/crypto", "SignedHash", "KeyString is Blank")
		keyToSign = key
	}

	h := hmac.New(sha256.New, []byte(keyToSign))
	h.Write([]byte(message))
	return h.Sum(nil)
}

//Checks whether a token is valid for a Secure Entity.
func IsTokenValid(secureEntity SecureEntity, token string) error {
	tracelog.STARTED("Utils", "IsValidToken")

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		tracelog.ERRORf(err, "Utils", "Utils.IsValidToken", "Error Decoding Passed In Token, %s", token)
		return err
	}

	entityToken, tErr := secureEntity.TokenBytes()
	if tErr != nil {
		tracelog.ERRORf(tErr, "Utils", "Utils.IsValidToken", "Error Generating Token for Entity")
		return tErr
	}

	if hmac.Equal(decodedToken, entityToken) == false {
		tracelog.ERRORf(err, "Utils", "Utils.IsValidToken", "Invalid Token Comparison,Tokens Are not the same, Invalid Token, entity[%s], decoded[%s]", string(entityToken), string(decodedToken))
		return errors.New("Invalid Token")
	}

	tracelog.COMPLETED("Utils", "IsValidToken, Token Is Valid")

	return nil

}
