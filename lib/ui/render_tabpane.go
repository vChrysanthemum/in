package ui

import (
	"github.com/gizak/termui"
	"github.com/gizak/termui/extra"
)

func (p *Page) _renderBodyTabpaneOneTab(nodeTab *Node) {
	var (
		nodeTabChild *Node
	)

	for nodeTabChild = nodeTab.FirstChild; nodeTabChild != nil; nodeTabChild = nodeTabChild.NextSibling {
		if false == *nodeTabChild.Display {
			continue
		}

		uiBuffer := nodeTab.uiBuffer.(*extra.Tab)
		uiBuffer.AddBlocks(nodeTabChild.uiBuffer.(termui.Bufferer))
	}
}

func (p *Page) renderBodyTabpane(node *Node) {
	nodeDataTabpane := node.Data.(*NodeTabpane)

	uiBuffer := node.uiBuffer.(*extra.Tabpane)

	p.normalRenderNodeBlock(node)
	node.UIBlock.X = 0
	node.UIBlock.Y = 0

	nodeDataTabpane.Tabs = []*extra.Tab{}
	index := 0
	for nodeTab := node.FirstChild; nodeTab != nil; nodeTab = nodeTab.NextSibling {
		nodeTab.Data.(*NodeTabpaneTab).Index = index
		p._renderBodyTabpaneOneTab(nodeTab)
		nodeDataTabpane.Tabs = append(nodeDataTabpane.Tabs, nodeTab.uiBuffer.(*extra.Tab))
		index += 1
	}

	uiBuffer.SetTabs(nodeDataTabpane.Tabs...)

	p.BufferersAppend(node, node.uiBuffer.(termui.Bufferer))

	return
}
