package implementations

import (
	"encoding/hex"
	"reson8-learning-api/util"

	"github.com/zeebo/blake3"
)

type Blake3HashingManager struct {
	hasher *blake3.Hasher
}

func NewBlake3HashingManager(configuration *util.Configuration) (*Blake3HashingManager, error) {
	key := make([]byte, 32)
	blake3.DeriveKey(configuration.GetSecretKey(), []byte(configuration.Salt), key)
	hasher, err := blake3.NewKeyed(key)
	if err != nil {
		return nil, err
	}
	return &Blake3HashingManager{hasher: hasher}, nil
}

func (h *Blake3HashingManager) GetHash(s string) (string, error) {
	count, err := h.hasher.Write([]byte(s))
	if err != nil || count == 0 {
		return "", err
	}
	res := hex.EncodeToString(h.hasher.Sum(nil))
	h.hasher.Reset()
	return res, nil
}

func (h *Blake3HashingManager) Verify(hashed string, new string) (bool, error) {
	res, err := h.GetHash(new)
	if err != nil {
		return false, err
	}
	return res == hashed, nil
}
