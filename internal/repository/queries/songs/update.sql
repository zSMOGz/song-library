UPDATE songs 
SET title = $1, artist = $2, album = $3, release_date = $4, 
    text = $5, link = $6, genre = $7, duration = $8
WHERE id = $9
RETURNING id; 