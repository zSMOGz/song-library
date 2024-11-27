SELECT id, song_id, verse_number, content, created_at
FROM verses
WHERE song_id = $1
ORDER BY verse_number
LIMIT $2 OFFSET $3; 