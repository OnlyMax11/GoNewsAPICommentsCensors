package api

import (
	"GoNewsAPICommentsCensors/pkg/storage"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	router *mux.Router       // Маршрутизатор запросов
	db     storage.Interface // база данных
}

// New Конструктор API.
func New(db storage.Interface) *API {
	api := API{
		router: mux.NewRouter(),
		db:     db,
	}
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.router.HandleFunc("/comments/check", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/comments/stop", api.addListHandler).Methods(http.MethodPost, http.MethodOptions)
}

// Получает из базы данных все слова из стоп листа.
func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	text := struct {
		Content string
	}{}
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stoplist, err := api.db.AllList()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, stopWord := range stoplist {
		matched, err := regexp.MatchString(stopWord.StopList, text.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if matched {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// Добавление слов в стоп лист базы данных.
func (api *API) addListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c storage.Stop
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	err = api.db.AddList(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.ResponseWriter.WriteHeader(w, http.StatusCreated)
}
