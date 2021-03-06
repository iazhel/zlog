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

func Test_WriteAllLog(t *testing.T) {
	logDir := []string{"/tmp/" /*"/tmp/1/", "/tmp/log_test","/tmp//"*/}
	logPath := []string{}
	logFilename := "zlog_autosave.log"

	// should rewrite files
	for i, dir := range logDir {
		logPath = append(logPath, path.Join(dir, logFilename))

		log := NewZL(logPath[i])
		log.CreateLogFile()
		fmt.Println("File created.")
		oldSize := getSize(logPath[i])

		log = log.NewStep("Step 1")
		n, err := log.WriteAllLog()
		if err != nil {
			t.Errorf("rewrite file: %v, %v ", logPath[i], err)
		}

		if size := getSize(logPath[i]); size-oldSize != n {
			fmt.Println(logPath[i], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n)
		} else {
			fmt.Println("Total size:", size)
		}
	}
	// should remove files
	for i := range logDir {
		err := deleteIfExist(logPath[i])
		if err != nil {
			t.Errorf("delete file: %v", err)
		}
	}
	// should create and write into files
	for i := range logDir {
		log := NewZL(logPath[i])
		log.NewStep("Step 1")
		n, err := log.WriteAllLog()
		if err != nil {
			t.Errorf("write file: %v", err)
		}
		if size := getSize(logPath[i]); size != n {
			fmt.Println(logPath[i], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n)
		} else {
			fmt.Println("Total size:", size)
		}
	}

	// should append to file
	log := NewZL(logPath[0])
	for _ = range logDir {
		log.NewStep("Step ")
		oldSize := getSize(logPath[0])
		n, err := log.WriteAllLog()
		if err != nil {
			t.Errorf("write file: ", err)
		}
		if size := getSize(logPath[0]); size != n+oldSize {
			fmt.Println(logPath[0], "FAIL")
			t.Errorf("001. Writed file size is not correct get: %d, need %d ", size, n+oldSize)
		} else {
			fmt.Println("Total size:", size)
		}

	}

	// should remove files
	for _, fileName := range logPath {
		fmt.Printf("Deleting file '%s'... ", fileName)
		err := deleteIfExist(fileName)
		if err != nil {
			t.Errorf("delete file: ", err)
		}
		fmt.Println(" Done.")
	}
	return

}
