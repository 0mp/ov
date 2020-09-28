package oviewer

import (
	"log"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

// main is manages and executes events in the main routine.
func (root *Root) main() {
	go root.countTimer()

loop:
	for {
		root.draw()
		ev := root.Screen.PollEvent()
		switch ev := ev.(type) {
		case *eventAppQuit:
			if root.input.mode == Help || root.input.mode == LogDoc {
				root.toNormal()
				continue
			}
			break loop
		case *eventTimer:
			root.updateEndNum()
		case *eventDocument:
			root.setDocument(ev.m)
		case *eventCopySelect:
			root.putClipboard()
		case *eventPaste:
			root.getClipboard()
		case *eventSearch:
			root.nextSearch()
		case *eventBackSearch:
			root.nextBackSearch()
		case *searchInput:
			root.search(ev.value)
		case *backSearchInput:
			root.backSearch(ev.value)
		case *gotoInput:
			root.goLine(ev.value)
		case *headerInput:
			root.setHeader(ev.value)
		case *delimiterInput:
			root.setDelimiter(ev.value)
		case *tabWidthInput:
			root.setTabWidth(ev.value)
		case *tcell.EventResize:
			root.resize()
		case *tcell.EventMouse:
			root.mouseEvent(ev)
		case *tcell.EventKey:
			root.setMessage("")
			switch root.input.mode {
			case Normal, Help, LogDoc:
				root.keyCapture(ev)
			default:
				root.inputEvent(ev)
			}
		}
	}
}

func (root *Root) checkScreen() bool {
	return root.Screen != nil
}

// eventAppQuit represents a quit event.
type eventAppQuit struct {
	tcell.EventTime
}

// Quit executes a quit event.
func (root *Root) Quit() {
	if !root.checkScreen() {
		return
	}
	ev := &eventAppQuit{}
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// WriteQuit sets the write flag and executes a quit event.
func (root *Root) WriteQuit() {
	if !root.checkScreen() {
		return
	}
	root.AfterWrite = true
	ev := &eventAppQuit{}
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// eventTimer represents a timer event.
type eventTimer struct {
	tcell.EventTime
}

// runOnTime runs at time.
func (root *Root) runOnTime() {
	if !root.checkScreen() {
		return
	}
	ev := &eventTimer{}
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// countTimer fires events periodically until it reaches EOF.
func (root *Root) countTimer() {
	timer := time.NewTicker(time.Millisecond * 500)
loop:
	for {
		<-timer.C
		root.runOnTime()
		if root.Doc.BufEOF() {
			break loop
		}
	}
	timer.Stop()
}

// MoveLine fires an event that moves to the specified line.
func (root *Root) MoveLine(num int) {
	if !root.checkScreen() {
		return
	}
	ev := &gotoInput{}
	ev.value = strconv.Itoa(num)
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// MoveTop fires the event of moving to top.
func (root *Root) MoveTop() {
	root.MoveLine(0)
}

// MoveBottom fires the event of moving to bottom.
func (root *Root) MoveBottom() {
	root.MoveLine(root.Doc.endNum)
}

// eventSearch represents search event.
type eventSearch struct {
	tcell.EventTime
}

func (root *Root) eventNextSearch() {
	ev := &eventSearch{}
	ev.SetEventNow()
	go func() {
		err := root.Screen.PostEvent(ev)
		if err != nil {
			log.Println(err)
		}
	}()
}

// eventBackSearch represents backward search event.
type eventBackSearch struct {
	tcell.EventTime
}

func (root *Root) eventNextBackSearch() {
	ev := &eventBackSearch{}
	ev.SetEventNow()
	go func() {
		err := root.Screen.PostEvent(ev)
		if err != nil {
			log.Println(err)
		}
	}()
}

// Search fires a forward search event.
// This is for calling Search from the outside.
// Normally, the event is executed from Confirm.
func (root *Root) Search(str string) {
	if !root.checkScreen() {
		return
	}
	if str == "" {
		root.input.reg = nil
		return
	}
	root.input.value = str
	root.input.reg = regexpComple(str, root.CaseSensitive)
	ev := &eventSearch{}
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// BackSearch fires a backward search event.
// This is for calling Search from the outside.
// Normally, the event is executed from Confirm.
func (root *Root) BackSearch(str string) {
	if !root.checkScreen() {
		return
	}
	if str == "" {
		root.input.reg = nil
		return
	}
	root.input.value = str
	root.input.reg = regexpComple(str, root.CaseSensitive)
	ev := &eventBackSearch{}
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}

// eventDocument represents a set model event.
type eventDocument struct {
	m *Document
	tcell.EventTime
}

// SetDocument fires a set document event.
func (root *Root) SetDocument(m *Document) {
	if !root.checkScreen() {
		return
	}
	ev := &eventDocument{}
	ev.m = m
	ev.SetEventNow()
	go func() {
		root.Screen.PostEventWait(ev)
	}()
}
