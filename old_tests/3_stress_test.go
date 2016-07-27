// USE -test.short flag to skip huge tests.

package zlog

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
)

type tTask struct {
	routCount int
	loopDeep  int
	msgsWant  int
	operation string
}

// Print memory size in in MB.
func memPrint(mem runtime.MemStats) {
	fmt.Printf("%4d/", int(mem.Alloc/1048576))
	fmt.Printf("%4d/", int(mem.HeapAlloc/1048576))
	fmt.Printf("%4d/", int(mem.HeapSys/1048576))
	fmt.Printf("%4d|\n\b", int(mem.HeapReleased/1048576))
}

func (self tTask) printInfo() {
	fmt.Printf("%10d|%10d|", self.routCount, self.loopDeep)
	fmt.Printf("%10d|", self.msgsWant)
}

func (self tTask) Control(t *testing.T, pass bool, get int, mem runtime.MemStats) {
	OK, ERR := "      [\033[32mOK\033[0m]", " "+"[ \033[31mERROR\033[0m ]"
	if curOS := runtime.GOOS; curOS == "windows" {
		OK, ERR = "[OK]", "[Error]"
	}

	if !pass {
		fmt.Printf("%10s|%10d|", ERR, get)
		memPrint(mem)
		t.Error("Getted messages count is not correct. ")

	} else {
		fmt.Printf("%10s|%10d|", OK, get)
		memPrint(mem)
	}
}

// TODO: Control by OS.
func (self tTask) writeControl(t *testing.T, pass error, want, nBytes int, mem runtime.MemStats) {
	OK, ERR := suffixOK, suffixError
	if curOS := runtime.GOOS; curOS == "windows" {
		OK, ERR = "[OK]", "[Error]"
	}
	if pass != nil {
		fmt.Printf("%40s|. Error:%v\n\b", ERR, pass)
		memPrint(mem)
	} else {
		fmt.Printf(" %51s|%10s|", OK, "???")
		memPrint(mem)
	}
}

// skip test if it is huge.
func checkShortMode(t *testing.T, msgsWant, limit int) {
	if msgsWant > limit {
		//		t.Skip("Skipping huge test in short mode")
		if testing.Short() {
			t.Skip("Skipping huge test in short mode")
		} else {
			fmt.Print("\n\b Huge test is started. It take about minute. \n\bUse '-test.short' flag to skip.")
		}
	}
}
func (task tTask) Do_A(t *testing.T) {
	var ops int32 = 0
	ch := make(chan bool)
	// running routines
	log := make([]*ZL, task.routCount)
	log[0] = NewZL()
	log[0].SetRemoveBeforeGet(true)

	for rc := 0; rc < task.routCount; rc++ {
		go func(n int) {
			//			log[n] = log[0].NewStep("NewStep %d", n)
			log[n] = log[0].NewStep("Step")
			//			fmt.Println("n=", n, "rc=", rc)
			for j := 0; j < task.loopDeep; j++ {
				log[n].Info(`In the loop.-----------------------------------------------------------------------`)
			}
			log[n].Error("End of goroutine")
			// increase counter at goroutine end.
			atomic.AddInt32(&ops, 1)
			runtime.Gosched()
			// find the last goroutine.
			if atomic.LoadInt32(&ops) == int32(task.routCount) {
				log[0].Error("Last test message!")
				ch <- true
			}
		}(rc)
	}
	// wait last routine signal.
	<-ch
	// do memory print.
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	// result control or writing.
	switch task.operation {
	case "get":
		msgs := log[0].GetAllLog()
		//		fmt.Println(msgs)
		msgsCount := len(msgs)
		pass := msgsCount >= task.msgsWant
		task.Control(t, pass, msgsCount, mem)
	case "write":
		//		nBytes, pass := log[0].WriteLog()
		//		task.writeControl(t, pass, task.msgsWant, nBytes, mem)
	case "rewrite":
		//		nBytes, pass := log[0].ReWriteLog()
		//		task.writeControl(t, pass, task.msgsWant, nBytes, mem)

	}
}

func (task tTask) Do_1(t *testing.T) {
	var ops int32 = 0
	ch := make(chan bool)
	// running routines
	log := NewZL()
	log.SetRemoveBeforeGet(true)

	for rc := 0; rc < task.routCount; rc++ {
		go func(n int) {
			log = log.NewStep("Step")
			//			fmt.Println("n=", n, "rc=", rc)
			for j := 0; j < task.loopDeep; j++ {
				log.Info(`In the loop.-----------------------------------------------------------------------`)
			}
			log.Error("End of goroutine")
			// increase counter at goroutine end.
			atomic.AddInt32(&ops, 1)
			//	runtime.Gosched()
			// find the last goroutine.
			if atomic.LoadInt32(&ops) == int32(task.routCount) {
				log.Error("Last test message!")
				ch <- true
			}
		}(rc)
	}
	// wait last routine signal.
	<-ch
	// do memory print.
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	// result control or writing.
	switch task.operation {
	case "get":
		msgs := log.GetAllLog()
		//		fmt.Println(msgs)
		msgsCount := len(msgs)
		pass := msgsCount >= task.msgsWant
		task.Control(t, pass, msgsCount, mem)
	case "write":
		nBytes, pass := log.WriteAllLog()
		task.writeControl(t, pass, task.msgsWant, nBytes, mem)
	case "rewrite":
		log.CreateLogFile()
		nBytes, pass := log.WriteAllLog()
		task.writeControl(t, pass, task.msgsWant, nBytes, mem)

	}
}

