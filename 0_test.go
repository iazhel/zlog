package zlog

import (
	"fmt"
	"strings"
	"testing"
	//	"runtime"
	. "github.com/pranavraja/zen"
)

func Test_OldLog(t *testing.T) {

	Desc(t, "zlog_Chek", func(it It) {

		it("shold create var zLog (type *ZLogger) and verify initial data", func(expect Expect) {
			zLog := NewZLog()
			zLog.Error("Part 1.")
			logs := zLog.GetLog()

			zLog.Info("Event 1")
			expect(zLog).ToExist() // verify initial data
			expect(len(logs)).ToEqual(2)

		})

		return
		it("shold get only name of the Step with suffix '[OK]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Error("Part 1.")
			zLog.Info("Event 1")
			zLog.Info("Part 1.")
			zLog.Info("Event 1")

			zLog.Step("Part 2.")
			zLog.Info("Event 2")
			zLog.Warning("Part x2.")
			zLog.Info("Event 2")
			zLog.Info("Part 2.")
			zLog.Info("Part 2.")

			logs := zLog.GetLog()
			expect(len(logs)).ToEqual(1)
			expect(strings.Index(logs[0], suffixOK) >= 0).ToEqual(true)

		})

		return

		it("shold get only name of the Step with suffix '[OK]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Warning("Part 1.")
			zLog.Info("Event 1")
			zLog.Info("Part 0.")
			zLog.Warning("Part x1.")
			zLog.Info("Event 2")

			logs := zLog.GetLog()
			expect(len(logs)).ToEqual(1)
			expect(strings.Index(logs[0], suffixOK) >= 0).ToEqual(true)

		})
		it("shold get only name of the Step with suffix '[OK]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Info("Event 1")
			zLog.Info("Part 0.")
			zLog.Info("Event 2")

			logs := zLog.GetLog()
			expect(len(logs)).ToEqual(1)
			expect(strings.Index(logs[0], suffixOK) >= 0).ToEqual(true)

		})
		it("shold get only name of the Step with suffix '[OK]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Error("Part 1.")
			zLog.Info("Event 1")
			zLog.Info("Part 0.")
			zLog.Info("Event 2")

			logs := zLog.GetLog()
			for _, w := range logs {
				fmt.Printf("\n%v", w)
			}
			//			expect(len(logs)).ToEqual(1)
			//			expect(strings.Index(logs[0], suffixOK) >= 0).ToEqual(true)

		})
		return

		it("shold get all blocks log and name of the block with suffix '[WARNING]'", func(expect Expect) {

			zLog := NewZLog()
			zLog.Step("Part 2.")
			zLog.Info("Event 1")
			zLog.Info("Event 2")
			zLog.Warning("Event 3")

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(4)
			expect(strings.Index(logs[0], suffixWarning) >= 0).ToEqual(true)

		})

		it("shold get all blocks log and name of the block with suffix '[ERROR]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Step("Part 3.")
			zLog.Info("Event 1")
			zLog.Info("Event 2")
			zLog.Error("Event 3")

			zLog.SetWarningLenght(1)

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(4)
			expect(strings.Index(logs[0], suffixError) >= 0).ToEqual(true)

		})

		it("shold get 34 from 49 logs and name of the block with suffix '[WARNING]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Step("Part 5.")
			for i := 0; i < 3; i++ {
				for j := 0; j < 15; j++ {
					zLog.Info("Event %d", i*15+j)
				}
				zLog.Warning("Result= %t", i < 2)
			}

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(34)
			expect(strings.Index(logs[0], suffixWarning) >= 0).ToEqual(true)

		})

		it("shold get 37 from 52 logs and name of the block with suffix '[WARNING]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Step("Part 6.")
			for i := 0; i < 3; i++ {
				for j := 0; j < 15; j++ {
					zLog.Info("Event %d", i*15+j)
				}
				zLog.Warning("Result= %t", i < 2)
				zLog.Step("Part 6.%d", i)
			}

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(37)
			expect(strings.Index(logs[0], suffixWarning) >= 0).ToEqual(true)

		})

		it("shold get 22 logs and name of the block with suffix '[WARNING]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Step("Part 7.")
			for i := 0; i < 20; i++ {
				zLog.Info("Event %d", i)
			}
			zLog.Warning("Result= %d", 2)

			zLog.SetWarningLenght(100)
			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(22)
			expect(strings.Index(logs[0], suffixWarning) >= 0).ToEqual(true)

		})

		it("shold get 29 logs and name of the block with suffix '[ERROR]'", func(expect Expect) {

			zLog := NewZLog()

			zLog.Step("Part 8.")
			for i := 0; i < 6; i++ {

				if i == 3 {
					zLog.Error("i== %d", i)
				}

				for j := 0; j < 15; j++ {
					zLog.Info("Event %d", i*15+j)
				}

				if i == 5 {
					zLog.SetWarningLenght(5)
					zLog.Warning("i== %d", i)
				}

				zLog.Step("Part 8.%d", i)
			}

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(29)
			expect(strings.Index(logs[0], suffixOK) >= 0).ToEqual(true)

		})

		it("shold verify output format", func(expect Expect) {

			i := 1
			f := 55.55555555
			t := true
			fileNames := []string{
				"data01.log",
				"data02.log",
				"out_data01.log"}

			message := []string{
				"\n\n Output examples: Block #04",
				"Opening files: %v ",
				"Reading file #%d",
				"Calculation finished: %t",
				"Calculation result: %.3f",
				"Cann't create file: %v\n\n"}

			zLog := NewZLog()

			zLog.Step(message[0])
			zLog.Info(message[1], fileNames)
			zLog.Info(message[2], i)
			zLog.Info(message[3], t)
			zLog.Warning(message[4], f)
			zLog.Error(message[5], fileNames[2])

			zLog.SetWarningLenght(1)

			logs := zLog.GetLog()

			expect(len(logs)).ToEqual(6)

			for j, w := range logs {
				fmt.Printf("\n%v", w)
				switch j {
				case 0:
					s := prefixStep + fmt.Sprintf(message[0])
					s = fmt.Sprintf(suffixFormat, s, suffixError)
					expect(logs[j]).ToEqual(s)
				case 3:
					s := prefixInfo + fmt.Sprintf(message[3], t)
					expect(logs[j]).ToEqual(s)
				case 4:
					s := prefixWarning + fmt.Sprintf(message[4], f)
					expect(logs[j]).ToEqual(s)
				case 5:
					s := prefixError + fmt.Sprintf(message[5], fileNames[2])
					expect(logs[j]).ToEqual(s)

				}

			}

		})

		fmt.Printf("\n")
		return

	})

}
