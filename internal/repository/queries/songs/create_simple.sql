INSERT INTO songs (title, artist) 
VALUES ($1, $2)
RETURNING id; 