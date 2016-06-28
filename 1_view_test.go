package zlog

import (
	"fmt"
	. "github.com/pranavraja/zen"
	"strings"
	"testing"
)

// for print string slice
func printlogs(self []string) {

	for _, msg := range self {
		fmt.Println(msg)
	}
	return
}

// tests
func Test_NewZLog(t *testing.T) {

	Desc(t, "zlog_Chek", func(it It) {

		it("should create var logger (type *ZLogger) and get empty log list.", func(expect Expect) {
			log := NewLogger()
			msgs := log.GetLog()
			logger := NewLogger()
			msgs = logger.GetLog()
			expect(logger).ToExist()     // verify creation
			expect(len(msgs)).ToEqual(0) // verify empty log list

		})

		it("should create var logger and make first Step.", func(expect Expect) {

			log := NewLogger()
			log.Step("Step 1")
			msgs := log.GetLog()
			expect(log).ToExist()        // verify initialization
			expect(len(msgs)).ToEqual(1) // verify first log

		})

		it("should not fail on NewStep and GetAllLog.(For old version compatibility)", func(expect Expect) {

			log := NewLogger()
			log = log.NewStep("Next Step")
			log.Info("Message ...")
			msgs := log.GetLog()
			expect(log).ToExist()        // verify initialization
			expect(len(msgs)).ToEqual(1) // verify first log

		})

		it("should get one message with step name and suffix [OK].", func(expect Expect) {

			log1 := NewLogger()
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Info("Message ...")

			msgs := log1.GetLog()

			expect(len(msgs)).ToEqual(1)                              // verify messages number
			expect(strings.Contains(msgs[0], suffixOK)).ToEqual(true) // verify suffix [OK]

		})
		it("should get 3 messages with suffix [Warning].", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Warning("Message ...")

			msgs := log1.GetLog()

			expect(len(msgs)).ToEqual(3)                                   // verify messages number
			expect(strings.Contains(msgs[0], suffixWarning)).ToEqual(true) // verify suffix

		})
		it("should get 3 messages with suffix [Warning].", func(expect Expect) {
			log := NewLogger()
			log.SetWarningLenght(2)

			log.Step("Step 1.")
			log.Info("Message ..1i.")
			log.Info("Message ..2i.")
			log.Info("Message ..3i.")
			log.Info("Message ..4i.")
			log.Warning("Message ..W.")

			log.Step("Step 2.")
			log.Info("Message ..1i.")
			log.Info("Message ..2i.")
			log.Info("Message ..3i.")
			log.Warning("Message ..W.")

			msgs := log.GetLog()

			expect(len(msgs)).ToEqual(8) // verify messages number
		})

		it("should get 3 messages with suffix [Error].", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Error("Message ...")

			msgs := log1.GetLog()

			expect(len(msgs)).ToEqual(3)                                 // verify messages number
			expect(strings.Contains(msgs[0], suffixError)).ToEqual(true) // verify suffix

		})
		it("should get 4 messages suffix [Error].", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Warning("Message ...")
			log1.Error("Message ...")

			msgs := log1.GetLog()

			expect(len(msgs)).ToEqual(4)                                 // verify messages number
			expect(strings.Contains(msgs[0], suffixError)).ToEqual(true) // verify suffix

		})

		it("should get no logs names, if GetLog had run before.", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")

			msgs := log1.GetLog()
			msgs1 := log1.GetLog()
			expect(len(msgs)).ToEqual(1)  // verify first log
			expect(len(msgs1)).ToEqual(0) // verify second log

		})

		it("should get two step names with suffixes [OK].", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Step("Step 1.2.")

			log1.Info("Message ...")
			log1.Info("Message ...")

			msgs := log1.GetLog()
			expect(len(msgs)).ToEqual(2)                              // verify messages number
			expect(strings.Contains(msgs[0], suffixOK)).ToEqual(true) // verify suffix [OK]
			expect(strings.Contains(msgs[1], suffixOK)).ToEqual(true) // verify suffix [OK]

		})

		it("should get 6 messages with suffixes [OK], [Warning], [Error].", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")

			log1.Step("Step 1.1.")
			log1.Warning("Message ...")

			log1.Step("Step 1.2.")
			log1.Error("Message ...")

			log1.Step("Step 1.3.")
			log1.Info("Message ...")

			msgs := log1.GetLog()

			expect(len(msgs)).ToEqual(6)                                   // verify messages number
			expect(strings.Contains(msgs[0], suffixOK)).ToEqual(true)      // verify suffix
			expect(strings.Contains(msgs[1], suffixWarning)).ToEqual(true) // verify suffix
			expect(strings.Contains(msgs[3], suffixError)).ToEqual(true)   // verify suffix
			expect(strings.Contains(msgs[5], suffixOK)).ToEqual(true)      // verify suffix

		})

		it("should create unknown step, if Info has been run before Step.", func(expect Expect) {

			log := NewLogger()
			log.Info("Message ...")

			msgs := log.GetLog()
			expect(len(msgs)).ToEqual(1)                                         // verify messages number
			expect(strings.Contains(msgs[0], unknownStepName[:6])).ToEqual(true) // verify step name

		})

		it("should create unknown step, if Error has been run before Step.", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			msgs := log1.GetLog()

			log1.Error("Message ...")

			msgs = log1.GetLog()
			expect(len(msgs)).ToEqual(2)                                         // verify messages number
			expect(strings.Contains(msgs[0], unknownStepName[:6])).ToEqual(true) // verify step name

		})
		it("should get 4 messages and save to reserve file.", func(expect Expect) {

			log1 := NewLogger()
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Warning("Message ...")
			log1.Error("Message ...")

			b, err := log1.WriteLog()

			expect(err).ToEqual(nil)     // verify messages number
			expect(b > 50).ToEqual(true) // verify messages number

		})
		it("should get 4 messages and save to reserve file.", func(expect Expect) {

			log1 := NewLogger("/tmp/test.log")
			log1.Step("Step 1.")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Warning("Message ...")

			b, err := log1.WriteLog()

			expect(err).ToEqual(nil)     // verify messages number
			expect(b > 70).ToEqual(true) // verify messages number

		})
		it("should create unknown step, if Warning has been run before Step.", func(expect Expect) {

			log1 := NewLogger("/tmp/test.log")
			log1.Warning("Message warning...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			msgs := log1.GetLog()
			fmt.Println(msgs)
			expect(len(msgs)).ToEqual(2)                                         // verify messages number
			expect(strings.Contains(msgs[0], unknownStepName[:6])).ToEqual(true) // verify step name

		})
		it("should create unknown step, if Warning has been run before Step.", func(expect Expect) {

			log1 := NewLogger("/tmp/test.log")
			log1.Error("Message error...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			log1.Info("Message ...")
			msgs := log1.GetLog()
			fmt.Println(msgs)
			expect(len(msgs)).ToEqual(5)                                         // verify messages number
			expect(strings.Contains(msgs[0], unknownStepName[:6])).ToEqual(true) // verify step name

		})

	})
}
