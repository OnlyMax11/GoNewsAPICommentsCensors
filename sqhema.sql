DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS stop;
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT, -- Заголовок новости.
    content TEXT NOT NULL UNIQUE, --Содержание новости.
    publishedAt BIGINT DEFAULT 0, --Время публикации новости.
    link TEXT --Ссылка на опубликованную новость.
);
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    news_id INT,
    content TEXT NOT NULL DEFAULT 'empty',
    PubTime BIGINT NOT NULL DEFAULT extract (epoch from now())
);
CREATE TABLE IF NOT EXISTS stop (
    id SERIAL PRIMARY KEY,
    stop_list TEXT
);
INSERT INTO comments(news_id,content)  VALUES (1,'комментарий');
INSERT INTO comments(news_id,content)  VALUES (2,'ups  проверка');
INSERT INTO stop (stop_list) VALUES ('qwerty');
INSERT INTO stop (stop_list) VALUES ('йцукен');
INSERT INTO stop (stop_list) VALUES ('zxvbnm');