package util

import (
	"strings"

	"github.com/satori/go.uuid"
)

func GetUUIDV1() string {
	return uuid.NewV1().String()
}

func GetUUIDV1WithoutLine() string {
	return strings.Replace(GetUUIDV1(), "-", "", -1)
}
