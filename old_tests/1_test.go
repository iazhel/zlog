package zlog

import (
	"fmt"
	"testing"
)

const InfoLogMsg = "info:----------------------------------------------------"

func checkTest(t *testing.T, b bool, s string) {
	if !b {
		t.Errorf("%s FAIL", s)
	}
}

func createAndFillZL(step, info int) []*ZL {
	logs := make([]*ZL, step)
	logs[0] = NewZL()
//	logs[0].Step("Step")
//	logs[0].Info("Info")


	// NewStep cycle
	for j := 0; j < step; j++ {
		logs[j] = logs[0].NewStep("step")
		// Info cycle
		for k := 0; k < info; k++ {
			logs[j].Info(InfoLogMsg)
		//	logs[j].Step(InfoLogMsg)
		}
	}
	return logs
}

func Test_Step(t *testing.T) {
	log := NewZL()
	log.NewStep("Step A")
	log.Info("Info")
	log.Error("EEEEEE")
	fmt.Printf("%#v", log)
	msgs := log.GetAllLog()
	msgsSize := len(msgs)
	fmt.Println(msgs)
	checkTest(t, msgsSize == 1, "000.Logs are created: ")

return
	log = NewZL()
	log.Step("Step B")
	log.Error("Error 2 !")

	msgs = log.GetAllLog()
	msgsSize = len(msgs)
	fmt.Println(msgs)
	checkTest(t, msgsSize == 2, "000.Logs are created: ")

	log = NewZL()
	log.Step("Step C")
	log.Warning("Warning 3 !")

	msgs = log.GetAllLog()
	msgsSize = len(msgs)
	fmt.Println(msgs)
	checkTest(t, msgsSize == 2, "000.Logs are created: ")


}


func Test_GetLog_clear(t *testing.T) {
	steps, infos := 10, 1000
	logs := createAndFillZL(steps, infos)

	msgsSize := len(logs[0].GetAllLog())
	checkTest(t, msgsSize > steps*10, "001.Logs are created: ")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "002.Removing logs after GetAllLog: ")
}

func Test_Compress_clear(t *testing.T) {
	steps, infos := 10, 1000
	logs := createAndFillZL(steps, infos)

	logs[0].compressAll()
	msgs := logs[0].GetAllLog()
	msgsSize := len(msgs)
	checkTest(t, msgsSize == 0, "003.Logs are compressed: ")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "004.Removing logs after GetAllLog: ")
}

func Test_Error_size(t *testing.T) {

	steps, infos := 2, 5
	logs := createAndFillZL(steps, infos)
	logs[0].Error("An error...")

	msgs := logs[0].GetAllLog()
	msgsSize := len(msgs)

	checkTest(t, msgsSize < 2*infos*len(InfoLogMsg), "005. Logs size is corect: ")
	checkTest(t, msgsSize > 1*infos*len(InfoLogMsg), "006. Logs size is corect: ")

	steps, infos = 2, 1000
	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "007.Removing logs after GetAllLog: ")

	logs = createAndFillZL(steps, infos)
	logs[0].Error("An error...")
	logs[1].Error("An error...")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize > steps*infos*len(InfoLogMsg), "008. Logs size is corect: ")

}
func Test_Warning_size(t *testing.T) {

	steps, infos := 2, 1000
	logs := createAndFillZL(steps, infos)
	logs[0].Warning("An warning.")

	msgsSize := len(logs[0].GetAllLog())
	checkTest(t, msgsSize < 2*10*len(InfoLogMsg), "009. Logs size is corect: ")
	checkTest(t, msgsSize > 1*10*len(InfoLogMsg), "010. Logs size is corect: ")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "011.Removing logs after GetAllLog: ")

	logs = createAndFillZL(steps, infos)
	logs[0].Warning("An warning.")
	logs[1].Warning("An warning.")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize > 2*10*len(InfoLogMsg), "012. Logs size is corect: ")

}
func Test_SetRemoveBeforeGet(t *testing.T) {

	steps, infos := 1000, 0

	logs := createAndFillZL(steps, infos)
	logs[0].SetRemoveBeforeGet(true)
	msgs := logs[0].GetAllLog()
	msgsSize := len(msgs)
	checkTest(t, msgsSize == 0, "013. Logs size is corect: ")

	logs = createAndFillZL(steps, infos)
	logs[0].SetRemoveBeforeGet(false)
	msgs = logs[0].GetAllLog()
	msgsSize = len(msgs)

	checkTest(t, msgsSize > 1*steps*10, "014. Logs size is corect: ")
}

func Test_SetWarningLines(t *testing.T) {

	steps, infos := 2, 1000
	logs := createAndFillZL(steps, infos)
	logs[0].Warning("An warning.")
	logs[0].SetWarningLines(100)

	msgsSize := len(logs[0].GetAllLog())
	checkTest(t, msgsSize < 2*100*len(InfoLogMsg), "015. Logs size is corect: ")
	checkTest(t, msgsSize > 1*100*len(InfoLogMsg), "016. Logs size is corect: ")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "017.Removing logs after GetAllLog: ")

	logs = createAndFillZL(steps, infos)
	logs[0].Warning("An warning.")
	logs[1].Warning("An warning.")
	logs[0].SetWarningLines(2)

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize > 2*2*len(InfoLogMsg), "018. Logs size is corect: ")
	checkTest(t, msgsSize < 2*3*len(InfoLogMsg), "019. Logs size is corect: ")

}
func Test_Stress_clear(t *testing.T) {
	fmt.Println()
	steps, infos := 1, 500
	logs := createAndFillZL(steps, infos)

	logs[0].compressAll()
	msgs := logs[0].GetAllLog()
	msgsSize := len(msgs)
	checkTest(t, msgsSize == 0, "020.Logs are compressed: ")

	msgsSize = len(logs[0].GetAllLog())
	checkTest(t, msgsSize == 0, "021.Removing logs after GetAllLog: ")
}
