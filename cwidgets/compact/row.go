package compact

import (
	"github.com/lordoverlord/ctop/config"
	"github.com/lordoverlord/ctop/logging"
	"github.com/lordoverlord/ctop/models"

	ui "github.com/gizak/termui"
)

const rowPadding = 1

var log = logging.Init()

type RowBufferer interface {
	SetX(int)
	SetY(int)
	SetWidths(int, []int)
	SetVisibleColumns(offset, count int)
	GetHeight() int
	Buffer() ui.Buffer
}

type CompactRow struct {
	Bg            *RowBg
	Cols          []CompactCol
	X, Y          int
	Height        int
	widths        []int // column widths
	visibleOffset int   // starting column index for visible columns
	visibleCount  int   // number of visible columns
}

func NewCompactRow() *CompactRow {
	row := &CompactRow{
		Bg:            NewRowBg(),
		Cols:          newRowWidgets(),
		X:             rowPadding,
		Height:        1,
		visibleOffset: 0,
		visibleCount:  0, // Will be set by SetVisibleColumns
	}

	return row
}

// SetVisibleColumns sets which columns should be rendered
func (row *CompactRow) SetVisibleColumns(offset, count int) {
	row.visibleOffset = offset
	row.visibleCount = count
}

// visibleCols returns the slice of columns that should be rendered
func (row *CompactRow) visibleCols() []CompactCol {
	if row.visibleCount == 0 {
		return row.Cols // Show all if not set
	}
	end := row.visibleOffset + row.visibleCount
	if end > len(row.Cols) {
		end = len(row.Cols)
	}
	return row.Cols[row.visibleOffset:end]
}

func (row *CompactRow) SetMeta(m models.Meta) {
	for _, w := range row.Cols {
		w.SetMeta(m)
	}
}

func (row *CompactRow) SetMetrics(m models.Metrics) {
	for _, w := range row.Cols {
		w.SetMetrics(m)
	}
}

// Set gauges, counters, etc. to default unread values
func (row *CompactRow) Reset() {
	for _, w := range row.Cols {
		w.Reset()
	}
}

func (row *CompactRow) GetHeight() int { return row.Height }

func (row *CompactRow) SetX(x int) { row.X = x }

func (row *CompactRow) SetY(y int) {
	if y == row.Y {
		return
	}

	row.Bg.Y = y
	for _, w := range row.Cols {
		w.SetY(y)
	}
	row.Y = y
}

func (row *CompactRow) SetWidths(totalWidth int, widths []int) {
	x := row.X // Use the X position set by SetX

	row.Bg.SetX(rowPadding) // Background always starts at padding
	row.Bg.SetWidth(totalWidth)

	visibleCols := row.visibleCols()
	for n, w := range visibleCols {
		if n >= len(widths) {
			break
		}
		w.SetX(x)
		w.SetWidth(widths[n])
		x += widths[n]
		// Only add spacing if not the last column
		if n < len(visibleCols)-1 && n < len(widths)-1 {
			x += colSpacing
		}
	}
}

func (row *CompactRow) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(row.Bg.Buffer())
	for _, w := range row.visibleCols() {
		buf.Merge(w.Buffer())
	}
	return buf
}

func (row *CompactRow) Highlight() {
	row.Cols[1].Highlight()
	if config.GetSwitchVal("fullRowCursor") {
		for _, w := range row.Cols {
			w.Highlight()
		}
	}
}

func (row *CompactRow) UnHighlight() {
	row.Cols[1].UnHighlight()
	if config.GetSwitchVal("fullRowCursor") {
		for _, w := range row.Cols {
			w.UnHighlight()
		}
	}
}

type RowBg struct {
	*ui.Par
}

func NewRowBg() *RowBg {
	bg := ui.NewPar("")
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ThemeAttr("par.text.bg")
	return &RowBg{bg}
}

func (w *RowBg) Highlight()   { w.Bg = ui.ThemeAttr("par.text.fg") }
func (w *RowBg) UnHighlight() { w.Bg = ui.ThemeAttr("par.text.bg") }
