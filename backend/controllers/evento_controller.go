package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEventos(c *gin.Context) {
	nombre := c.Query("nombre")
	estado := c.Query("estado")

	eventos, err := services.GetEventos(nombre, estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener eventos",
		})
		return
	}

	var response []dto.EventoDetailResponse

	for _, evento := range eventos {
		response = append(response, dto.EventoDetailResponse{
			ID:          evento.ID,
			Nombre:      evento.Nombre,
			Descripcion: evento.Descripcion,
			Fecha:       evento.Fecha,
			Lugar:       evento.Lugar,
			Estado:      evento.Estado,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetEventoByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	evento, err := services.GetEventoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Evento no encontrado",
		})
		return
	}

	response := dto.EventoDetailResponse{
		ID:          evento.ID,
		Nombre:      evento.Nombre,
		Descripcion: evento.Descripcion,
		Fecha:       evento.Fecha,
		Lugar:       evento.Lugar,
		Estado:      evento.Estado,
	}

	c.JSON(http.StatusOK, response)
}

func GetTiposEntradaByEvento(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	tipos, err := services.GetTiposEntradaByEvento(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tipos de entrada"})
		return
	}

	type tipoEntradaResponse struct {
		ID              uint    `json:"id"`
		Nombre          string  `json:"nombre"`
		Precio          float64 `json:"precio"`
		StockDisponible int     `json:"stock_disponible"`
	}

	var response []tipoEntradaResponse
	for _, t := range tipos {
		response = append(response, tipoEntradaResponse{
			ID:              t.ID,
			Nombre:          t.Nombre,
			Precio:          t.Precio,
			StockDisponible: t.StockDisponible,
		})
	}

	c.JSON(http.StatusOK, response)
}
