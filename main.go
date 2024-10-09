package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-peyrard/swim-spot-checker/internal/http"
	"github.com/a-peyrard/swim-spot-checker/internal/llm"
	"github.com/a-peyrard/swim-spot-checker/internal/swim"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

const (
	swimSchoolURL   = "https://andersonswim.com/"
	previousContent = "/tmp/previous_content"
	startMarker     = "Welcome to Anderson’s Swim School"
	endMarker       = "Average email response time"
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

		log.Info().Msg("Initialize model...")
		var apiKey = os.Getenv("LLM_API_KEY")
		cfg := llm.Config{
			ModelName: "gemini-1.5-flash",
			ApiKey:    apiKey,
		}
		var model *llm.Model
		model, err = llm.NewGoogleAIModel(context.Background(), cfg)
		if err != nil {
			log.Error().Err(err).Msg("unable to initialize model")
			return
		}

		log.Info().Msg("Checking for spots...")
		var (
			startExecTime = time.Now()
			foundSpot     bool
			explanation   string
		)
		foundSpot, explanation, err = checkAvailability(swimSchoolURL, model)
		if err != nil {
			return
		}

		if foundSpot {
			log.Info().Msgf("Spot found (in %s)", time.Since(startExecTime))
			if !skipNotifications {
				notifySpotFound(explanation)
			}
		} else {
			log.Info().Msgf("No spot found (in %s)", time.Since(startExecTime))
		}

		return
	},
}

func notifySpotFound(explanation string) {
	log.Info().Msgf("Sending notification, we found availability: %s", explanation)
	// Send a notification to your phone
	// fixme: use twilio or something
}

func checkAvailability(url string, model *llm.Model) (foundSpot bool, explanation string, err error) {
	var (
		oldContent string
		newContent string
	)
	newContent, err = extractContentFromURL(url)
	if err != nil {
		return
	}

	oldContent, err = loadPreviousContent()
	if err != nil {
		return
	}

	if oldContent == "" {
		log.Info().Msgf("No previous content found, saving current content")
		err = storePreviousContent(newContent)
		return
	}

	if oldContent == newContent {
		log.Info().Msg("No change in content")
		return
	}

	foundSpot, explanation, err = swim.CheckAvailability(context.Background(), model, oldContent, newContent)
	if err == nil {
		log.Info().Msgf("Availability check result: %t", foundSpot)
		log.Info().Msgf("Explanation: %s", explanation)

		err = storePreviousContent(newContent)
	}

	return
}

func extractContentFromURL(url string) (content string, err error) {
	content, err = http.ScrapURL(url)
	if err != nil {
		return
	}

	startIndex := strings.Index(content, startMarker)
	if startIndex == -1 {
		err = fmt.Errorf("start marker not found")
		return
	}
	endIndex := strings.Index(content, endMarker)
	if endIndex == -1 {
		err = fmt.Errorf("end marker not found")
		return
	}
	content = strings.TrimSpace(content[startIndex+len(startMarker) : endIndex])

	return
}

func loadPreviousContent() (content string, err error) {
	var b []byte
	b, err = os.ReadFile(previousContent)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return
	}
	content = string(b)
	return
}

func storePreviousContent(content string) (err error) {
	err = os.WriteFile(previousContent, []byte(content), 0644)
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
