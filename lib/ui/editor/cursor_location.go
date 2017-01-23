package editor

import (
	uiutils "fin/ui/utils"

	"github.com/gizak/termui"
)

type CursorLocation struct {
	IsDisplay     bool
	OffXCellIndex int
	Editor        *Editor
}

func NewCursorLocation(editor *Editor) *CursorLocation {
	ret := &CursorLocation{
		IsDisplay:     false,
		OffXCellIndex: 0,
		Editor:        editor,
	}
	return ret
}

func (p *CursorLocation) ResetLocation() {
	uiutils.UISetCursor(p.Editor.Block.InnerArea.Min.X, p.Editor.Block.InnerArea.Min.Y)
}

func (p *CursorLocation) ResumeCursor() {
	//uiutils.UISetCursor(p.Location.X, p.Location.Y)
}

func (p *CursorLocation) MoveCursorNRuneLeft(n int) {
	if len(p.Editor.CurrentLine.Cells) == 0 {
		p.OffXCellIndex = 0
		return
	}

	p.OffXCellIndex -= n
	if p.OffXCellIndex < 0 {
		p.OffXCellIndex = 0
	}

	cell := p.Editor.CurrentLine.Cells[p.OffXCellIndex]
	uiutils.UISetCursor(cell.X, cell.Y)
	uiutils.UIRender(p.Editor)
}

func (p *CursorLocation) MoveCursorNRuneRight(n int) {
	if len(p.Editor.CurrentLine.Cells) == 0 {
		p.OffXCellIndex = 0
		return
	}

	p.OffXCellIndex += n
	if p.OffXCellIndex >= len(p.Editor.CurrentLine.Cells) {
		p.OffXCellIndex = len(p.Editor.CurrentLine.Cells) - 1
	}

	cell := p.Editor.CurrentLine.Cells[p.OffXCellIndex]
	uiutils.UISetCursor(cell.X, cell.Y)
	uiutils.UIRender(p.Editor)
}

func (p *CursorLocation) MoveCursorAfterWrite(line *Line) {
	if nil == line {
		uiutils.UISetCursor(0, 0)
		return
	}

	if 0 == len(line.Cells) {
		uiutils.UISetCursor(line.ContentStartX, line.ContentStartY)
		return
	}

	if p.OffXCellIndex >= len(line.Cells) {
		p.OffXCellIndex = len(line.Cells)
	}

	var cell termui.Cell
	if p.OffXCellIndex == len(line.Cells) {
		cell = line.Cells[p.OffXCellIndex-1]
		width := cell.Width()
		if cell.X+width >= p.Editor.Block.InnerArea.Max.X {
			uiutils.UISetCursor(line.ContentStartX, cell.Y+1)
		} else {
			uiutils.UISetCursor(cell.X+width, cell.Y)
		}
	} else {
		cell = line.Cells[p.OffXCellIndex]
		uiutils.UISetCursor(cell.X, cell.Y)
	}
}
