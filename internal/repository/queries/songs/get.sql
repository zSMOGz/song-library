SELECT id, title, artist, album, genre, duration, release_date, text, link
FROM songs
WHERE id = $1; 