-- Сначала удаляем зависимые триггеры
DROP TRIGGER IF EXISTS update_songs_updated_at ON songs;
DROP TRIGGER IF EXISTS update_verses_updated_at ON verses;

-- Затем удаляем таблицы в правильном порядке (сначала зависимые)
DROP TABLE IF EXISTS verses;
DROP TABLE IF EXISTS verse_types;
DROP TABLE IF EXISTS songs;

-- В конце удаляем общую функцию
DROP FUNCTION IF EXISTS update_updated_at_column();