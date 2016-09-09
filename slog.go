package zlog

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	SLWarning       = "[WARNING]: "
	SLError         = "[ERROR]: "
	SLInfo          = "[info]: "
	SLInfoSuffix    = "[ok]"
	SLWarningSuffix = "[warning]"
	SLErrorSuffix   = "[error]"
	SLStep          = "   Step: "
)

// lines sepatator
var SLSeparator string

// SL means simle or strem logger.
type SL struct {
	logs       []string // current logs place
	storage    string   // storage of comresed steps.
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
	self.storage = ""
	defer debug.FreeOSMemory()
	defer runtime.GC()
	return msgs
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

	self.storage += strings.Join(self.logs, SLSeparator) + SLSeparator
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
		filledMsg := fill(self.logs[0], 53)
		self.logs[0] = fmt.Sprintf("%-15s%-53s%s", "   Unknown Step: ", filledMsg, suffix)
		return
	}
	// Filling "................"
	// Do not fill steps with OK.
	filledMsg := self.logs[0]
	if suffix != SLInfoSuffix {
		filledMsg = fill(self.logs[0], 70)
	}

	self.logs[0] = fmt.Sprintf("%-70s%s", filledMsg, suffix)
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
