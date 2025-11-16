-- migrate:up
CREATE TABLE tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    title TEXT,
    artist TEXT,
    album TEXT,
    duration REAL NOT NULL,
    has_plain_lyrics BOOLEAN NOT NULL DEFAULT 0,
    has_synced_lyrics BOOLEAN NOT NULL DEFAULT 0
);

CREATE INDEX idx_tracks_title ON tracks(title);
CREATE INDEX idx_tracks_artist ON tracks(artist);
CREATE INDEX idx_tracks_album ON tracks(album);
CREATE INDEX idx_tracks_path ON tracks(path);

-- migrate:down
DROP INDEX IF EXISTS idx_tracks_path;
DROP INDEX IF EXISTS idx_tracks_album;
DROP INDEX IF EXISTS idx_tracks_artist;
DROP INDEX IF EXISTS idx_tracks_title;
DROP TABLE IF EXISTS tracks;
