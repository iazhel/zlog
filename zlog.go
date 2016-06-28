package zlog

// Package zlog implements a logging package. /It defines a type, Log.
/*
	// initialiation
	log := NewLogger()
	//  or
	log := NewLogger('/tmp/zlog.log')

	// working code
	log.Step("Step 1.")
	log.Info("Msg ...")
	log.Warning("Msg ...%d...", intVar, ...)
	log.Error("Msg ... %v"..., errVar, ...)

	// OUTPUT
	messages := log.GetLog() // get only messages and clear log.
	writtenByte, err := log.WriteLog() // write to source file, and clear all logs.

*/

// Messages in this logger are written with methods Error, Warning, Info, Step.
// Method GetLog moves all saved information into string slice.
// Method WriteLog moves all saved information into string slice.

// Log adds prefixes to all messages.
// These methods support formatting like fmt.Sprintf.
// Logs divide messages into the blocks by Step method.
// This method creates and saves name for block.
// It adds suffixe to block name.
// That Step method deletes all info logs, if these logs contain
// no warning or error message.
// When error is in the block,
// Step saves all logs and name of current step to stringRes.
// When warning is in the block,
// zlog saves logs, that were written before warning message.
// Number of saved logs euivalents warningLen.
// WarningLen may be set with SetWarningLen function.
// WarningLen is 10, when ZLoger is initialaized.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

const (
	reserveLogFile_Lin = "/tmp/zlog_autosave.log"
	reserveLogFile_Win = "C:\\go_test.log"
)

type ZLogger struct {
	logThisStep      []string // logs of the block
	stringRes        []string // for saving valuable logs
	flagError        bool     // error availability flag
	posWarning       []int    // warning availability position
	notFirstMsg      bool     // the not first message in the step flag
	warningLen       int      // the numbers of lines to save
	linesAlloc       int      // the numbers of lines had saved, after memory clear.
	sync.Mutex                //
	ReserveFile      string
	OutSource        string
	pOK, pWarn, pErr string // perfixex
	sOK, sWarn, sErr string // suffixes
}

// NewLogger creates a new object.
func NewLogger(out ...interface{}) *ZLogger {
	self := &ZLogger{}
	currentOS := runtime.GOOS
	if len(out) > 0 {
		value := reflect.ValueOf(out[0])
		self.OutSource = value.String()
	}
	switch currentOS {
	case "linux":
		self.ReserveFile = reserveLogFile_Lin
		self.pWarn = prefixWarning
		self.pErr = prefixError
		self.sOK = suffixOK
		self.sWarn = suffixWarning
		self.sErr = suffixError

	case "windows":
		self.ReserveFile = reserveLogFile_Win
		self.pWarn = prefixWarning_Win
		self.pErr = prefixError_Win
		self.sOK = suffixOK_Win
		self.sWarn = suffixWarning_Win
		self.sErr = suffixError_Win
	}
	self.warningLen = 10
	return self
}

// Function creates log with error description.
func (self *ZLogger) Error(format string, v ...interface{}) {
	self.Lock()
	self.checkMsgPos("")
	self.flagError = true
	self.logThisStep = append(self.logThisStep, fmt.Sprintf(self.pErr+format, v...))
	self.Unlock()
}

// Function creates log with warning description.
func (self *ZLogger) Warning(format string, v ...interface{}) {
	self.Lock()
	self.checkMsgPos("")
	self.posWarning = append(self.posWarning, len(self.logThisStep))
	self.logThisStep = append(self.logThisStep, fmt.Sprintf(self.pWarn+format, v...))
	self.Unlock()
}

// Function creates log with info message.
func (self *ZLogger) Info(format string, v ...interface{}) {
	self.Lock()
	self.checkMsgPos(firstInfoMsg + format)
	self.logThisStep = append(self.logThisStep, fmt.Sprintf(prefixInfo+format, v...))
	self.Unlock()
}

// Function Step analysys previous step message.
// When there are no warning or error in the block,
// Step always deletes unsaved logs. It saves only
// previous step head with suffix [OK].
// When error is in the block,
// Step always saves all logs.
// When warning is in the block,
// Step saves logs, that were written before warning message.
// Number of this logs euivalents warningLen.

