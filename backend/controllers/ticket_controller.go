package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	idFloat, ok := value.(float64)
	if !ok {
		return 0, false
	}
	return uint(idFloat), true
}

func ComprarTicket(c *gin.Context) {
	usuarioID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no identificado"})
		return
	}

	var request dto.ComprarConCuponRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	var ticket interface{}
	var err error

	if request.CodigoCupon != "" {
		ticket, err = services.ComprarTicketConCupon(usuarioID, request.TipoEntradaID, request.CodigoCupon)
	} else {
		ticket, err = services.ComprarTicket(usuarioID, request.TipoEntradaID)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Compra realizada con éxito",
		"ticket":  ticket,
	})
}

func GetMisTickets(c *gin.Context) {
	usuarioID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no identificado"})
		return
	}

	tickets, err := services.GetMisTickets(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tickets"})
		return
	}

	var response []dto.TicketResponse
	for _, ticket := range tickets {
		response = append(response, dto.TicketResponse{
			ID:           ticket.ID,
			EventoNombre: ticket.TipoEntrada.Evento.Nombre,
			TipoEntrada:  ticket.TipoEntrada.Nombre,
			PrecioPagado: ticket.PrecioPagado,
			FechaCompra:  ticket.FechaCompra,
			Estado:       ticket.Estado,
		})
	}

	c.JSON(http.StatusOK, response)
}

func CancelarTicket(c *gin.Context) {
	usuarioID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no identificado"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := services.CancelarTicket(usuarioID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Ticket cancelado con éxito"})
}

func TransferirTicket(c *gin.Context) {
	usuarioID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no identificado"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request dto.TransferirTicketRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	if err := services.TransferirTicket(usuarioID, uint(id), request.EmailDestino); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Ticket transferido con éxito"})
}
