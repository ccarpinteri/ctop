package compact

import (
	ui "github.com/gizak/termui"
)

type CompactGrid struct {
	ui.GridBufferer
	header  *CompactHeader
	cols    []CompactCol // reference columns
	Rows    []RowBufferer
	X, Y    int
	Width   int
	Height  int
	Offset  int // starting row offset (vertical)
	HOffset int // starting column offset (horizontal)
}

func NewCompactGrid() *CompactGrid {
	cg := &CompactGrid{header: NewCompactHeader()}
	cg.rebuildHeader()
	return cg
}

func (cg *CompactGrid) Align() {
	y := cg.Y

	if cg.Offset >= len(cg.Rows) || cg.Offset < 0 {
		cg.Offset = 0
	}

	// Rebuild header with visible columns and scroll indicators
	cg.updateHeader()

	// Calculate visible columns
	visibleCols := cg.visibleCols()
	visibleCount := len(visibleCols)

	// Calculate starting X position for data rows (account for left scroll indicator)
	startX := rowPadding
	if cg.HasLeftScroll() {
		startX += 2 // Space for "<" indicator
	}

	// update row ypos, width recursively
	colWidths := cg.calcWidths()
	for _, r := range cg.pageRows() {
		r.SetVisibleColumns(cg.HOffset, visibleCount)
		// Header manages its own X through indicator Pars, data rows need explicit X
		if _, isHeader := r.(*CompactHeader); !isHeader {
			r.SetX(startX)
		}
		r.SetY(y)
		y += r.GetHeight()
		r.SetWidths(cg.Width, colWidths)
	}
}

func (cg *CompactGrid) Clear() {
	cg.Rows = []RowBufferer{}
	cg.rebuildHeader()
}

func (cg *CompactGrid) GetHeight() int { return len(cg.Rows) + cg.header.Height }
func (cg *CompactGrid) SetX(x int)     { cg.X = x }
func (cg *CompactGrid) SetY(y int)     { cg.Y = y }
func (cg *CompactGrid) SetWidth(w int) { cg.Width = w }
func (cg *CompactGrid) MaxRows() int   { return ui.TermHeight() - cg.header.Height - cg.Y }
func (cg *CompactGrid) MaxCols() int   { return len(cg.cols) }

// ScrollLeft moves the horizontal view left by one column
func (cg *CompactGrid) ScrollLeft() {
	if cg.HOffset > 0 {
		cg.HOffset--
	}
}

// ScrollRight moves the horizontal view right by one column
func (cg *CompactGrid) ScrollRight() {
	visibleCols := cg.visibleCols()
	if cg.HOffset+len(visibleCols) < len(cg.cols) {
		cg.HOffset++
	}
}

// HasLeftScroll returns true if there are columns hidden to the left
func (cg *CompactGrid) HasLeftScroll() bool {
	return cg.HOffset > 0
}

// HasRightScroll returns true if there are columns hidden to the right
func (cg *CompactGrid) HasRightScroll() bool {
	visibleCols := cg.visibleCols()
	return cg.HOffset+len(visibleCols) < len(cg.cols)
}

// visibleCols returns the slice of columns that fit within the current width
// starting from HOffset
func (cg *CompactGrid) visibleCols() []CompactCol {
	if cg.HOffset >= len(cg.cols) {
		cg.HOffset = 0
	}

	// Fast path: if at start, check if all columns fit at minimum widths
	if cg.HOffset == 0 {
		totalMinWidth := 0

		// Use minimum widths for all columns
		for _, col := range cg.cols {
			if col.FixedWidth() > 0 {
				totalMinWidth += col.FixedWidth()
			} else {
				totalMinWidth += col.MinWidth()
			}
		}

		// Add spacing (between columns, not after last)
		if len(cg.cols) > 1 {
			totalMinWidth += colSpacing * (len(cg.cols) - 1)
		}

		// Account for rowPadding and right padding
		// Reserve 3 chars to ensure visual padding after last column
		availableWidth := cg.Width - rowPadding - 3

		// If all columns fit at minimum widths, show all columns
		// calcWidths() will distribute extra space according to priorities
		if totalMinWidth <= availableWidth {
			return cg.cols
		}
	}

	// Need to scroll - calculate which columns are visible
	var widthUsed int
	visible := []CompactCol{}

	// Account for rowPadding and right padding
	// Reserve 3 chars to ensure visual padding after last column
	availableWidth := cg.Width - rowPadding - 3

	// Reserve space for scroll indicators
	indicatorWidth := 0
	if cg.HOffset > 0 {
		indicatorWidth += 2 // "< " indicator
	}

	for i := cg.HOffset; i < len(cg.cols); i++ {
		col := cg.cols[i]
		// Use minimum width to determine if column fits
		// (columns will shrink to min if needed)
		colWidth := col.FixedWidth()
		if colWidth == 0 {
			colWidth = col.MinWidth()
		}

		// Add spacing before this column (except for first visible column)
		spacingNeeded := 0
		if len(visible) > 0 {
			spacingNeeded = colSpacing
		}

		// Check if this column fits
		potentialWidth := widthUsed + spacingNeeded + colWidth

		// Reserve space for right indicator if more columns remain
		if i+1 < len(cg.cols) {
			potentialWidth += 2
		}

		if potentialWidth+indicatorWidth > availableWidth && len(visible) > 0 {
			break
		}

		visible = append(visible, col)
		widthUsed += spacingNeeded + colWidth
	}

	// Ensure at least one column is visible
	if len(visible) == 0 && len(cg.cols) > 0 {
		visible = []CompactCol{cg.cols[cg.HOffset]}
	}

	return visible
}

