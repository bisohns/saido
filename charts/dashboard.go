package charts

import (
	"context"
	"fmt"
	"image"
	"math/rand"
	"sync"
	"time"

	"github.com/bisohns/saido/config"
	"github.com/bisohns/saido/inspector"
	log "github.com/sirupsen/logrus"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/sparkline"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

// widgets holds the widgets used by this demo.
type widgets struct {
	segDist      *segmentdisplay.SegmentDisplay
	input        *textinput.TextInput
	rollT        *text.Text
	spGreen      *sparkline.SparkLine
	barChart     *barchart.BarChart
	donut        *donut.Donut
	hosts        [][]grid.Element
	goToMainPage *button.Button
}

var (
	hostsPerPage     int = 5
	currentHostPage  int = 0
	currentHost      string
	currentMetric    string
	inspectorWidgets map[string]map[string]widgetapi.Widget = map[string]map[string]widgetapi.Widget{}
	globalLogWidget  *text.Text                             = nil
	globalLogFunc    func(string) error
	globalContainer  *container.Container
	dashboardInfo    *config.DashboardInfo
	globalCtx        context.Context
	globalCancel     context.CancelFunc
	globalTerminal   terminalapi.Terminal
)

// getOrSetGlobalLog gets or sets the global log widget
// in a persistent matter despite window refreshes
func getOrSetGlobalLog() (*text.Text, error) {
	var (
		rollT              *text.Text
		logToDashBoardfunc func(string) error
		err                error
	)
	if globalLogWidget == nil {
		rollT, logToDashBoardfunc, err = newRollText()
		if err != nil {
			return nil, err
		}
		globalLogWidget = rollT
		globalLogFunc = logToDashBoardfunc
	} else {
		rollT = globalLogWidget
		logToDashBoardfunc = globalLogFunc
	}
	return rollT, nil
}

// newWidgets creates all widgets used by this demo.
func newWidgets() (*widgets, error) {
	sd, err := newRotatingSegmentDisplay(dashboardInfo.Title)
	if err != nil {
		return nil, err
	}

	rollT, err := getOrSetGlobalLog()
	if err != nil {
		return nil, err
	}
	spGreen, err := newSparkLines()
	if err != nil {
		return nil, err
	}

	bc, err := newBarChart()
	if err != nil {
		return nil, err
	}

	don, err := newDonut()
	if err != nil {
		return nil, err
	}

	paginatedHosts := Paginate(dashboardInfo.Hosts, hostsPerPage)
	goTo, err := addGoToMain()
	if err != nil {
		return nil, err
	}
	constantWidgets := &widgets{
		segDist:      sd,
		rollT:        rollT,
		goToMainPage: goTo,
	}
	hosts, err := newHostButtons(paginatedHosts, dashboardInfo.Metrics, constantWidgets)
	if err != nil {
		return nil, err
	}

	return &widgets{
		segDist:      constantWidgets.segDist,
		rollT:        constantWidgets.rollT,
		spGreen:      spGreen,
		barChart:     bc,
		donut:        don,
		hosts:        hosts,
		goToMainPage: goTo,
	}, nil
}

// layoutType represents the possible layouts the buttons switch between.
type layoutType int

const (
	// mainDashboard displays all the widgets.
	mainDashboard layoutType = iota
	// layoutSingle shows a single inspector of a host
	hostMetric
	// layoutText focuses onto the text widget.
	layoutText
	// layoutSparkLines focuses onto the sparklines.
	layoutSparkLines
	// layoutLineChart focuses onto the linechart.
	layoutLineChart
)

// gridLayout prepares container options that represent the desired screen layout.
// This function demonstrates the use of the grid builder.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType mainDashboard.
func gridLayout(w *widgets, lt layoutType) ([]container.Option, error) {
	leftRows := []grid.Element{
		grid.RowHeightPerc(25,
			grid.Widget(w.segDist,
				container.Border(linestyle.Light),
				container.BorderTitle("Press Esc to quit"),
			),
		),
	}
	switch lt {
	case mainDashboard:
		leftRows = append(leftRows,
			grid.RowHeightPerc(20,
				grid.ColWidthPerc(20,
					grid.Widget(w.rollT,
						container.Border(linestyle.Light),
						container.BorderTitle("Log reports"),
					),
				),
				grid.ColWidthPercWithOpts(60,
					[]container.Option{
						container.Border(linestyle.Light),
						container.BorderTitle("Hosts"),
					},
					w.hosts[currentHostPage]...,
				),
			),
		)
	case hostMetric:
		headerWidget, err := newSegmentDisplay(currentHost)
		if err != nil {
			return nil, err
		}
		header := []grid.Element{
			grid.RowHeightPerc(25,
				grid.Widget(headerWidget,
					container.Border(linestyle.Light),
					container.BorderTitle("Press Esc to quit"),
				),
			),
			grid.RowHeightPerc(10,
				grid.Widget(w.goToMainPage,
					container.Border(linestyle.Light),
					container.BorderTitle("Press Esc to quit"),
				),
			),
		}

		leftRows = append(header,
			grid.RowHeightPerc(65,
				grid.ColWidthPerc(20,
					grid.Widget(w.rollT,
						container.Border(linestyle.Light),
						container.BorderTitle("Log reports"),
					),
				),
				grid.ColWidthPercWithOpts(80,
					[]container.Option{
						container.Border(linestyle.Light),
						container.BorderTitle(fmt.Sprintf("%s-%s", currentHost, currentMetric)),
						container.PlaceWidget(inspectorWidgets[currentHost][currentMetric]),
					},
				),
			),
		)
	case layoutText:
		leftRows = append(leftRows,
			grid.RowHeightPerc(65,
				grid.Widget(w.rollT,
					container.Border(linestyle.Light),
					container.BorderTitle("A rolling text"),
				),
			),
		)

	}

	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc(80, leftRows...),
	)

	innergridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return innergridOpts, nil
}

