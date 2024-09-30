package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	tDate := Library{
		Group:       "Nirvana",
		Song:        "Lithium",
		ReleaseDate: time.Date(1999, 5, 11, 0, 0, 0, 0, time.UTC),
		Text:        "Angie, Angie, when will those dark clouds disappear",
		Link:        "http://www.ops.com",
	}

	ctx := context.Background()

	// Создаём песню в базе данных в соответствии с AddParam и проверяем результат.
	addParam := AddParams{
		Group: tDate.Group,
		Song:  tDate.Song,
	}

	add, err := testQueries.Add(ctx, addParam)
	require.NoError(t, err)
	require.NotEmpty(t, add)

	assert.Equal(t, addParam.Group, add.Group)
	assert.Equal(t, addParam.Song, add.Song)

	require.NotZero(t, add.ID)
	require.NotZero(t, add.ReleaseDate)

	// Обновляем параметры песни в базе данных в соответствии с FetchParam.
	fetchParam := FetchParams{
		ID:          add.ID,
		ReleaseDate: tDate.ReleaseDate,
		Text:        tDate.Text,
		Link:        tDate.Link,
	}

	err = testQueries.Fetch(ctx, fetchParam)
	require.NoError(t, err)

	// Получаем песню и проверяем обновлённые параметры.
	res, err := testQueries.GetOne(ctx, add.ID)
	require.NoError(t, err)

	assert.Equal(t, fetchParam.ReleaseDate, res.ReleaseDate)
	assert.Equal(t, fetchParam.Text, res.Text)
	assert.Equal(t, fetchParam.Link, res.Link)

	// Получаем текст песни и проверяем результат.
	txt, errTxt := testQueries.GetText(ctx, res.ID)
	require.NoError(t, errTxt)

	assert.Equal(t, res.Text, txt.Text)

	// Настраиваем фильтры и делаем запрос в базу данных для вывода песен с их учётом.
	listFilter := ListWithFiltersParams{
		Column1: sql.NullString{String: tDate.Group, Valid: true},
		Column4: sql.NullString{String: "will", Valid: true},
		Limit:   10,
		Offset:  0,
	}

	list, errList := testQueries.ListWithFilters(ctx, listFilter)
	require.NoError(t, errList)

	require.NotNil(t, list)

	// Обновляем параметры песни.
	updParam := UpdateParams{
		ID:          res.ID,
		ReleaseDate: time.Date(2004, 4, 8, 0, 0, 0, 0, time.UTC),
	}

	err = testQueries.Update(ctx, updParam)
	require.NoError(t, err)

	// Получаем песню и проверяем обновлённые параметры.
	res, err = testQueries.GetOne(ctx, add.ID)
	require.NoError(t, err)

	assert.Equal(t, updParam.ReleaseDate, res.ReleaseDate)
	assert.Equal(t, updParam.Text, res.Text)
	assert.Equal(t, updParam.Link, res.Link)

	// Удаляем песню из базы данных.
	err = testQueries.Delete(ctx, res.ID)
	require.NoError(t, err)

	// Получаем песню и проверяем её удаление.
	res, err = testQueries.GetOne(ctx, add.ID)
	require.Error(t, err)

	assert.Equal(t, sql.ErrNoRows, err)
}
