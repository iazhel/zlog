package zlog

import (
	"strings"
	"testing"
)

func BenchmarkFillString(b *testing.B) {
	// NOOOOOOO
	out := ""
	logs := []string{}
	msg := "123456789012345678901234567890123456789012345678901234567890"
	for i := 0; i < 20; i++ {
		logs = append(logs, msg)
	}
	b.N = 200
	for n := 0; n < b.N; n++ {
		out += strings.Join(logs, "\n")
	}
	b.Log("len:", len(out))
}
func BenchmarkFillSlice(b *testing.B) {
	// YESSSSSSSSSSS
	out := []string{}
	logs := []string{}
	msg := "123456789012345678901234567890123456789012345678901234567890"
	for i := 0; i < 20; i++ {
		logs = append(logs, msg)
	}
	b.N = 200
	for n := 0; n < b.N; n++ {
		out = append(out, logs...)
	}
	out1 := strings.Join(out, "\n")
	b.Log("len:", len(out1))
}
