package api

import (
	"math/rand"
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

//func makeRelative(src, dst string) (string, error) {
//	path := src
//	base := dst
//
//	if !filepath.IsAbs(path) || !filepath.IsAbs(base) {
//		return "", errors.New("paths must be absolute paths only")
//	}
//
//}

func symlinkOrRename(src, dst string) error {
	if info, err := os.Stat(dst); err == nil && info != nil {
		return nil
	}

	//relSrc := makeRelative(src, dst)

	err := os.Symlink(src, dst)
	if err != nil && !os.IsExist(err) {
		err = os.Rename(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}
