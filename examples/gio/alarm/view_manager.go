package main

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/profile"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"git.sr.ht/~whereswaldon/sprig/core"
	"git.sr.ht/~whereswaldon/sprig/icons"
	sprigTheme "git.sr.ht/~whereswaldon/sprig/widget/theme"
)

type ViewManager interface {
	RequestViewSwitch(viewId)
	SetView(viewId)
	RegisterView(id viewId, view View)
	RegisterIntentHandler(intent Intent) bool
	RequestInvalidate()
	HandleBackNavigation()
	RequestContextualBar(
		gtx layout.Context,
		title string,
		actions []component.AppBarAction,
		overflow []component.OverflowAction,
	)
	DismissContextualBar(gtx layout.Context)
	DismissOverflow(gtx layout.Context)
	SelectedOverflowTag() interface{}
	Layout(gtx layout.Context) layout.Dimensions
	SetProfiling(bool)
	SetThemeing(bool)
	ApplySettings(core.SettingsService)
}

type viewManager struct {
	views       map[viewId]View
	currentView viewId
	window      *app.Window

	core.AP

	*component.AppBar
	*component.ModalNavDrawer
	*component.ModalLayer
	component.NavDrawer
	navigationAnimator component.VisibilityAnimation

	intentToView map[IntentID]viewId

	SelectedOverflowTagFunc func() interface{}

	viewStack []viewId

	dockDrawer bool

	profilingEnabled bool
	profile          profile.Event
	lastMallocs      uint64

	themeingEnabled bool
	themeView       View
}
