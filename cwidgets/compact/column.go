package compact

import (
	"github.com/lordoverlord/ctop/config"
	"github.com/lordoverlord/ctop/models"

	ui "github.com/gizak/termui"
)

var (
	allCols = map[string]NewCompactColFn{
		"status":  NewStatus,
		"name":    NewNameCol,
		"id":      NewCIDCol,
		"image":   NewImageCol,
		"ports":   NewPortsCol,
		"IPs":     NewIpsCol,
		"created": NewCreatedCol,
		"cpu":     NewCPUCol,
		"cpus":    NewCpuScaledCol,
		"mem":     NewMemCol,
		"net":     NewNetCol,
		"io":      NewIOCol,
		"pids":    NewPIDCol,
		"uptime":  NewUptimeCol,
	}
)

type NewCompactColFn func() CompactCol

func newRowWidgets() []CompactCol {
	enabled := config.EnabledColumns()
	cols := make([]CompactCol, len(enabled))

	for n, name := range enabled {
		wFn, ok := allCols[name]
		if !ok {
			panic("no such widget name: %s" + name)
		}
		cols[n] = wFn()
	}

	return cols
}

type CompactCol interface {
	ui.GridBufferer
	Reset()
	Header() string       // header text to display for column
	FixedWidth() int      // fixed width size. if == 0, width is automatically calculated
	MinWidth() int        // minimum width for growable columns
	MaxWidth() int        // maximum width for growable columns
	GrowPriority() int    // priority for growing (1=highest, 0=fixed)
	Highlight()
	UnHighlight()
	SetMeta(models.Meta)
	SetMetrics(models.Metrics)
}
