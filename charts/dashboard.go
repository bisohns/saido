package charts

import (
	"context"
	"fmt"
	"image"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/bisohns/saido/config"
	log "github.com/sirupsen/logrus"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/sparkline"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

// widgets holds the widgets used by this demo.
type widgets struct {
	segDist  *segmentdisplay.SegmentDisplay
	input    *textinput.TextInput
	rollT    *text.Text
	spGreen  *sparkline.SparkLine
	barChart *barchart.BarChart
	donut    *donut.Donut
	leftB    *button.Button
	rightB   *button.Button
	sineLC   *linechart.LineChart
	hosts    []grid.Element

	buttons *layoutButtons
}

var logToDashBoard func(string) error

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, t terminalapi.Terminal, c *container.Container, dashboardInfo *config.DashboardInfo) (*widgets, error) {
	sd, err := newSegmentDisplay(ctx, t, dashboardInfo.Title)
	if err != nil {
		return nil, err
	}

	rollT, logToDashBoardfunc, err := newRollText(ctx)
	if err != nil {
		return nil, err
	}
	logToDashBoard = logToDashBoardfunc
	spGreen, err := newSparkLines(ctx)
	if err != nil {
		return nil, err
	}

	bc, err := newBarChart(ctx)
	if err != nil {
		return nil, err
	}

	don, err := newDonut(ctx)
	if err != nil {
		return nil, err
	}

	leftB, rightB, sineLC, err := newSines(ctx)
	if err != nil {
		return nil, err
	}

	hosts, err := newHostButtonGrid(ctx, dashboardInfo.Hosts)
	if err != nil {
		return nil, err
	}

	return &widgets{
		segDist:  sd,
		rollT:    rollT,
		spGreen:  spGreen,
		barChart: bc,
		donut:    don,
		leftB:    leftB,
		rightB:   rightB,
		sineLC:   sineLC,
		hosts:    hosts,
	}, nil
}

// layoutType represents the possible layouts the buttons switch between.
type layoutType int

