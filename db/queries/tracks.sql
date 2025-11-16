-- name: GetAllTracks :many
SELECT * FROM tracks;

-- name: GetTrackByPath :one
SELECT * FROM tracks WHERE path = ? LIMIT 1;

-- name: GetTrackByID :one
SELECT * FROM tracks WHERE id = ? LIMIT 1;

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

-- name: GetTracksByArtist :many
SELECT * FROM tracks WHERE artist = ? ORDER BY album, title;

-- name: GetTracksByAlbum :many
SELECT * FROM tracks WHERE album = ? ORDER BY title;

-- name: SearchTracks :many
SELECT * FROM tracks
WHERE title LIKE ? OR artist LIKE ? OR album LIKE ?
ORDER BY artist, album, title;
