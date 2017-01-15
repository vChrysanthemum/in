package ui

import (
	uiutils "fin/ui/utils"
	"github.com/gizak/termui"
	"golang.org/x/net/html"
)

func (p *Page) parseBodyLineChart(parentNode *Node, htmlNode *html.Node) (ret *Node, isFallthrough bool) {
	ret = p.newNode(htmlNode)
	parentNode.addChild(ret)
	isFallthrough = true

	ret.InitNodeLineChart()

	return
}

func (p *NodeLineChart) NodeDataParseAttribute(attr []html.Attribute) (isUIChange, isNeedReRenderPage bool) {
	isUIChange = false
	isNeedReRenderPage = false

	uiBuffer := p.Node.uiBuffer.(*termui.LineChart)

	for _, v := range attr {
		switch v.Key {
		case "mode":
			isUIChange = true
			switch v.Val {
			case "braille":
				uiBuffer.Mode = v.Val

			case "dot":
				uiBuffer.Mode = v.Val
			}

		case "axescolor":
			isUIChange = true
			uiBuffer.AxesColor = uiutils.ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_LINE_CHART_AXES)

		case "linecolor":
			isUIChange = true
			uiBuffer.LineColor = uiutils.ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_LINE_CHART_LINE)
		}
	}

	return
}
