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
	// LINUX
	lineSep        = "\n"
	reserveLogFile = "/tmp/zlog_autosave.log"

	// WINDOWS
	//lineSep           = "\r\n"
	//reserveLogFile = "c:\\zlog_autosave.log"

)

// chain relationship implementing
type ZL struct {
	sync.Mutex //
	isRoot     bool
	parent     *ZL              // pointer to parent(root)
	key        int              // key for map child
	nextKey    int              //
	children   map[int]*ZL      // pointers to childs
	logs       map[int][]string // current logs place
	//	logsCounter        map[int]int      // logs counter
	storage         map[int]string // storage of comresed steps.
	warningPls      []int          // log positions what contains warnings
	errorPls        []int          // log positions what contains errors
	warningLines    int            // number lines before warning to save
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
	n := newNode.savePoint()
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixStep+format, v...))
	root.Unlock()
	return newNode
}

func (self *ZL) Step(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.key
	if root.removeBeforeGet {
		self.compress()
	}
	_ = root.processChild(n)
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixStep+format, v...))
	root.Unlock()
	return
}

func (self *ZL) Info(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.savePoint()
	// when msgs coount == 0
	if len(root.logs[n]) == 0 {
		format = lineSep + "Unknown step: " + format
		root.logs[n] = append(root.logs[n], fmt.Sprintf(format, v...))
	}
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixInfo+format, v...))
	root.Unlock()
}

func (self *ZL) Warning(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.savePoint()
	self.warningPls = append(self.warningPls, len(root.logs[n]))
	root.logs[n] = append(root.logs[n], fmt.Sprintf(prefixWarning+format, v...))
	root.Unlock()
}

func (self *ZL) Error(format string, v ...interface{}) {
	root := self.goRoot()
	root.Lock()
	n := self.savePoint()
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

	if self.isRoot {
		if root.removeBeforeGet {
			root.compressAll()
		}
		// get all by sorted keys
		keys := []int{}
		for key := range root.children {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		for _, key := range keys {
			_ = root.processChild(key)
			if msg, e := root.storage[key]; e == true {
				msgs = append(msgs, msg)
				root.storage[key] = string("")
			}
		}
	} else { // get by one key
		key := self.key
		if root.removeBeforeGet {
			self.compress()
		}
		msg, e := root.storage[key]
		if e == true {
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
// If *ZL don't contains logs with attentions,
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

// Do like compress
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
// returns moveing code.
func (root *ZL) processChild(n int) bool {

	if !root.isRoot {
		root = root.goRoot()
	}
	// there are no logs
	if len(root.logs[n]) == 0 {
		return false
	}
	// check error existence
	if len(root.children[n].errorPls) != 0 {
		// make step header
		root.storage[n] += fmt.Sprintf(suffixFormat, root.logs[n][0], suffixError)
		// copy data to the storage
		root.storage[n] += strings.Join(root.logs[n][1:], lineSep)
		// claear logs
		root.logs[n] = []string{}
		root.children[n].errorPls = []int{}
		root.children[n].warningPls = []int{}
		return true
	}
	// check warning existence
	if len(root.children[n].warningPls) != 0 {
		// save step header
		root.storage[n] += fmt.Sprintf(suffixFormat, root.logs[n][0], suffixWarning)
		// get lines ranges to save
		saveRange := getRange(root.children[n].warningPls, root.warningLines)
		msgs := []string{}
		// join lines in blocks
		for _, pairs := range saveRange {
			s, e := pairs.begin, pairs.end
			msgs = append(msgs, strings.Join(root.logs[n][s:e], lineSep))
		}
		// join all blocks in one string
		root.storage[n] += strings.Join(msgs, lineSep)
		// clear logs
		root.logs[n] = []string{}
		root.children[n].warningPls = []int{}
		return true
	}
	// at the end, if all OK. I make step header:
	root.storage[n] += fmt.Sprintf(suffixFormat, root.logs[n][0], suffixOK)
	// clear logs
	root.logs[n] = []string{}
	return false
}

// returns key for logs & storage.
func (self *ZL) savePoint() int {
	saveKey := self.key
	root := self.goRoot()

	// when is used root as enter point
	if saveKey == 0 {
		// and root has not children
		if len(root.children) == 0 {
			// add new step
			_ = self.NewStep("Unknown New step")
			fmt.Println("savePoint: Unknown New step!")
		}
		// search first child key
		for key := range self.children {
			saveKey = key
			fmt.Println("savePoint: KEY from cycle!!!", key)
			break
		}
	}
	return saveKey
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
	startLine := fmt.Sprintf(endOutputLine, time.Now())[:65] + lineSep
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
		//		fmt.Println("reWrite", err, fileName)
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
		//		fmt.Print(err)
		return 0, err
	}
	fmt.Printf("%d bytes. Done.\n", w)
	return w, nil
}
