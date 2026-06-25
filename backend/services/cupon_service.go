package services

import (
	"backend/dao"
	"backend/domain"
	"backend/dto"
	"errors"
	"time"
)

func ValidarCupon(codigo string, tipoEntradaID uint) (*dto.ValidarCuponResponse, error) {
	cupon, err := dao.GetCuponByCodigo(codigo)
	if err != nil {
		return nil, errors.New("cupón no encontrado")
	}

	if cupon.Estado != "ACTIVO" {
		return nil, errors.New("el cupón no está activo")
	}

	if time.Now().After(cupon.FechaVencimiento) {
		return nil, errors.New("el cupón está vencido")
	}

	if cupon.UsosActuales >= cupon.LimiteUsos {
		return nil, errors.New("el cupón alcanzó su límite de usos")
	}

	// Si el cupon tiene evento asignado, verificar que corresponda al evento del tipo de entrada
	if cupon.EventoID != 0 {
		tipoEntrada, err := dao.GetTiposEntradaByID(tipoEntradaID)
		if err != nil {
			return nil, errors.New("tipo de entrada no encontrado")
		}
		if tipoEntrada.EventoID != cupon.EventoID {
			return nil, errors.New("el cupón no es válido para este evento")
		}
	}

	// Obtener precio original
	tipoEntrada, err := dao.GetTiposEntradaByID(tipoEntradaID)
	if err != nil {
		return nil, errors.New("tipo de entrada no encontrado")
	}

	precioOriginal := tipoEntrada.Precio
	var descuento float64

	if cupon.TipoDescuento == "PORCENTAJE" {
		descuento = precioOriginal * cupon.ValorDescuento / 100
	} else {
		descuento = cupon.ValorDescuento
	}

	precioFinal := precioOriginal - descuento
	if precioFinal < 0 {
		precioFinal = 0
	}

	return &dto.ValidarCuponResponse{
		CuponID:        cupon.ID,
		Codigo:         cupon.Codigo,
		TipoDescuento:  cupon.TipoDescuento,
		ValorDescuento: cupon.ValorDescuento,
		PrecioOriginal: precioOriginal,
		PrecioFinal:    precioFinal,
		Descuento:      descuento,
	}, nil
}

func GetCupones() ([]domain.Cupon, error) {
	return dao.GetCupones()
}

func CrearCupon(req dto.CrearCuponRequest) (*domain.Cupon, error) {
	if req.TipoDescuento != "PORCENTAJE" && req.TipoDescuento != "MONTO_FIJO" {
		return nil, errors.New("tipo de descuento inválido: use PORCENTAJE o MONTO_FIJO")
	}

	if req.ValorDescuento <= 0 {
		return nil, errors.New("el valor del descuento debe ser mayor a 0")
	}

	if req.FechaVencimiento.Before(time.Now()) {
		return nil, errors.New("la fecha de vencimiento debe ser futura")
	}

	cupon := &domain.Cupon{
		Codigo:           req.Codigo,
		TipoDescuento:    req.TipoDescuento,
		ValorDescuento:   req.ValorDescuento,
		FechaVencimiento: req.FechaVencimiento,
		LimiteUsos:       req.LimiteUsos,
		UsosActuales:     0,
		Estado:           "ACTIVO",
		EventoID:         req.EventoID,
	}

	if err := dao.CrearCupon(cupon); err != nil {
		return nil, errors.New("error al crear el cupón")
	}

	return cupon, nil
}

func ActualizarCupon(id uint, req dto.ActualizarCuponRequest) (*domain.Cupon, error) {
	cupon, err := dao.GetCuponByID(id)
	if err != nil {
		return nil, errors.New("cupón no encontrado")
	}

	if req.TipoDescuento != "" {
		cupon.TipoDescuento = req.TipoDescuento
	}
	if req.ValorDescuento > 0 {
		cupon.ValorDescuento = req.ValorDescuento
	}
	if !req.FechaVencimiento.IsZero() {
		cupon.FechaVencimiento = req.FechaVencimiento
	}
	if req.LimiteUsos > 0 {
		cupon.LimiteUsos = req.LimiteUsos
	}
	if req.Estado != "" {
		cupon.Estado = req.Estado
	}

	if err := dao.ActualizarCupon(cupon); err != nil {
		return nil, errors.New("error al actualizar el cupón")
	}

	return cupon, nil
}

func DesactivarCupon(id uint) error {
	cupon, err := dao.GetCuponByID(id)
	if err != nil {
		return errors.New("cupón no encontrado")
	}

	cupon.Estado = "INACTIVO"
	return dao.ActualizarCupon(cupon)
}
