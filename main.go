package zlog

import (
	"fmt"
	"math"
	"time"
)

func routine(n int, log *ZL) {
	//	fmt.Println("main", n, log.key)
	for i := 0; i < math.MaxInt32; i++ {
		time.Sleep(time.Millisecond)
		f := time.Now().Nanosecond()
		c := f % 10000000
		switch {
		case c >= 9999990:
			log.Error("rand number =%d, routine %d", c, n)
			fmt.Printf("\n%v", time.Now())
			fmt.Printf("\nError number %d, routine %d.\n", c, n)
		case c > 9999900:
			//			log.WriteAllLog()
			log.WriteStep()
			log.Step("Step: Writed step, cmd number =%d, routine %d", c, n)
		case c < 999:
			log.Step("Step Next: rand number =%d, routine %d", c, n)
		case c < 1999:
			log.Warning("rand number =%d, routine %d", c, n)
		case c < 9999990:
			log.Info("rand number =%d, routine %d", c, n)
		}
	}
}

func main() {
	n := 10
	log := make([]*ZL, n)
	log[0] = NewZL("/tmp/long_test.log")
	log[0].CreateLogFile()
	//	log[0].SetRemoveBeforeGet(true)
	log[0].SetWarningLines(5)
	//	log[0].SetRemoveBeforeGet(true)
	// bench cycle
	for j := 0; j < n; j++ {
		log[j] = log[0].NewStep("The New Step %d", j)
		go routine(j, log[j])
	}
	fmt.Printf("Started <%d> goroutines, %v \n", n, time.Now())
	time.Sleep(time.Minute * 1200)
	log[0].WriteAllLog()
	fmt.Println("END")
}