func (self *ZLogger) Step(format string, v ...interface{}) {
	self.Lock()
	self.linesAlloc += len(self.logThisStep)

	// When error is in the block...
	if self.flagError == true {
		// add suffix to step name, copy all logs, clear self variable.
		self.logThisStep[0] = fmt.Sprintf(suffixFormat, self.logThisStep[0], self.sErr)
		self.stringRes = append(self.stringRes, self.logThisStep...)
		self.clearStep()
	}

	// When warning is in the block...
	if len(self.posWarning) > 0 {
		var nfirst, nlast int
		// add suffix to step name
		self.stringRes = append(self.stringRes, fmt.Sprintf(suffixFormat, self.logThisStep[0], self.sWarn))
		for _, pos := range self.posWarning {
			if pos-self.warningLen < nlast {
				nfirst = nlast + 1
			} else {
				nfirst = pos - self.warningLen
			}
			// copy selected logs
			self.stringRes = append(self.stringRes, self.logThisStep[nfirst:pos+1]...)
			nlast = pos
		}
		// clear self variable and return memory to OS.
		self.clearStep()
	}

	// When there are no warning or error in the block...
	if len(self.logThisStep) > 0 {
		// add suffix to step name, clear self variable.
		self.stringRes = append(self.stringRes, fmt.Sprintf(suffixFormat, self.logThisStep[0], self.sOK))
		self.clearStep()
	}

	// clear self variable and return to OS.
	if self.linesAlloc > linesToFreeOsMem {
		runtime.GC()
		debug.FreeOSMemory()
		self.linesAlloc = 0

	}
	self.notFirstMsg = true
	// append step name
	self.logThisStep = append(self.logThisStep, fmt.Sprintf(lineSep+format, v...))
	self.Unlock()
	return
}

// NewStep creates new Step, and return himself(for compatibility).
func (self *ZLogger) NewStep(format string, v ...interface{}) *ZLogger {
	self.Step(fmt.Sprintf(format, v...))
	return self
}

// This clearStep method deletes all unsaved logs.
// Metod clears all flags.
// Variable warningLen desn't change only.
func (self *ZLogger) clearStep() {

	self.flagError = false
	self.notFirstMsg = false
	self.posWarning = []int{}
	self.logThisStep = []string{}
}

// GetLog gets string slice of all logs copied by Step function.
// After that, it removes all logs from memory.
// GetLog saves only argument warningLen.
func (self *ZLogger) GetLog() (msgs []string) {

	self.Step("")
	msgs = self.stringRes
	self.clearStep()
	self.stringRes = []string{}
	//	runtime.GC()
	//	debug.FreeOSMemory()
	return msgs
}

// append logs to the file
func (self *ZLogger) WriteLog() (n int, err error) {
	outPath := self.OutSource
	// verify file existence
	if _, err = os.Stat(outPath); err == nil {
		return self.writeFile("add")
	}
	err = dirPrepare(outPath)
	if err == nil {
		return self.writeFile("rewrite")
	}
	// if can't create dir, writeFile should write ReserveFile.
	return self.writeFile("rewrite")
}

// creates new file and write into it
func (self *ZLogger) ReWriteLog() (n int, err error) {
	err = dirPrepare(self.OutSource)
	return self.writeFile("rewrite")
}

// makes dir and waits until it is has been created
func dirPrepare(outPath string) (err error) {
	dir, _ := filepath.Split(outPath)
	if _, err := os.Stat(dir); err == nil {
		return nil
	}
	err = os.MkdirAll(dir, 0777)
	// waiting for the dir creating
	if err == nil {
		for attempt := 0; attempt < 10; attempt++ {
			if _, err = os.Stat(dir); err == nil {
				return nil
			}
		}
	}
	return err
}

// writeFile should write result into outSource,
// if it can't open it, writeFile should write into ReserveFile.
func (self *ZLogger) writeFile(operation string) (n int, err error) {
	logs := self.GetLog()
	var outFile *os.File

	switch operation {
	case "rewrite":
		outFile, err = os.Create(self.OutSource)
	case "add":
		outFile, err = os.OpenFile(self.OutSource, os.O_APPEND|os.O_WRONLY, 0600)
	}

	// if fail to write, try to write into ReserveFile.
	if err != nil {
		outFile, err = os.OpenFile(self.ReserveFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			outFile, err = os.Create(self.ReserveFile)
			if err != nil {
				return n, err
			}
		}

	}
	defer outFile.Close()
	fmt.Printf("\nLogs is writing into '%s'...", outFile.Name())
	for _, log := range logs {
		w, err := outFile.WriteString(log)
		n += w
		if err != nil {
			fmt.Print(err)
			return n, err
		}
	}
	fmt.Printf("%d bytes. Done.\n", n)
	endLine := fmt.Sprintf(endOutputLine, time.Now())
	_, _ = outFile.WriteString(endLine)
	return n, nil
}

// NewZlog creates new Log, and return himself(for compatibility).
func NewZlog(v ...interface{}) *ZLogger {
	return NewLogger()
}

func (self *ZLogger) GetAllLog() (msgs []string) {
	return self.GetLog()
}

// Method checks Step creation. If msg position is first,
// this method creates and saves name for Step.
func (self *ZLogger) checkMsgPos(msg string) {
	if !self.notFirstMsg {
		// add first message.
		self.logThisStep = append(self.logThisStep, fmt.Sprintf(unknownStepName, msg))
		self.notFirstMsg = true
	}
}

// Function set  warningLen.
func (self *ZLogger) SetWarningLenght(n int) {
	self.warningLen = n
}
