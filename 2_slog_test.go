package zlog

import (
	"fmt"
	"testing"
)

func Test_1100(t *testing.T) {
	sl := NewSL()

	sl.Info("I")

	sl.Step("step .....xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")

	sl.Step("step .....")

	sl.Step("step i....")
	sl.Info("I")

	sl.Step("step iw...")
	sl.Info("I")
	sl.Warning("W")

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

	sl.Step("step ieixxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	sl.Info("I")
	sl.Error("E")
	sl.Info("I")

	logs := sl.GetAllLog()

	sl.Info("I")
	logs += sl.GetAllLog()

	sl.Warning("W")
	logs += sl.GetAllLog()

	sl.Error("E")
	logs += sl.GetAllLog()

	sl.Step("s.....")
	logs += sl.GetAllLog()

	fmt.Println("******SLOG**********")
	fmt.Println(logs)

}
