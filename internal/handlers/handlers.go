package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/Ra1nz0r/effective_mobile-1/db/sqlc"
	cfg "github.com/Ra1nz0r/effective_mobile-1/internal/config"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
	"github.com/Ra1nz0r/effective_mobile-1/internal/models"
	"github.com/Ra1nz0r/effective_mobile-1/internal/services"
)

type HandleQueries struct {
	*db.Queries
	cfg.Config
}

func NewHandlerQueries(connect *sql.DB, cfg cfg.Config) *HandleQueries {
	return &HandleQueries{
		db.New(connect),
		cfg,
	}
}

// AddSongInLibrary добавляет песню в библиотеку. Обрабатывает POST запрос в формате
// JSON {"group": "Muse", "song": "Supermassive Black Hole"}, полученные данные добавляются
// в базу данных. Далее делается GET запрос во внешнее API для получения дополнительной
// информации о добавленной песне. Если данные не найдены или сервер недоступен, то дополнительные
// поля песни не заполняются и работа завершается. В случае успеха, делается запрос в базу данных
// для добавления дополнительных сведений о песне.
func (hq *HandleQueries) AddSongInLibrary(w http.ResponseWriter, r *http.Request) {
	// Получаем group и song из запроса, и помещаем данные в структуру.
	var baseParam db.AddParams
	if err := json.NewDecoder(r.Body).Decode(&baseParam); err != nil {
		logger.Zap.Error(err)
		ErrReturn(fmt.Errorf("invalid request"), http.StatusBadRequest, w)
		return
	}

	// Добавляем group и song в базу данных, без дополнительной информации.
	insert, err := hq.Add(r.Context(), baseParam)
	if err != nil {
		logger.Zap.Error(fmt.Errorf("error adding song: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Делаем запрос во внешний API для получения дополнительной информации о песне.
	// Если запрос завершился неудачей, то песня добавляется без дополнительных данных.
	details, errDet := services.FetchSongDetails(baseParam.Group, baseParam.Song, hq.ExternalApiURL)
	if errDet != nil {
		logger.Zap.Error(errDet)

		res := fmt.Sprintf(
			`
			Song ID: %d
			Unable to get additional information about the song.
            There is no data or the server is unavailable.
			The song will be added to the database without additional information.
			`,
			insert.ID,
		)

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		w.WriteHeader(http.StatusCreated)

		if _, err = w.Write([]byte(res)); err != nil {
			logger.Zap.Error(fmt.Errorf("failed attempt WRITE response: %w", err))
			return
		}
		return
	}

	// Добавляем в песню дополнительные параметры, полученные из внешнего API.
	fetch := db.FetchParams{
		ID:   insert.ID,
		Text: details.Text,
		Link: details.Link,
	}

	// Приводим дату к нужному формату и обновляем в FetchParams.
	fetch.ReleaseDate, err = time.Parse("02.01.2006", details.ReleaseDate)
	if err != nil {
		logger.Zap.Error(fmt.Errorf("error parsing date: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Делаем update песни в базе данных, заполняя поля releaseDate, text, link
	if err = hq.Fetch(r.Context(), fetch); err != nil {
		logger.Zap.Error(fmt.Errorf("error updating song: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := map[string]int32{
		"id": insert.ID,
	}

	resJSON, errJSON := json.Marshal(result)
	if errJSON != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(resJSON); err != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt WRITE response: %w", err))
		return
	}
}

// DeleteSong обрабатывает DELETE запрос и удаляет песню из библиотеки по указанному ID: "?id=21".
func (hq *HandleQueries) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		logger.Zap.Error("invalid string to number conversion or ID number")
		ErrReturn(fmt.Errorf("invalid string to number conversion or ID number"), http.StatusBadRequest, w)
		return
	}

	// Проверям существование песни и возвращаем ошибку, если её нет в базе данных.
	if _, err = hq.GetOne(r.Context(), int32(id)); err != nil {
		logger.Zap.Error("ID does not exist")
		ErrReturn(fmt.Errorf("ID does not exist"), http.StatusBadRequest, w)
		return
	}

	// Удаляем задачу из базы данных.
	if err = hq.Delete(r.Context(), int32(id)); err != nil {
		logger.Zap.Error("Delete request failed.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write([]byte(`{}`)); err != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt WRITE response: %w", err))
		return
	}
}

// ListAllSongsWithFilters обрабатывает GET запрос, получает данные из базы данных и
// выводит весь список песен из библиотеки в соответствии с фильтрами.
// Формат запроса "?group=Pink Floyd&releaseDate=11.11.2022&limit5&offset=0".
func (hq *HandleQueries) ListSongsWithFilters(w http.ResponseWriter, r *http.Request) {
	// Чтение параметров запроса из URL.
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releaseDate")
	text := r.URL.Query().Get("text")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = hq.PaginationLimit
	}

	offset, errOffset := strconv.Atoi(r.URL.Query().Get("offset"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}

	// Если полученные параметры не пусты, то записываем их в структуру запроса к базе данных.
	params := db.ListWithFiltersParams{
		Column1: sql.NullString{String: group, Valid: group != ""},
		Column2: sql.NullString{String: song, Valid: song != ""},
		Column4: sql.NullString{String: text, Valid: text != ""},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	if releaseDate != "" {
		params.ReleaseDate, err = time.Parse("02.01.2006", releaseDate)
		if err != nil {
			logger.Zap.Error(fmt.Errorf("error parsing date: %w", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Делаем запрос в базу данных с учётом указанных параметров фильтра.
	res, errUpdate := hq.ListWithFilters(r.Context(), params)
	if errUpdate != nil {
		logger.Zap.Error("Request could not be processed based on the specified filters.")
		w.WriteHeader(http.StatusInternalServerError)
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

// TextSongWithPagination обрабатывает GET запрос и выводит текст песни по указанному ID,
// разбитый на куплеты по страницам. Текст разделяется на куплеты по символу "\n\n".
// Формат запроса: "?id=16&page=1".
func (hq *HandleQueries) TextSongWithPagination(w http.ResponseWriter, r *http.Request) {
	songID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || songID < 1 {
		logger.Zap.Error("invalid string to number conversion or ID number")
		ErrReturn(fmt.Errorf("invalid string to number conversion or ID number"), http.StatusBadRequest, w)
		return
	}

	page, errPage := strconv.Atoi(r.URL.Query().Get("page"))
	if errPage != nil {
		logger.Zap.Error("invalid string to number conversion or PAGE number")
		ErrReturn(fmt.Errorf("invalid string to number conversion or PAGE number"), http.StatusBadRequest, w)
		return
	}

	// Получаем данные песни из базы данных.
	song, errSG := hq.GetText(r.Context(), int32(songID))
	if errSG != nil {
		logger.Zap.Error("Unable to retrieve song data.")
		ErrReturn(fmt.Errorf("invalid ID number"), http.StatusInternalServerError, w)
		return
	}

	// Разбиваем текст на куплеты по символу '\n\n'.
	couplet := strings.Split(song.Text, "\n\n")

	// Проверяем, не выходит ли запрашиваемая страница за пределы.
	if page > len(couplet) || page < 1 {
		logger.Zap.Error("Page out of range")
		ErrReturn(fmt.Errorf("page out of range"), http.StatusBadRequest, w)
		return
	}

	// Заменяем все символы '\n' в куплетах для корректного вывода.
	verse := strings.ReplaceAll(couplet[page-1], "\n", "\n")

	// Конфигурируем выходной результат.
	result := fmt.Sprintf("Group: %s, Song: %s\n\n%s", song.Group, song.Song, verse)

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write([]byte(result)); err != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}

// UpdateSong обрабатывает PUT запрос в формате JSON и обновляет параметры песни в базе данных.
// Формат запроса: {"id": 3, "releaseDate": "11.04.2022", "text": "You set my soul alight", "link": "ops link"}.
func (hq *HandleQueries) UpdateSong(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var sd models.SongDetail
	if err := json.NewDecoder(r.Body).Decode(&sd); err != nil {
		logger.Zap.Error(err)
		ErrReturn(fmt.Errorf("invalid request"), http.StatusBadRequest, w)
		return
	}

	// Заполняем данные для запроса в базу данных в соответствии с полученными данными.
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
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errParse), http.StatusBadRequest, w)
		return
	}

	// Делаем запрос и обновляем releaseDate, text, link песни в базе данных, в соответствии с полученными.
	if errUpdate := hq.Update(r.Context(), upd); errUpdate != nil {
		ErrReturn(fmt.Errorf("can't update task scheduler: %w", errUpdate), http.StatusBadRequest, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}

// WithRequestDetails (middleware) добавляет дополнительный код для регистрации сведений о запросе.
func (hs *HandleQueries) WithRequestDetails(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		logger.Zap.Info(
			"Method:", r.Method,
			"Duration:", time.Since(start),
			"URI:", r.RequestURI,
		)
	})
}

// WithResponseDetails (middleware) добавляет дополнительный код для регистрации сведений об ответе.
func (hs *HandleQueries) WithResponseDetails(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := logginResponseWriter{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		h.ServeHTTP(&lw, r)

		logger.Zap.Info(
			"Status:", lw.status,
			"Size:", lw.size,
		)
	})
}

// Переопределение методов для выведения дополнительной информации о запросах и ответах.
type logginResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *logginResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *logginResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}
