-- Таблица для хранения песен
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    genre VARCHAR(100),
    duration INTEGER NOT NULL,
    release_date DATE,  
    text TEXT,
    link VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Добавим триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_songs_updated_at
    BEFORE UPDATE ON songs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
    
-- Создаем индексы для оптимизации наиболее частых запросов:

-- Индекс для быстрого поиска песен по названию
-- Оптимизирует запросы вида: WHERE title LIKE '%search%' или WHERE title = 'song_name'
CREATE INDEX idx_songs_title ON songs(title);

-- Индекс для поиска всех песен определенного исполнителя
-- Оптимизирует запросы вида: WHERE artist = 'artist_name'
CREATE INDEX idx_songs_artist ON songs(artist);

-- Индекс для фильтрации песен по году выпуска
-- Оптимизирует запросы вида: WHERE release_date BETWEEN 01.01.1999 AND 01.12.2000
CREATE INDEX idx_songs_release_date ON songs(release_date); 