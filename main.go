package main

import (
	"fmt"
	"github.com/a-peyrard/swim-spot-checker/internal/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

const (
	swimSchoolURL                 = "https://andersonswim.com/"
	noAvailabilityPatternsDefault = "No Single Lessons"
)

var (
	skipNotifications bool
)

var swimSpotCheckerCmd = &cobra.Command{
	Use:   "swim-spot-checker",
	Short: "Check if there are some opening for swimming lessons",
	Long:  `Check if there are some opening for swimming lessons`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) > 0 && args[0] == "completion" {
			shell := "zsh"
			if len(args) > 1 {
				shell = args[1]
			}
			return handleCompletion(shell, cmd)
		}

		noAvailabilityPatternsRaw := os.Getenv("NO_AVAILABILITY_PATTERNS")
		if noAvailabilityPatternsRaw == "" {
			noAvailabilityPatternsRaw = noAvailabilityPatternsDefault
		}
		noAvailabilityPatterns := strings.Split(noAvailabilityPatternsRaw, ",")

		log.Info().Msg("Checking for spots...")
		var startExecTime = time.Now()
		var foundSpot bool
		foundSpot, err = checkAvailability(swimSchoolURL, noAvailabilityPatterns)
		if err != nil {
			return
		}

		if foundSpot {
			log.Info().Msgf("Spot found (in %s)", time.Since(startExecTime))
			if !skipNotifications {
				notifySpotFound()
			}
		} else {
			log.Info().Msgf("No spot found (in %s)", time.Since(startExecTime))
		}

		return
	},
}

func notifySpotFound() {
	log.Info().Msg("Sending notification")
	// Send a notification to your phone
	// fixme: use twilio or something
}

func checkAvailability(url string, patterns []string) (foundSpot bool, err error) {
	var page string
	page, err = http.ScrapURL(url)
	if err != nil {
		return
	}

	foundSpot = true
	for _, pattern := range patterns {
		if strings.Contains(page, pattern) {
			foundSpot = false
			break
		}
	}

	return
}

func handleCompletion(shell string, cmd *cobra.Command) error {
	switch shell {
	case "bash":
		return cmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.GenFishCompletion(os.Stdout, true)
	default:
		return cmd.Help()
	}
}

func init() {
	swimSpotCheckerCmd.Flags().BoolVar(
		&skipNotifications,
		"skip-notifications",
		false,
		"don't send notifications, just print result",
	)
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	if err := swimSpotCheckerCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
