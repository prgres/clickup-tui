package main

import (
	"fmt"
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

	f, err := tea.LogToFileWith("debug.log", logger.GetPrefix(), logger)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger.Info("Starting up...")

	logger.Info("Initializing config...")
	var cfg config.Config
	if err := fig.Load(&cfg,
		fig.File(*flagConfigFilename),
		fig.Dirs(
			".",
			"/etc/myapp",
			"/home/user/myapp",
			"$HOME/.config/clickup-tui",
		),
	); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Initializing cache...")
	cache := cache.NewCache(
		slog.New(logger.WithPrefix("Cache")),
		*flagCachePath,
	)

	defer func() {
		if err := cache.Dump(); err != nil {
			logger.Fatal(err)
		}
	}()

	if *flagCleanCache || *flagCleanCacheOnly {
		logger.Info("Cleaning cache...")
		if err := cache.Invalidate(); err != nil {
			logger.Fatal(err)
		}

		if *flagCleanCacheOnly {
			return
		}
	}

	if err := cache.Load(); err != nil {
		logger.Fatal(err)
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
		logger.Fatal(err)
	}
}
