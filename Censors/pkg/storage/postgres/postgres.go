package postgres

import (
	"GoNewsAPICommentsCensors/pkg/storage"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store Хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор объекта хранилища
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
	s := Store{
		db: db,
	}

	return &s, nil
}

// AllList Выводит все комментарии.
func (p *Store) AllList() ([]storage.Stop, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM stop")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []storage.Stop
	for rows.Next() {
		var c storage.Stop
		err = rows.Scan(&c.ID, &c.StopList)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, rows.Err()
}

// AddList Добавляет комментарии в стоп лист.
func (p Store) AddList(c storage.Stop) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO stop (stop_list) VALUES ($1);", c.StopList)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
