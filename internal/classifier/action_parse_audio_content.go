package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
)

const parseAudioContentName = "parse_audio_content"

type parseAudioContentAction struct{}

func (parseAudioContentAction) name() string {
	return parseAudioContentName
}

var parseAudioContentPayloadSpec = payloadLiteral[string]{
	literal:     parseAudioContentName,
	description: "Parse audio-related attributes (codec) from the files of the current torrent",
}

func (parseAudioContentAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := parseAudioContentPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed, err := parsers.ParseAudioContent(ctx.torrent, ctx.result)
			cl := ctx.result
			if err != nil {
				return cl, err
			}
			cl.Merge(parsed)
			return cl, nil
		},
	}, nil
}

func (parseAudioContentAction) JSONSchema() JSONSchema {
	return parseAudioContentPayloadSpec.JSONSchema()
}
