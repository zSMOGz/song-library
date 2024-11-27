INSERT INTO songs (title, artist, album, release_date, text, link, genre, duration)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id; 