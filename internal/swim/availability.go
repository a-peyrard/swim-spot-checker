package swim

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-peyrard/swim-spot-checker/internal/llm"
	"github.com/rs/zerolog/log"
)

const query = `I'm looking for a spot for a swimming lesson, the school is having individual single lessons and also X weeks session.
I captured the website's content a while ago and now it seems to be different. Can you compare the two version, and answer those two questions:
- Is there any spot available single or weeks session since last capture?
- Can you summarize what is newly available in a short sentence?
In order to be able to parse your answer, can you use a JSON format like this:
{
  "available": <boolean response to the first question>,
  "explanation": <string response to the second question>
}
Please can you only output JSON, no extra characters around it, I need to be able to get the response as it and parse it.

Here is the old content:
------------------------
%s
------------------------

Here is the new content:
------------------------
%s
------------------------
`

func CheckAvailability(
	ctx context.Context,
	model *llm.Model,
	oldContent string,
	newContent string,
) (found bool, explanation string, err error) {

	var (
		rawResponse string
		response    map[string]any
		varFound    bool
		varOk       bool
	)
	rawResponse, err = model.Answer(ctx, fmt.Sprintf(query, oldContent, newContent))
	if err != nil {
		return
	}

	log.Debug().Msgf("Model response: <%s>", rawResponse)

	err = json.Unmarshal([]byte(rawResponse), &response)
	if err != nil {
		err = fmt.Errorf("unable to parse model response: %s, because of %w", rawResponse, err)
		return
	}

	found, varFound, varOk = extractVar[bool](response, "available")
	if !varFound || !varOk {
		err = fmt.Errorf("unable to find 'available' field in model response (found %t, ok %t)", varFound, varOk)
		return
	}
	explanation, varFound, varOk = extractVar[string](response, "explanation")
	if !varFound || !varOk {
		err = fmt.Errorf("unable to find 'explanation' field in model response (found %t, ok %t)", varFound, varOk)
		return
	}

	return
}

func extractVar[T any](m map[string]any, field string) (value T, found bool, ok bool) {
	var anyValue any
	anyValue, found = m[field]
	if !found {
		return
	}

	value, ok = anyValue.(T)
	return
}
