INSERT INTO verses (song_id, verse_number, content)
VALUES ($1, $2, $3)
RETURNING id; 