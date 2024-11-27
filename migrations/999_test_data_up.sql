-- Добавляем тестовые песни
INSERT INTO songs (title, artist, album, genre, duration, release_date, text) VALUES
    ('Группа крови', 'Кино', 'Группа крови', 'Рок', 285, '1988-01-01', NULL),
    ('Все идет по плану', 'Гражданская Оборона', 'Все идет по плану', 'Панк-рок', 260, '1988-01-01', NULL),
    ('Районы-кварталы', 'Звери', 'Районы-кварталы', 'Рок', 210, '2004-01-01', NULL);

-- Добавляем куплеты для "Группа крови"
INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    1,
    (SELECT id FROM verse_types WHERE name = 'verse'),
    'Теплое место, но улицы ждут
Отпечатков наших ног.
Звездная пыль - на сапогах.
Мягкое кресло, клетчатый плед,
Не нажатый вовремя курок.'
FROM songs s WHERE s.title = 'Группа крови';

INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    2,
    (SELECT id FROM verse_types WHERE name = 'chorus'),
    'Пожелай мне удачи в бою,
Пожелай мне:
Не остаться в этой траве,
Не остаться в этой траве.
Пожелай мне удачи,
Пожелай мне удачи!'
FROM songs s WHERE s.title = 'Группа крови';

-- Добавляем куплеты для "Все идет по плану"
INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    1,
    (SELECT id FROM verse_types WHERE name = 'verse'),
    'А при коммунизме все будет заебись,
Он наступит скоро, надо только подождать.
Там все будет бесплатно, там все будет в кайф,
Там наверное вообще не надо будет умирать!'
FROM songs s WHERE s.title = 'Все идет по плану';

INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    2,
    (SELECT id FROM verse_types WHERE name = 'chorus'),
    'Все идет по плану,
Все идет по плану!'
FROM songs s WHERE s.title = 'Все идет по плану';

-- Добавляем куплеты для "Районы-кварталы"
INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    1,
    (SELECT id FROM verse_types WHERE name = 'verse'),
    'Районы, кварталы, жилые массивы,
Я ухожу, ухожу красиво.
Районы, кварталы, жилые массивы,
Я ухожу, ухожу красиво.'
FROM songs s WHERE s.title = 'Районы-кварталы';

INSERT INTO verses (song_id, verse_number, verse_type_id, content)
SELECT 
    s.id,
    2,
    (SELECT id FROM verse_types WHERE name = 'chorus'),
    'И никто не знает, что будет дальше,
Все мои надежды - не напрасны.
Районы, кварталы, жилые массивы,
Я ухожу, ухожу красиво.'
FROM songs s WHERE s.title = 'Районы-кварталы';
  