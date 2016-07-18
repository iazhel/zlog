package zlog

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

func deleteIfExist(name string) error {
	if _, err := os.Stat(name); err == nil {
		return os.Remove(name)
		fmt.Println("exsists")
		// delete file
	}
	return nil
}

func getSize(name string) int {
	for i := 1; i < 20; i++ {
		if s, err := os.Stat(name); err == nil {
			return int(s.Size())
		}
		time.Sleep(time.Millisecond)
	}
	return -1
}

func checkSize(name string, n int, t *testing.T) bool {
	if size := getSize(name); size != n {
		t.Errorf("Writed file size is not correct get: %d, need %d ", size, n)
		return false
	}
	return true
}

func Test_logFile_creating(t *testing.T) {
	pass_logPath := []string{"/tmp/1/1/1/1", "/tmp/zlog_test", "/tmp/1", "/tmp/log_test", "/tmp//"}
	outPath := []string{}
	logFilename := "zlog_autosave.log"

	// rewrite files
	for i, logpath := range pass_logPath {
		outPath = append(outPath, path.Join(logpath, logFilename))

		log := NewZL(outPath[i])
		log.CreateLogFile()
		oldSize := getSize(outPath[i])

		log = log.NewStep("Step 1")
		n, err := log.WriteAllLog()

		if err != nil {
			t.Errorf("rewrite file: %v, %v ", outPath[i], err)
		}

		if size := getSize(outPath[i]); size-oldSize != n {
			fmt.Println(outPath[i], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n)
		}
	}
	// remove files
	for i := range pass_logPath {
		err := deleteIfExist(outPath[i])
		if err != nil {
			t.Errorf("delete file: %v", err)
		}
	}
	// create and write into files
	for i := range pass_logPath {
		log := NewZL(outPath[i])
		log.NewStep("Step 1")
		n, err := log.WriteAllLog()
		if err != nil {
			t.Errorf("write file: %v", err)
		}
		if size := getSize(outPath[i]); size != n {
			fmt.Println(outPath[i], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n)
		}
	}

	// append to file
	log := NewZL(outPath[0])
	for _ = range pass_logPath {
		log.NewStep("Step ")
		oldSize := getSize(outPath[0])
		n, err := log.WriteAllLog()
		if err != nil {
			t.Errorf("write file: ", err)
		}
		if size := getSize(outPath[0]); size != n+oldSize {
			fmt.Println(outPath[0], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n+oldSize)
		}

	}

	return
	// remove files
	for i := range pass_logPath {
		err := deleteIfExist(outPath[i])
		if err != nil {
			t.Errorf("delete file: ", err)
		}
	}

}
