package api

import (
	"GoNewsAPICommentsCensors/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API Программный интерфейс сервера.
type API struct {
	router *mux.Router       // Маршрутизатор запросов
	db     storage.Interface // база данных
}

// New Конструктор объекта API.
func New(db storage.Interface) *API {
	api := API{
		router: mux.NewRouter(),
		db:     db,
	}
	api.router = mux.NewRouter()
	api.endpoints()

	return &api
}

// Router Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {

	return api.router
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// обработчиков новостей
	api.router.HandleFunc("/comments", api.commentsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/comments/del", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)
}

// commentsHandler, который выводит заданное кол-во новостей.
// Требуемое количество публикаций указывается в пути запроса
func (api *API) commentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	parseId := r.URL.Query().Get("news_id")
	newsId, err := strconv.Atoi(parseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comments, err := api.db.AllComments(newsId)
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// addCommentHandler Добавление комментария.
func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c storage.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	err = api.db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.ResponseWriter.WriteHeader(w, http.StatusCreated)
}

// deletePostHandler Удаление комментария.
func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c storage.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
