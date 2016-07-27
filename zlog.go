package zlog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	reserveLogFile = "/tmp/zlog_autosave.log"
)

// It is var for Windows
var (
	linesSep      = "\r\n"
	prefixWarning = "  [warning]: "
	prefixError   = "  [error]  : "
	suffixOK      = "[OK]" + linesSep
	suffixWarning = "[WARNING]" + linesSep
	suffixError   = "[ERROR]" + linesSep

	endOutputLine = "\r\n################ Zlog session ############### %s"
	prefixStep    = linesSep + "Step: "
)

// It is var for all OS
var (
	prefixInfo    = "     [info]:"
	suffixFormat  = "%-65s %s" // indent for [OK][ERROR][WARNING] is 65 points.
)

// chain relationship implementing
type ZL struct {
	sync.Mutex      //
	isRoot          bool
	parent          *ZL              // pointer to parent(root)
	key             int              // key for map child
	nextKey         int              //
	children        map[int]*ZL      // pointers to childs
	logs            map[int][]string // current logs place
	storage         map[int]string   // storage of comresed steps.
	warningPls      []int            // log positions what contains warnings
	errorPls        []int            // log positions what contains errors
	warningLines    int              // number lines before warning to save
	removeBeforeGet bool
	OutSource       string
	processed       bool
}

type pair struct {
	begin, end int
}

func NewZL(out ...interface{}) *ZL {
	curOS := runtime.GOOS
	switch curOS {
	case "linux":
		linesSep = "\n"
		prefixWarning = "  [\033[35mwarning\033[0m]: "
		prefixError = "  [ \033[31merror\033[0m ]: "
		suffixOK = "[\033[32mOK\033[0m]" // + linesSep
		suffixWarning = "[\033[35mWARNING\033[0m]" + linesSep
		suffixError = "[ \033[31mERROR\033[0m ]" + linesSep
		prefixStep = linesSep + "Step: "
	    endOutputLine = linesSep + "################ Zlog session ############### %s"
	case "darwin":
		linesSep = "\n"
	    endOutputLine = linesSep + "################ Zlog session ############### %s"
	}

	zl := &ZL{
		isRoot:          true,
		removeBeforeGet: false,
		warningLines:    10,
		logs:            make(map[int][]string, 0),
		//		logsCounter:        make(map[int]int, 0),
		storage:  make(map[int]string, 0),
		children: make(map[int]*ZL, 0),
	}
	if len(out) > 0 {
		value := reflect.ValueOf(out[0])
		if name := value.String(); len(name) > 0 {
			zl.OutSource = name
			return zl
		}
	}
	zl.OutSource = reserveLogFile
	return zl
}

func (self *ZL) NewStep(format string, v ...interface{}) *ZL {
	root := self.goRoot()
	root.Lock()
	newNode := self.clone()
	n := newNode.key
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixStep+format, v...))
	root.Unlock()
	self = newNode
	return newNode
}

func (self *ZL) Step(format string, v ...interface{}) {
	root := self.goRoot()
	n := self.key
	root.Lock()
	if root.removeBeforeGet {
		self.compress()
	}
	n = self.key
	_ = root.processChild(n)
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixStep+format, v...))
	root.Unlock()
	return
}

func (self *ZL) Info(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.key
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixInfo+format, v...))
	root.Unlock()
}

func (self *ZL) Warning(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.key
	self.warningPls = append(self.warningPls, len(root.logs[n]))
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixWarning+format, v...))
	root.Unlock()
}

func (self *ZL) Error(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.key
	self.errorPls = append(self.errorPls, len(root.logs[n]))
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixError+format, v...))
	root.Unlock()
}

func (self *ZL) GetStep() string {
	return self.getLog()
}

func (self *ZL) GetAllLog() string {
	root := self.goRoot()
	return root.getLog()
}

