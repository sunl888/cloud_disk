package go_file_uploader

import (
	"io"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"errors"
)

type HashFunc func(file io.ReadSeeker) (string, error)

func (hf HashFunc) Hash(file io.ReadSeeker) (string, error) {
	return hf(file)
}

type Hasher interface {
	Hash(file io.ReadSeeker) (string, error)
}

func MD5HashFunc(file io.ReadSeeker) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type Hash2StorageName interface {
	Convent(hash string) (storageName string, err error)
}

type Hash2StorageNameFunc func(hash string) (storageName string, err error)

func (f Hash2StorageNameFunc) Convent(hash string) (storageName string, err error) {
	return f(hash)
}

func DefaultHash2StorageNameFunc(hash string) (storageName string, err error) {
	return hash, nil
}

func TwoCharsPrefixHash2StorageNameFunc(hash string) (storageName string, err error) {
	if len(hash) <= 2 {
		return "", errors.New("hash length must greater than 2 chars")
	}
	sb := strings.Builder{}
	sb.WriteString(hash[:2])
	sb.WriteString("/")
	sb.WriteString(hash[2:])
	return sb.String(), nil
}
