-- Таблица для хранения куплетов песен
CREATE TABLE IF NOT EXISTS verses (
    id SERIAL PRIMARY KEY,
    song_id INTEGER NOT NULL,
    verse_number INTEGER NOT NULL,
    verse_type_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Ограничение внешнего ключа для связи с таблицей songs
    CONSTRAINT fk_verses_song
        FOREIGN KEY (song_id)
        REFERENCES songs(id)
        ON DELETE CASCADE,
    -- Ограничение внешнего ключа для связи с таблицей verse_types
    CONSTRAINT fk_verses_type
        FOREIGN KEY (verse_type_id)
        REFERENCES verse_types(id),
    -- Обеспечиваем уникальность комбинации song_id и verse_number
    CONSTRAINT unique_verse_number_per_song 
        UNIQUE(song_id, verse_number)
);

-- Триггер для updated_at
CREATE TRIGGER update_verses_updated_at
    BEFORE UPDATE ON verses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- Индекс для быстрого поиска всех куплетов конкретной песни
CREATE INDEX idx_verses_song_id ON verses(song_id); 