func (self *ZL) getLog() string {
	root := self.goRoot()
	root.Lock()
	msgs := []string{}
	keys := []int{}

	if self.isRoot {
		// get all by sorted keys
		keys = append(keys, 0)
		for key := range root.children {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		if root.removeBeforeGet {
			root.compressAll()
		}

	} else {
		// get by one key
		key := self.key
		keys = append(keys, key)
		if root.removeBeforeGet {
			self.compress()
		}
	}

	for _, key := range keys {
		_ = root.processChild(key)
		// check map on key existense
		if msg, e := root.storage[key]; e == true {
			msgs = append(msgs, msg)
			root.storage[key] = string("")
		}
	}

	defer debug.FreeOSMemory()
	defer runtime.GC()

	root.Unlock()
	return strings.Join(msgs, "")
}

// Checks logs on attentions messages.
// If *ZL do not contains logs with attentions,
// removes logs.
func (self *ZL) compress() {
	n := self.key
	// check on attentions.
	withoutAttention := len(self.errorPls) == 0 && len(self.warningPls) == 0
	root := self.goRoot()
	// means no errors, warnings and the storage is empty
	if withoutAttention && len(root.storage[n]) == 0 {
		root.logs[self.key] = []string{}
		//		delete(root.children, n)
		//		delete(root.logs, n)
		//		delete(root.storage, n)
	}
}

func (self *ZL) compressAll() {
	root := self.goRoot()
	for _, zl := range root.children {
		zl.compress()
	}
	return
}

func (self *ZL) SetWarningLines(n int) {
	self.goRoot().warningLines = n
}

func (self *ZL) SetRemoveBeforeGet(s bool) {
	self.goRoot().removeBeforeGet = s
}

func (self *ZL) goRoot() *ZL {
	if !self.isRoot {
		return self.parent
	}
	return self
}

func (self *ZL) clone() *ZL {
	theParent := self.goRoot()
	theParent.nextKey++
	key := theParent.nextKey
	// create child
	newNode := &ZL{
		parent: theParent,
		key:    key,
	}
	theParent.logs[key] = []string{}
	theParent.storage[key] = string("")
	theParent.children[key] = newNode
	return newNode
}

// returns range of lines, what should be saved with Warning log
func getRange(warningPls []int, warningLines int) (pairs []pair) {
	var saved, begin int
	for _, p := range warningPls {
		if p-warningLines < saved {
			begin = saved + 1
		} else {
			begin = p - warningLines
		}
		pairs = append(pairs, pair{begin, p + 1})
		saved = p
	}
	return pairs
}

// moves important msgs from logs to storage
// returns code.
func (root *ZL) processChild(n int) bool {
	var ePls, wPls *[]int

	if n == 0 {
		ePls = &root.errorPls
		wPls = &root.warningPls
	} else {
		root = root.goRoot()
		ePls = &root.children[n].errorPls
		wPls = &root.children[n].warningPls
	}

	// there are no logs
	if len(root.logs[n]) == 0 {
		return false
	}

	// check error existence
	if len(*ePls) != 0 {
		// make step header
		root.makeCaption(n, suffixError)
		// copy data to the storage
		root.storage[n] += strings.Join(root.logs[n][1:], linesSep)
		// claear logs
		root.logs[n] = []string{}
		*ePls = []int{}
		*wPls = []int{}
		return true
	}
	// check warning existence
	if len(*wPls) != 0 {
		// save step header
		root.makeCaption(n, suffixWarning)
		// get lines ranges to save
		saveRange := getRange(*wPls, root.warningLines)
		msgs := []string{}
		// join lines in blocks
		for _, pairs := range saveRange {
			s, e := pairs.begin, pairs.end
			msgs = append(msgs, strings.Join(root.logs[n][s:e], linesSep))
		}
		// join all blocks in one string
		root.storage[n] += strings.Join(msgs, linesSep)
		// clear logs
		root.logs[n] = []string{}
		*wPls = []int{}
		return true
	}
	// at the end, if all OK. I make step header:
	root.makeCaption(n, suffixOK)
	// clear logs
	root.logs[n] = []string{}
	return false
}

func (root *ZL) makeCaption(n int, suffix string) {
	caption := root.logs[n][0]
	if !strings.Contains(caption, prefixStep) {
		caption = prefixStep +"(unknown)"+ caption
	}
	root.storage[n] += fmt.Sprintf(suffixFormat, caption, suffix)
}

// append logs to the file
func (self *ZL) WriteAllLog() (n int, err error) {
	outPath := self.goRoot().OutSource
	// verify file existence
	if _, err = os.Stat(outPath); err == nil {
		return self.writeFile(self.GetAllLog(), "add")
	}
	_ = dirPrepare(outPath)
	return self.writeFile(self.GetAllLog(), "rewrite")
}

// append logs to the file
func (self *ZL) WriteStep() (n int, err error) {
	outPath := self.goRoot().OutSource
	// verify file existence
	if _, err = os.Stat(outPath); err == nil {
		return self.writeFile(self.GetStep(), "add")
	}
	_ = dirPrepare(outPath)
	return self.writeFile(self.GetStep(), "rewrite")
}

// creates new file and write start line into it.
func (self *ZL) CreateLogFile() (n int, err error) {
	err = dirPrepare(self.goRoot().OutSource)
	startLine := fmt.Sprintf(endOutputLine, time.Now())[:65] + linesSep
	return self.writeFile(startLine, "rewrite")
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
		for attempt := 0; attempt < 20; attempt++ {
			if _, err = os.Stat(dir); err == nil {
				return nil
			}
		}
	}
	return err
}

// writeFile should write result into outSource,
// if it can't open it, writeFile should write into ReserveFile.
func (self *ZL) writeFile(logs, operation string) (n int, err error) {
	//	logs := self.GetAllLog()
	var outFile *os.File
	fileName := self.goRoot().OutSource
	switch operation {
	case "rewrite":
		outFile, err = os.Create(fileName)
	case "add":
		outFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	}
	// if fail to write, try to write into ReserveFile.
	if err != nil {
		outFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			outFile, err = os.Create(fileName)
			if err != nil {
				return 0, err
			}
		}

	}
	defer outFile.Close()
	fmt.Printf("Logs is writing into '%s'...", outFile.Name())
	w, err := outFile.WriteString(logs)
	if err != nil {
		return 0, err
	}
	fmt.Printf("%d bytes. Done.\n", w)
	return w, nil
}
