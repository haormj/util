package util

import (
	"log"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	log.Println(RandStringRunes(10))
}

func BenchmarkRandStringRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringRunes(60)
	}
}

func TestRandStringBytes(t *testing.T) {
	log.Println(RandStringBytes(10))
}

func BenchmarkRandStringBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytes(60)
	}
}

func TestRandStringBytesRmndr(t *testing.T) {
	log.Println(RandStringBytesRmndr(10))
}

func BenchmarkRandStringBytesRmndr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesRmndr(60)
	}
}

func TestRandStringBytesMask(t *testing.T) {
	log.Println(RandStringBytesMask(10))
}

func BenchmarkRandStringBytesMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(60)
	}
}

func TestRandStringBytesMaskImpr(t *testing.T) {
	log.Println(RandStringBytesMaskImpr(10))
}

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImpr(60)
	}
}

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	log.Println(RandStringBytesMaskImprSrc(10))
}

func BenchmarkRandStringBytesMaskImprSrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImpr(60)
	}
}

func TestRandomString(t *testing.T) {
	log.Println(RandomString("ab", 10))
}

func TestRandomStringByLetterAndDigital(t *testing.T) {
	log.Println(RandomStringByLetterAndDigital(43))
}

func TestRandomStringByDigital(t *testing.T) {
	log.Println(RandomStringByDigital(6))
}