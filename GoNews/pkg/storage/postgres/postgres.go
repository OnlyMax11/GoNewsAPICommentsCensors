package postgres

import (
	"GoNewsAPICommentsCensors/pkg/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор объекта хранилища.
func New(ctx context.Context, constr string) (*Store, error) {
	for {
		_, err := pgxpool.New(ctx, constr)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.New(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := Store{db: db}

	return &s, err
}

// PostsCreation Создание n-ого кол-ва публикаций
func (p *Store) PostsCreation(posts []storage.Post) error {
	for _, post := range posts {
		err := p.AddPost(post)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddPost создаёт новую запись.
func (s *Store) AddPost(t storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO news (title, content, publishedAt, link)
		VALUES ($1, $2, $3, $4);
		`,
		t.Title,
		t.Content,
		t.PublishedAt,
		t.Link,
	).Scan()

	return err
}

// PostSearchILIKE Поиск по заголовку
func (p *Store) PostSearchILIKE(pattern string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	pattern = "%" + pattern + "%"
	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	row := p.db.QueryRow(context.Background(), "SELECT count(*) FROM news WHERE title ILIKE $1;", pattern)
	err := row.Scan(&pagination.NumOfPages)
	if pagination.NumOfPages%limit > 0 {
		pagination.NumOfPages = pagination.NumOfPages/limit + 1
	} else {
		pagination.NumOfPages /= limit
	}
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	rows, err := p.db.Query(context.Background(), "SELECT * FROM news WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;", pattern, limit, offset)
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PublishedAt, &p.Link)
		if err != nil {
			return nil, storage.Pagination{}, err
		}
		posts = append(posts, p)
	}

	return posts, pagination, rows.Err()
}

// Posts Получение странице с определенным номером
func (s *Store) Posts(limit, offset int) ([]storage.Post, error) {
	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	rows, err := s.db.Query(context.Background(), `
	SELECT * FROM news
	ORDER BY pubtime DESC LIMIT $1 OFFSET $2
	`,
		pagination.Limit, pagination.Page,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []storage.Post
	// итерированние по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PublishedAt,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}

	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// PostDetail Получение публикаций по id
func (p *Store) PostDetail(id int) (storage.Post, error) {
	row := p.db.QueryRow(context.Background(), `
	SELECT * FROM news 
    WHERE id =$1;
	`, id)
	var post storage.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PublishedAt,
		&post.Link)
	if err != nil {
		return storage.Post{}, err
	}

	return post, nil
}
