// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func randStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)

	var result []byte
	rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(100))))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func shouldRedirect(statusCode int) bool {
	redirectStatusCodes := http.StatusMovedPermanently |
		http.StatusFound |
		http.StatusSeeOther |
		http.StatusTemporaryRedirect |
		http.StatusPermanentRedirect
	return (redirectStatusCodes & statusCode) != 0
}

func GetSHA256FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
