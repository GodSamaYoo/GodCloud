package main

import (
	"crypto/md5"
	"encoding/hex"
)


func md5_(s string) string {
	m := md5.New()
	m.Write([]byte (s))
	return hex.EncodeToString(m.Sum(nil))[8:24]
}
