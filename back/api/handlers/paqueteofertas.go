package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ciudadPaquete struct {
	Ciudad string `json:"ciudad"`
}

func ObtenerPaquetesOfertas(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data ciudadPaquete

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Variable 'ciudad' recibida: %s\n", data.Ciudad)

	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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
		AND oferta_vuelo > 0
		AND nombre_ciudad_origen = $1
		`, data.Ciudad)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var paquetesOfertas []models.PaqueteOferta
	//Itera a través de los resultados y agrega a la slice
	for rows.Next() {
		var PaqueteOferta models.PaqueteOferta
		var infoPaqueteOferta models.PaqueteInfoAdicionalOferta
		var hotelInfoOferta models.HotelInfoOferta

		err := rows.Scan(
			&PaqueteOferta.ID,
			&PaqueteOferta.Nombre,
			&PaqueteOferta.TotalPersonas,
			&PaqueteOferta.FechaInit,
			&PaqueteOferta.FechaFin,
			&PaqueteOferta.IdOrigen,
			&PaqueteOferta.IdDestino,
			&PaqueteOferta.NombreCiudadOrigen,
			&PaqueteOferta.NombreCiudadDestino,
			&PaqueteOferta.PrecioOfertaVuelo,
			&PaqueteOferta.PrecioVuelo,
			&PaqueteOferta.PrecioNoche,
			&PaqueteOferta.Descripcion,
			&PaqueteOferta.Detalles,
			&PaqueteOferta.Imagenes,
			&infoPaqueteOferta.NombreOpcionHotel,
			&infoPaqueteOferta.DescripcionHabitacion,
			&infoPaqueteOferta.ServiciosHabitacion,
			&hotelInfoOferta.NombreHotel,
			&hotelInfoOferta.DireccionHotel,
			&hotelInfoOferta.ValoracionHotel,
			&hotelInfoOferta.DescripcionHotel,
			&hotelInfoOferta.ServiciosHotel,
			&hotelInfoOferta.TelefonoHotel,
			&hotelInfoOferta.CorreoElectronico,
			&hotelInfoOferta.SitioWeb,
			&infoPaqueteOferta.RowNum,
		)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear los resultados", http.StatusInternalServerError)
			return
		}
		fechaInicio, err := time.Parse("2006-01-02T15:04:05Z", PaqueteOferta.FechaInit)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de inicio", http.StatusInternalServerError)
			return
		}

		fechaFin, err := time.Parse("2006-01-02T15:04:05Z", PaqueteOferta.FechaFin)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al parsear la fecha de fin", http.StatusInternalServerError)
			return
		}

		// Formatea las fechas como solo la parte de la fecha (sin la hora)
		PaqueteOferta.FechaInit = fechaInicio.Format("2006-01-02")
		PaqueteOferta.FechaFin = fechaFin.Format("2006-01-02")

		infoPaqueteOferta.HotelInfo = hotelInfoOferta
		PaqueteOferta.InfoPaquete = infoPaqueteOferta

		paquetesOfertas = append(paquetesOfertas, PaqueteOferta)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetesOfertas); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
