package zlog

import (
	"fmt"
	"os"
	"testing"
)

func Test_1100(t *testing.T) {
	sl := NewSL()

	sl.Info("I0")
	sl.Info("I1")
	sl.Info("I2")

	sl.Step("step ii..")
	sl.Info("I1")
	sl.Info("I1")

	sl.Step("")
	sl.Step("step 333333333333333333333333333333333333333333333333333333333333333333333332")

	sl.Step("step iwi..")
	sl.Info("I")
	sl.Warning("W")
	sl.Info("I")

	sl.Step("step iewi.")
	sl.Error("E")
	sl.Info("I")
	sl.Error("E")
	sl.Warning("W")
	sl.Info("I")

	sl.Step("step iwei.")
	sl.Warning("W")
	sl.Info("I")
	sl.Warning("W")
	sl.Error("E")
	sl.Info("I")

	sl.Step("step i")
	sl.Info("I")
	sl.Error("E")
	sl.Info("I")

	logs := sl.GetAllLog()

	sl.Info("I")
	sl.Info("I")
	sl.Info("I")
	logs += sl.GetAllLog()

	sl.Warning("W")
	sl.Info("I")
	sl.Info("I")
	logs += sl.GetAllLog()

	sl.Error("E")
	sl.Info("I")
	sl.Info("I")
	logs += sl.GetAllLog()

	sl.Step("s.....")
	logs += sl.GetAllLog()

	fmt.Println("******SLOG**********")
	fmt.Println(logs)

	f, _ := os.Create("out.log")
	n, _ := f.WriteString(logs)
	fmt.Printf("Writed %d bytes\n", n)

}
