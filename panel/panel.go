package panel

import (
	"dcd/line"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type PanelType uint32

const (
	InputType PanelType = iota
	OutputType
)

type Panel struct {
	X      int
	Y      int
	Width  int
	Height int
	Bg     termbox.Attribute
	Fg     termbox.Attribute

	Type PanelType
	PosX int
	PosY int

	FilterRegexp string

	buffers    []line.Line
	lines      []line.Line
	startLine  int
	selectLine int
}

func (p *Panel) drawEdge() {
	for i := p.X; i < p.X+p.Width; i++ {
		termbox.SetCell(i, p.Y+p.Height-1, '-', termbox.ColorGreen, p.Bg)
	}
	for i := p.Y; i < p.Y+p.Height; i++ {
		termbox.SetCell(p.X+p.Width-1, i, '|', termbox.ColorGreen, p.Bg)
	}
}

func (p *Panel) drawCursor() {
	if p.Type == InputType {
		termbox.SetCursor(p.X+p.PosX, p.Y+p.PosY)
	} else if p.Type == OutputType {
	} else {
	}
}

func (p *Panel) filter() {
	p.lines = make([]line.Line, 0)
	if p.FilterRegexp == "" {
		p.lines = append(p.lines, p.buffers...)
	}
}

func (p *Panel) drawLines() {
	p.filter()
	if p.selectLine < p.startLine {
		p.startLine = p.selectLine
	}
	minStartLine := 0
	linesHeight := 0
	for i := p.selectLine; i >= 0 && i < len(p.lines); i-- {
		if p.lines[i].GetHeight(p.Width)+linesHeight >= p.Height {
			break
		} else {
			linesHeight += p.lines[i].GetHeight(p.Width)
			minStartLine = i
		}
	}
	if minStartLine > p.startLine {
		p.startLine = minStartLine
	}
	endLine := 0
	linesHeight = 0
	for i := p.startLine; i < len(p.lines); i++ {
		if p.lines[i].GetHeight(p.Width)+linesHeight >= p.Height {
			break
		} else {
			linesHeight += p.lines[i].GetHeight(p.Width)
			endLine = i
		}
	}
	yIndex := 0
	xIndex := 0
	for i := p.startLine; i <= endLine && i < len(p.lines); i++ {
		xIndex = 0
		for _, v := range p.lines[i].Cs {
			if xIndex+runewidth.RuneWidth(v.Ch) < p.Width {
			} else {
				yIndex++
				xIndex = 0
			}
			if i == p.selectLine && p.Type == OutputType {
				termbox.SetCell(p.X+xIndex, p.Y+yIndex, v.Ch, v.Fg|p.lines[i].Fg, v.Bg|termbox.ColorCyan)
			} else {
				termbox.SetCell(p.X+xIndex, p.Y+yIndex, v.Ch, v.Fg|p.lines[i].Fg, v.Bg|p.lines[i].Bg)
			}
			xIndex += runewidth.RuneWidth(v.Ch)
		}
		yIndex++
	}
}

func (p *Panel) Init(x, y, w, h int, fg, bg termbox.Attribute, t PanelType, px, py int) {
	p.X = x
	p.Y = y
	p.Width = w
	p.Height = h
	p.Bg = bg
	p.Fg = fg
	p.Type = t
	p.PosX = px
	p.PosY = py
	p.selectLine = 0
}

func (p *Panel) PushLine(b []byte) {
	l := line.Line{Fg: p.Fg, Bg: p.Bg}
	l.PushBytes(b)
	p.buffers = append(p.buffers, l)
}

func (p *Panel) Draw() {
	p.drawEdge()
	p.drawCursor()
	p.drawLines()
}

func (p *Panel) Up() {
	if p.selectLine > 0 && p.Type == OutputType {
		p.selectLine--
	}
}

func (p *Panel) Down() {
	if p.selectLine < len(p.lines)-1 && p.Type == OutputType {
		p.selectLine++
	}
}

func (p *Panel) Push(Ch rune) {
	if p.PosX >= p.Width-2 || p.Type == OutputType {
		return
	}
	p.PosX += runewidth.RuneWidth(Ch)
	if len(p.buffers) == 0 {
		l := line.Line{Fg: p.Fg, Bg: p.Bg}
		p.buffers = append(p.buffers, l)
	}
	p.buffers[0].PushCell(termbox.Cell{Ch: Ch, Fg: p.Fg, Bg: p.Bg})
}

func (p *Panel) Pop() {
	if p.PosX <= 0 || p.Type == OutputType {
		return
	}
	p.PosX -= p.buffers[0].PopCell()
}