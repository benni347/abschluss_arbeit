package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/x/profiling"
	"git.sr.ht/~whereswaldon/sprig/core"
	sprigTheme "git.sr.ht/~whereswaldon/sprig/widget/theme"
	utils "github.com/benni347/messengerutils"
	"github.com/inkeliz/giohyperlink"
	"github.com/pkg/profile"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	log.SetFlags(log.Lshortfile | log.Flags())
	go func() {
		w := app.NewWindow(app.Title("Messenger"))
		if err := run(w); err != nil {
			utils.PrintError("During the app running an error ocured", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var (
		dataDir        string
		invalidate     bool
		profileOptions string
	)
	dataDir, err := getDataDir("messenger")
	if err != nil {
		utils.PrintError("During the data dir fetching an error ocured", err)
		return err
	}

	flag.StringVar(
		&profileOptions,
		"profile",
		"none",
		"enable profiling mode, one of [none, cpu, mem, mutex, block, goroutine, trace, gio]",
	)
	flag.BoolVar(
		&invalidate,
		"invalidate",
		false,
		"invalidate on each frame, only useful for profiling",
	)
	flag.StringVar(&dataDir, "data-dir", dataDir, "data directory")
	flag.Parse()

	profiler := ProfileOptions(profileOptions).NewProfiler()
	profiler.Start()
	defer profiler.Stop()

	app, err := core.NewApp(w, dataDir)
	if err != nil {
		utils.PrintError("During the app creation an error ocured", err)
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	viewManager := NewViewManager(w, app)
	return nil
}

type viewId int

const (
	ConnectFormID viewId = iota
	IdentityFormID
	SettingsID
	ReplyViewID
	ConsentViewID
	SubscriptionViewID
	SubscriptionSetupFormViewID
	DynamicChatViewID
)

func getDataDir(folderName string) (string, error) {
	data, err := app.DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(data, folderName), nil
}

type ProfileOptions string

const (
	None      ProfileOptions = "none"
	CPU       ProfileOptions = "cpu"
	Memory    ProfileOptions = "mem"
	Mutex     ProfileOptions = "mutex"
	Block     ProfileOptions = "block"
	Goroutine ProfileOptions = "goroutine"
	Trace     ProfileOptions = "trace"
	Gio       ProfileOptions = "gio"
)

func (p ProfileOptions) NewProfiler() Profiler {
	switch p {
	case "", None:
		return Profiler{}
	case CPU:
		return Profiler{Starter: profile.CPUProfile}
	case Memory:
		return Profiler{Starter: profile.MemProfile}
	case Mutex:
		return Profiler{Starter: profile.MutexProfile}
	case Block:
		return Profiler{Starter: profile.BlockProfile}
	case Goroutine:
		return Profiler{Starter: profile.GoroutineProfile}
	case Trace:
		return Profiler{Starter: profile.TraceProfile}
	case Gio:
		var (
			recorder *profiling.CSVTimingRecorder
			err      error
		)

		return Profiler{
			Starter: func(*profile.Profile) {
				recorder, err = profiling.NewRecorder(nil)
				if err != nil {
					log.Fatal(err)
				}
			},
			Stopper: func() {
				if recorder == nil {
					return
				}
				if err := recorder.Stop(); err != nil {
					utils.PrintError("During the recorder stopping an error ocured", err)
				}
			},
			Recorder: func(gtx C) {
				if recorder == nil {
					return
				}
				recorder.Profile(gtx)
			},
		}
	}
	return Profiler{}
}

type Profiler struct {
	Starter  func(*profile.Profile)
	Stopper  func()
	Recorder func(gtx C)
}

func (p *Profiler) Start() {
	if p.Starter != nil {
		p.Stopper = profile.Start(p.Starter).Stop
	}
}

func (p *Profiler) Stop() {
	if p.Stopper != nil {
		p.Stopper()
	}
}

func (p *Profiler) Record(gtx C) {
	if p.Recorder != nil {
		p.Recorder(gtx)
	}
}
