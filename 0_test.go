package zlog

import (
	"fmt"
	"testing"
	//    "strings"
)

func generateLogn(n int) (msgs []string) {
	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("%10d", i)
		msgs = append(msgs, msg)
		fmt.Println(msg)
	}
	return msgs

}
func Test_Step_0(t *testing.T) {
	log := NewZL()
	log.NewStep("001.        ")
	log.Info("002. Info    ")
	log.Error("003. EEEEE    ")
	log.NewStep("004.        ")
	log.Info("005. Info   ")
	log.Warning("006. WWWww    ")
	log.NewStep("007.        ")
	log.Error("008. Error")
	log.Warning("009. WWWww    ")
	log = log.NewStep("010.        ")
	log.Info("011. Info")

	msgs := log.GetAllLog()
	//	msgs := log.GetStep()
	fmt.Print(msgs)
	log.Error("012. Error")
	log.Warning("013. WWWww    ")
	log.Info("014. Info")
	log.Info("015. Info")
	msgs = log.GetAllLog()
	fmt.Print(msgs)

	return
}

func Test_New_Step_1(t *testing.T) {
	log := NewZL()
	//	log.Step("001. Step 00")
	log.Info("001. info")
	log.Info("002. Info 00 ")
	log.Error("003. EEEEE 00 ")
	log.Step("004. Step 01")
	log.Info("005. Info 01")
	log.Warning("006. WWWww 011")
	log.Step("007. Step 01")
	log.Error("008. Error")
	log.Warning("009. WWWww 011")
	log.Step("010. Step 01")
	log.Info("011. Info")

	msgs := log.GetAllLog()
	fmt.Print(msgs)
	return
}
