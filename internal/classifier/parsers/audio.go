package parsers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

// ParseAudioContent infers audio-related attributes (currently AudioCodec)
// from the torrent's files. Unlike ParseVideoContent (which parses the name),
// this inspects file extensions — because music releases are multi-file albums
// whose format (FLAC/MP3/...) lives in the track extensions, not the folder name.
//
// The caller (classifier workflow) must ensure this only runs for torrents
// already classified as contentType.music — otherwise non-music torrents with
// incidental audio files (e.g. a movie with an MP3 soundtrack sample) would be
// misattributed.
func ParseAudioContent(torrent model.Torrent, result classification.Result) (classification.ContentAttributes, error) {
	attrs := classification.ContentAttributes{}
	attrs.InferAudioAttributes(torrent.Files)
	return attrs, nil
}
