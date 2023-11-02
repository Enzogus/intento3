package models

type PaqueteOferta struct {
	ID                  int                        `json:"id"`
	Nombre              string                     `json:"nombre"`
	TotalPersonas       int                        `json:"total_personas"`
	FechaInit           string                     `json:"fecha_init"`
	FechaFin            string                     `json:"fecha_fin"`
	IdOrigen            int                        `json:"id_ciudad_origen"`
	IdDestino           int                        `json:"id_ciudad_destino"`
	NombreCiudadOrigen  string                     `json:"nombre_ciudad_origen"`
	NombreCiudadDestino string                     `json:"nombre_ciudad_destino"`
	PrecioOfertaVuelo   float64                    `json:"precio_oferta_vuelo"`
	PrecioVuelo         float64                    `json:"precio_vuelo"`
	PrecioNoche         float64                    `json:"precio_noche"`
	Descripcion         string                     `json:"descripcion"`
	Detalles            string                     `json:"detalles"`
	Imagenes            string                     `json:"imagenes"`
	InfoPaquete         PaqueteInfoAdicionalOferta `json:"info_paquete"`
}

type PaqueteInfoAdicionalOferta struct {
	NombreOpcionHotel     string          `json:"nombre_opcion_hotel"`
	DescripcionHabitacion string          `json:"descripcion_habitacion"`
	ServiciosHabitacion   string          `json:"servicios_habitacion"`
	HotelInfo             HotelInfoOferta `json:"hotel_info"`
	RowNum                int             `json:"row_num"`
}

type HotelInfoOferta struct {
	NombreHotel       string  `json:"nombre_hotel"`
	DireccionHotel    string  `json:"direccion_hotel"`
	ValoracionHotel   float64 `json:"valoracion_hotel"`
	DescripcionHotel  string  `json:"descripcion_hotel"`
	ServiciosHotel    string  `json:"servicios_hotel"`
	TelefonoHotel     string  `json:"telefono_hotel"`
	CorreoElectronico string  `json:"correo_electronico_hotel"`
	SitioWeb          string  `json:"sitio_web_hotel"`
}
