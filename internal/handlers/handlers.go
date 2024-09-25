package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	db "github.com/Ra1nz0r/effective_mobile-1/db/sqlc"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
	"github.com/jackc/pgx/v5"
)

type HandleQueries struct {
	*db.Queries
}

func NewHandleQueries(conn *pgx.Conn) *HandleQueries {
	return &HandleQueries{db.New(conn)}
}

func (hq *HandleQueries) AddSongInLibrary(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		//logerr.ErrEvent("cannot read from BODY", errBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var task db.AddParams
	if errUnm := json.Unmarshal(result, &task); errUnm != nil {

		//ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	fmt.Println(task)

	// Если данные корректны, то создаём запись в датабазе.
	insertedTask, errCreate := hq.Add(r.Context(), task)
	if errCreate != nil {
		fmt.Println(errCreate)
		//logerr.ErrEvent("cannot create task in DB", errCreate)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Создание мапы и выведение последнего ID добавленного в датабазу, ответ в виде: {"id":"186"}.
	respResult := make(map[string]int32)
	respResult["id"] = insertedTask.ID
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
	if errID != nil {
		log.Fatal(errID)
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

// Собирает все метрики метрики из локального хранилища и выводит их в
// результирующей карте при получении GET запроса.
// Вызывает метод интерфейса, который возвращает копию локального хранилища.
// Формат JSON, в виде {"Alloc":146464,"Frees":10,...}.
func (hq *HandleQueries) ListAllSongs(w http.ResponseWriter, r *http.Request) {

	limit, errID := strconv.Atoi(r.URL.Query().Get("limit"))
	if errID != nil {
		log.Fatal(errID)
	}

	offset, errID := strconv.Atoi(r.URL.Query().Get("offset"))
	if errID != nil {
		log.Fatal(errID)
	}

	var rr db.ListParams
	rr.Limit = int32(limit)
	rr.Offset = int32(offset)

	res, errUpdate := hq.List(r.Context(), rr)
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

func (hq *HandleQueries) UpdateSong(w http.ResponseWriter, r *http.Request) {
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		//logerr.ErrEvent("cannot read from BODY", errBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var task db.UpdateParams
	if errUnm := json.Unmarshal(result, &task); errUnm != nil {

		//ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	if errUpdate := hq.Update(r.Context(), task); errUpdate != nil {
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
