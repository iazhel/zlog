package zlog

import (
	"fmt"
	"strings"
	"testing"
)

// this message should be added in log message,
// what will disapear in result.
const EXCLD = "Should be excluded from log!"

func generateLog(n int) (msgs []string) {
	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("LOG message:%5d", i)
		msgs = append(msgs, msg)
	}
	return msgs
}

func checkLogExist(log string, msgs []string, n int) (bool, string) {
	return checkLog(log, msgs, n, false)
}
func checkLogOrder(log string, msgs []string, n int) (bool, string) {
	return checkLog(log, msgs, n, true)
}

func checkLog(log string, msgs []string, n int, order bool) (bool, string) {
	var last int
	for i := 0; i <= n; i++ {
		pos := strings.Index(log, msgs[i])
		if pos < 0 {
			return false, "'" + msgs[i] + "' not found."
		}
		if pos < last && order {
			return false, "'" + msgs[i] + "' is not in true order!"
		}
		last = pos
	}
	if strings.Contains(log, EXCLD) {
		return false, "Log contains excluded message!"
	}
	return true, ""
}

// Check NewZL
// Logs should be empty
func Test_0000(t *testing.T) {
	log := NewZL()
	msgs := log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check NewStep
// Should be all Step messages
// Logs order should be correct
func Test_0100(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := 0; i < n; i++ {
		log = log.NewStep(ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check NewStep
// Should be all root messages
// Should not mix log1 and log(root and child)
func Test_0101(t *testing.T) {
	n := 10
	ms := generateLog(n)
	msgs := ""
	log := NewZL()
	// fill root
	for i := 0; i < n; i++ {
		log.Step(ms[i] + "A")
		i++
		log.Error(ms[i] + "A")
	}
	// create and fill log1
	log1 := log.NewStep(EXCLD + "B")
	for i := 0; i < n; i++ {
		log1 = log1.NewStep(EXCLD + ms[i] + "B")
		i++
		log1.Warning(EXCLD + ms[i] + "B")
	}
	// get only root,
	// check on excluding message from log1
	msgs += log.GetStep()
	//
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check NewStep
// Should be all log1 messages
// Should not mix log1 and log(root and child)
// Logs order should be correct
func Test_0102(t *testing.T) {
	n := 10
	ms := generateLog(n)
	msgs := ""
	log := NewZL()
	// fill root
	for i := 0; i < n; i++ {
		log.Step(EXCLD + ms[i] + "A")
		i++
		log.Error(EXCLD + ms[i] + "A")
	}
	log1 := log.NewStep(ms[0] + "B")

	// fill log1
	for i := 0; i < n; i++ {
		log1.Step(ms[i] + "B")
		i++
		log1.Warning(ms[i] + "B")
	}
	// get only log1,
	// check on excluding message from root
	msgs += log1.GetStep()
	//
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Step
// Should be writed all Step msgs
// Logs order should be uncorrect
func Test_0200(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := n - 1; i >= 0; i-- {
		log.Step(ms[i])
	}
	msgs := log.GetAllLog()
	// should be
	pass, err := checkLogExist(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// should be inverse order
	pass, err = checkLogOrder(msgs, ms, n-1)
	if pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Step
// Should be all Step messages
// Should be cleared after GetAllLog
// Should be writed all Step msgs
// Should be cleared after GetStep
// Should be writed all Step msgs
// Logs order should be correct
func Test_0201(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := 0; i < n; i++ {
		log.Step(ms[i])
	}
	// should be
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// should not be
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
	for i := 0; i < n; i++ {
		log.Step(ms[i])
	}
	// should be
	msgs = log.GetStep()
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// should not be
	msgs = log.GetStep()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
	for i := 0; i < n; i++ {
		log.Step(ms[i])
	}
	// should be
	msgs = log.GetStep()
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Step
// Logs should be saved after many Step.
// Logs order should be correct
func Test_0202(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	//for i := n-1; i >0; i--{
	for i := 0; i < n; i++ {
		log.Step(ms[i])
		i++
		log.Error(ms[i])
	}
	// should be
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// should not be
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check Info
// Logs should be in log result
func Test_0300(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n-1; i++ {
		log.Info(ms[i])
	}
	log.Error(ms[n-1])
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Info
// Logs should be excluded in log result
func Test_0301(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n; i++ {
		log.Step(ms[i])
		log.Info(EXCLD + ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Info
// should be created Step
func Test_0302(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	log.Info(ms[0])
	for i := 1; i < n; i++ {
		log.Info(EXCLD + ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, 0)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Warning
// Logs should be in log result
func Test_0400(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n; i++ {
		log.Warning(ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Warning
// Logs should be 10 info logs in the result
// Logs should set info logs count in the result
func Test_0401(t *testing.T) {
	n := 11
	ms := generateLog(n)
	log := NewZL()
	log.Step("step")
	for i := 0; i < n; i++ {
		log.Info(ms[i])
	}
	log.Warning("ms")
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	// test should be fail
	if pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// increase message count
	log.SetWarningLines(15)
	log.Step("step")
	for i := 0; i < n; i++ {
		log.Info(ms[i])
	}
	log.Warning("ms")
	msgs = log.GetAllLog()
	pass, err = checkLogOrder(msgs, ms, n-1)
	// test should be pass
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}

}

// Check Warning
// Only warning logs should be in log result
// Should be missed any excluded message
// Should be getted all messages

func Test_0402(t *testing.T) {
	n := 12
	ms := generateLog(n)
	log := NewZL()
	log.SetWarningLines(0)
	log.Step("step")
	for i := 0; i < n; i++ {
		log.Info(EXCLD)
		log.Warning(ms[i])
	}
	msgs := log.GetAllLog()
	// Only warning logs should be in log result
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}

	log.SetWarningLines(1)
	log.Step("step")
	for i := 0; i < n; i++ {
		log.Info(EXCLD)
		log.Info(ms[i])
		i++
		log.Warning(ms[i])
	}
	msgs = log.GetAllLog()
	// Should be missed any excluded messages
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}

	log.SetWarningLines(2)
	log.Step("step")
	for i := 0; i < n; i++ {
		log.Info(ms[i])
		i++
		log.Info(ms[i])
		i++
		log.Warning(ms[i])
	}
	msgs = log.GetAllLog()
	// Should be getted all messages
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Warning
// After many Step, Logs should be in log result

func Test_0403(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n; i++ {
		log.Info(ms[i])
		i++
		log.Warning(ms[i])
		i++
		log.Step(ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Error
// All error should be in log result
func Test_0500(t *testing.T) {
	n := 20
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n; i++ {
		log.Error(ms[i])
	}
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Error
// All Info should be in log result
func Test_0501(t *testing.T) {
	n := 20
	ms := generateLog(n)
	log := NewZL()
	for i := 0; i < n-1; i++ {
		log.Info(ms[i])
	}
	log.Error(ms[n-1])
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Error
// All Info & Warning & Error should be in log result
func Test_0502(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	log.SetWarningLines(2)
	for i := 0; i < n-1; i++ {
		log.Info(ms[i])
		i++
		log.Info(ms[i])
		i++
		log.Warning(ms[i])
	}
	log.Error(ms[n-1])
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check Error
// All logs should be in log result
// Logs shoud be saved after Step
func Test_0503(t *testing.T) {
	n := 40
	ms := generateLog(n)
	log := NewZL()
	log.SetWarningLines(2)
	for i := 0; i < n-1; i++ {
		log.Info(ms[i])
		i++
		log.Info(ms[i])
		i++
		log.Warning(ms[i])
		i++
		log.Step(ms[i])
	}
	log.Error(ms[n-1])
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check GetStep
// Should get all root messages
// Should clear log
// Should get all log1 messages
// Should clear log1
func Test_0600(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	// fill root(log)
	for i := 0; i < n; i++ {
		log.Step(ms[i] + "A")
		i++
		log.Error(ms[i] + "A")
	}
	// create and fill log1
	log1 := log.NewStep("B")
	for i := 0; i < n; i++ {
		log1.Step(ms[i] + "B")
		i++
		log1.Warning(ms[i] + "B")
	}
	// check message from log
	msgs := log.GetStep()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// Should not be messages after
	msgs = log.GetStep()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
	// check message from log
	msgs = log1.GetStep()
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// Should not be logs
	msgs = log1.GetStep()
	if len(msgs) != 0 {
		fmt.Println(msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}

}

// Check GetStep
// Should get all root messages
// Should not mix  log(root) and log1 (childs)
func Test_0601(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()
	// fill root
	for i := 0; i < n; i++ {
		log.Step(ms[i] + "A")
		i++
		log.Error(ms[i] + "A")
	}
	// create and fill log1
	log1 := log.NewStep(EXCLD + "B")
	for i := 0; i < n; i++ {
		log1 = log1.NewStep(EXCLD + ms[i] + "B")
		i++
		log1.Warning(EXCLD + ms[i] + "B")
	}
	// get only root,
	// check on excluding message from log1
	msgs := log.GetStep()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check GetStep
// Should not mix log1(child) and log(root and any chids)
// Should be all log1 messages
// Logs order should be correct
func Test_0602(t *testing.T) {
	n := 10
	ms := generateLog(n)
	msgs := ""
	log := NewZL()
	// fill root
	for i := 0; i < n; i++ {
		log = log.NewStep(EXCLD + ms[i] + "A")
		i++
		log.Error(EXCLD + ms[i] + "A")
	}
	log1 := log.NewStep(ms[0] + "B")

	// fill log1
	for i := 0; i < n; i++ {
		log1.Step(ms[i] + "B")
		i++
		log1.Warning(ms[i] + "B")
	}
	// get only log1,
	// check on excluding message from root
	msgs += log1.GetStep()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check GetStep
// Should get not all messages
// Should be in order
func Test_0603(t *testing.T) {
	n := 9
	ms := generateLog(n)
	log := NewZL()
	// fill root(log)
	for i := 0; i < n; i++ {
		log = log.NewStep(ms[i])
		i++
		log.Step(ms[i])
		i++
		log.Warning(ms[i])
	}
	// check message from log
	msgs := log.GetStep()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
}

// Check GetAllStep
// Should get all messages from all loggers
// Should be in order
// Should clear log
func Test_0700(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	// fill root(log)
	for i := 0; i < n; i++ {
		log = log.NewStep(ms[i])
		i++
		log.Step(ms[i])
		i++
		log.Warning(ms[i])
	}
	// check message from log
	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// Should clear log
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Printf("GetStep: '%s'", msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check GetAllStep
// Should get all messages from all loggers
// Should be in order
// Should clear log
func Test_0701(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	log = log.NewStep(ms[0])
	log1 := log.NewStep(ms[1])
	log2 := log.NewStep(ms[2])
	// fill root(log)
	for i := 3; i < n; i++ {
		log.Warning(ms[i])
		i++
		log1.Warning(ms[i])
		i++
		log2.Warning(ms[i])
	}
	// check message from log
	msgs := log.GetAllLog()
	pass, err := checkLogExist(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	pass, err = checkLogOrder(msgs, ms, n-1)
	if pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// Should clear log
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Printf("GetStep: '%s'", msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check GetAllStep
// Should get all messages from all loggers
// Should not mix step
func Test_0702(t *testing.T) {
	n := 30
	ms := generateLog(n)
	log := NewZL()
	log = log.NewStep(ms[0])
	log1 := log.NewStep(ms[10])
	log2 := log.NewStep(ms[20])
	// fill root(log)
	for i := 1; i < 10; i++ {
		log.Warning(ms[i])
		log1.Warning(ms[10+i])
		log2.Warning(ms[20+i])
	}
	// check message from log
	msgs := log.GetAllLog()
	pass, err := checkLogExist(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	pass, err = checkLogOrder(msgs, ms, n-1)
	if !pass {
		fmt.Println(msgs)
		t.Errorf(err)
	}
	// Should clear log
	msgs = log.GetAllLog()
	if len(msgs) != 0 {
		fmt.Printf("GetStep: '%s'", msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}

	msgs = log.GetStep()

	if len(msgs) != 0 {
		fmt.Printf("GetStep: '%s'", msgs)
		t.Errorf("Log size is not zero,  %d bytes", len(msgs))
	}
}

// Check Output view
// Should get all messages from all loggers
func Test_0800(t *testing.T) {
	n := 42
	ms := generateLog(n)
	log := NewZL()
	log = log.NewStep(ms[0])
	log1 := log.NewStep(ms[10])
	log2 := log.NewStep(ms[20])
	// fill root(log)
	for i := 12; i < 20; i++ {
		log.Info(ms[i])
		log1.Info(ms[10+i])
		log2.Info(ms[20+i])
	}
	log.Info("Last msg")
	log1.Warning("Last msg")
	log2.Error("Last msg")
	log.Step("step")
	log1.Step("step")
	log2.Step("step")
	msgs := log.GetAllLog()
	fmt.Println(msgs)
}
