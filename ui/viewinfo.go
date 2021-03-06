package ui

import (
	"gosrg/ui/base"
	"gosrg/utils"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jessewkun/gocui"
)

var iView *InfoView

type InfoView struct {
	base.GView
}

func init() {
	iView = new(InfoView)
	iView.Name = "info"
	iView.Title = " Info "
	iView.ShortCuts = []base.ShortCut{
		base.ShortCut{Key: gocui.KeyCtrlY, Level: base.SC_LOCAL_Y, Handler: iView.copy},
	}
}

func (i *InfoView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(i.Name, maxX-29, 0, maxX-1, maxY-15, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = i.Title
		v.Wrap = true
		i.View = v
	}
	return nil
}

func (i *InfoView) formatOutput(argv interface{}) {
	if info, ok := argv.([]string); ok {
		i.Outputln(utils.Yellow(strings.ToLower(info[0])+":") + info[1])
	} else {
		opView.error("argv does not contain a variable of type []string")
	}
}

func (i *InfoView) copy(g *gocui.Gui, v *gocui.View) error {
	if err := clipboard.WriteAll(v.ViewBuffer()); err != nil {
		opView.error(err.Error())
		return err
	}
	opView.info("Copy key info success")
	return nil
}
