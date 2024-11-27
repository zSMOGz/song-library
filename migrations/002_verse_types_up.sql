CREATE TABLE verse_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Предзаполним базовыми типами
INSERT INTO verse_types (name) VALUES
    ('verse'),      -- куплет
    ('chorus'),     -- припев
    ('bridge'),     -- бридж
    ('intro'),      -- вступление
    ('outro'),      -- заключение
    ('pre_chorus'); -- предприпев 