package zlog

import (
	"fmt"
	"testing"
	//	"strings"
)

func Benchmark_Info(b *testing.B) {

	message := "Ready bytes: %07d................................"
	size := 400 //Size of Step, kB
	sizeMessage := len(message) + len(prefixInfo) + 7
	m := int(size * 1024 / sizeMessage) // Number of Info logs.
	counter := 0

	zLog := NewZLog()
	zLog.Step("Making Step with size %d Kb", size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < m; j++ {
			counter += sizeMessage
			zLog.Info(message, counter)
		}
	}

	fmt.Printf("\nMade %d messages, per lenght: %d bytes, total: %d Kbytes", m*b.N, sizeMessage, int(counter/(1024)))

	_ = zLog.GetLog()
	//n	fmt.Printf("\nLast message in Getlog:  %v\n", logs[len(logs)-1])

}

func Benchmark_Step_Warning(b *testing.B) {

	size := 400 //Size of Step, kB
	message := "Event bytes: %07d................................"
	sizeMessage := len(message) + len(prefixInfo) + 7
	m := int(size * 1024 / sizeMessage)
	counter := 0

	zLog := NewZLog()
	zLog.Step("Making Step(block) with size %d Kb and Warning", size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < m; j++ {
			counter += sizeMessage
			zLog.Info(message, counter)
		}
	}

	zLog.Warning("")
	fmt.Printf("\nMade %d messages with WARNING, per lenght: %d bytes, total: %d Kbytes", m*b.N, sizeMessage, int(counter/(1024)))

	_ = zLog.GetLog()
	//	fmt.Printf("\nLast message in Getlog:  %v\n", logs[len(logs)-1])

}

func Benchmark_Step_Error(b *testing.B) {

	size := 400 //Size of Step, kB
	message := "Event bytes: %07d................................"
	sizeMessage := len(message) + len(prefixInfo) + 7
	m := int(size * 1024 / sizeMessage)
	counter := 0

	zLog := NewZLog()
	zLog.Step("Making Step(block) with size %d Kb and ERROR", size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < m; j++ {
			counter += sizeMessage
			zLog.Info(message, counter)
		}
	}

	zLog.Error("")
	fmt.Printf("\nMade %d messages with ERROR lenght: %d bytes, total: %d Kbytes", m*b.N, sizeMessage, int(counter/(1024)))

	_ = zLog.GetLog()
	//	fmt.Printf("\nLast message in Getlog:  %v\n", logs[len(logs)-1])

}

func Benchmark_Step_v1(b *testing.B) {

	size := 102400 //Size of Step, kB
	message := "Event bytes: %07d................................"
	sizeMessage := len(message) + len(prefixInfo) + 7
	m := int(size * 1024 / sizeMessage)
	counter := 0
	counterSteps := 0
	zLog := NewZLog()
	zLog.Step("Making Step(block) with size %d Kb and ERROR", size)
	fmt.Printf("\nMade %d messages with ERROR lenght: %d bytes, total: %d Kbytes", m, sizeMessage, int(counter/(1024)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counterSteps++
		zLog.Step(" Step:  %07d", counterSteps)
	}
	fmt.Printf("\nAdd %d Steps.", counterSteps)

}
