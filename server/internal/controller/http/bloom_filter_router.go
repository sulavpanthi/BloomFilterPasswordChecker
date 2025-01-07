package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/internal/usecase"
)

type BloomFilterHandler struct {
	useCase *usecase.BloomFilterUseCase
}

func NewBloomFilterHandler(useCase *usecase.BloomFilterUseCase) *BloomFilterHandler {
	return &BloomFilterHandler{useCase: useCase}
}

func (h *BloomFilterHandler) AddPassword(c *gin.Context) {
	var request struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.useCase.AddPassword(request.Password)
	c.JSON(201, gin.H{"message": "Password added to Bloom Filter"})
}

func (h *BloomFilterHandler) CheckPassword(c *gin.Context) {
	var request struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	match := h.useCase.IsPasswordCommon(request.Password)

	var responseMessage string
	if match {
		responseMessage = "Password is probably present in common password list"
	} else {
		responseMessage = "Password is not present in common password list"
	}
	c.JSON(200, gin.H{"message": responseMessage})
}

func (h *BloomFilterHandler) GetBloomFilter(c *gin.Context) {
	jsonBF := h.useCase.SerializeAsJSON()
	c.JSON(200, jsonBF)
}
