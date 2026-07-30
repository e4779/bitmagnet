package search

import (
	"database/sql/driver"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

var VideoResolutionCriteria = torrentContentAttributeCriteria[model.VideoResolution](videoResolutionField)

var Video3DCriteria = torrentContentAttributeCriteria[model.Video3D](video3dField)

// audioCodecField returns the DAO field for torrent_contents.audio_codec.
// Defined here (not in a facet file) because the GraphQL/facet layer for
// AudioCodec is not yet implemented — but the SQL filter is needed by the
// torznab adapter (cat=3010/3040 filtering).
func audioCodecField(q *dao.Query) field.Field {
	return q.TorrentContent.AudioCodec
}

var AudioCodecCriteria = torrentContentAttributeCriteria[model.AudioCodec](audioCodecField)

func torrentContentAttributeCriteria[T attribute](getFld func(*dao.Query) field.Field) func(...T) query.Criteria {
	return func(values ...T) query.Criteria {
		return query.DaoCriteria{
			Conditions: func(ctx query.DBContext) ([]field.Expr, error) {
				q := ctx.Query()
				fld := getFld(q)
				valuers := make([]driver.Valuer, 0, len(values))
				for _, v := range values {
					valuers = append(valuers, v)
				}
				return []field.Expr{fld.In(valuers...)}, nil
			},
			Joins: maps.NewInsertMap(
				maps.MapEntry[string, struct{}]{Key: model.TableNameTorrentContent},
			),
		}
	}
}