// rootID is the ID assigned to the root container.
const rootID = "root"

// Terminal implementations
const (
	termboxTerminal = "termbox"
	tcellTerminal   = "tcell"
)

func Main(cfg *config.Config) {
	var err error
	dashboardInfo = config.GetDashboardInfoConfig(cfg)
	log.Errorf("%v", dashboardInfo)
	log.Debugf("Starting %s", dashboardInfo.Title)
	globalTerminal, err = tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		panic(err)
	}
	defer globalTerminal.Close()

	globalContainer, err = container.New(globalTerminal, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	globalCtx, globalCancel = context.WithCancel(context.Background())
	w, err := newWidgets()
	if err != nil {
		panic(err)
	}
	err = addNextPrevButtons(w)
	if err != nil {
		panic(err)
	}

	gridOpts, err := gridLayout(w, mainDashboard)
	if err != nil {
		panic(err)
	}

	if err := globalContainer.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			globalCancel()
		}
	}
	if err := termdash.Run(globalCtx, globalTerminal, globalContainer, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(RedrawInterval)); err != nil {
		panic(err)
	}
}

// textState creates a rotated state for the text we are displaying.
func textState(text string, capacity, step int) []rune {
	if capacity == 0 {
		return nil
	}

	var state []rune
	for i := 0; i < capacity; i++ {
		state = append(state, ' ')
	}
	state = append(state, []rune(text)...)
	step = step % len(state)
	return rotateRunes(state, step)
}

func newSegmentDisplay(text string) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	chunks := []*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(text,
			segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorNumber(33))),
		),
	}
	if err := sd.Write(chunks); err != nil {
		return nil, err
	}
	return sd, nil
}

// newRotatingSegmentDisplay creates a new SegmentDisplay that shows the dashboard name
func newRotatingSegmentDisplay(text string) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	colors := []cell.Color{
		cell.ColorNumber(33),
		cell.ColorRed,
		cell.ColorYellow,
		cell.ColorNumber(33),
		cell.ColorGreen,
		cell.ColorRed,
		cell.ColorGreen,
		cell.ColorRed,
	}

	step := 0

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		capacity := 0
		t := globalTerminal
		termSize := t.Size()
		for {
			select {
			case <-ticker.C:
				if capacity == 0 {
					// The segment display only knows its capacity after both
					// text size and terminal size are known.
					capacity = sd.Capacity()
				}
				if t.Size().Eq(image.ZP) || !t.Size().Eq(termSize) {
					// Update the capacity initially the first time the
					// terminal reports a non-zero size and then every time the
					// terminal resizes.
					//
					// This is better than updating the capacity on every
					// iteration since that leads to edge cases - segment
					// display capacity depends on the length of text and here
					// we are trying to adjust the text length to the capacity.
					termSize = t.Size()
					capacity = sd.Capacity()
				}

				state := textState(text, capacity, step)
				var chunks []*segmentdisplay.TextChunk
				for i := 0; i < capacity; i++ {
					if i >= len(state) {
						break
					}

					color := colors[i%len(colors)]
					chunks = append(chunks, segmentdisplay.NewChunk(
						string(state[i]),
						segmentdisplay.WriteCellOpts(cell.FgColor(color)),
					))
				}
				if len(chunks) == 0 {
					continue
				}
				if err := sd.Write(chunks); err != nil {
					panic(err)
				}
				step++

			case <-globalCtx.Done():
				return
			}
		}
	}()
	return sd, nil
}

