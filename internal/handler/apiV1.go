package handler

import (
	"awesomeProjectRentaTeam/internal/logger"
	"awesomeProjectRentaTeam/internal/model"
	"awesomeProjectRentaTeam/pkg/erx"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogrusInstance()
}

func (h *Handler) GetPostsWithLimit(c *gin.Context) {
	limit := c.DefaultQuery("limit", "25")
	offset := c.DefaultQuery("offset", "0")
	tag := c.DefaultQuery("tag", "")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}

	res, err := h.service.GetPostsWithLimit(tag, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get data"})
		log.Error(erx.New(err))
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) InsertPost(c *gin.Context) {

	reqInputs := model.Post{}

	err := c.ShouldBindJSON(&reqInputs)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse data"})
		return
	}

	err = h.service.InsertPostData(reqInputs)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to write data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}

func (h *Handler) GenerateRandomPosts(c *gin.Context) {

	count := c.DefaultQuery("count", "30")

	countInt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}

	err = h.service.GenerateRandomPosts(countInt)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to write data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}
