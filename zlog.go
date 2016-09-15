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
	// will save logs in reserveLogFile, when NewZL have not an argument,
	// or saving failed into this argument.
	reserveLogFile = "/tmp/zlog_autosave.log"
)

type ZL struct {
	sync.Mutex
	isRoot          bool
	parent          *ZL              // pointer to parent(root)
	key             int              // key for map child
	nextKey         int              //
	children        map[int]*ZL      // pointers to childs
	logs            map[int][]string // current logs place
	storage         map[int][]string // storage of comresed steps.
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
	zl := &ZL{
		isRoot:          true,
		removeBeforeGet: false,
		warningLines:    10,
		logs:            make(map[int][]string, 0),
		storage:         make(map[int][]string, 0),
		children:        make(map[int]*ZL, 0),
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
	return self.getLog(false)
}

func (self *ZL) GetAllLog() string {
	root := self.goRoot()
	return root.getLog(true)
}

func (self *ZL) getLog(all bool) string {
	root := self.goRoot()
	root.Lock()
	msgs := []string{}
	keys := []int{}

	if all {
		// get all by sorted keys
		keys = append(keys, 0)
		for theKey := range root.children {
			keys = append(keys, theKey)
		}
		sort.Ints(keys)
	} else {
		// get by one key
		key := self.key
		keys = append(keys, key)
	}

	for _, key := range keys {
		_ = root.processChild(key)

		// Add newline ??
		if msg, e := root.storage[key]; e == true {
			msgs = append(msgs, msg...)
			root.storage[key] = []string{}
		}
	}

	defer debug.FreeOSMemory()
	defer runtime.GC()

	root.Unlock()
	if len(msgs) == 0 {
		return ""
	}
	return strings.Join(msgs, linesSep) + linesSep
}

func (self *ZL) SetWarningLines(n int) {
	self.goRoot().warningLines = n
}

func (self *ZL) goRoot() *ZL {
	if self.isRoot {
		return self
	}
	return self.parent
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
	theParent.storage[key] = []string{}
	theParent.children[key] = newNode
	return newNode
}

// returns range of lines, what should be saved with Warning log
func getWarningRange(warningPls []int, warningLines int) (pairs []pair) {
	var saved, begin int
	pairs = append(pairs, pair{0, 1})
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
	var errLines, warnLines *[]int

	// there are no logs
	if len(root.logs[n]) == 0 {
		return false
	}
	// define lines number with error and warning in
	if n == 0 {
		errLines = &root.errorPls
		warnLines = &root.warningPls
	} else {
		root = root.goRoot()
		errLines = &root.children[n].errorPls
		warnLines = &root.children[n].warningPls
	}

	// check error existence
	if len(*errLines) != 0 {
		// make step header
		root.makeCaption(n, suffixError)
		// copy data to the storage
		root.storage[n] = append(root.storage[n], root.logs[n]...)
		// claear logs
		root.logs[n] = []string{}
		*errLines = []int{}
		*warnLines = []int{}
		return true
	}
	// check warning existence
	if len(*warnLines) != 0 {
		// save step header
		root.makeCaption(n, suffixWarning)
		// get lines ranges to save
		warningSaveRanges := getWarningRange(*warnLines, root.warningLines)
		// join lines in blocks
		for _, pair := range warningSaveRanges {
			root.storage[n] = append(root.storage[n], root.logs[n][pair.begin:pair.end]...)
		}
		// clear logs
		root.logs[n] = []string{}
		*warnLines = []int{}
		return true
	}

	// at the end, if all OK. I make step header:
	root.makeCaption(n, suffixOK)
	if len(root.logs[n]) > 0 {
		root.storage[n] = append(root.storage[n], root.logs[n][0])
		// clear logs
		root.logs[n] = []string{}
	}
	return false
}

// It makes the step header and save it.
// If step was not defined, marks it 'Unknown step'
func (root *ZL) makeCaption(n int, suffix string) {
	if len(root.logs[n]) == 0 {
		return
	}
	// when step header isn't exist
	if !strings.Contains(root.logs[n][0], prefixStep) {
		// add step header to storage with first msg included in
		//root.logs[n][0] = fmt.Sprintf("%-15s%-52s%s"+linesSep+"%s", + "" after "   UnknownStep"
		root.logs[n][0] = fmt.Sprintf("%-69s%s"+linesSep+"%s", "    Unknown Step: ", suffix, root.logs[n][0])
		return
	}
	// add step header only
	root.logs[n][0] = fmt.Sprintf("%-69s%s", root.logs[n][0], suffix)
}

// It writes (appends) logs into the file,
// what was argument the NewZL() function.
// Else function saves into reserveLogFile.
func (self *ZL) WriteAllLog() (n int, err error) {
	outPath := self.goRoot().OutSource
	// verify file existence
	if _, err = os.Stat(outPath); err == nil {
		return self.writeFile(self.GetAllLog(), "add")
	}
	_ = dirPrepare(outPath)
	return self.writeFile(self.GetAllLog(), "rewrite")
}

// Function appends current step logs to the file.
// It clear only this step.
func (self *ZL) WriteStep() (n int, err error) {
	outPath := self.goRoot().OutSource
	// verify file existence
	if _, err = os.Stat(outPath); err == nil {
		return self.writeFile(self.GetStep(), "add")
	}
	_ = dirPrepare(outPath)
	return self.writeFile(self.GetStep(), "rewrite")
}

// It creates(rewrites) new file and write start line into it.
func (self *ZL) CreateLogFile() (n int, err error) {
	err = dirPrepare(self.goRoot().OutSource)
	startLine := fmt.Sprintf("Started at %v %s", time.Now(), linesSep)
	return self.writeFile(startLine, "rewrite")
}

// It makes dir and waits until it is has been created
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

// Functions writeFile should write result into NewZL() argument,
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