//buttonChunks creates a button chunk with design
func buttonChunks(text string) []*button.TextChunk {
	if len(text) == 0 {
		return nil
	}
	// TODO: Customize outlook of Button
	return []*button.TextChunk{
		button.NewChunk(text,
			button.TextCellOpts(cell.FgColor(cell.ColorWhite)),
			button.FocusedTextCellOpts(cell.FgColor(cell.ColorBlack)),
			button.PressedTextCellOpts(cell.FgColor(cell.ColorBlack)),
		),
	}
}

// newRollText creates a new Text widget that displays rolling text.
func newRollText() (*text.Text, func(string) error, error) {
	t, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		return nil, nil, err
	}

	logToDashBoard := func(message string) error {
		if err := t.Write(fmt.Sprintf("%s\n", message), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(142)))); err != nil {
			return err
		}
		return nil
	}
	return t, logToDashBoard, nil
}

// newSparkLines creates two new sparklines displaying random values.
func newSparkLines() (*sparkline.SparkLine, error) {
	spGreen, err := sparkline.New(
		sparkline.Color(cell.ColorGreen),
	)
	if err != nil {
		return nil, err
	}

	const max = 100
	go Periodic(globalCtx, 250*time.Millisecond, func() error {
		v := int(rand.Int31n(max + 1))
		return spGreen.Add([]int{v})
	})

	return spGreen, nil

}

// newDonut creates a demo Donut widget.
func newDonut() (*donut.Donut, error) {
	d, err := donut.New(donut.CellOpts(
		cell.FgColor(cell.ColorNumber(33))),
	)
	if err != nil {
		return nil, err
	}

	const start = 35
	progress := start

	go Periodic(globalCtx, 500*time.Millisecond, func() error {
		if err := d.Percent(progress); err != nil {
			return err
		}
		progress++
		if progress > 100 {
			progress = start
		}
		return nil
	})
	return d, nil
}

// newBarChart returns a BarcChart that displays random values on multiple bars.
func newBarChart() (*barchart.BarChart, error) {
	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorNumber(33),
			cell.ColorNumber(39),
			cell.ColorNumber(45),
			cell.ColorNumber(51),
			cell.ColorNumber(81),
			cell.ColorNumber(87),
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
		}),
		barchart.ShowValues(),
	)
	if err != nil {
		return nil, err
	}

	const (
		bars = 6
		max  = 100
	)
	values := make([]int, bars)
	go Periodic(globalCtx, 1*time.Second, func() error {
		for i := range values {
			values[i] = int(rand.Int31n(max + 1))
		}

		return bc.Values(values, max)
	})
	return bc, nil
}

// distance is a thread-safe int value used by the newSince method.
// Buttons write it and the line chart reads it.
type distance struct {
	v  int
	mu sync.Mutex
}

// add adds the provided value to the one stored.
func (d *distance) add(v int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.v += v
}

// get returns the current value.
func (d *distance) get() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.v
}

// addGoToMain adds a button to go to main page
func addGoToMain() (*button.Button, error) {
	backToMain, err := button.New("<< Back To Main", func() error {
		w, err := newWidgets()
		if err != nil {
			return err
		}
		err = addNextPrevButtons(w)
		if err != nil {
			return err
		}
		refreshPage(w)
		globalLogFunc("Back to Main clicked")
		return nil
	},
		nextandPrevButtonStyle...,
	)
	if err != nil {
		return nil, err
	}

	return backToMain, nil
}

