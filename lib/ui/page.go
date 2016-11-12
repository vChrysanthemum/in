package ui

import (
	"container/list"
	"log"

	"github.com/gizak/termui"
	"golang.org/x/net/html"
)

type Page struct {
	Title string

	IdToNodeMap map[string]*Node

	Bufferers []termui.Bufferer

	script         *Script
	parseAgentMap  []*ParseAgent
	renderAgentMap []*RenderAgent

	doc                     *html.Node
	parsingNodesStack       *list.List
	FirstChildNode          *Node
	FocusNode               *list.Element
	WorkingNodes            *list.List
	ActiveNode              *Node
	ActiveNodeAfterRerender *Node

	renderingX int
	renderingY int
}

func newPage() *Page {
	ret := new(Page)

	ret.IdToNodeMap = make(map[string]*Node, 0)

	ret.parsingNodesStack = list.New()
	ret.WorkingNodes = list.New()

	ret.prepareScript()
	ret.prepareParse()
	ret.prepareRender()

	return ret
}

func (p *Page) dumpNodesHtmlData(node *Node) {
	log.Println(node.HtmlData)
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		p.dumpNodesHtmlData(childNode)
	}
}

func (p *Page) DumpNodesHtmlData() {
	p.dumpNodesHtmlData(p.FirstChildNode)
}

func (p *Page) RemoveNode(node *Node) {
	if nil != node.OnRemove {
		node.OnRemove()
	}

	delete(p.IdToNodeMap, node.Id)

	if nil != node.PrevSibling {
		node.PrevSibling.NextSibling = node.NextSibling
	}

	if nil != node.NextSibling {
		node.NextSibling.PrevSibling = node.PrevSibling
	}

	if nil != node.Parent {
		node.Parent.ChildrenCount -= 1
		if node.Parent.FirstChild == node {
			node.Parent.FirstChild = node.NextSibling
		}
		if node.Parent.LastChild == node {
			node.Parent.LastChild = node.PrevSibling
		}
	}

	p.Rerender()
}

func (p *Page) Refresh() {
	if nil != p.FocusNode {
		focusNode := p.FocusNode.Value.(*Node)
		if nil != focusNode.UnFocusMode {
			focusNode.UnFocusMode()
		}
	}

	if nil != p.ActiveNode && nil != p.ActiveNode.UnActiveMode {
		p.ActiveNode.UnActiveMode()
	}

	uiClear()

	if nil != p.ActiveNodeAfterRerender {
		p.SetActiveNode(p.ActiveNodeAfterRerender)
		p.FocusNode = nil
	} else if nil != p.FocusNode {
		p.SetActiveNode(p.FocusNode.Value.(*Node))
	}

	p.uiRender()
}

func (p *Page) nodeAfterRenderHandle(node *Node) {
	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		p.nodeAfterRenderHandle(childNode)
	}

	if nil != node.afterRenderHandle {
		node.afterRenderHandle()
	}
}

func (p *Page) uiRender() {
	if 0 == len(p.Bufferers) {
		return
	}
	termui.Render(p.Bufferers...)
	p.nodeAfterRenderHandle(p.FirstChildNode)
}

func (p *Page) Rerender() {
	p.Render()
	p.Refresh()
}

func (p *Page) Serve() {
	defer termui.Close()

	p.Refresh()

	p.registerHandles()
	go func() {
		p.script.Run()
	}()

	termui.Loop()
}
