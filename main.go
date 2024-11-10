package main

import (
	"context"
	"fmt"
	"github.com/a-peyrard/swim-spot-checker/internal/llm"
	"github.com/a-peyrard/swim-spot-checker/internal/notification"
	"github.com/a-peyrard/swim-spot-checker/internal/swim"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	swimSchoolURL = "https://andersonswim.com/"
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

		notifier := notification.NewEmailNotifier()
		recipient := notification.Recipient{
			PhoneNumber: os.Getenv("TO_PHONE_NUMBER"),
			Carrier:     os.Getenv("TO_PHONE_CARRIER"),
		}

		c := cron.New()
		rawCron := os.Getenv("SCHEDULE")
		schedule, err := cron.ParseStandard(rawCron)
		if err != nil {
			log.Err(err).Msg("Failed to parse schedule")
			return
		}
		c.Schedule(schedule, cron.FuncJob(func() {
			err := check(model, notifier, recipient)
			if err != nil {
				log.Err(err).Msg("Failed to check for spots")
			}
			log.Info().Msgf("Next execution at %s", schedule.Next(time.Now()))
		}))
		log.Info().Msgf("Scheduled to run with cron %s, next execution at %s", rawCron, schedule.Next(time.Now()))
		c.Run()

		log.Info().Msg("Bye!")
		return
	},
}

func check(model *llm.Model, notifier notification.Notifier, recipient notification.Recipient) (err error) {
	log.Info().Msg("Checking for spots...")
	var (
		startExecTime = time.Now()
		foundSpot     bool
		explanation   string
	)
	foundSpot, explanation, err = swim.CheckAvailability(swimSchoolURL, model)
	if err != nil {
		return
	}

	if foundSpot {
		log.Info().Msgf("Spot found (in %s)", time.Since(startExecTime))
		if !skipNotifications {
			err = notifier.Text(
				notification.Sms{Body: fmt.Sprintf("%s\n\nGo check it out: %s", explanation, swimSchoolURL)},
				recipient,
			)
			if err != nil {
				log.Err(err).Msg("Failed to send SMS")
			} else {
				log.Info().Msg("SMS sent successfully!")
			}
		}
	} else {
		log.Info().Msgf("No spot found (in %s)", time.Since(startExecTime))
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