func tableHeader() {
	fmt.Printf("\n[\033[34m%9s|%10s|%10s|%10s|%10s|%20s|\033[0m\n", "Routines", "Loops", "Msgs want", "result", "Get msgs", "Memory MB A/HA/HS/HR")
}

// All Info messages must be saved in file.
// Should return memory to OS after finish.
func Test_Error__method(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep, msgsWant(will recalculate)}
		// light
		{1, 100000, 0, "get"},
		{10, 10000, 0, "get"},
		// huge
		{100000, 30, 0, "write"},
		{100000, 40, 0, "write"},
		// standart
		{100, 1000, 0, "get"},
		{1000, 100, 0, "write"},
		{10000, 10, 0, "rewrite"},
		{100000, 1, 0, "write"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount*task.loopDeep + 2*task.routCount + 2
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000000)
		task.Do_1(t)
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000000)
		task.Do_A(t)
	}
}

/*
// Should get Warning messages and control it min count.
func Test_Warning_method(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep}
		{1, 100000, 0, "get"},
		{10, 100000, 0, "get"},
		{100, 10000, 0, "get"},
		{1000, 100, 0, "rewrite"},
		{10000, 10, 0, "get"}}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount*2 + 2
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000000)
		log := NewZL("/tmp/go_tests/test2.log")
		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.NewStep, log.Info, log.Warning, log.Info)
	}
}

// All steps messages must be saved. It count is routines count.
// Should return memory to programm.
func Test_NewStep_method(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep}
		{1, 1000000, 0, "write"},
		{10, 100000, 0, "get"},
		{100, 10000, 0, "get"},
		{10000, 1000, 0, "get"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount + 1
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000)
		log := NewZL()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.NewStep, log.Info, log.Info, log.Info)
	}
}

// All Info messages must be saved. It count is routines*loops + 2.
// Should return memory to OS after finish.
func Test_Error_Heap(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep}
		{1, 100000, 0, "write"},
		{10, 10000, 0, "get"},
		{100, 1000, 0, "get"},
		{1000, 100, 0, "get"},
		{10000, 10, 0, "get"},
		{100000, 1, 0, "get"},
		{100, 100000, 0, "get"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount + 1
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000)
		log := NewZL()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Info, log.Info, log.Info, log.Error)
	}

}

// All NewSteps messages must be saved into file.
// Test should return memory to OS after finish.
// Free OS memory should be increased.
func Test_NewSteps_Heap(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep}
		{1, 100000, 0, "get"},
		{10, 10000, 0, "rewrite"},
		{100, 1000, 0, "get"},
		{1000, 100, 0, "get"},
		{10000, 10, 0, "get"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount + 1
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000)
		log := NewZL()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.NewStep, log.Info, log.Info, log.Info)
	}
}
// A huge test for OS memory. All error messages must be saved.
//It count is more then routines*loops + 1.
func Test_GetAllLog_Heap(t *testing.T) {
	taskCases := []tTask{
		//{routCount, loopDeep}
		{1, 100000, 0, "get"},
		{10, 10000, 0, "write"},
		{100, 1000, 0, "get"},
		{333, 1000, 0, "get"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount*task.loopDeep + 1
		task.printInfo()
		checkShortMode(t, task.msgsWant, 333000)
		var ops, msgsGet int32 = 0, 0
		ch := make(chan bool)
		log := NewZL()
		log.NewStep("NewStep 1")
		for rc := 0; rc < task.routCount; rc++ {
			go func(n int) {
				for j := 0; j < task.loopDeep; j++ {
					log.Error(`--------------------------------------
					----------------------------------------------`)
				}
				runtime.Gosched()
				//		log.Error("Last message is test error.")
				l := int32(len(log.GetAllLog()))
				atomic.AddInt32(&msgsGet, l)
				atomic.AddInt32(&ops, 1)
				if atomic.LoadInt32(&ops) == int32(task.routCount) {
					ch <- true
				}
			}(rc)
		}
		<-ch
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		atomic.AddInt32(&msgsGet, int32(len(log.GetAllLog())))
		msgsCountInt := int(atomic.LoadInt32(&msgsGet))
		pass := task.msgsWant <= msgsCountInt
		task.Control(t, pass, msgsCountInt, mem)

	}
	fmt.Println()
}

*/
