package kdf

import (
	"crypto/sha512"
	"fmt"
	"hash"
)

const (
	ipad = 0x36
	opad = 0x5c
)

const sizeHash = 64

type kdf struct {
}

func New() kdf {
	return kdf{}
}

func (k kdf) KDF(S, T []byte) ([]byte, error) {
	var err error

	h := sha512.New()

	tempKey, err := k1(h, S, T)
	if err != nil {
		return nil, fmt.Errorf("k1: %v", err)
	}

	result, err := k2(h, tempKey, nil)
	if err != nil {
		return nil, fmt.Errorf("k2: %v", err)
	}

	defer func() { // очистка ключевой информации
		tempKey = nil
		result = nil
		S = nil
		T = nil
	}()

	return result, nil
}

func k1(h hash.Hash, S, T []byte) ([]byte, error) {
	res, err := hmac(h, S, T)
	if err != nil {
		return nil, err
	}
	resLB := lb(res)
	return resLB, nil
}

func k2(h hash.Hash, K, info []byte) ([]byte, error) {
	var err error
	z := make([]byte, sizeHash)
	count := numIteration(len(K), sizeHash)

	finZ := make([]byte, count*sizeHash)
	c := 0

	for i := 0; i < count; i++ {
		z, err = hmac(h, format(z, c), K)
		if err != nil {
			return nil, err
		}
		c++
		addSlice(finZ, z, i*sizeHash)
	}
	return lb(finZ), nil
}

func lb(data []byte) []byte {
	return data[:32]
}

func hmac(h hash.Hash, S, T []byte) ([]byte, error) {
	hash, err := internalHmac(h, S, T)
	if err != nil {
		return nil, err
	}
	res, err := sha(h, hash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func xorBlocks(left []byte, right byte) []byte {
	out := make([]byte, sizeHash)
	for i := range left {
		out[i] = left[i] ^ right
	}
	return out
}

func internalHmac(h hash.Hash, S, T []byte) ([]byte, error) {
	internalS := make([]byte, sizeHash+sizeHash+len(S))
	newT := make([]byte, sizeHash)

	addSlice(newT, T, 0) // расширение T до 512 бит
	T = newT

	addSlice(internalS, xorBlocks(T, opad), 0)
	internalHash, err := sha(h, xorBlocks(T, ipad))
	if err != nil {
		return nil, err
	}
	addSlice(internalS, internalHash, sizeHash)
	addSlice(internalS, S, len(T)+sizeHash)

	return internalS, nil
}

func addSlice(receiver, source []byte, position int) {
	for i := range source {
		receiver[position+i] = source[i]
	}
}

func sha(h hash.Hash, data []byte) ([]byte, error) {
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func numIteration(len, n int) int {
	return len * 8 / n
}

func format(z []byte, c int) []byte {
	newZ := make([]byte, len(z)+1)
	copy(newZ, z)
	newZ[len(z)] = byte(c)
	return newZ
}
