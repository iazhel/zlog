package zlog

// type SL means SIMPLE or ONE STREAM logger.

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	SLWarning       = "[WARNING]: "
	SLError         = "[ERROR]: "
	SLInfo          = "[info]: "
	SLInfoSuffix    = " [ok]"
	SLWarningSuffix = " [WARNING]"
	SLErrorSuffix   = " [ERROR]"
	SLStep          = "    Step: "
)

// lines sepatator
var SLSeparator string

type SL struct {
	logs       []string // current logs place
	storage    []string // storage of comresed steps.
	warningPls []int    // log positions what contains warnings
	errorPls   []int    // log positions what contains errors
}

func NewSL(out ...interface{}) *SL {
	SLSeparator = "\r\n"
	if runtime.GOOS != "windows" {
		SLSeparator = "\n"
	}
	return &SL{
		logs: make([]string, 0),
	}
}

func (self *SL) Step(format string, v ...interface{}) {
	self.processLogs()
	self.logs = append(self.logs, fmt.Sprintf(SLStep+format, v...))
	return
}

func (self *SL) Info(format string, v ...interface{}) {
	self.logs = append(self.logs, fmt.Sprintf(SLInfo+format, v...))
}

func (self *SL) Warning(format string, v ...interface{}) {
	self.warningPls = append(self.warningPls, len(self.logs))
	self.logs = append(self.logs, fmt.Sprintf(SLWarning+format, v...))
}

func (self *SL) Error(format string, v ...interface{}) {
	self.errorPls = append(self.errorPls, len(self.logs))
	self.logs = append(self.logs, fmt.Sprintf(SLError+format, v...))
}

func (self *SL) GetAllLog() string {
	return self.GetLog()
}

func (self *SL) GetLog() string {
	self.processLogs()
	msgs := self.storage
	self.storage = []string{}
	defer debug.FreeOSMemory()
	defer runtime.GC()
	if len(msgs) == 0 {
		return ""
	}
	return strings.Join(msgs, SLSeparator) + SLSeparator
}

func (self *SL) processLogs() {
	if len(self.logs) == 0 {
		return
	}
	switch {
	case len(self.errorPls) != 0:
		self.makeCaption(SLErrorSuffix)
	case len(self.warningPls) != 0:
		self.makeCaption(SLWarningSuffix)
	default:
		self.makeCaption(SLInfoSuffix)
	}

	self.storage = append(self.storage, self.logs...)
	// clear logs
	self.logs = []string{}
	self.errorPls = []int{}
	self.warningPls = []int{}
	return
}

// It makes the step header and save it.
// If step was not defined, marks it 'Unknown step'
func (self *SL) makeCaption(suffix string) {
	if len(self.logs) == 0 {
		return
	}
	// when step header isn't exist
	if !strings.Contains(self.logs[0], SLStep) {
		// add step header to storage with first msg included in
		filledMsg := fill("    Unknown Step: ", 69)
		self.logs[0] = fmt.Sprintf("%-69s%s"+SLSeparator+"%s", filledMsg, suffix, self.logs[0])
		return
	}
	// Do not fill steps with OK.
	if suffix == SLInfoSuffix {
		self.logs[0] = fmt.Sprintf("%-69s%s", self.logs[0], suffix)
		return
	}

	filledMsg := fill(self.logs[0], 69)
	self.logs[0] = fmt.Sprintf("%-69s%s", filledMsg, suffix)
	return
}

func fill(line string, n int) string {
	filler := " ....................................................................."
	f := n - len(line)
	switch {
	case f <= 0:
		return line
	case f > 5 && f < len(filler):
		return line + filler[:f]
	case f >= len(filler):
		return line + filler
	}
	return line
}

func (self *SL) WriteLog(filePath string) (n int, err error) {
	return self.Write(self.GetAllLog(), filePath)
}

func (self *SL) Write(text, filePath string) (n int, err error) {
	dir, file := path.Split(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return 0, err
	}
	f, err := os.OpenFile(path.Join(dir, file), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	defer f.Close()
	if err != nil {
		return 0, err
	}
	b, err := f.WriteString(text)
	if err != nil {
		return 0, err
	}
	return b, err

}
