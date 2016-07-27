package zlog

import (
	"fmt"
	"strings"
	"testing"
)

const EXCLD = "Should be excluded from log!"

func generateLog(n int) (msgs []string) {
	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("LOG message:%5d", i)
		msgs = append(msgs, msg)
	}
	return msgs
}

func checkLogOrder(log string, msgs []string, n int) (bool, string) {
	var last int
	for i := 0; i <= n; i++ {
		pos := strings.Index(log, msgs[i])
		if pos < 0 {
			return false, "'" + msgs[i] + "' not found."
		}
		if pos < last {
			return false, "'" + msgs[i] + "' not in order."
			last = pos
		}
	}
	if strings.Contains(log, EXCLD) {
		return false, "Log contains excluded message!"
	}
	return true, ""
}

func Test_01(t *testing.T) {
	m := generateLog(16)
	log := NewZL()
	log.Info(m[0])
	log.Info(m[1])
	log.Error(m[2])

	log.Step(m[3])
	log.Info(m[4])
	log.Warning(m[5])

	log.Step(m[6])
	log.Error(m[7])
	log.Error(m[8])
	log.Warning(m[9])

	log.Step(m[10])
	log.Info(EXCLD)

	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, m, 10)
	if !pass {
		t.Errorf(err)
	}
	fmt.Println(msgs)
	return
}

func Test_02(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := 0; i < n; i++ {
		log.Step(ms[i])
	}

	msgs := log.GetAllLog()
	fmt.Println(msgs)
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		t.Errorf(err)
	}
}

func Test_03(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := 0; i < n; i++ {
		log.NewStep(ms[i])
	}

	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		t.Errorf(err)
	}
	fmt.Println(msgs)
}

func Test_04(t *testing.T) {
	n := 10
	ms := generateLog(n)
	log := NewZL()

	for i := 0; i < n; i++ {
    	//log = log.NewStep(ms[i])
    	 log.NewStep(ms[i])
		log.Info(EXCLD + ms[i])
	}

	msgs := log.GetAllLog()
	pass, err := checkLogOrder(msgs, ms, n-1)
	if !pass {
		t.Errorf(err)
	}
	fmt.Println(msgs)
}

func Test_05(t *testing.T) {
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
		t.Errorf(err)
	}
	fmt.Println(msgs)
}




