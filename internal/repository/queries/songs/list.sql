SELECT 
    s.*, 
    COUNT(*) OVER() as total_count
FROM songs s
WHERE 
    ($1 = '' OR LOWER(title) LIKE LOWER('%' || $1 || '%')) AND
    ($2 = '' OR LOWER(artist) LIKE LOWER('%' || $2 || '%')) AND
    ($3 = '' OR LOWER(album) LIKE LOWER('%' || $3 || '%')) AND
    ($4 = 0 OR EXTRACT(YEAR FROM release_date) = $4) AND
    ($5 = '' OR LOWER(genre) LIKE LOWER('%' || $5 || '%'))
ORDER BY title
LIMIT $6 OFFSET $7; 