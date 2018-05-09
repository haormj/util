# Util

Golang util tools

## Install

```shell
go get github.com/haormj/util
```

## Usage

```go
package main

import (
	"log"

	"github.com/haormj/util"
)

func main() {
	dst, err := util.ByteToUint32([]byte{1, 2, 3, 4})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(dst)
}
```