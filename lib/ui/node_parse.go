package ui

import (
	. "in/ui/utils"
	"strconv"

	"golang.org/x/net/html"
)

func (p *Node) ParseAttribute(attr []html.Attribute) (isNeedRerenderPage bool) {
	isNeedRerenderPage = false

	if nil == p.uiBlock {
		return
	}
	p.uiBlock.BorderLabelFg = COLOR_DEFAULT_BORDER_LABEL_FG
	p.uiBlock.BorderFg = COLOR_DEFAULT_BORDER_FG

	for _, v := range attr {
		switch v.Key {
		case "borderlabelfg":
			p.uiBlock.BorderLabelFg = ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_BORDER_LABEL_FG)
		case "borderlabel":
			p.uiBlock.BorderLabel = v.Val

		case "borderfg":
			p.uiBlock.BorderFg = ColorToTermuiAttribute(v.Val, COLOR_DEFAULT_BORDER_FG)

		case "border":
			p.uiBlock.Border = StringToBool(v.Val, p.uiBlock.Border)
		case "borderleft":
			p.uiBlock.BorderLeft = StringToBool(v.Val, p.uiBlock.BorderLeft)
		case "borderright":
			p.uiBlock.BorderRight = StringToBool(v.Val, p.uiBlock.BorderRight)
		case "bordertop":
			p.uiBlock.BorderTop = StringToBool(v.Val, p.uiBlock.BorderTop)
		case "borderbottom":
			p.uiBlock.BorderBottom = StringToBool(v.Val, p.uiBlock.BorderBottom)

		case "height":
			isNeedRerenderPage = true
			p.uiBlock.Height, _ = strconv.Atoi(v.Val)
			if p.uiBlock.Height < 0 {
				p.uiBlock.Height = 0
			}
			p.isShouldCalculateHeight = false
		case "width":
			isNeedRerenderPage = true
			p.uiBlock.Width, _ = strconv.Atoi(v.Val)
			if p.uiBlock.Width < 0 {
				p.uiBlock.Width = 0
			}
			p.isShouldCalculateWidth = false
		}
	}

	return
}
