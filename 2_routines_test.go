package zlog

import (
	"fmt"
	. "github.com/pranavraja/zen"
	"strings"
	"testing"
)

// Function creates 3 Info messages.
func routInfo(logs *ZLogger, msg string, ready chan<- bool) {

	logs.Info("Goroutine routInfo is running...")
	logs.Info(msg)
	logs.Info("")
	ready <- true

}

// Function creates 1 warning and 2 Info messages.
func routWarning(logs *ZLogger, msg string, ready chan<- bool) {

	readyW := make(chan bool)
	logs.Info("Goroutine routWarning is running...")
	logs.Info("Goroutine routWarning call routInfo...")
	go routInfo(logs, "", ready)
	logs.Warning(msg)
	<-readyW

}

// Function creates 1 error and 5 Info messages.
func routError(logs *ZLogger, msg string, ready chan<- bool) {

	readyE := make(chan bool)
	logs.Info("Goroutine routError is running...")
	logs.Info("Goroutine routError call routInfo...")
	routInfo(logs, "", ready)
	logs.Error(msg)
	<-readyE

}

// Tests...
func Test_zlogRoutine(t *testing.T) {

	Desc(t, "zlog_ChekRoutine", func(it It) {
		it("should create two independent logs and fill their by goroutines.", func(expect Expect) {

			log1 := NewZLog()
			ready := make(chan bool)

			log1.Step("Step 1.")
			go func() {
				log1.Info("Msg 1..")
				ready <- true
			}()
			<-ready
			go func(msg string) {
				log1.Info("Msg 2..")
				log1.Warning("Msg 2..%s", msg)
				ready <- true
			}("Incorrect format  X")
			<-ready

			msgs1 := log1.GetLog() // get log1 messages
			msgs2 := log1.GetLog() // get log2 messages

			expect(len(msgs1)).ToEqual(4) // msgs1 len must be
			expect(len(msgs2)).ToEqual(0) // msgs2 len must be
		})

		it("should create two independent logs variables by one logger.", func(expect Expect) {

			logger := NewZLog()
			logger.Step("Step 1.")

			msgs1 := logger.GetLog()

			expect(msgs1).ToExist()       // verify initialization
			expect(len(msgs1)).ToEqual(1) // verify len log

		})

		it("should create two independent logs and fill it.", func(expect Expect) {

			logger := NewZLog()

			logger.Step("Step 1.")
			logger.Error("Msg 1..")
			msgs1 := logger.GetLog()

			logger.Step("Step 2.")
			logger.Error("Msg 2..")
			msgs2 := logger.GetLog()

			expect(strings.Contains(msgs1[1], "Msg 1..")).ToEqual(true)  // must contain
			expect(strings.Contains(msgs1[1], "Msg 2..")).ToEqual(false) // must not contain
			expect(strings.Contains(msgs2[1], "Msg 2..")).ToEqual(true)  // must contain
			expect(strings.Contains(msgs2[1], "Msg 1..")).ToEqual(false) // must not contain
		})

		it("should create 3 independent logs and fill their by goroutines.", func(expect Expect) {

			logger := NewZLog()
			ready := make(chan bool, 3)

			logger.Step("The fmt.Print(GetLog()) output example. Step 1.")
			logger.Step("The output example. Step 2.")
			logger.Step("The output example. Step 3.")

			go routWarning(logger, "Rout 2..", ready)
			go routError(logger, "Rout 3..", ready)
			go routInfo(logger, "Rout 1..", ready)

			<-ready
			<-ready
			<-ready

			allmsgs := logger.GetLog() // get all result table, and clear all.
			fmt.Println(allmsgs)
			expect(len(allmsgs)).ToEqual(18) // allmsgs len must be
		})

	})
}
