package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Hasher struct {
	salt string
}

func NewHasher(cfg Config) *Hasher {
	return &Hasher{salt: cfg.Salt}
}

func (h *Hasher) Hash(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