// addNextPrevButtons adds next and previous buttons
// to every host page
func addNextPrevButtons(w *widgets) error {
	pageLength := len(w.hosts)
	next, err := button.New(">>>", func() error {
		currentHostPage = Next(currentHostPage, pageLength)
		refreshPage(w)
		globalLogFunc(fmt.Sprintf("Moving on to next Page %d", currentHostPage))
		return nil
	},
		nextandPrevButtonStyle...,
	)
	prev, err := button.New("<<<", func() error {
		currentHostPage = Prev(currentHostPage, pageLength)
		refreshPage(w)
		globalLogFunc(fmt.Sprintf("Moving on to previous Page %d", currentHostPage))
		return nil
	},
		nextandPrevButtonStyle...,
	)
	if err != nil {
		return err
	}
	paginationBar := []grid.Element{
		grid.RowHeightPerc(10,
			grid.ColWidthPerc(50,
				grid.Widget(prev),
			),
			grid.ColWidthPerc(50,
				grid.Widget(next),
			),
		),
	}
	// prepend next and previous buttons to each
	// page
	for ind := range w.hosts {
		w.hosts[ind] = append(paginationBar, w.hosts[ind]...)
	}
	return nil
}

// newHostButtons returns all the pages of host buttons
func newHostButtons(paginatedHosts [][]config.Host, metrics []string, w *widgets) ([][]grid.Element, error) {
	buttonGrid := [][]grid.Element{}
	for _, page := range paginatedHosts {
		pageButtons, err := newHostButtonPage(page, metrics, w)
		if err != nil {
			return nil, err
		}
		buttonGrid = append(buttonGrid, pageButtons)
	}
	// create next and previous buttons
	return buttonGrid, nil
}

// newHostButtonPage returns a group of buttons that displays each individual host
// for expansion upon click
func newHostButtonPage(hosts []config.Host, metrics []string, w *widgets) ([]grid.Element, error) {
	buttonGrid := []grid.Element{}
	percentage := 90 / hostsPerPage
	for _, host := range hosts {
		// freeze variables for the closure
		address := host.Address
		inspectorWidgets[address] = map[string]widgetapi.Widget{}
		driver := host.Connection.ToDriver()
		for _, metric := range metrics {
			i, _ := inspector.Init(metric, &driver)
			inspectorWidgets[address][metric] = i.GetWidget()
			go Periodic(globalCtx, 500*time.Millisecond, i.UpdateWidget)
			currentMetric = metric
		}
		aliasText := host.Alias
		if aliasText == "" {
			aliasText = "None"
		}

		hostButton, err := button.NewFromChunks(buttonChunks(host.Address), func() error {
			currentHost = address
			globalLogFunc(fmt.Sprintf("View %s - %s", currentHost, currentMetric))
			setLayout(w, hostMetric)
			return nil
		},
			buttonStyles...,
		)
		driverText, err := text.New()
		driverText.Write(host.Connection.Type)
		alias, err := text.New()
		alias.Write(aliasText)

		if err != nil {
			return nil, err
		}
		buttonGrid = append(buttonGrid,
			grid.RowHeightPerc(percentage,
				grid.ColWidthPerc(34,
					grid.Widget(hostButton,
						singleGridStyle...,
					),
				),
				grid.ColWidthPerc(33,
					grid.Widget(driverText,
						singleGridStyle...,
					)),
				grid.ColWidthPerc(33,
					grid.Widget(alias,
						singleGridStyle...,
					)),
			))
	}
	return buttonGrid, nil
}

// setLayout sets the specified layout.
func setLayout(w *widgets, lt layoutType) error {
	gridOpts, err := gridLayout(w, lt)
	if err != nil {
		return err
	}
	return globalContainer.Update(rootID, gridOpts...)
}

func refreshPage(w *widgets) error {
	gridOpts, err := gridLayout(w, mainDashboard)
	if err != nil {
		return err
	}
	return globalContainer.Update(rootID, gridOpts...)
}

// rotateFloats returns a new slice with inputs rotated by step.
// I.e. for a step of one:
//   inputs[0] -> inputs[len(inputs)-1]
//   inputs[1] -> inputs[0]
// And so on.
func rotateFloats(inputs []float64, step int) []float64 {
	return append(inputs[step:], inputs[:step]...)
}

// rotateRunes returns a new slice with inputs rotated by step.
// I.e. for a step of one:
//   inputs[0] -> inputs[len(inputs)-1]
//   inputs[1] -> inputs[0]
// And so on.
func rotateRunes(inputs []rune, step int) []rune {
	return append(inputs[step:], inputs[:step]...)
}