// calculate and return per-column width for visible columns only
// Uses smart distribution: allocates minimums first, then distributes extra space by priority
func (cg *CompactGrid) calcWidths() []int {
	visibleCols := cg.visibleCols()
	if len(visibleCols) == 0 {
		return []int{}
	}

	// Account for rowPadding and right padding
	// Reserve 3 chars to ensure visual padding after last column
	availableWidth := cg.Width - rowPadding - 3
	colWidths := make([]int, len(visibleCols))

	// Reserve space for scroll indicators
	if cg.HasLeftScroll() {
		availableWidth -= 2 // "< " indicator
	}
	if cg.HasRightScroll() {
		availableWidth -= 2 // " >" indicator
	}

	// Reserve space for column spacing (between columns, not after last)
	if len(visibleCols) > 1 {
		spacing := colSpacing * (len(visibleCols) - 1)
		availableWidth -= spacing
	}

	// Step 1: Allocate minimum widths (or fixed widths for fixed columns)
	for n, col := range visibleCols {
		if col.FixedWidth() > 0 {
			// Fixed width columns use their fixed width
			colWidths[n] = col.FixedWidth()
			availableWidth -= col.FixedWidth()
		} else {
			// Growable columns start at minimum width
			colWidths[n] = col.MinWidth()
			availableWidth -= col.MinWidth()
		}
	}

	// Prevent negative available width
	if availableWidth < 0 {
		availableWidth = 0
	}

	// Step 2: Distribute extra space by priority (1 → 2)
	for priority := 1; priority <= 2 && availableWidth > 0; priority++ {
		// Find growable columns at this priority level that can still grow
		canGrow := make(map[int]int) // map of column index → max additional width
		for n, col := range visibleCols {
			if col.GrowPriority() == priority && col.FixedWidth() == 0 {
				maxWidth := col.MaxWidth()
				currentWidth := colWidths[n]

				if maxWidth > 0 {
					remaining := maxWidth - currentWidth
					if remaining > 0 {
						canGrow[n] = remaining
					}
				} else {
					// No max limit - can grow indefinitely
					canGrow[n] = availableWidth
				}
			}
		}

		if len(canGrow) == 0 {
			continue // No columns can grow at this priority
		}

		// Distribute available width equally, respecting max widths
		// Keep distributing until no space left or all columns hit max
		for len(canGrow) > 0 && availableWidth > 0 {
			perCol := availableWidth / len(canGrow)
			if perCol == 0 {
				perCol = 1 // Give at least 1 width if space available
			}

			distributed := 0
			for n, maxGrowth := range canGrow {
				actualGrowth := perCol
				if actualGrowth > maxGrowth {
					actualGrowth = maxGrowth
				}
				if actualGrowth > availableWidth {
					actualGrowth = availableWidth
				}

				colWidths[n] += actualGrowth
				distributed += actualGrowth

				// Remove from canGrow if hit max
				if actualGrowth >= maxGrowth {
					delete(canGrow, n)
				} else {
					canGrow[n] -= actualGrowth
				}
			}

			availableWidth -= distributed
			if distributed == 0 {
				break // Couldn't distribute any more
			}
		}
	}

	return colWidths
}

func (cg *CompactGrid) pageRows() (rows []RowBufferer) {
	rows = append(rows, cg.header)
	rows = append(rows, cg.Rows[cg.Offset:]...)
	return rows
}

func (cg *CompactGrid) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	for _, r := range cg.pageRows() {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (cg *CompactGrid) AddRows(rows ...RowBufferer) {
	cg.Rows = append(cg.Rows, rows...)
}

func (cg *CompactGrid) rebuildHeader() {
	cg.cols = newRowWidgets()
	cg.updateHeader()
}

// updateHeader rebuilds the header with only visible columns
func (cg *CompactGrid) updateHeader() {
	cg.header.clearFieldPars()

	// Add left scroll indicator if needed
	if cg.HasLeftScroll() {
		cg.header.addLeftIndicator()
	}

	// Add visible columns
	for _, col := range cg.visibleCols() {
		cg.header.addFieldPar(col.Header())
	}

	// Add right scroll indicator if needed
	if cg.HasRightScroll() {
		cg.header.addRightIndicator()
	}
}
