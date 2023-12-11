// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:build !windows

package api

import (
	"os"
	"path/filepath"
)

const useANSICodes = true

func symlinkOrRename(src, dst string) error {
	if info, err := os.Stat(dst); err == nil && info != nil {
		return nil
	}

	absSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	err = os.Symlink(absSrc, absDst)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
