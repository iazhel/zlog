package zlog

import (
	"fmt"
	"os"
	"testing"
)

func Test_1100(t *testing.T) {
	zl := NewZL()
	return
	zl.Info("I0")
	zl.Info("I1")

	zl.Step("step ii..")
	zl.Info("I0")

	zl.Step("")
	zl.Step("step Long log ------------------------------------------------------------------------")

	zl.Step("step iwi..")
	zl.Info("I")
	zl.Warning("W")
	zl.Info("I")

	zl.Step("step iewi.")
	zl.Error("E")
	zl.Info("I")
	zl.Error("E")
	zl.Warning("W")
	zl.Info("I")

	zl.Step("step i")
	zl.Info("I")
	zl.Error("E")
	zl.Info("I")

	logs := zl.GetAllLog()

	zl.Info("I0")
	zl.Info("I1")
	logs += zl.GetAllLog()

	zl.Warning("W")
	zl.Info("I0")
	logs += zl.GetAllLog()

	zl.Error("E")
	zl.Info("I0")
	zl.Info("I1")
	logs += zl.GetAllLog()

	zl.Step("s.....")
	logs += zl.GetAllLog()

	fmt.Println("******zlOG**********")
	fmt.Println(logs)

	f, _ := os.Create("outz.log")
	n, _ := f.WriteString(logs)
	fmt.Printf("Writed %d bytes\n", n)

}
