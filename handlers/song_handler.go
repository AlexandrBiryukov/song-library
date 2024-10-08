// handlers/song_handler.go
package handlers

import (
	"net/http"
	"song-library/config"
	"song-library/database"
	"song-library/externalapi"
	"song-library/logger"
	"song-library/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type AddSongRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

// Добавление новой песни
func AddSong(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req AddSongRequest
		if err := c.Bind(&req); err != nil {
			logger.Log.Error(err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
		}

		songDetail, err := externalapi.FetchSongDetail(cfg, req.Group, req.Song)
		if err != nil {
			logger.Log.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch song details"})
		}

		song := models.Song{
			GroupName:   req.Group,
			Song:        req.Song,
			ReleaseDate: songDetail.ReleaseDate,
			Text:        songDetail.Text,
			Link:        songDetail.Link,
		}

		if err := database.DB.Create(&song).Error; err != nil {
			logger.Log.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to save song"})
		}

		logger.Log.Infof("Song added: %s - %s", req.Group, req.Song)
		return c.JSON(http.StatusCreated, song)
	}
}

// Получение списка песен с фильтрацией и пагинацией
func GetSongs(c echo.Context) error {
	var songs []models.Song
	query := database.DB

	// Фильтрация по полям
	if group := c.QueryParam("group"); group != "" {
		query = query.Where("group ILIKE ?", "%"+group+"%")
	}
	if song := c.QueryParam("song"); song != "" {
		query = query.Where("song ILIKE ?", "%"+song+"%")
	}

	// Пагинация
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch songs"})
	}

	return c.JSON(http.StatusOK, songs)
}

// Удаление песни
func DeleteSong(c echo.Context) error {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Song{}, id).Error; err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete song"})
	}
	logger.Log.Infof("Song deleted: %s", id)
	return c.NoContent(http.StatusNoContent)
}

// Изменение данных песни
func UpdateSong(c echo.Context) error {
	id := c.Param("id")
	var song models.Song
	if err := database.DB.First(&song, id).Error; err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Song not found"})
	}

	if err := c.Bind(&song); err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := database.DB.Save(&song).Error; err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update song"})
	}

	logger.Log.Infof("Song updated: %s", id)
	return c.JSON(http.StatusOK, song)
}

// Получение текста песни с пагинацией по куплетам
func GetSongText(c echo.Context) error {
	id := c.Param("id")
	var song models.Song
	if err := database.DB.First(&song, id).Error; err != nil {
		logger.Log.Error(err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Song not found"})
	}

	verses := splitIntoVerses(song.Text)

	// Пагинация
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 1
	}
	offset := (page - 1) * limit

	if offset >= len(verses) {
		return c.JSON(http.StatusOK, []string{})
	}

	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := verses[offset:end]
	return c.JSON(http.StatusOK, paginatedVerses)
}

// Помощная функция для разделения текста на куплеты
func splitIntoVerses(text string) []string {
	return strings.Split(text, "\n\n")
}
