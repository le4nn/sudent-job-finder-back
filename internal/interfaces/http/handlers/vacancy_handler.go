package handlers

import (
	"net/http"

	"github.com/albkvv/student-job-finder-back/internal/application/usecases"
	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

type VacancyHandler struct {
	Service *usecases.VacancyService
}

func NewVacancyHandler(service *usecases.VacancyService) *VacancyHandler {
	return &VacancyHandler{Service: service}
}

// CreateVacancy создает новую вакансию
// POST /api/vacancies
func (h *VacancyHandler) CreateVacancy(c *gin.Context) {
	var req entities.Vacancy
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.Service.CreateVacancy(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "vacancy created successfully",
		"data": req,
	})
}

// GetVacancy получает вакансию по ID
// GET /api/vacancies/:id
func (h *VacancyHandler) GetVacancy(c *gin.Context) {
	id := c.Param("id")
	
	vacancy, err := h.Service.GetVacancy(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "vacancy not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vacancy,
	})
}

// GetAllVacancies получает все вакансии с фильтрацией
// GET /api/vacancies?status=Активна
func (h *VacancyHandler) GetAllVacancies(c *gin.Context) {
	status := c.Query("status")
	
	vacancies, err := h.Service.GetAllVacancies(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vacancies,
		"count": len(vacancies),
	})
}

// UpdateVacancy обновляет вакансию
// PUT /api/vacancies/:id
func (h *VacancyHandler) UpdateVacancy(c *gin.Context) {
	id := c.Param("id")
	
	var req entities.Vacancy
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	req.ID = id

	if err := h.Service.UpdateVacancy(c.Request.Context(), &req); err != nil {
		if err.Error() == "vacancy not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vacancy updated successfully",
		"data": req,
	})
}

// UpdateVacancyStatus обновляет статус вакансии
// PATCH /api/vacancies/:id/status
func (h *VacancyHandler) UpdateVacancyStatus(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.Service.UpdateVacancyStatus(c.Request.Context(), id, req.Status); err != nil {
		if err.Error() == "vacancy not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vacancy status updated successfully",
		"status": req.Status,
	})
}

// DeleteVacancy удаляет вакансию
// DELETE /api/vacancies/:id
func (h *VacancyHandler) DeleteVacancy(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.DeleteVacancy(c.Request.Context(), id); err != nil {
		if err.Error() == "vacancy not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vacancy deleted successfully",
	})
}
