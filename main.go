package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/lordoverlord/ctop/config"
	"github.com/lordoverlord/ctop/connector"
	"github.com/lordoverlord/ctop/container"
	"github.com/lordoverlord/ctop/cwidgets/compact"
	"github.com/lordoverlord/ctop/logging"
	"github.com/lordoverlord/ctop/widgets"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
)

var (
	build     = "none"
	version   = "dev-build"
	goVersion = runtime.Version()

	log     *logging.CTopLogger
	cursor  *GridCursor
	cGrid   *compact.CompactGrid
	header  *widgets.CTopHeader
	status  *widgets.StatusLine
	errView *widgets.ErrorView

	versionStr = fmt.Sprintf("ctop version %v, build %v %v", version, build, goVersion)
)

func main() {
	defer panicExit()

	// Disable default help flag to avoid conflict with -h
	flag.Usage = printHelp

	// parse command line arguments
	var (
		// Connection options
		hostFlag      string
		contextFlag   string
		connectorFlag string

		// Filtering options
		filterFlag     string
		activeOnlyFlag bool

		// Display options
		sortFieldFlag   string
		reverseSortFlag bool
		invertFlag      bool

		// General options
		versionFlag bool
		helpFlag    bool
	)

	// Connection flags
	flag.StringVar(&hostFlag, "H", "", "")
	flag.StringVar(&hostFlag, "host", "", "")
	flag.StringVar(&contextFlag, "context", "", "")
	flag.StringVar(&connectorFlag, "c", "docker", "")
	flag.StringVar(&connectorFlag, "connector", "docker", "")

	// Filtering flags
	flag.BoolVar(&activeOnlyFlag, "a", false, "")
	flag.BoolVar(&activeOnlyFlag, "all", false, "")
	flag.StringVar(&filterFlag, "f", "", "")
	flag.StringVar(&filterFlag, "filter", "", "")

	// Display flags
	flag.StringVar(&sortFieldFlag, "s", "", "")
	flag.StringVar(&sortFieldFlag, "sort", "", "")
	flag.BoolVar(&reverseSortFlag, "r", false, "")
	flag.BoolVar(&reverseSortFlag, "reverse", false, "")
	flag.BoolVar(&invertFlag, "i", false, "")
	flag.BoolVar(&invertFlag, "invert", false, "")

	// General flags
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.BoolVar(&helpFlag, "help", false, "")

	flag.Parse()

	if versionFlag {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if helpFlag {
		printHelp()
		os.Exit(0)
	}

	// init logger
	log = logging.Init()

	// init global config and read config file if exists
	config.Init()
	if err := config.Read(); err != nil {
		log.Warningf("reading config: %s", err)
	}

	// override default config values with command line flags
	if filterFlag != "" {
		config.Update("filterStr", filterFlag)
	}

	if activeOnlyFlag {
		config.Toggle("allContainers")
	}

	if sortFieldFlag != "" {
		validSort(sortFieldFlag)
		config.Update("sortField", sortFieldFlag)
	}

	if reverseSortFlag {
		config.Toggle("sortReversed")
	}

	// init ui
	if invertFlag {
		InvertColorMap()
	}
	ui.ColorMap = ColorMap // override default colormap
	if err := ui.Init(); err != nil {
		panic(err)
	}
	tm.SetInputMode(tm.InputAlt)

	defer Shutdown()
	// init grid, cursor, header
	cSuper, err := connector.ByNameWithConfig(connectorFlag, hostFlag, contextFlag)
	if err != nil {
		panic(err)
	}
	cursor = &GridCursor{cSuper: cSuper}
	cGrid = compact.NewCompactGrid()
	header = widgets.NewCTopHeader("v" + version)
	status = widgets.NewStatusLine()
	errView = widgets.NewErrorView()

	for {
		exit := Display()
		if exit {
			return
		}
	}
}

func Shutdown() {
	log.Notice("shutting down")
	log.Exit()
	if tm.IsInit {
		ui.Close()
	}
}

// ensure a given sort field is valid
func validSort(s string) {
	if _, ok := container.Sorters[s]; !ok {
		fmt.Printf("invalid sort field: %s\n", s)
		os.Exit(1)
	}
}

func panicExit() {
	if r := recover(); r != nil {
		Shutdown()
		panic(r)
		fmt.Printf("error: %s\n", r)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`ctop - interactive container viewer

USAGE:
  ctop [OPTIONS]

CONNECTION OPTIONS:
  -H, --host HOST         Docker daemon socket or TCP address
                          Examples: tcp://192.168.1.100:2376, unix:///var/run/docker.sock
      --context NAME      Docker context to use
  -c, --connector TYPE    Container connector (default: docker)

FILTERING OPTIONS:
  -a, --all               Show all containers (default: running only)
  -f, --filter PATTERN    Filter containers by name

DISPLAY OPTIONS:
  -s, --sort FIELD        Sort by: name, cpu, mem, net, io (default: name)
  -r, --reverse           Reverse sort order
  -i, --invert            Invert default colours

GENERAL OPTIONS:
  -h, --help              Display this help
  -v, --version           Show version information

`)
	fmt.Printf("Available connectors: %s\n", strings.Join(connector.Enabled(), ", "))
}
