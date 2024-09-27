package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/Ra1nz0r/effective_mobile-1/db/sqlc"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
	"github.com/Ra1nz0r/effective_mobile-1/internal/models"
	"github.com/Ra1nz0r/effective_mobile-1/internal/services"
)

type HandleQueries struct {
	*db.Queries
}

func NewHandleQueries(conn *sql.DB) *HandleQueries {
	return &HandleQueries{
		db.New(conn),
	}
}

func (hq *HandleQueries) AddSongInLibrary(w http.ResponseWriter, r *http.Request) {
	// Получаем group и song из запроса, и помещаем данные в структуру.
	var baseParam db.AddParams
	if err := json.NewDecoder(r.Body).Decode(&baseParam); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Добавляем group и song в базу данных, без дополнительной информации.
	insert, errCreate := hq.Add(r.Context(), baseParam)
	if errCreate != nil {
		//logerr.ErrEvent("cannot create task in DB", errCreate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Делаем запрос во внешний API для получения дополнительной информации о песне.
	details, errDetail := services.FetchSongDetails(baseParam.Group, baseParam.Song)
	if errDetail != nil {
		logger.Zap.Error(errDetail)
		// Ошибку для клиента исправить на то что сервер недоступен или невозможно получить данные, добавлено с базовыми
		http.Error(w, "Error fetching song details", http.StatusInternalServerError)
		return
	}

	// Добавляем в песню дополнительные параметры, полученные из внешнего API.
	fetch := db.FetchParams{
		ID:   insert.ID,
		Text: details.Text,
		Link: details.Link,
	}

	// Приводим дату к нужному формату и обновляем дату в FetchParams.
	var errParse error
	fetch.ReleaseDate, errParse = time.Parse("02.01.2006", details.ReleaseDate)
	if errParse != nil {
		logger.Zap.Error("Error parsing date: %w", errParse)
	}

	// Делаем update песни в базе данных, заполняя поля releaseDate, text, link
	if errFetch := hq.Fetch(r.Context(), fetch); errFetch != nil {
		http.Error(w, "Error updating song", http.StatusInternalServerError)
		return
	}

	// Создание мапы и выведение последнего ID добавленного в датабазу, ответ в виде: {"id":"186"}.
	respResult := make(map[string]int32)
	respResult["id"] = insert.ID
	jsonResp, errJSON := json.Marshal(respResult)
	if errJSON != nil {
		//logerr.ErrEvent("failed attempt json-marshal response", errJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusCreated)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		//logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}

func (hq *HandleQueries) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id, errID := strconv.Atoi(r.URL.Query().Get("id"))
	if errID != nil || id < 1 {
		logger.Zap.Error(fmt.Errorf("invalid string to number conversion or ID number: DeleteSong"))
		ErrReturn(fmt.Errorf("invalid string to number conversion or ID number"), 404, w)
		return
	}

	// Проверям существование задачи и возвращаем ошибку, если её нет в базе данных.
	_, errGeted := hq.GetOne(r.Context(), int32(id))
	if errGeted != nil {
		ErrReturn(fmt.Errorf("the ID you entered does not exist: %w", errGeted), 404, w)
		return
	}

	// Удаляем задачу из базы данных, при DELETE запросе в виде "/api/task?id=185".
	if errDel := hq.Delete(r.Context(), int32(id)); errDel != nil {
		//ErrReturn(fmt.Errorf("failed delete: %w", errDel), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		//logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}

// Собирает все метрики метрики из локального хранилища и выводит их в
// результирующей карте при получении GET запроса.
// Вызывает метод интерфейса, который возвращает копию локального хранилища.
// Формат JSON, в виде {"Alloc":146464,"Frees":10,...}.
func (hq *HandleQueries) ListAllSongs(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	releaseDateStr := r.URL.Query().Get("releaseDate")
	text := r.URL.Query().Get("text")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Преобразуем limit и offset в int
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // дефолтное значение
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // дефолтное значение
	}

	// Преобразуем releaseDate
	var releaseDate sql.NullTime
	if releaseDateStr != "" {
		date, err := time.Parse("2006-01-02", releaseDateStr)
		if err == nil {
			releaseDate = sql.NullTime{Time: date, Valid: true}
		}
	}

	// Создаем параметры для SQL-запроса
	params := db.ListLibraryParams{
		Group:       sql.NullString{String: group, Valid: group != ""},
		Song:        sql.NullString{String: song, Valid: song != ""},
		ReleaseDate: releaseDate,
		Text:        sql.NullString{String: text, Valid: text != ""},
		Limit:       int32(limit),
		Offset:      int32(offset),
	}

	// Выполняем запрос через SQLC
	libraries, err := hq.ListLibrary(context.Background(), params)
	if err != nil {
		http.Error(w, "Failed to query libraries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libraries)
}

func (hq *HandleQueries) UpdateSong(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var sd models.SongDetail
	if err := json.NewDecoder(r.Body).Decode(&sd); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Обновляем параметры песни в соответствии с полученными данными.
	upd := db.UpdateParams{
		ID:   sd.ID,
		Text: sd.Text,
		Link: sd.Link,
	}

	// Приводим дату к нужному формату и обновляем дату в UpdateParams.
	var errParse error
	upd.ReleaseDate, errParse = time.Parse("02.01.2006", sd.ReleaseDate)
	if errParse != nil {
		logger.Zap.Error("Error parsing date: %w", errParse)
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errParse), 404, w)
		return
	}

	// Делаем запрос и обновляем параметры песни в базе данных, в соответствии с полученными.
	if errUpdate := hq.Update(r.Context(), upd); errUpdate != nil {
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errUpdate), 404, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusAccepted)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		//logerr.ErrEvent("failed attempt WRITE response", errWrite)
		return
	}
}

// =========================================================================

func (hq *HandleQueries) GetAll(w http.ResponseWriter, r *http.Request) {
	res, errUpdate := hq.ListAll(r.Context())
	if errUpdate != nil {
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errUpdate), 404, w)
		return
	}

	ans, errJSON := json.Marshal(res)
	if errJSON != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", errJSON))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(ans)); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}
