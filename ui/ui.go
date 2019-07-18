package ui

import (
	"fmt"
	"gosrg/config"
	"gosrg/utils"

	"github.com/jessewkun/gocui"
)

var Ui UI

type GHandler interface {
	Layout(g *gocui.Gui) error
	initialize() error
	bindShortCuts() error
	focus(arg ...interface{}) error
	clear() error
	setCurrent(v GHandler, arg ...interface{}) error
	getCurrentLine() string
	output(arg ...interface{}) error
	outputln(arg ...interface{}) error
}

type UI struct {
	G        *gocui.Gui
	AllView  map[string]GHandler
	TabNo    int
	NextView GHandler
}

type GView struct {
	Name      string
	Title     string
	View      *gocui.View
	ShortCuts []ShortCut
}

const (
	GLOBAL_N = iota // global and cannnot unbind
	GLOBAL_Y        // global and can unbind
	LOCAL_N         // local and cannot unbind
	LOCAL_Y         // local and can unbind
)

type ShortCut struct {
	Key     interface{}
	Level   int
	Handler func(*gocui.Gui, *gocui.View) error
}

type ButtonWidget struct {
	Name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

var confirmBtn *ButtonWidget
var cancelBtn *ButtonWidget

func NewButtonWidget(name string, x, y int, label string, handler func(g *gocui.Gui, v *gocui.View) error) *ButtonWidget {
	return &ButtonWidget{Name: name, x: x, y: y, w: len(label) + 1, label: label, handler: handler}
}

func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	if v, err := Ui.G.SetView(w.Name, w.x, w.y, w.x+w.w, w.y+2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		Ui.G.Cursor = false
		if err := Ui.G.SetKeybinding(w.Name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return err
		}
		fmt.Fprint(v, w.label)
	}
	return nil
}

func (gv *GView) Layout(g *gocui.Gui) error {
	return nil
}

func (gv *GView) bindShortCuts() error {
	for _, sc := range gv.ShortCuts {
		vName := gv.Name
		if sc.Level == GLOBAL_Y || sc.Level == GLOBAL_N {
			vName = ""
		}
		if err := Ui.G.SetKeybinding(vName, sc.Key, gocui.ModNone, sc.Handler); err != nil {
			utils.Logger.Fatalln(err)
			return err
		}
	}
	return nil
}

func (gv *GView) unbindShortCuts() error {
	for _, sc := range gv.ShortCuts {
		if sc.Level == GLOBAL_N || sc.Level == LOCAL_N {
			continue
		}
		vName := gv.Name
		if sc.Level == GLOBAL_Y {
			vName = ""
		}
		if err := Ui.G.DeleteKeybinding(vName, sc.Key, gocui.ModNone); err != nil {
			utils.Logger.Fatalln(err)
			return err
		}
	}
	return nil
}

func (gv *GView) initialize() error {
	return nil
}

func (gv *GView) focus(arg ...interface{}) error {
	Ui.G.Cursor = false
	if tip, ok := config.TipsMap[gv.Name]; ok {
		tView.output(tip)
	} else {
		tView.output("")
	}
	return nil
}

func (gv *GView) clear() error {
	gv.View.Clear()
	return nil
}

func (gv *GView) output(arg ...interface{}) error {
	if _, err := fmt.Fprint(gv.View, arg...); err != nil {
		utils.Logger.Fatalln(err)
		return err
	}

	return nil
}

func (gv *GView) outputln(arg ...interface{}) error {
	if _, err := fmt.Fprintln(gv.View, arg...); err != nil {
		utils.Logger.Fatalln(err)
		return err
	}

	return nil
}

func (gv *GView) setCurrent(v GHandler, arg ...interface{}) error {
	if _, err := Ui.G.SetCurrentView(gv.Name); err != nil {
		utils.Logger.Fatalln(err)
		return err
	}
	v.focus(arg...)
	return nil
}

func setNextView() {
	Ui.TabNo++
	next := Ui.TabNo % len(config.TabView)
	Ui.NextView = Ui.AllView[config.TabView[next]]
}

func (gv *GView) getCurrentLine() string {
	var line string
	var err error

	_, cy := gv.View.Cursor()
	if line, err = gv.View.Line(cy); err != nil {
		utils.Logger.Println(err)
		return ""
	}
	return line
}

func (gv *GView) deleteCursorLine() error {
	_, cy := gv.View.Cursor()
	return gv.deleteLine(cy)
}

func (gv *GView) deleteLine(y int) error {
	if err := gv.View.DeleteLine(y); err != nil {
		utils.Logger.Fatalln(err)
		return err
	}
	return nil
}

func (gv *GView) cursorUp() error {
	ox, oy := gv.View.Origin()
	cx, cy := gv.View.Cursor()
	if err := gv.View.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := gv.View.SetOrigin(ox, oy-1); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (gv *GView) cursorDown() error {
	cx, cy := gv.View.Cursor()
	lineHeight := gv.View.LinesHeight()
	if cy+1 >= lineHeight {
		return nil
	}
	if err := gv.View.SetCursor(cx, cy+1); err != nil {
		ox, oy := gv.View.Origin()
		if err := gv.View.SetOrigin(ox, oy+1); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (gv *GView) command(str string) {
	gv.outputln(utils.Now() + utils.Bule("[COMMAND]") + str)
}

func (gv *GView) info(str string) {
	gv.outputln(utils.Now() + utils.Tianqing("[INFO]") + str)
}

func (gv *GView) res(str string) {
	gv.outputln(utils.Now() + utils.Green("[RESULT]") + str)
}

func (gv *GView) error(str string) {
	gv.outputln(utils.Now() + utils.Red("[ERROR]") + str)
}

func (gv *GView) debug(arg ...interface{}) {
	if config.DEBUG {
		arg = append([]interface{}{utils.Now() + utils.Orange("[DEBUG]")}, arg...)
		gv.outputln(arg...)
	}
}

func (gv *GView) cursorBegin() error {
	if err := gv.View.SetCursor(0, 0); err != nil {
		if err := gv.View.SetOrigin(0, 0); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (gv *GView) cursorLast() error {
	lineHeight := gv.View.LinesHeight()
	lineHeight--
	lastLine, _ := gv.View.Line(lineHeight)
	lastLineWidth := len(lastLine)
	opView.debug("lineHeight:", lineHeight, "lastLine:", lastLine, "lastLineWidth:", lastLineWidth)
	if err := gv.View.SetCursor(lastLineWidth, lineHeight); err != nil {
		if err := gv.View.SetOrigin(lastLineWidth, lineHeight); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (gv *GView) cursorEnd() error {
	lineHeight := gv.View.LinesHeight()
	lineHeight--
	lastLine, _ := gv.View.Line(lineHeight)
	lastLineWidth := len(lastLine)
	opView.debug("lineHeight:", lineHeight, "lastLine:", lastLine, "lastLineWidth:", lastLineWidth)
	if err := gv.View.SetCursor(0, lineHeight); err != nil {
		if err := gv.View.SetOrigin(0, lineHeight); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func InitUI() {
	Ui.AllView = map[string]GHandler{
		"global": gView,
		"server": sView,
		"info":   iView,
		"key":    kView,
		// "keydel":  kdView,
		"detail":  dView,
		"output":  opView,
		"tip":     tView,
		"project": pView,
		// "help":    hView,
		// "db":      dbView,
	}
	Ui.NextView = sView
	Ui.G.SetManager(iView, tView, pView, opView, dView, sView, kView)
	for _, item := range Ui.AllView {
		item.bindShortCuts()
	}
}
