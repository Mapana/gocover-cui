package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"golang.org/x/tools/cover"
	"io"
	"io/ioutil"
	"math"
	"time"
)

// https://github.com/golang/go/blob/master/src/cmd/cover/html.go#L24
//
// imitate
func cuiOutput(profile string, logfile string) error {
	var (
		cuiData templateData
		logData []byte
	)

	if len(profile) > 0 {
		profiles, err := cover.ParseProfiles(profile)
		if err != nil {
			return err
		}

		dirs, err := findPkgs(profiles)
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			fn := profile.FileName
			file, err := findFile(dirs, fn)
			if err != nil {
				return err
			}
			src, err := ioutil.ReadFile(file)
			if err != nil {
				return fmt.Errorf("can't read %q: %v", fn, err)
			}
			var buf bytes.Buffer
			err = cuiGen(&buf, src, profile.Boundaries(src))
			if err != nil {
				return err
			}
			cuiData.Files = append(cuiData.Files, &templateFile{
				Name:     fn,
				Body:     buf.String(),
				Coverage: percentCovered(profile),
			})
		}
	}
	if len(logfile) > 0 {
		var err error

		logData, err = ioutil.ReadFile(logfile)
		if err != nil {
			return err
		}
	}

	cui := &GoCoverCui{
		tview.NewApplication(),
		&MainView{
			tview.NewFlex(),
			tview.NewFlex(),
			tview.NewPages(),
			tview.NewFlex(),
		},
	}
	if err := cui.view(&cuiData, logData); err != nil {
		return err
	}
	return cui.Application.SetRoot(cui.Main.Flex, true).SetFocus(cui.Main.Top).Run()
}

// https://github.com/golang/go/blob/master/src/cmd/cover/html.go#L112
//
// imitate
// modify the output
// change html to bash
func cuiGen(w io.Writer, src []byte, boundaries []cover.Boundary) error {
	color := defaultCov
	dst := bufio.NewWriter(w)
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			b := boundaries[0]
			if b.Start {
				n := 0
				if b.Count > 0 {
					n = int(math.Floor(b.Norm*9)) + 1
				}

				if n == 0 {
					color = cov0
				} else if n == 8 {
					color = cov8
				}
			} else {
				color = defaultCov
			}
			dst.WriteString(color)
			boundaries = boundaries[1:]
		}
		dst.WriteByte(src[i])
		if src[i] == '\n' {
			dst.WriteString(color)
		}
	}
	return dst.Flush()
}

func (cui *GoCoverCui) view(cuiData *templateData, logData []byte) error {
	var showPage string
	rowView := tview.NewFlex().SetDirection(tview.FlexRow)
	item := tview.NewDropDown().SetLabel(" Select File: ").SetCurrentOption(0)

	// generate log ui
	if len(logData) > 0 {
		err := cui.logView(logData, item)
		if err != nil {
			return err
		}
		showPage = logger
	}

	// generate cover ui
	if cuiData != nil && len(cuiData.Files) > 0 {
		err := cui.cuiView(cuiData.Files, item)
		if err != nil {
			return err
		}
		showPage = cuiData.Files[0].Name
	}

	// help view
	helpView, err := cui.helpView()
	if err != nil {
		return err
	}
	rowView.AddItem(helpView, 4, 0, false)

	// top view
	cui.Main.Top.SetBorder(true).SetTitle(" Cover Files ").SetTitleAlign(tview.AlignCenter)
	item.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			cui.SetFocus(cui.Main.Pages)
		}
	})
	cui.Main.Top.AddItem(item, 0, 3, true)
	for k, v := range statusMap {
		cui.Main.Top.AddItem(statusColorGen(fmt.Sprintf("%s%s", v, k)), 0, 1, false)
	}
	rowView.AddItem(cui.Main.Top, 3, 1, false)

	// pages view
	cui.Main.Pages.SwitchToPage(showPage)
	rowView.AddItem(cui.Main.Pages, 0, 1, false)

	// main view
	cui.Main.AddItem(rowView, 0, 1, false)

	return nil
}

