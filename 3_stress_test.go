// USE -test.v flag to see output at PASS.
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
	OK, ERR := "      "+suffixOK, " "+suffixError
	if curOS := runtime.GOOS; curOS == "windows" {
		OK, ERR = "[OK]", "[Error]"
	}

	if !pass {
		fmt.Printf("%10s|%10d\n\b", ERR, get)
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
		if testing.Short() {
			t.Skip("Skipping huge test in short mode")
		} else {
			fmt.Print("\n\b Huge test is started. It take about minute. \n\bUse '-test.short' flag to skip.  ")
		}
	}
}

// It is working routine.
func (task tTask) Do(t *testing.T, log *ZLogger, f0, f1, f2, f3 func(string, ...interface{})) {
	var ops int32 = 0
	ch := make(chan bool)
	// running routines
	log.Step("Step 1")
	for rc := 0; rc < task.routCount; rc++ {
		go func(n int) {
			f0("Start routine %d", n)
			for j := 0; j < task.loopDeep; j++ {
				f1(`In the loop.------------------------------------------------------------------------------------`)
			}
			f2("End of goroutine %d.", rc)
			// increase counter at goroutine end.
			atomic.AddInt32(&ops, 1)
			runtime.Gosched()
			// find the last goroutine.
			if atomic.LoadInt32(&ops) == int32(task.routCount) {
				f3("Last test message!")
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
		msgsCount := len(log.GetLog())
		pass := msgsCount >= task.msgsWant
		task.Control(t, pass, msgsCount, mem)
	case "write":
		nBytes, pass := log.WriteLog()
		task.writeControl(t, pass, task.msgsWant, nBytes, mem)
	case "rewrite":
		nBytes, pass := log.ReWriteLog()
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
		{1, 100000, 0, "get"},
		{10, 10000, 0, "get"},
		{100, 1000, 0, "get"},
		{1000, 100, 0, "write"},
		{10000, 10, 0, "rewrite"},
		{1000, 1000, 0, "get"},
	}
	tableHeader()
	for _, task := range taskCases {
		task.msgsWant = task.routCount*task.loopDeep + 2*task.routCount + 2
		task.printInfo()
		checkShortMode(t, task.msgsWant, 1000000)
		log := NewZLog("/tmp/go_tests/test1.log")

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Info, log.Info, log.Info, log.Error)
	}
}

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
		log := NewZLog("/tmp/go_tests/test2.log")
		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Step, log.Info, log.Warning, log.Info)
	}
}

// All steps messages must be saved. It count is routines count.
// Should return memory to programm.
func Test_Step_method(t *testing.T) {
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
		log := NewZLog()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Step, log.Info, log.Info, log.Info)
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
		log := NewZLog()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Info, log.Info, log.Info, log.Error)
	}

}

// All Steps messages must be saved into file.
// Test should return memory to OS after finish.
// Free OS memory should be increased.
func Test_Steps_Heap(t *testing.T) {
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
		log := NewZLog()

		// Do( .., .., at start routine, in loop, at end each routine, at end of all)
		task.Do(t, log, log.Step, log.Info, log.Info, log.Info)
	}
}

// A huge test for OS memory. All error messages must be saved.
//It count is more then routines*loops + 1.
func Test_GetLog_Heap(t *testing.T) {
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
		log := NewZLog()
		log.Step("Step 1")
		for rc := 0; rc < task.routCount; rc++ {
			go func(n int) {
				for j := 0; j < task.loopDeep; j++ {
					log.Error(`--------------------------------------
					----------------------------------------------`)
				}
				runtime.Gosched()
				//		log.Error("Last message is test error.")
				l := int32(len(log.GetLog()))
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
		atomic.AddInt32(&msgsGet, int32(len(log.GetLog())))
		msgsCountInt := int(atomic.LoadInt32(&msgsGet))
		pass := task.msgsWant <= msgsCountInt
		task.Control(t, pass, msgsCountInt, mem)

	}
	fmt.Println()
}
