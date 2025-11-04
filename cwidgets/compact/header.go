package compact

import (
	ui "github.com/gizak/termui"
)

type CompactHeader struct {
	X, Y   int
	Width  int
	Height int
	cols   []CompactCol
	widths []int
	pars   []*ui.Par
}

func NewCompactHeader() *CompactHeader {
	return &CompactHeader{
		X:      rowPadding,
		Height: 2,
	}
}

func (row *CompactHeader) GetHeight() int {
	return row.Height
}

func (row *CompactHeader) SetWidths(totalWidth int, widths []int) {
	x := rowPadding // Header always starts at padding
	widthIdx := 0

	for _, w := range row.pars {
		w.SetX(x)
		// Scroll indicators get fixed width of 2
		if w.TextFgColor == ui.ColorYellow {
			w.SetWidth(2)
			x += 2
		} else {
			// Regular columns use calculated widths
			if widthIdx < len(widths) {
				w.SetWidth(widths[widthIdx])
				x += widths[widthIdx]
				// Only add spacing if not the last column
				if widthIdx < len(widths)-1 {
					x += colSpacing
				}
				widthIdx++
			}
		}
	}
	row.Width = totalWidth
}

func (row *CompactHeader) SetX(x int) {
	row.X = x
}

func (row *CompactHeader) SetY(y int) {
	for _, p := range row.pars {
		p.SetY(y)
	}
	row.Y = y
}

func (row *CompactHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, p := range row.pars {
		buf.Merge(p.Buffer())
	}
	return buf
}

func (row *CompactHeader) clearFieldPars() {
	row.pars = []*ui.Par{}
}

func (row *CompactHeader) addFieldPar(s string) {
	p := ui.NewPar(s)
	p.Height = row.Height
	p.Border = false
	row.pars = append(row.pars, p)
}

func (row *CompactHeader) addLeftIndicator() {
	p := ui.NewPar("<")
	p.Height = row.Height
	p.Border = false
	p.TextFgColor = ui.ColorYellow
	row.pars = append(row.pars, p)
}

func (row *CompactHeader) addRightIndicator() {
	p := ui.NewPar(">")
	p.Height = row.Height
	p.Border = false
	p.TextFgColor = ui.ColorYellow
	row.pars = append(row.pars, p)
}

// SetVisibleColumns is a no-op for headers (header manages its own visibility)
func (row *CompactHeader) SetVisibleColumns(offset, count int) {
	// Header manages its own column visibility in updateHeader()
}
