package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type extractedArgonParams struct {
	params       *argonParams
	salt         []byte
	passwordHash []byte
}

func genRandomSalt(size uint32) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func CreateHash(str string) ([]byte, error) {
	params := &argonParams{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := genRandomSalt(params.saltLength)
	if err != nil {
		return nil, err
	}
	hashStr := argon2.IDKey([]byte(str), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashStrBase64 := base64.RawStdEncoding.EncodeToString(hashStr)
	argonHash := []byte(fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.memory, params.iterations, params.parallelism, saltBase64, hashStrBase64))
	return argonHash, nil
}

func extractArgonParams(hashData []byte) (*extractedArgonParams, error) {
	tokens := strings.Split(string(hashData), "$")
	if len(tokens) != 6 {
		return nil, errors.New("invalid hash")
	}
	var p argonParams
	var argonVersion int
	var ep extractedArgonParams
	_, err := fmt.Sscanf(string(hashData), "$argon2id$v=%d$m=%d,t=%d,p=%d$", &argonVersion, &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, err
	}
	salt, err := base64.RawStdEncoding.Strict().DecodeString(tokens[4])
	if err != nil {
		return nil, err
	}
	passwordHash, err := base64.RawStdEncoding.Strict().DecodeString(tokens[5])
	if err != nil {
		return nil, err
	}
	p.keyLength = uint32(len(passwordHash))
	ep.salt = salt
	ep.passwordHash = passwordHash
	ep.params = &p
	return &ep, nil
}

func VerifyHash(password string, hashData []byte) (bool, error) {
	ep, err := extractArgonParams(hashData)
	if err != nil {
		return false, err
	}
	regeneratedHash := argon2.IDKey([]byte(password), ep.salt, ep.params.iterations, ep.params.memory, ep.params.parallelism, ep.params.keyLength)
	if subtle.ConstantTimeCompare(ep.passwordHash, regeneratedHash) == 1 {
		return true, nil
	}
	return false, nil
}
