package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/kkyr/fig"
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/internal/config"
	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/ui"
	"github.com/prgrs/clickup/ui/context"
	"github.com/spf13/pflag"
)

const (
	// Version is the current version
	AppVersion = "0.0.1-pre-alpha"
	// Name is the name of the application
	AppName = "clickup-tui"
	// Description is the description of the application
	AppDescription = "A terminal user interface for ClickUp"
	// DefaultConfigFilename is the default config filename
	DefaultConfigFilename = "config.yaml"
	// DefaultCachePath is the default cache path
	DefaultCachePath = "./cache"

	ErrMissingToken = `ClickUp API token is required. Follow the steps: [ClickUp API docs: Generate your personal API token](https://clickup.com/api/developer-portal/authentication/#generate-your-personal-api-token) and please set it in the config file. See https://docs.clickup.com/en/articles/1367130-getting-started-with-the-clickup-api")`
)

var (
	flag               *pflag.FlagSet = pflag.NewFlagSet(AppName, pflag.ContinueOnError)
	flagDebug          *bool          = flag.Bool("debug", false, "Enable debug mode")
	flagDebugDeep      *bool          = flag.Bool("debug-deep", false, "Enable deep debug mode")
	flagHelp           *bool          = flag.BoolP("help", "h", false, "Show help")
	flagVersion        *bool          = flag.BoolP("version", "v", false, "Show version")
	flagConfigFilename *string        = flag.StringP("config", "c", DefaultConfigFilename, "A config filename")
	flagCleanCache     *bool          = flag.Bool("clean-cache", false, "Cleans cache data")
	flagCleanCacheOnly *bool          = flag.Bool("clean-cache-only", false, "Cleans cache data and exits")
	flagCachePath      *string        = flag.String("cache-path", DefaultCachePath, "The path to the cache directory")

	flagUsage func() string = func() string {
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("%s - %s\n", AppName, AppDescription))
		s.WriteString("Usage:\n")
		s.WriteString(fmt.Sprintf("  %s [flags]\n", AppName))
		s.WriteString("Flags:\n")
		s.WriteString(flag.FlagUsages())
		return s.String()
	}

	termLogger *log.Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: false,
		Level:        log.InfoLevel,
	})
)

func main() {
	if err := flag.Parse(os.Args[1:]); err != nil {
		fmt.Println(flagUsage())
		os.Exit(2)
	}

	if *flagHelp {
		fmt.Println(flagUsage())
		return
	}

	if *flagVersion {
		fmt.Printf("%s %s\n", AppName, AppVersion)
		return
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: *flagDebugDeep,
		Level: func() log.Level {
			lvl := log.InfoLevel
			if *flagDebug {
				lvl = log.DebugLevel
			}
			return lvl
		}(),
		ReportTimestamp: true,
	})

	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()
	logger.SetOutput(f)

	termLogger.SetOutput(io.MultiWriter(os.Stderr, f))

	logger.Info("Starting up...")

	logger.Info("Initializing config...")
	var cfg config.Config
	if err := fig.Load(&cfg,
		fig.File(*flagConfigFilename),
		fig.Dirs(
			".",
			"/etc/clickup-tui",
			"/home/user/clickup-tui",
			"$HOME/.config/clickup-tui",
		),
	); err != nil {
		termLogger.Fatal(err)
	}

	if cfg.Token == "" {
		termLogger.Fatal(ErrMissingToken)
	}

	logger.Info("Initializing cache...")
	cache := cache.NewCache(
		slog.New(logger.WithPrefix("Cache")),
		*flagCachePath,
	)

	defer func() {
		if err := cache.Dump(); err != nil {
			termLogger.Fatal(err)
		}
	}()

	if *flagCleanCache || *flagCleanCacheOnly {
		logger.Info("Cleaning cache...")
		if err := cache.Invalidate(); err != nil {
			termLogger.Fatal(err)
		}

		if *flagCleanCacheOnly {
			return
		}
	}

	if err := cache.Load(); err != nil {
		termLogger.Fatal(err)
	}

	logger.Info("Initializing api...")
	api := api.NewApi(logger, cache, cfg.Token)

	logger.Info("Initializing user context...")
	ctx := context.NewUserContext(logger, &api, &cfg)

	logger.Info("Initializing main model...")
	mainModel := ui.InitialModel(&ctx, logger)

	logger.Info("Initializing program...")
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		termLogger.Fatal(err)
	}
}
