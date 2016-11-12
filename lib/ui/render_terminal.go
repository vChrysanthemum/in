package ui

import (
	. "in/ui/utils"

	"github.com/gizak/termui"
)

func (p *Page) renderBodyTerminal(node *Node) (isFallthrough bool) {
	isFallthrough = false
	uiBuffer := node.Data.(*NodeTerminal).Editor

	p.normalRenderNodeBlock(node)

	if "" != node.ColorFg {
		uiBuffer.TextFgColor = ColorToTermuiAttribute(node.ColorFg, termui.ColorDefault)
	}
	if "" != node.ColorBg {
		uiBuffer.TextBgColor = ColorToTermuiAttribute(node.ColorBg, termui.ColorDefault)
	}

	p.BufferersAppend(node, uiBuffer)

	p.pushWorkingNode(node)

	p.renderingY = uiBuffer.Block.Y + uiBuffer.Block.Height

	return
}
