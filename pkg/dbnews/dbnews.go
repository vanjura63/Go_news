package dbnews

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Data base.
type DB struct {
	db *pgxpool.Pool
}

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

/*func (db *DB) Delete() error {

	_, err := db.db.Exec(context.Background(), `
		DELETE  from news

		`,
	)
	return err
}*/

// Constructor, accepts a DB connection string.
func New(constr string) (*DB, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := DB{
		db: db,
	}
	return &s, nil
}

func (db *DB) StoreNews(news []Post) error {
	for _, post := range news {
		_, err := db.db.Exec(context.Background(), `
		INSERT INTO news(title, content, public_time, link)
		VALUES ($1, $2, $3, $4)`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// News ruturns news from DB.
func (db *DB) News(n int) ([]Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.db.Query(context.Background(), `
	SELECT id, title, content, public_time, link FROM news
	ORDER BY public_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var news []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}
