package models

type PaquetesDestacados struct {
	Vistas              int                           `json:"total_vistas"`
	ID                  int                           `json:"id"`
	Nombre              string                        `json:"nombre"`
	IdOrigen            int                           `json:"id_ciudad_origen"`
	IdDestino           int                           `json:"id_ciudad_destino"`
	Descripcion         string                        `json:"descripcion"`
	Detalles            string                        `json:"detalles"`
	PrecioVuelo         float64                       `json:"precio_vuelo"`
	ListaHH             string                        `json:"id_hh"`
	Imagenes            string                        `json:"imagenes"`
	TotalPersonas       int                           `json:"total_personas"`
	NombreCiudadOrigen  string                        `json:"nombre_ciudad_origen"`
	NombreCiudadDestino string                        `json:"nombre_ciudad_destino"`
	InfoPaquete         PaqueteInfoAdicionalDestacado `json:"info_paquete"`
	PrecioNoche         float64                       `json:"precio_noche"`
}

type PaqueteInfoAdicionalDestacado struct {
	HabitacionId          int                `json:"habitacion_id"`
	OpcionHotelId         int                `json:"opcion_hotel_id"`
	NombreOpcionHotel     string             `json:"nombre_opcion_hotel"`
	DescripcionHabitacion string             `json:"descripcion_habitacion"`
	ServiciosHabitacion   string             `json:"servicios_habitacion"`
	HotelInfo             HotelInfoDestacado `json:"hotel_info"`
	RowNum                int                `json:"row_num"`
}

type HotelInfoDestacado struct {
	ID                int     `json:"id_hotel"`
	NombreHotel       string  `json:"nombre_hotel"`
	CiudadIdHotel     int     `json:"ciudad_id_hotel"`
	DireccionHotel    string  `json:"direccion_hotel"`
	ValoracionHotel   float64 `json:"valoracion_hotel"`
	DescripcionHotel  string  `json:"descripcion_hotel"`
	ServiciosHotel    string  `json:"servicios_hotel"`
	TelefonoHotel     string  `json:"telefono_hotel"`
	CorreoElectronico string  `json:"correo_electronico_hotel"`
	SitioWeb          string  `json:"sitio_web_hotel"`
}
