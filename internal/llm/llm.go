package llm

import (
	"context"
	"fmt"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/googleai"
	"github.com/rs/zerolog/log"
)

type (
	Model struct {
		model ai.Model
	}

	Config struct {
		ModelName string
		ApiKey    string
	}
)

func NewGoogleAIModel(ctx context.Context, cfg Config) (model *Model, err error) {
	err = googleai.Init(ctx, &googleai.Config{
		APIKey: cfg.ApiKey,
	})
	if err != nil {
		log.Error().Err(err).Msg("unable to initialize Google AI")
		return
	}
	aiModel := googleai.Model(cfg.ModelName)
	if aiModel == nil {
		log.Error().Msg("unable to get Google AI model")
		err = fmt.Errorf("unable to get Google AI model")
		return
	}
	model = &Model{model: aiModel}
	return
}

func (m *Model) Answer(ctx context.Context, question string) (response string, err error) {
	var genResponse *ai.GenerateResponse
	genResponse, err = ai.Generate(ctx, m.model, ai.WithTextPrompt(question))
	if err != nil {
		log.Error().Err(err).Msgf("unable to get answer from model: %v", err.Error())
		return
	}

	response = genResponse.Text()
	return
}
