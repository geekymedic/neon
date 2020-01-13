package tool

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"sync"

	"github.com/google/uuid"
	"github.com/zentures/cityhash"
)

var syncPool sync.Pool

func init() {
	syncPool.New = func() interface{} {
		return cityhash.New64()
	}
}

func MD5UUID4() string {
	hasher := md5.New()
	txt := uuid.New()
	hasher.Write(txt[:])
	return hex.EncodeToString(hasher.Sum(nil))
}

func RandomUint64(input ...[]byte) uint64 {
	input = append(input, []byte(MD5UUID4()))
	hash64 := syncPool.Get().(hash.Hash64)
	for _, b := range input {
		hash64.Write(b)
	}
	id := hash64.Sum64()
	hash64.Reset()
	return id
}

func ConvertBytesToUint64(input ...[]byte) uint64 {
	input = append(input, []byte(MD5UUID4()))
	hash64 := syncPool.Get().(hash.Hash64)
	id := hash64.Sum64()
	hash64.Reset()
	return id
}
