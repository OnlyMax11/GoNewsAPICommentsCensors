package main

import (
	"GoNewsAPICommentsCensors/config"
	"GoNewsAPICommentsCensors/pkg/api"
	"GoNewsAPICommentsCensors/pkg/middl"
	"GoNewsAPICommentsCensors/pkg/storage"
	"GoNewsAPICommentsCensors/pkg/storage/postgres"
	"GoNewsAPICommentsCensors/supply"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

// сервер
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
	// объект сервера
	var srv server
	cfg := config.New()
	// Адрес базы данных
	dbURL := cfg.Comments.URLdb
	// Порт по умолчанию.
	port := cfg.Censor.AdrPort
	// Можно сменить Порт при запуске флагом < --censor-port= >
	portFlag := flag.String("censor-port", port, "Порт для censor сервиса")
	flag.Parse()
	portCensor := *portFlag
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// объект базы данных postgresql
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Println(err)
	}
	// Получение списка для стоп листа из файла words.txt
	stop, err := supply.StopList()
	if err != nil {
		log.Println(err)
	}
	// Добавление в таблицу stop полученного списка
	for _, v := range stop {
		err := db.AddList(v)
		if err != nil {
			log.Println(err)
		}
	}
	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db
	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)
	srv.api.Router().Use(middlewareeware.Middlewaree)
	log.Print("Запуск сервера на http://127.0.0.1" + portCensor)
	err = http.ListenAndServe(portCensor, srv.api.Router())
	if err != nil {
		log.Println("Не удалось запустить сервер шлюза. Error:", err)
	}
}
