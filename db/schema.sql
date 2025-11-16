CREATE TABLE schema_migrations (version uint64,dirty bool);
CREATE UNIQUE INDEX version_unique ON schema_migrations (version);
CREATE TABLE backlite_tasks (
    id text PRIMARY KEY,
    created_at integer NOT NULL,
    queue text NOT NULL,
    task blob NOT NULL,
    wait_until integer,
    claimed_at integer,
    last_executed_at integer,
    attempts integer NOT NULL DEFAULT 0
) STRICT;
CREATE TABLE backlite_tasks_completed (
    id text PRIMARY KEY NOT NULL,
    created_at integer NOT NULL,
    queue text NOT NULL,
    last_executed_at integer,
    attempts integer NOT NULL,
    last_duration_micro integer,
    succeeded integer,
    task blob,
    expires_at integer,
    error text
) STRICT;
CREATE INDEX backlite_tasks_wait_until ON backlite_tasks (wait_until) WHERE wait_until IS NOT NULL;
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
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  (20251116073135);
