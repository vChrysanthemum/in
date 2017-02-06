package ui

import (
	"fin/ui/utils"

	"github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

type EditorMode int

type Editor struct {
	termui.Block
	Buf *termui.Buffer

	// CommandMode
	CommandModeBufAreaY      int
	CommandModeBufAreaHeight int
	CommandModeBuf           *EditorLine
	CommandModeCursor        *EditorCommandCursor

	KeyEvents                       chan string
	KeyEventsResultIsQuitActiveMode chan bool

	Views []*EditorView
	*EditorView
}

func NewEditor() *Editor {
	ret := &Editor{
		Block: *termui.NewBlock(),
	}

	ret.PrepareCommandMode()

	ret.CommandModeCursor = NewEditorCommandCursor(ret)

	ret.KeyEvents = make(chan string, 200)
	ret.KeyEventsResultIsQuitActiveMode = make(chan bool)
	ret.RegisterKeyEventHandlers()

	ret.Views = append(ret.Views, NewEditorView(ret))
	ret.EditorView = ret.Views[0]

	return ret
}

func (p *Editor) Close() {
	close(p.KeyEvents)
	close(p.KeyEventsResultIsQuitActiveMode)
}

func (p *Editor) Buffer() termui.Buffer {
	if nil == p.Buf {
		if 0 == len(p.Lines) {
			p.EditModeAppendNewLine(p.EditModeCursor)
		}
		buf := p.Block.Buffer()
		p.Buf = &buf
		p.Buf.IfNotRenderByTermUI = true
		p.CommandModeBufAreaY = p.Block.InnerArea.Max.Y - 1
		p.CommandModeBufAreaHeight = 1
		p.EditModeBufAreaHeight = p.Block.InnerArea.Dy() - p.CommandModeBufAreaHeight
		p.isShouldRefreshEditModeBuf = true
		p.isShouldRefreshCommandModeBuf = true
	} else {
		p.isShouldRefreshEditModeBuf = true
		p.isShouldRefreshCommandModeBuf = true
	}

	if true == p.Block.Border {
		p.Block.DrawBorder(*p.Buf)
		p.Block.DrawBorderLabel(*p.Buf)
		for point, c := range p.Buf.CellMap {
			termbox.SetCell(point.X, point.Y, c.Ch, toTmAttr(c.Fg), toTmAttr(c.Bg))
		}
	}

	p.RefreshBuf()
	p.RefreshCursorByEditorLine()
	p.RefreshBuf()

	return *p.Buf
}

func (p *Editor) RefreshCursorByEditorLine() {
	switch p.Mode {
	case EditorEditMode:
		p.EditModeCursor.RefreshCursorByEditorLine(p.EditModeCursor.Line())
	case EditorNormalMode:
		p.EditModeCursor.RefreshCursorByEditorLine(p.EditModeCursor.Line())
	case EditorCommandMode:
		p.CommandModeCursor.RefreshCursorByEditorLine(p.CommandModeBuf)
	}
}

func (p *Editor) ActiveMode() {
	p.EditModeEnter(p.EditModeCursor)
}

func (p *Editor) UnActiveMode() {
	p.Mode = EditorModeNone
	utils.UISetCursor(-1, -1)
}
