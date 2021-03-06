package zlog

import (
	"fmt"
	"testing"
)

func Test_Z1100(t *testing.T) {
	sl := NewSL()
	sl.Info("Should not preset in output !!!!")
	filePath := "/tmp/logs/slog/s1.log"
	n, err := sl.WriteNewLog(filePath)
	if err != nil {
		t.Errorf("Can't possible to write into file %v", err)
	}

	fmt.Printf("Writed %d bytes\n", n)
	sl.Info("I0")
	sl.Info("I1")

	sl.Step("step ii..")
	sl.Info("I0")

	sl.Step("")
	sl.Step("step Long log ------------------------------------------------------------------------")

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

	sl.Step("step i")
	sl.Info("I")
	sl.Error("E")
	sl.Info("I")

	logs := sl.GetAllLog()

	sl.Info("I0")
	sl.Info("I1")
	logs += sl.GetAllLog()

	sl.Warning("W")
	sl.Info("I0")
	logs += sl.GetAllLog()

	sl.Error("E")
	sl.Info("I0")
	sl.Info("I1")
	logs += sl.GetAllLog()

	logs += sl.GetAllLog()

	fmt.Println("******SLOG**********")
	fmt.Print(logs)

	n, err = sl.WriteNew(logs, filePath)
	sl.Step("WriteNew()")
	if err != nil {
		t.Errorf("Can't possible to write into file %v", err)
		sl.Error("Not Correctly.")
	}

	logs = sl.GetAllLog()
	n, err = sl.Write(logs, filePath)
	fmt.Print(logs)

	sl.Step("Append .....")
	if err != nil {
		t.Errorf("Can't possible to append into file %v", err)
		sl.Error("Not Correctly.")
	}

	n, err = sl.WriteLog(filePath)
	if err != nil {
		t.Errorf("Can't possible to write into file %v", err)
	}
	fmt.Printf("Writed %d bytes\n", n)

}
