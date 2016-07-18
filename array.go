package zlog

import (
	//	"fmt"
	"strings"
)

// chain relationship implementing
type ZLA struct {
	isRoot   bool
	parent   *ZLA       // pointer to parent(root)
	order    int        // plase logs in the storage
	children []*ZLA     // pointers to childs
	logs     [][]string // current logs plase
	storage  []string   // storage of comresed steps.
}

type processing interface {
	NewStep() *ZLA
	Info() *ZLA
	GetAllLog() ZLA
	clone() *ZLA
	goRoot() *ZLA
}

func NewZLA() *ZLA {
	return &ZLA{isRoot: true, storage: make([]string, 0)}
}

func (self *ZLA) goRoot() *ZLA {
	if !self.isRoot {
		return self.parent
	}
	return self
}

func (self *ZLA) clone() *ZLA {
	theParent := self.goRoot()
	nextOrder := len(theParent.children) + 1
	// create child
	newNode := &ZLA{
		parent: theParent,
		order:  nextOrder,
	}
	newNode.parent.logs = append(newNode.parent.logs, []string{})
	newNode.parent.storage = append(newNode.parent.storage, string(""))
	theParent.children = append(theParent.children, newNode)
	//	fmt.Println("clone return", newNode.order)
	return newNode
}

func (self *ZLA) processStep() {
	// n is the storage index
	n := self.order
	root := self.parent
	//	fmt.Println("n=", n, len(root.storage), len(root.logs))
	// move all logs to storage
	root.storage[n-1] += strings.Join(root.logs[n-1], "\n")
	root.logs[n-1] = []string{""}
}

func (self *ZLA) NewStep(s string) *ZLA {
	newNode := self.clone()
	newNode.Info("\n     " + s)
	//	fmt.Println("Step return order", newNode.order)
	return newNode
}

func (self *ZLA) Info(s string) {
	//	fmt.Println("Info order", self.order)
	n := self.order
	// when use root as logger
	if n == 0 {
		// root has not children
		if len(self.goRoot().children) == 0 {
			self = self.NewStep("Unknown step")
		} else {
			// save to first logger
			self = self.children[0]
		}
		n = self.order
	}
	self.parent.logs[n-1] = append(self.parent.logs[n-1], s+".")
}
func (self *ZLA) Error(s string) {
	//	fmt.Println("Info order", self.order)
	n := self.order
	// when use root as logger
	if n == 0 {
		// root has not children
		if len(self.goRoot().children) == 0 {
			self = self.NewStep("Unknown step")
		} else {
			// save to first logger
			self = self.children[0]
		}
		n = self.order
	}
	self.parent.logs[n-1] = append(self.parent.logs[n-1], s+".")
}

func (self *ZLA) GetAllLog() string {
	root := self.goRoot()
	for _, node := range root.children {
		node.processStep()
	}
	return strings.Join(root.storage, "/n")
}
