package main

import (
	"GoNewsAPICommentsCensors/config"
	"GoNewsAPICommentsCensors/pkg/api"
	"GoNewsAPICommentsCensors/pkg/middl"
	"GoNewsAPICommentsCensors/pkg/rss"
	"GoNewsAPICommentsCensors/pkg/storage"
	"GoNewsAPICommentsCensors/pkg/storage/postgres"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"path/filepath"
	"time"
)


// Сервер.
type server struct {
	db  storage.Interface
	api *api.API
}

// init вызывается перед main()
func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	configURL := filepath.Join("aggregator", "cmd", "server", "config.json")
	// Создаём объект сервера.
	var srv server
	cfg := config.New()
	// Адрес базы данных
	dbURL := cfg.News.URLdb
	// Порт по умолчанию.
	port := cfg.News.AdrPort
	// Можно сменить Порт при запуске флагом < --news-port= >
	portFlag := flag.String("news-port", port, "Порт для news сервиса")
	flag.Parse()
	portNews := *portFlag
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// объект базы данных postgresql
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Println(err)
	}
	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(db)
	// Создаём каналы для новостей и ошибок.
	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)
	// Чтение RSS-лент из конфига с заданным интервалом
	go func() {
		err := rss.GoNews(configURL, chanPosts, chanErrs)
		if err != nil {
			log.Println(err)
		}
	}()
	// вывод ошибок
	go func() {
		for err := range chanErrs {
			log.Println(err)
		}
	}()
	srv.api.Router().Use(middlewareeware.Middlewaree)
	log.Print("Запуск сервера на http://127.0.0.1" + portNews)
	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(portNews, srv.api.Router())
	if err != nil {
		log.Println("Не удалось запустить сервер. Ошибка:", err)
	}
}