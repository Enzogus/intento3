package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/api/models"
	"backend/api/utils"
)

type SolicitudPaquete struct {
	OrigenID  int    `json:"origen_id"`
	DestinoID int    `json:"destino_id"`
	FechaInit string `json:"fechaInit"`
	FechaFin  string `json:"fechaFin"`
	Personas  int    `json:"personas"`
}

type SolicitudPaqueteMes struct {
	OrigenID  int `json:"origen_id"`
	DestinoID int `json:"destino_id"`
	Mes       int `json:"mes"`
	Personas  int `json:"personas"`
}

func ObtenerPaquetes(w http.ResponseWriter, r *http.Request) {
	var solicitud SolicitudPaquete
	err := json.NewDecoder(r.Body).Decode(&solicitud)
	if err != nil {
		http.Error(w, "Error al parsear la solicitud JSON", http.StatusBadRequest)
		return
	}

	fechaInicio, err := time.Parse("2006-01-02", solicitud.FechaInit)
	if err != nil {
		http.Error(w, "Error al parsear la fecha de inicio", http.StatusBadRequest)
		return
	}

	fechaFin, err := time.Parse("2006-01-02", solicitud.FechaFin)
	if err != nil {
		http.Error(w, "Error al parsear la fecha de fin", http.StatusBadRequest)
		return
	}

	// Abre la conexión con la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza la consulta SQL con los parámetros
	rows, err := db.Query(`
	WITH ranked_dates AS (
		SELECT
			fechapaquete.id,
			paquete.nombre,
			COALESCE(total_personas, 0) AS total_personas,
			fechapaquete.fechainit,
			fechapaquete.fechafin,
			ciudad_origen.id AS id_ciudad_origen,
			ciudad_destino.id AS id_ciudad_destino,
			ciudad_origen.nombre AS nombre_ciudad_origen,
			ciudad_destino.nombre AS nombre_ciudad_destino,
			fechapaquete.precio_oferta_vuelo as oferta_vuelo,
			paquete.precio_vuelo,
			habitacionhotel.precio_noche,
			paquete.descripcion,
			paquete.detalles,
			paquete.imagenes,
			opcionhotel.nombre AS nombre_opcion_hotel,
			habitacionhotel.descripcion AS descripcion_habitacion,
			habitacionhotel.servicios AS servicios_habitacion,
			hotel.nombre AS nombre_hotel,
			hotel.direccion AS direccion_hotel,
			hotel.valoracion AS valoracion_hotel,
			hotel.descripcion AS descripcion_hotel,
			hotel.servicios AS servicios_hotel,
			hotel.telefono AS telefono_hotel,
			hotel.correo_electronico AS correo_electronico_hotel,
			hotel.sitio_web AS sitio_web_hotel,
			ROW_NUMBER() OVER (PARTITION BY fechapaquete.id ORDER BY fechapaquete.id) AS row_num
		FROM
			paquete
			INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
			INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
			INNER JOIN hotel ON habitacionhotel.hotel_id = hotel.id
			INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
			INNER JOIN ciudad ciudad_origen ON paquete.id_origen = ciudad_origen.id
			INNER JOIN ciudad ciudad_destino ON paquete.id_destino = ciudad_destino.id
			INNER JOIN fechapaquete ON paquete.id = fechapaquete.id_paquete
			LEFT JOIN (
				SELECT
					paquete.id AS paquete_id,
					SUM(opcionhotel.cantidad) AS total_personas
				FROM
					paquete
					INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
					INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
					INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
				GROUP BY
					paquete.id
			) AS subquery ON paquete.id = subquery.paquete_id
	)
	SELECT *
	FROM ranked_dates
	WHERE 
		row_num = 1 	
	AND id_ciudad_origen = $1
	AND id_ciudad_destino = $2
	AND fechainit = $3
	AND fechafin = $4
	AND total_personas = $5
    `, solicitud.OrigenID, solicitud.DestinoID, fechaInicio, fechaFin, solicitud.Personas)
	if err != nil {
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Procesa los resultados y crea una estructura para la respuesta JSON
	var paquetes []models.PaqueteInfo
	for rows.Next() {
		var paquete models.PaqueteInfo
		var infoPaquete models.PaqueteInfoAdicional
		var hotelInfo models.HotelInfo
		if err := rows.Scan(
			&paquete.ID,
			&paquete.Nombre,
			&paquete.TotalPersonas,
			&paquete.FechaInit,
			&paquete.FechaFin,
			&paquete.IdOrigen,
			&paquete.IdDestino,
			&paquete.NombreCiudadOrigen,
			&paquete.NombreCiudadDestino,
			&paquete.PrecioOfertaVuelo,
			&paquete.PrecioVuelo,
			&paquete.PrecioNoche,
			&paquete.Descripcion,
			&paquete.Detalles,
			&paquete.Imagenes,
			&infoPaquete.NombreOpcionHotel,
			&infoPaquete.DescripcionHabitacion,
			&infoPaquete.ServiciosHabitacion,
			&hotelInfo.NombreHotel,
			&hotelInfo.DireccionHotel,
			&hotelInfo.ValoracionHotel,
			&hotelInfo.DescripcionHotel,
			&hotelInfo.ServiciosHotel,
			&hotelInfo.TelefonoHotel,
			&hotelInfo.CorreoElectronico,
			&hotelInfo.SitioWeb,
			&infoPaquete.RowNum,
		); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaInit)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paquete.FechaInit = fechaInicio.Format("2006-01-02")
		paquete.FechaFin = fechaFin.Format("2006-01-02")

		infoPaquete.HotelInfo = hotelInfo
		paquete.InfoPaquete = infoPaquete

		paquetes = append(paquetes, paquete)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}

func ObtenerPaquetesMes(w http.ResponseWriter, r *http.Request) {
	var solicitudMes SolicitudPaqueteMes
	err := json.NewDecoder(r.Body).Decode(&solicitudMes)
	if err != nil {
		http.Error(w, "Error al parsear la solicitud JSON", http.StatusBadRequest)
		return
	}
	// Abre la conexión con la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza la consulta SQL con los parámetros
	rows, err := db.Query(`
	WITH ranked_dates AS (
		SELECT
			fechapaquete.id,
			paquete.nombre,
			COALESCE(total_personas, 0) AS total_personas,
			fechapaquete.fechainit,
			fechapaquete.fechafin,
			ciudad_origen.id AS id_ciudad_origen,
			ciudad_destino.id AS id_ciudad_destino,
			ciudad_origen.nombre AS nombre_ciudad_origen,
			ciudad_destino.nombre AS nombre_ciudad_destino,
			fechapaquete.precio_oferta_vuelo as oferta_vuelo,
			paquete.precio_vuelo,
			habitacionhotel.precio_noche,
			paquete.descripcion,
			paquete.detalles,
			paquete.imagenes,
			opcionhotel.nombre AS nombre_opcion_hotel,
			habitacionhotel.descripcion AS descripcion_habitacion,
			habitacionhotel.servicios AS servicios_habitacion,
			hotel.nombre AS nombre_hotel,
			hotel.direccion AS direccion_hotel,
			hotel.valoracion AS valoracion_hotel,
			hotel.descripcion AS descripcion_hotel,
			hotel.servicios AS servicios_hotel,
			hotel.telefono AS telefono_hotel,
			hotel.correo_electronico AS correo_electronico_hotel,
			hotel.sitio_web AS sitio_web_hotel,
			ROW_NUMBER() OVER (PARTITION BY fechapaquete.id ORDER BY fechapaquete.id) AS row_num
		FROM
			paquete
			INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
			INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
			INNER JOIN hotel ON habitacionhotel.hotel_id = hotel.id
			INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
			INNER JOIN ciudad ciudad_origen ON paquete.id_origen = ciudad_origen.id
			INNER JOIN ciudad ciudad_destino ON paquete.id_destino = ciudad_destino.id
			INNER JOIN fechapaquete ON paquete.id = fechapaquete.id_paquete
			LEFT JOIN (
				SELECT
					paquete.id AS paquete_id,
					SUM(opcionhotel.cantidad) AS total_personas
				FROM
					paquete
					INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
					INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
					INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
				GROUP BY
					paquete.id
			) AS subquery ON paquete.id = subquery.paquete_id
	)
	SELECT *
	FROM ranked_dates
	WHERE 
    row_num = 1
	AND id_ciudad_origen = $1
    AND id_ciudad_destino = $2 	
    AND EXTRACT(MONTH FROM fechainit) = $3
    AND EXTRACT(MONTH FROM fechafin) = $3
    AND total_personas = $4;
    `, solicitudMes.OrigenID, solicitudMes.DestinoID, solicitudMes.Mes, solicitudMes.Personas)
	if err != nil {
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var paquetes []models.PaqueteInfo
	for rows.Next() {
		var paquete models.PaqueteInfo
		var infoPaquete models.PaqueteInfoAdicional
		var hotelInfo models.HotelInfo
		if err := rows.Scan(
			&paquete.ID,
			&paquete.Nombre,
			&paquete.TotalPersonas,
			&paquete.FechaInit,
			&paquete.FechaFin,
			&paquete.IdOrigen,
			&paquete.IdDestino,
			&paquete.NombreCiudadOrigen,
			&paquete.NombreCiudadDestino,
			&paquete.PrecioOfertaVuelo,
			&paquete.PrecioVuelo,
			&paquete.PrecioNoche,
			&paquete.Descripcion,
			&paquete.Detalles,
			&paquete.Imagenes,
			&infoPaquete.NombreOpcionHotel,
			&infoPaquete.DescripcionHabitacion,
			&infoPaquete.ServiciosHabitacion,
			&hotelInfo.NombreHotel,
			&hotelInfo.DireccionHotel,
			&hotelInfo.ValoracionHotel,
			&hotelInfo.DescripcionHotel,
			&hotelInfo.ServiciosHotel,
			&hotelInfo.TelefonoHotel,
			&hotelInfo.CorreoElectronico,
			&hotelInfo.SitioWeb,
			&infoPaquete.RowNum,
		); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaInit)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", paquete.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		paquete.FechaInit = fechaInicio.Format("2006-01-02")
		paquete.FechaFin = fechaFin.Format("2006-01-02")

		infoPaquete.HotelInfo = hotelInfo
		paquete.InfoPaquete = infoPaquete

		paquetes = append(paquetes, paquete)
	}

	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
