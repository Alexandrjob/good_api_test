package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"good_api_test/models"
	"good_api_test/service"
)

func (s *Server) createGoodHandler(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project Id"})
		return
	}

	var good models.Good
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	good.ProjectId = projectId

	if good.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	id, err := s.service.Create(c.Request.Context(), &good)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	good.Id = id

	c.JSON(http.StatusOK, good)
}

func (s *Server) getGoodsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	if offset <= 0 {
		offset = 1
	}

	goodsResponse, err := s.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goodsResponse)
}

func (s *Server) updateGoodHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project Id"})
		return
	}

	var good models.Good
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	good.Id = id
	good.ProjectId = projectId

	if good.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	if err := s.service.Update(c.Request.Context(), id, &good); err != nil {
		if errors.Is(err, service.ErrGoodNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) deleteGoodHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project Id"})
		return
	}

	if err := s.service.Delete(c.Request.Context(), id, projectId); err != nil {
		if errors.Is(err, service.ErrGoodNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "project_id": projectId, "removed": true})
}

func (s *Server) reprioritizeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project Id"})
		return
	}

	var payload struct {
		NewPriority int `json:"newPriority"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goods, err := s.service.Reprioritize(c.Request.Context(), id, projectId, payload.NewPriority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"priorities": goods})
}

func (s *Server) getGoodHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project Id"})
		return
	}

	good, err := s.service.Get(c.Request.Context(), id, projectId)
	if err != nil {
		if errors.Is(err, service.ErrGoodNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, good)
}
