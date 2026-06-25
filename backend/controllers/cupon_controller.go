package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Cliente: validar un cupón antes de comprar
func ValidarCupon(c *gin.Context) {
	var request dto.ValidarCuponRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	response, err := services.ValidarCupon(request.Codigo, request.TipoEntradaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Admin: listar todos los cupones
func GetCupones(c *gin.Context) {
	cupones, err := services.GetCupones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cupones"})
		return
	}
	c.JSON(http.StatusOK, cupones)
}

// Admin: crear cupón
func CrearCupon(c *gin.Context) {
	var request dto.CrearCuponRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	cupon, err := services.CrearCupon(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cupon)
}

// Admin: actualizar cupón
func ActualizarCupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request dto.ActualizarCuponRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	cupon, err := services.ActualizarCupon(uint(id), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cupon)
}

// Admin: desactivar cupón
func DesactivarCupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := services.DesactivarCupon(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Cupón desactivado con éxito"})
}
