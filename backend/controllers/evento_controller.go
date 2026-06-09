package controllers

import (
	"backend/dto"
	"backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetEventos(c *gin.Context) {

	eventos, err := services.GetEventos()

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