func (cui *GoCoverCui) helpView() (*tview.Flex, error) {
	var proportion = 1

	for i, data := range helpData {
		tv := tview.NewTextView()
		tv.SetDynamicColors(true).SetWrap(true).SetScrollable(true)
		if _, err := fmt.Fprint(tview.ANSIWriter(tv), data); err != nil {
			return cui.Main.Help, err
		}
		if i == len(helpData)-1 {
			proportion = 2
		}
		cui.Main.Help.AddItem(tv, 0, proportion, false)
	}
	return cui.Main.Help, nil
}

func (cui *GoCoverCui) logView(logData []byte, item *tview.DropDown) error {
	// new select option
	item.AddOption(logger, func() {
		cui.Main.Pages.SwitchToPage(logger)
	}).SetFieldWidth(len(logger) + 10)

	// new data view
	tv, err := dataViewGen(logger, string(logData), func(key tcell.Key) {
		if key == tcell.KeyTAB {
			cui.SetFocus(cui.Main.Top)
		}
	})
	if err != nil {
		return err
	}
	cui.Main.Pages.AddPage(logger, tv, true, false)

	return nil
}

func (cui *GoCoverCui) cuiView(files []*templateFile, item *tview.DropDown) error {
	for _, f := range files {
		ch := make(chan error)
		// generate main view for cover file
		go func(f *templateFile, ch chan error) {
			defer close(ch)

			tv, err := dataViewGen(f.Name, f.Body, func(key tcell.Key) {
				if key == tcell.KeyTAB {
					cui.SetFocus(cui.Main.Top)
				}
			})
			if err != nil {
				ch <- err
				return
			}
			cui.Main.Pages.AddPage(f.Name, tv, true, false)
		}(f, ch)

		select {
		case err := <-ch:
			if err != nil {
				return err
			}
		case <-time.After(time.Second * 10):
			return fmt.Errorf("\033[0;32mgenerate cui timeout 10 seconds")
		}

		// generate select view for cover file
		go func(f *templateFile) {
			item.AddOption(fmt.Sprintf("%s  %.1f%s", f.Name, f.Coverage, "%"), func() {
				cui.Main.Pages.SwitchToPage(f.Name)
			}).SetFieldWidth(len(f.Name) + 10)
		}(f)
	}

	return nil
}

func dataViewGen(name string, data string, handler func(key tcell.Key)) (*tview.TextView, error) {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true).SetWrap(true).SetScrollable(true).SetDoneFunc(handler).
		SetRegions(true).SetBorder(true).SetTitle(fmt.Sprintf(" Data View: [ %s ] ", name))

	if _, err := fmt.Fprint(tview.ANSIWriter(tv), data); err != nil {
		return nil, err
	}

	return tv, nil
}

func statusColorGen(text string) *tview.TextView {
	tv := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	if _, err := fmt.Fprint(tview.ANSIWriter(tv), text); err != nil {
		panic(err)
	}
	return tv
}

type GoCoverCui struct {
	*tview.Application
	Main *MainView
}

type MainView struct {
	*tview.Flex
	Top   *tview.Flex
	Pages *tview.Pages
	Help  *tview.Flex
}

var (
	statusMap = map[string]string{"not tracked": defaultCov, "not covered": cov0, "covered": cov8}
	helpData  = []string{
		`[ [yellow::b]global keyboard [white]]
- [green::b]TAB[white]: switch focus in "Select File" or "Data View"`,
		`[ [yellow::b]"Select File" keyboard [white]]
[green::b]0-9、a-Z、↑、↓、ENTER[white] or other: ctivation option
[green::b]↑ [white]or [green::b]↓[white]: switch the option
[green::b]ENTER[white]: determine the choice`,
		`[ [yellow::b]Data View keyboard [white]]
[green::b]↑ [white]or [green::b]k[white]: move up                            | [green::b]g [white]or [green::b]home[white]: move to the top
[green::b]↓ [white]or [green::b]j[white]: move down                          | [green::b]G [white]or [green::b]end[white]: move to the bottom
[green::b]CTRL [white]+ [green::b]F[white] or [green::b]PAGE UP[white]: move down by one page | [green::b]CTRL [white]+ [green::b]B[white] or [green::b]PAGE DOWN[white]: move up by one page`,
	}
)

const (
	defaultCov = "\033[0;37m"
	cov0       = "\033[0;31m"
	cov8       = "\033[0;32m"
	logger     = "logger"
)
