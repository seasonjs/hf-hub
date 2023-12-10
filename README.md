# hf-hub

[![Go Reference](https://pkg.go.dev/badge/github.com/seasonjs/hf-hub.svg)](https://pkg.go.dev/github.com/seasonjs/hf-hub)

golang client for the huggingface hub aiming for minimal subset of features over `huggingface-hub` python package

## Usage

Add the dependency

```bash
go get github.com/seasonjs/hf-hub
```

use this package

```golang
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
```

## Thanks

* [huggingface_hub](https://github.com/huggingface/huggingface_hub)
* [hf-hub](https://github.com/huggingface/hf-hub)

## License

Copyright (c) seasonjs. All rights reserved.
Licensed under the MIT License. See License.txt in the project root for license information.
