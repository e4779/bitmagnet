package model

// AudioCodec classifies a music torrent's encoding: lossless, lossy, or other.
// Unlike video attributes (inferred from the torrent name via regex), AudioCodec
// is inferred from the dominant audio file extension by size — because music is
// almost always multi-file (album folders with tracks + artwork + cue/log), so
// the torrent name rarely carries the format.
//
// ENUM(Lossless, Lossy, Other)
type AudioCodec string

// Label returns the human-readable label for the audio codec.
func (a AudioCodec) Label() string {
	return a.String()
}

// lossless audio extensions — preserved-bitrate formats (FLAC, APE, WAV, DSD).
var losslessAudioExtensions = map[string]struct{}{
	"flac": {},
	"ape":  {},
	"wav":  {},
	"dsf":  {},
	"dff":  {},
}

// lossy audio extensions — lossy-compressed formats (MP3, AAC, OGG, M4A/OPUS).
var lossyAudioExtensions = map[string]struct{}{
	"mp3":  {},
	"aac":  {},
	"ogg":  {},
	"m4a":  {},
	"opus": {},
	"wma":  {},
}

// InferAudioCodecFromFiles determines the dominant audio codec from the torrent's
// files. It sums the byte size of lossless vs lossy audio extensions (ignoring
// non-audio files like jpg/cue/log/txt/nfo) and returns the dominant one.
//
// Returns NullAudioCodec{Valid: false} when the torrent has no recognizable
// audio files, or when both lossless and lossy content is present without a
// clear majority (ambiguous mixes are classified as Other).
func InferAudioCodecFromFiles(files []TorrentFile) NullAudioCodec {
	var losslessBytes, lossyBytes uint
	for _, f := range files {
		if !f.Extension.Valid {
			continue
		}
		ext := f.Extension.String
		if _, ok := losslessAudioExtensions[ext]; ok {
			losslessBytes += f.Size
		} else if _, ok := lossyAudioExtensions[ext]; ok {
			lossyBytes += f.Size
		}
		// non-audio extensions (jpg, cue, log, txt, nfo, png, m3u, ...) are ignored
	}

	// Require a clear majority: the dominant type must exceed 2x the other.
	// This avoids misclassification of mixed-format torrents (e.g. a FLAC album
	// with a bonus MP3 preview) while still catching single-format albums.
	switch {
	case losslessBytes > 0 && losslessBytes >= 2*lossyBytes:
		return NewNullAudioCodec(AudioCodecLossless)
	case lossyBytes > 0 && lossyBytes >= 2*losslessBytes:
		return NewNullAudioCodec(AudioCodecLossy)
	case losslessBytes > 0 || lossyBytes > 0:
		// Both present without a clear majority, or single file — ambiguous.
		return NewNullAudioCodec(AudioCodecOther)
	default:
		// No recognizable audio files.
		return NullAudioCodec{}
	}
}
