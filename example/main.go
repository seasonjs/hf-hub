// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"github.com/seasonjs/hf-hub/api"
)

func main() {
	hapi, err := api.NewApi()
	if err != nil {
		print(err.Error())
		return
	}

	modelPath, err := hapi.Model("bert-base-uncased").Get("config.json")
	if err != nil {
		print(err.Error())
		return
	}

	print(modelPath)
}
