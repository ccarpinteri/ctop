package widgets

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

type CTopHeader struct {
	Time    *ui.Par
	Count   *ui.Par
	Filter  *ui.Par
	bg      *ui.Par
	version string
}

func NewCTopHeader(version string) *CTopHeader {
	return &CTopHeader{
		Time:    headerParWide(2, "", 40),  // Wider for version string
		Count:   headerPar(44, "-"),
		Filter:  headerPar(66, ""),
		bg:      headerBg(),
		version: version,
	}
}

func (c *CTopHeader) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	c.Time.Text = c.timeStr()
	buf.Merge(c.bg.Buffer())
	buf.Merge(c.Time.Buffer())
	buf.Merge(c.Count.Buffer())
	buf.Merge(c.Filter.Buffer())
	return buf
}

func (c *CTopHeader) Align() {
	c.bg.SetWidth(ui.TermWidth() - 1)
}

func (c *CTopHeader) Height() int {
	return c.bg.Height
}

func headerBgBordered() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Height = 3
	bg.Bg = ui.ThemeAttr("header.bg")
	return bg
}

func headerBg() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Height = 1
	bg.Border = false
	bg.Bg = ui.ThemeAttr("header.bg")
	return bg
}

func (c *CTopHeader) SetCount(val int) {
	c.Count.Text = fmt.Sprintf("%d containers", val)
}

func (c *CTopHeader) SetFilter(val string) {
	if val == "" {
		c.Filter.Text = ""
	} else {
		c.Filter.Text = fmt.Sprintf("filter: %s", val)
	}
}

func (c *CTopHeader) timeStr() string {
	ts := time.Now().Local().Format("15:04:05 MST")
	return fmt.Sprintf("ctop %s - %s", c.version, ts)
}

func headerPar(x int, s string) *ui.Par {
	return headerParWide(x, s, 20)
}

func headerParWide(x int, s string, width int) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.X = x
	p.Border = false
	p.Height = 1
	p.Width = width
	p.Bg = ui.ThemeAttr("header.bg")
	p.TextFgColor = ui.ThemeAttr("header.fg")
	p.TextBgColor = ui.ThemeAttr("header.bg")
	return p
}
