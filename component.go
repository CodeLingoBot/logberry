package logberry

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type Component struct {
	UID    uint64
	Parent Context
	Root   *Root
	Label  string

	Class ComponentClass

	mute bool
	highlight bool
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func newcomponent(parent Context, label string, data ...interface{}) *Component {

	var class = COMPONENT
	if data != nil && len(data) > 0 {
		if cc, ok := data[0].(ComponentClass); ok {
			class = cc
			data = data[1:]
		}
	}
	d := DAggregate(data)

	c := &Component{
		UID:    newcontextuid(),
		Parent: parent,
		Label:  label,
		Class:  class,
	}

	if parent != nil {
		c.Root = parent.GetRoot()
		d.Set(c.Root.FieldPrefix+"Parent", parent.GetUID())
	} else {
		c.Root = Std
	}

	if c.Class < 0 || c.Class >= componentclass_sentinel {
		c.Root.InternalError(NewError("ComponentClass out of range", c.UID, c.Class))
		d.Set(c.Root.FieldPrefix+"Class", c.Class)
	} else {
		d.Set(c.Root.FieldPrefix+"Class", ComponentClassText[c.Class])
	}

	if parent != nil {
		c.Root.ComponentEvent(c, COMPONENT_START, "Instantiate", d)
	}

	return c

}

func (x *Component) Component(label string, data ...interface{}) *Component {
	return newcomponent(x, label, data...)
}

func (x *Component) Task(activity string, data ...interface{}) *Task {
	return newtask(x, activity, data)
}

//----------------------------------------------------------------------
func (x *Component) GetLabel() string {
	return x.Label
}

func (x *Component) GetUID() uint64 {
	return x.UID
}

func (x *Component) GetParent() Context {
	return x.Parent
}

func (x *Component) GetRoot() *Root {
	return x.Root
}


func (x *Component) Mute() *Component {
	x.mute = true
	return x
}
func (x *Component) Unmute() *Component {
	x.mute = false
	return x
}
func (x *Component) IsMute() bool {
	return x.mute
}


func (x *Component) Highlight() *Component {
	x.highlight = true
	return x
}

func (x *Component) ClearHighlight() *Component {
	x.highlight = false
	return x
}

func (x *Component) IsHighlighted() bool {
	return x.highlight
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Component) Build(build *BuildMetadata) {
	x.Root.ComponentEvent(x, COMPONENT_CONFIGURATION, "Build", DBuild(build))
}

func (x *Component) Configuration(data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_CONFIGURATION, "Configuration", DBuild(data))
}

func (x *Component) CommandLine() {

	hostname, err := os.Hostname()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve hostname"))
		return
	}

	u, err := user.Current()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve user info"))
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve program path"))
		return
	}

	prog := path.Base(os.Args[0])

	d := D{
		"Host":    hostname,
		"User":    u.Username,
		"Path":    dir,
		"Program": prog,
		"Args":    os.Args[1:],
	}

	x.Root.ComponentEvent(x, COMPONENT_CONFIGURATION, "Command line", &d)

}

func (x *Component) Environment() {

	d := D{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		d[pair[0]] = pair[1]
	}
	x.Root.ComponentEvent(x, COMPONENT_CONFIGURATION, "Environment", &d)

}

func (x *Component) Process() {

	hostname, err := os.Hostname()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve hostname"))
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve working dir"))
		return
	}

	u, err := user.Current()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve user info"))
		return
	}

	d := D{
		"Host": hostname,
		"WD":   wd,
		"UID":  u.Uid,
		"User": u.Username,
		"PID":  os.Getpid(),
	}

	x.Root.ComponentEvent(x, COMPONENT_CONFIGURATION, "Process", &d)

}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Component) Info(msg string, data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_INFO, msg, DAggregate(data))
}

func (x *Component) Recovered(msg string, err error, data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_WARNING, msg,
		DAggregate(data).Set(x.Root.FieldPrefix+"Error", err.Error()))
}

func (x *Component) Warning(msg string, data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_WARNING, msg, DAggregate(data))
}

func (x *Component) Error(msg string, err error, data ...interface{}) error {

	// Note that this can't/shouldn't just throw err into the data blob
	// because the standard errors package error doesn't expose
	// anything, even the message.  So you basically have to reduce to a
	// string via Error().

	x.Root.ComponentEvent(x, COMPONENT_ERROR, msg,
		DAggregate(data).Set(x.Root.FieldPrefix+"Error", err.Error()))
	return WrapError(err, msg)

}

// Failure is the same as Error but doesn't take an error object.
func (x *Component) Failure(msg string, data ...interface{}) error {

	x.Root.ComponentEvent(x, COMPONENT_ERROR, msg, DAggregate(data))
	return NewError(msg)

}

// Generally only the top level should invoke fatal, not sub-components.
func (x *Component) Fatal(msg string, err error, data ...interface{}) {

	x.Root.ComponentEvent(x, COMPONENT_FATAL, msg,
		DAggregate(data).Set(x.Root.FieldPrefix+"Error", err.Error()))
	os.Exit(1)

}

func (x *Component) Abort(msg string, data ...interface{}) {

	x.Root.ComponentEvent(x, COMPONENT_FATAL, msg, DAggregate(data))
	os.Exit(1)

}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Component) Finalize(data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_FINISH, "Finalize", DAggregate(data))
}

func (x *Component) Ready(msg string, data ...interface{}) {
	x.Root.ComponentEvent(x, COMPONENT_READY, msg, DAggregate(data))
}
