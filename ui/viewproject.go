package ui

import (
	"gosrg/config"
	"gosrg/ui/base"
	"gosrg/utils"
	"strings"

	"github.com/jessewkun/gocui"
)

var pView *ProjectView

type ProjectView struct {
	base.GView
}

const MAX_LEN = 18

func init() {
	pView = new(ProjectView)
	pView.Name = "project"
	pView.ShortCuts = []base.ShortCut{
		base.ShortCut{Key: gocui.MouseLeft, Level: base.SC_LOCAL_N, Handler: pView.openGit},
	}
}

func (p *ProjectView) Layout(g *gocui.Gui) error {
	maxX, maxY := Ui.G.Size()
	if v, err := g.SetView(p.Name, maxX-20, maxY-2, maxX-1, maxY, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Frame = false
		p.View = v
		p.Initialize()
	}
	return nil
}

func (p *ProjectView) Initialize() error {
	str := config.PROJECT_NAME + " " + config.Version
	l := len(str)
	if MAX_LEN > l {
		str = strings.Repeat(" ", MAX_LEN-l) + utils.Pink(utils.UnderLine(str))
	} else {
		str = utils.Pink(utils.UnderLine(str))
	}
	p.Output(str)
	return nil
}

func (p *ProjectView) openGit(g *gocui.Gui, v *gocui.View) error {
	utils.OpenLink(config.PROJECT_URL)
	return nil
}
