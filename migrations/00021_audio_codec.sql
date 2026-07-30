-- +goose Up
-- +goose StatementBegin

alter table "torrent_contents" add column "audio_codec" text;
create index on "torrent_contents" (audio_codec);

alter table "torrent_hints" add column "audio_codec" text;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "torrent_hints" drop column "audio_codec";

drop index if exists "torrent_contents_audio_codec_idx";
alter table "torrent_contents" drop column "audio_codec";

-- +goose StatementEnd
