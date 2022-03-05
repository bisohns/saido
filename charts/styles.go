package charts

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/button"
)

var buttonStyles = []button.Option{
	button.Key(keyboard.KeyEnter),
	button.DisableShadow(),
	button.Height(1),
	button.TextHorizontalPadding(0),
	button.FillColor(cell.ColorBlack),
	button.FocusedFillColor(cell.ColorNumber(117)),
	button.PressedFillColor(cell.ColorNumber(220)),
}

var singleGridStyle = []container.Option{
	container.AlignHorizontal(align.HorizontalLeft),
	container.PaddingLeft(5),
	container.Border(linestyle.Round),
}

var nextandPrevButtonStyle = append(buttonStyles,
	button.TextColor(cell.ColorWhite),
	button.TextHorizontalPadding(10),
)
