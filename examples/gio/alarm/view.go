package main

import (
	"gioui.org/layout"
	"gioui.org/x/component"
)

type View interface {
	SetManager(ViewManager)
	AppBarData() (bool, string, []component.AppBarAction, []component.OverflowAction)
	NavItem() *component.NavItem
	BecomeVisible()
	HandleIntent(Intent)
	Update(gtx layout.Context)
	Layout(gtx layout.Context) layout.Dimensions
}
