package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const AudioCodecFacetKey = "audio_codec"

func audioCodecField(q *dao.Query) field.Field {
	return q.TorrentContent.AudioCodec
}

func AudioCodecFacet(options ...query.FacetOption) query.Facet {
	return audioCodecFacet{torrentContentAttributeFacet[model.AudioCodec]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(AudioCodecFacetKey),
				query.FacetHasLabel("Audio Codec"),
				query.FacetUsesOrLogic(),
				query.FacetTriggersCte(),
			}, options...)...,
		),
		field: audioCodecField,
		parse: model.ParseAudioCodec,
	}}
}

type audioCodecFacet struct {
	torrentContentAttributeFacet[model.AudioCodec]
}

func (audioCodecFacet) Values(query.FacetContext) (map[string]string, error) {
	acs := model.AudioCodecValues()
	values := make(map[string]string, len(acs))

	for _, ac := range acs {
		values[ac.String()] = ac.Label()
	}

	return values, nil
}
