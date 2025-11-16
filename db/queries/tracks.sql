-- name: GetAllTracks :many
SELECT
    id,
    path,
    CAST(title AS TEXT),
    CAST(artist AS TEXT),
    CAST(album AS TEXT),
    duration,
    has_plain_lyrics,
    has_synced_lyrics
FROM tracks;

-- name: GetTrackByPath :one
SELECT
    id,
    path,
    CAST(title AS TEXT),
    CAST(artist AS TEXT),
    CAST(album AS TEXT),
    duration,
    has_plain_lyrics,
    has_synced_lyrics
FROM tracks
WHERE path = ?
LIMIT 1;

-- name: GetTrackByID :one
SELECT
    id,
    path,
    CAST(title AS TEXT),
    CAST(artist AS TEXT),
    CAST(album AS TEXT),
    duration,
    has_plain_lyrics,
    has_synced_lyrics
FROM tracks
WHERE id = ?
LIMIT 1;

-- name: CreateTrack :exec
INSERT INTO tracks (
    path,
    title,
    artist,
    album,
    duration,
    has_plain_lyrics,
    has_synced_lyrics
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateTrack :exec
UPDATE tracks
SET
    title = ?,
    artist = ?,
    album = ?,
    duration = ?,
    has_plain_lyrics = ?,
    has_synced_lyrics = ?
WHERE path = ?;

-- name: DeleteTrack :exec
DELETE FROM tracks WHERE path = ?;

-- name: SearchTracks :many
SELECT
    id,
    path,
    CAST(title AS TEXT),
    CAST(artist AS TEXT),
    CAST(album AS TEXT),
    duration,
    has_plain_lyrics,
    has_synced_lyrics
FROM tracks
WHERE title LIKE ? OR artist LIKE ? OR album LIKE ?
ORDER BY artist, album, title;

-- name: GetStats :one
SELECT
    COUNT(*) AS total_tracks,
    CAST(SUM(CASE WHEN has_plain_lyrics = 1 AND has_synced_lyrics = 1 THEN 1 ELSE 0 END) AS INTEGER) AS tracks_with_both_lyrics,
    CAST(SUM(CASE WHEN has_plain_lyrics = 0 AND has_synced_lyrics = 1 THEN 1 ELSE 0 END) AS INTEGER) AS tracks_with_synced_only,
    CAST(SUM(CASE WHEN has_plain_lyrics = 1 AND has_synced_lyrics = 0 THEN 1 ELSE 0 END) AS INTEGER) AS tracks_with_plain_only,
    CAST(SUM(CASE WHEN has_plain_lyrics = 0 AND has_synced_lyrics = 0 THEN 1 ELSE 0 END) AS INTEGER) AS instrumental_tracks,
    CAST(SUM(
        CASE WHEN title IS NULL OR title = ''
              OR artist IS NULL OR artist = ''
              OR album IS NULL OR album = ''
        THEN 1 ELSE 0 END
    ) AS INTEGER) AS tracks_missing_metadata
FROM tracks;
