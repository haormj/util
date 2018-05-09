# Util [![Build Status](https://travis-ci.org/haormj/util.svg?branch=master)](https://travis-ci.org/haormj/util) [![GoDoc](https://godoc.org/github.com/haormj/util?status.svg)](https://godoc.org/github.com/haormj/util) [![Go Report Card](https://goreportcard.com/badge/github.com/haormj/util)](https://goreportcard.com/report/github.com/haormj/util)

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