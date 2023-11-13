package models

type PaqueteMasComprado struct {
	ID                  int     `json:"id_paquete"`
	CountPorPaquete     int     `json:"count_por_paquete"`
	NombrePaquete       string  `json:"nombre_paquete"`
	DescripcionPaquete  string  `json:"descripcion_paquete"`
	DetallesPaquete     string  `json:"detalles_paquete"`
	TotalPersonas       int     `json:"total_personas"`
	NombreCiudadOrigen  string  `json:"nombre_ciudad_origen"`
	NombreCiudadDestino string  `json:"nombre_ciudad_destino"`
	PrecioVueloMin      float64 `json:"precio_vuelo_min"`
	OfertaVueloMin      float64 `json:"oferta_vuelo_min"`
	OfertaVueloMax      float64 `json:"oferta_vuelo_max"`
}

// Otros modelos y funciones si es necesario...
