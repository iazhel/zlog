package zlog

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	SLWarning       = "[WARNING]: "
	SLError         = "[ERROR]: "
	SLInfo          = "[info]: "
	SLInfoSuffix    = "[ok]"
	SLWarningSuffix = "[warning]"
	SLErrorSuffix   = "[error]"
	SLStep          = "    Step: "
)

// lines sepatator
var SLSep string

// SL means simle or strem logger.
type SL struct {
	sync.Mutex
	logs         []string // current logs place
	storage      string   // storage of comresed steps.
	warningPls   []int    // log positions what contains warnings
	errorPls     []int    // log positions what contains errors
	warningLines int      // number lines before warning to save
	OutSource    string
}

type sPair struct {
	begin, end int
}

func NewSL(out ...interface{}) *SL {

	SLSep = "\r\n"
	if runtime.GOOS != "windows" {
		SLSep = "\n"
	}

	sl := &SL{
		warningLines: 10,
		logs:         make([]string, 0),
	}
	if len(out) > 0 {
		value := reflect.ValueOf(out[0])
		if name := value.String(); len(name) > 0 {
			sl.OutSource = name
		}
	}
	return sl
}

func (self *SL) Step(format string, v ...interface{}) {
	self.processChild()
	//	self.add(SLStep, v...)
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
	self.processChild()
	msgs := self.storage
	self.storage = ""
	defer debug.FreeOSMemory()
	defer runtime.GC()
	return msgs
}

// returns range of lines, what should be saved with Warning log
func sRange(warningPls []int, warningLines int) (pairs []sPair) {
	var saved, begin int
	for _, p := range warningPls {
		if p-warningLines < saved {
			begin = saved + 1
		} else {
			begin = p - warningLines
		}
		pairs = append(pairs, sPair{begin, p + 1})
		saved = p
	}
	return pairs
}

// moves important msgs from logs to storage
// returns code.
func (self *SL) processChild() {

	if len(self.logs) == 0 {
		return
	}

	// check error existence
	if len(self.errorPls) != 0 {
		// make step header
		self.makeCaption(SLErrorSuffix)
		// copy data to the storage
		self.storage += strings.Join(self.logs[1:], SLSep)
		// claear logs
		self.logs = []string{}
		self.errorPls = []int{}
		self.warningPls = []int{}
		return
	}

	// check warning existence
	if len(self.warningPls) != 0 {
		// save step header
		self.makeCaption(SLWarningSuffix)
		// get lines ranges to save
		saveRange := sRange(self.warningPls, self.warningLines)
		msgs := []string{}
		// join lines in blocks
		for _, pairs := range saveRange {
			s, e := pairs.begin, pairs.end
			msgs = append(msgs, strings.Join(self.logs[s:e], SLSep))
		}
		// join all blocks in one string
		self.storage += strings.Join(msgs, SLSep)
		// clear logs
		self.logs = []string{}
		self.warningPls = []int{}
		return
	}
	// at the end, if all OK. I make step header:
	self.makeCaption(SLInfoSuffix)
	// clear logs
	self.logs = []string{}
	return
}

// It makes the step header and save it in storage.
// If step was not defined, marks it 'Unknown step'
// and always save first message in unknown step.
func (self *SL) makeCaption(suffix string) {

	if len(self.logs) == 0 {
		return
	}
	firstLog := self.logs[0]
	// for unknown steps
	formatU := SLSep + "%-15s%-65s%-15s"
	// for known steps
	formatK := SLSep + "%-80s%-15s"

	// when step header isn't exist
	if !strings.Contains(firstLog, SLStep) {
		prefixUnknown := "Unknown Step"
		switch suffix {
		case SLInfoSuffix:
			// add step header to storage with first msg included in
			self.storage += fmt.Sprintf(formatU, prefixUnknown, firstLog, suffix)
		default:
			// add step header to storage
			self.storage += fmt.Sprintf(formatU, prefixUnknown, fill("", 65), suffix)
			self.storage += SLSep
			// add the first message
			self.storage += firstLog
		}
		return
	}

	// add step header only
	if suffix != SLInfoSuffix {
		firstLog = fill(firstLog, 80)
		self.storage += fmt.Sprintf(formatK, firstLog, suffix)
		self.storage += SLSep
		return
	}

	self.storage += fmt.Sprintf(formatK, firstLog, suffix)
	return

}

func fill(line string, n int) string {
	// bad step lineFiller
	badFiller := " ....................................................................."
	if f := n - len(line); f > 5 && f < len(badFiller) {
		return line + badFiller[:f]
	}
	return line + badFiller

}