const (
	// layoutAll displays all the widgets.
	layoutAll layoutType = iota
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
// both produce equivalent layouts for layoutType layoutAll.
func gridLayout(w *widgets, lt layoutType) ([]container.Option, error) {
	leftRows := []grid.Element{
		grid.RowHeightPerc(25,
			grid.Widget(w.segDist,
				container.Border(linestyle.Light),
				container.BorderTitle("Press Esc to quit"),
			),
		),
		grid.RowHeightPerc(5,
			grid.ColWidthPerc(25,
				grid.Widget(w.buttons.allB),
			),
			grid.ColWidthPerc(25,
				grid.Widget(w.buttons.textB),
			),
			grid.ColWidthPerc(25,
				grid.Widget(w.buttons.spB),
			),
			grid.ColWidthPerc(25,
				grid.Widget(w.buttons.lcB),
			),
		),
	}
	switch lt {
	case layoutAll:
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
					w.hosts...,
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

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

// contLayout prepares container options that represent the desired screen layout.
// This function demonstrates the direct use of the container API.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType layoutAll.
// contLayout only produces layoutAll.
func contLayout(w *widgets) ([]container.Option, error) {

	builder := grid.New()
	builder.Add(
		w.hosts...,
	)
	hostContainerOpts, err := builder.Build()
	hostContainerOpts = append(
		hostContainerOpts,
		container.Border(linestyle.Light),
		container.BorderTitle("Hosts"),
	)
	if err != nil {
		return nil, err
	}

	textAndSparks := []container.Option{
		container.SplitVertical(
			container.Left(
				container.Border(linestyle.Light),
				container.BorderTitle("Logs"),
				container.PlaceWidget(w.rollT),
			),
			container.Right(
				hostContainerOpts...,
			),
			container.SplitPercent(40),
		),
	}

	return []container.Option{
		container.SplitHorizontal(
			container.Top(
				container.Border(linestyle.Light),
				container.BorderTitle("Press Esc to quit"),
				container.PlaceWidget(w.segDist),
			),
			container.Bottom(
				container.SplitHorizontal(
					container.Top(),
					container.Bottom(textAndSparks...),
					container.SplitPercent(1),
				),
			),
			container.SplitPercent(20),
		),
	}, nil
}

// rootID is the ID assigned to the root container.
const rootID = "root"

// Terminal implementations
const (
	termboxTerminal = "termbox"
	tcellTerminal   = "tcell"
)

func Main(cfg *config.Config) {
	dashboardInfo := config.GetDashboardInfoConfig(cfg)
	log.Errorf("%v", dashboardInfo)
	log.Debugf("Starting %s", dashboardInfo.Title)
	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		panic(err)
	}
	defer t.Close()

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	w, err := newWidgets(ctx, t, c, dashboardInfo)
	if err != nil {
		panic(err)
	}
	lb, err := newLayoutButtons(c, w)
	if err != nil {
		panic(err)
	}
	w.buttons = lb

	gridOpts, err := gridLayout(w, layoutAll) // equivalent to contLayout(w)
	//  gridOpts, err := contLayout(w)
	if err != nil {
		panic(err)
	}

	if err := c.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(RedrawInterval)); err != nil {
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

// newSegmentDisplay creates a new SegmentDisplay that shows the dashboard name
func newSegmentDisplay(ctx context.Context, t terminalapi.Terminal, text string) (*segmentdisplay.SegmentDisplay, error) {
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

			case <-ctx.Done():
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
func newRollText(ctx context.Context) (*text.Text, func(string) error, error) {
	t, err := text.New(text.RollContent())
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
func newSparkLines(ctx context.Context) (*sparkline.SparkLine, error) {
	spGreen, err := sparkline.New(
		sparkline.Color(cell.ColorGreen),
	)
	if err != nil {
		return nil, err
	}

	const max = 100
	go Periodic(ctx, 250*time.Millisecond, func() error {
		v := int(rand.Int31n(max + 1))
		return spGreen.Add([]int{v})
	})

	return spGreen, nil

}

// newDonut creates a demo Donut widget.
func newDonut(ctx context.Context) (*donut.Donut, error) {
	d, err := donut.New(donut.CellOpts(
		cell.FgColor(cell.ColorNumber(33))),
	)
	if err != nil {
		return nil, err
	}

	const start = 35
	progress := start

	go Periodic(ctx, 500*time.Millisecond, func() error {
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
func newBarChart(ctx context.Context) (*barchart.BarChart, error) {
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
	go Periodic(ctx, 1*time.Second, func() error {
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

// newHostButtonGrid returns a group of buttons that displays each individual host
// for expansion upon click
func newHostButtonGrid(ctx context.Context, hosts []config.Host) ([]grid.Element, error) {
	buttonGrid := []grid.Element{}
	percentage := 100 / len(hosts)
	for _, host := range hosts {
		// freeze addresss for the closure
		address := host.Address
		aliasText := host.Alias
		if aliasText == "" {
			aliasText = "None"
		}
		hostButton, err := button.NewFromChunks(buttonChunks(host.Address), func() error {
			logToDashBoard(address)
			return nil
		},
			button.Key(keyboard.KeyEnter),
			button.DisableShadow(),
			button.Height(1),
			button.TextHorizontalPadding(0),
			button.FillColor(cell.ColorBlack),
			button.FocusedFillColor(cell.ColorNumber(117)),
			button.PressedFillColor(cell.ColorNumber(220)),
		)
		driver, err := text.New()
		driver.Write(host.Connection.Type)
		alias, err := text.New()
		alias.Write(aliasText)

		if err != nil {
			return nil, err
		} else {
			commonStyles := []container.Option{
				container.AlignHorizontal(align.HorizontalLeft),
				container.PaddingLeft(5),
				container.Border(linestyle.Round),
			}

			buttonGrid = append(buttonGrid,
				grid.RowHeightPerc(percentage,
					grid.ColWidthPerc(34,
						grid.Widget(hostButton,
							commonStyles...,
						),
					),
					grid.ColWidthPerc(33,
						grid.Widget(driver,
							commonStyles...,
						)),
					grid.ColWidthPerc(33,
						grid.Widget(alias,
							commonStyles...,
						)),
				))
		}
	}
	return buttonGrid, nil
}

// newSines returns a line chart that displays multiple sine series and two buttons.
// The left button shifts the second series relative to the first series to
// the left and the right button shifts it to the right.
func newSines(ctx context.Context) (left, right *button.Button, lc *linechart.LineChart, err error) {
	var inputs []float64
	for i := 0; i < 200; i++ {
		v := math.Sin(float64(i) / 100 * math.Pi)
		inputs = append(inputs, v)
	}

	sineLc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		return nil, nil, nil, err
	}
	step1 := 0
	secondDist := &distance{v: 100}
	go Periodic(ctx, RedrawInterval/3, func() error {
		step1 = (step1 + 1) % len(inputs)
		if err := lc.Series("first", rotateFloats(inputs, step1),
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(33))),
		); err != nil {
			return err
		}

		step2 := (step1 + secondDist.get()) % len(inputs)
		return lc.Series("second", rotateFloats(inputs, step2), linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite)))
	})

	// diff is the difference a single button press adds or removes to the
	// second series.
	const diff = 20
	leftB, err := button.New("(l)eft", func() error {
		secondDist.add(diff)
		return nil
	},
		button.GlobalKey('l'),
		button.WidthFor("(r)ight"),
		button.FillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	rightB, err := button.New("(r)ight", func() error {
		secondDist.add(-diff)
		return nil
	},
		button.GlobalKey('r'),
		button.FillColor(cell.ColorNumber(196)),
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return leftB, rightB, sineLc, nil
}

// setLayout sets the specified layout.
func setLayout(c *container.Container, w *widgets, lt layoutType) error {
	gridOpts, err := gridLayout(w, lt)
	if err != nil {
		return err
	}
	return c.Update(rootID, gridOpts...)
}

// layoutButtons are buttons that change the layout.
type layoutButtons struct {
	allB  *button.Button
	textB *button.Button
	spB   *button.Button
	lcB   *button.Button
}

// newLayoutButtons returns buttons that dynamically switch the layouts.
func newLayoutButtons(c *container.Container, w *widgets) (*layoutButtons, error) {
	opts := []button.Option{
		button.WidthFor("sparklines"),
		button.FillColor(cell.ColorNumber(220)),
		button.Height(1),
	}

	allB, err := button.New("all", func() error {
		return setLayout(c, w, layoutAll)
	}, opts...)
	if err != nil {
		return nil, err
	}

	textB, err := button.New("text", func() error {
		return setLayout(c, w, layoutText)
	}, opts...)
	if err != nil {
		return nil, err
	}

	spB, err := button.New("sparklines", func() error {
		return setLayout(c, w, layoutSparkLines)
	}, opts...)
	if err != nil {
		return nil, err
	}

	lcB, err := button.New("linechart", func() error {
		return setLayout(c, w, layoutLineChart)
	}, opts...)
	if err != nil {
		return nil, err
	}

	return &layoutButtons{
		allB:  allB,
		textB: textB,
		spB:   spB,
		lcB:   lcB,
	}, nil
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
