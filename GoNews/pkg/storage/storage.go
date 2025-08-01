package storage

// Post Публикация, получаемая из RSS.
type Post struct {
	ID          int    // Идентификатор записи.
	Title       string // Заголовок новости.
	Content     string // Содержание новости.
	PublishedAt int64  // Время публикации новости.
	Link        string // Ссылка на источник новости.
}

type Pagination struct {
	NumOfPages int `json:"numOfPages,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
}

// Interface Задаёт контракт на работу с БД.
type Interface interface {
	Posts(limit, offset int) ([]Post, error)                                       // Получение последних новостей.
	AddPost(p Post) error                                                          // Добавление новости в БД.
	PostSearchILIKE(keyWord string, limit, offset int) ([]Post, Pagination, error) // Поиск по заголовку
	PostsCreation([]Post) error                                                    // Создание n-ого кол-ва публикаций
	PostDetail(id int) (Post, error)                                               // Детальный вывод
}
