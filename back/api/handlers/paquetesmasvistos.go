package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"log"
	"net/http"
)

func ObtenerMasVistos(w http.ResponseWriter, r *http.Request) {
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
	WITH ranked_packages AS (
		SELECT
			COALESCE(SUM(lp.cantidad_vistas), 0) AS total_vistas,
			paquete.*,
			COALESCE(total_personas, 0) AS total_personas,
			ciudad_origen.nombre AS nombre_ciudad_origen,
			ciudad_destino.nombre AS nombre_ciudad_destino,
			habitacionhotel.id AS habitacion_id,
			habitacionhotel.opcion_hotel_id AS opcion_hotel_id,
			opcionhotel.nombre AS nombre_opcion_hotel,
			habitacionhotel.descripcion AS descripcion_habitacion,
			habitacionhotel.servicios AS servicios_habitacion,
			habitacionhotel.precio_noche AS precio_noche,
			hotel.id AS hotel_id,
			hotel.nombre AS nombre_hotel,
			hotel.ciudad_id AS ciudad_id_hotel,
			hotel.direccion AS direccion_hotel,
			hotel.valoracion AS valoracion_hotel,
			hotel.descripcion AS descripcion_hotel,
			hotel.servicios AS servicios_hotel,
			hotel.telefono AS telefono_hotel,
			hotel.correo_electronico AS correo_electronico_hotel,
			hotel.sitio_web AS sitio_web_hotel,
			ROW_NUMBER() OVER (PARTITION BY paquete.id ORDER BY paquete.id) AS row_num
		FROM
			paquete
			INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
			INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
			INNER JOIN hotel ON habitacionhotel.hotel_id = hotel.id
			INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
			INNER JOIN ciudad ciudad_origen ON paquete.id_origen = ciudad_origen.id
			INNER JOIN ciudad ciudad_destino ON paquete.id_destino = ciudad_destino.id
			INNER JOIN fechapaquete as fp ON paquete.id = fp.id_paquete
			INNER JOIN logs_paquetes as lp ON fp.id = fk_fechapaquete
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
		GROUP BY
		paquete.id,
		paquete.nombre,
		paquete.id_origen,
		paquete.id_destino,
		paquete.descripcion,
		paquete.detalles,
		paquete.precio_vuelo,
		paquete.id_hh,
		paquete.imagenes,
		ciudad_origen.nombre,
		ciudad_destino.nombre,
		habitacionhotel.id,
		habitacionhotel.opcion_hotel_id,
		opcionhotel.nombre,
		habitacionhotel.descripcion,
		habitacionhotel.servicios,
		habitacionhotel.precio_noche,
		hotel.id,
		hotel.nombre,
		hotel.ciudad_id,
		hotel.direccion,
		hotel.valoracion,
		hotel.descripcion,
		hotel.servicios,
		hotel.telefono,
		hotel.correo_electronico,
		hotel.sitio_web,
		total_personas,
		nombre_ciudad_origen,
		nombre_ciudad_destino
	)
	SELECT *
	FROM ranked_packages
	WHERE row_num = 1
	`)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var paquetesMasVistos []models.PaquetesMasVistos

	for rows.Next() {
		var paquetesMasVisto models.PaquetesMasVistos
		var infoPaqueteMasVisto models.PaqueteInfoAdicionalMasVisto
		var hotelInfoMasVisto models.HotelInfoMasVisto
		err := rows.Scan(
			&paquetesMasVisto.Vistas,
			&paquetesMasVisto.ID,
			&paquetesMasVisto.Nombre,
			&paquetesMasVisto.IdOrigen,
			&paquetesMasVisto.IdDestino,
			&paquetesMasVisto.Descripcion,
			&paquetesMasVisto.Detalles,
			&paquetesMasVisto.PrecioVuelo,
			&paquetesMasVisto.ListaHH,
			&paquetesMasVisto.Imagenes,
			&paquetesMasVisto.TotalPersonas,
			&paquetesMasVisto.NombreCiudadOrigen,
			&paquetesMasVisto.NombreCiudadDestino,
			&infoPaqueteMasVisto.HabitacionId,
			&infoPaqueteMasVisto.OpcionHotelId,
			&infoPaqueteMasVisto.NombreOpcionHotel,
			&infoPaqueteMasVisto.DescripcionHabitacion,
			&infoPaqueteMasVisto.ServiciosHabitacion,
			&paquetesMasVisto.PrecioNoche,
			&hotelInfoMasVisto.ID,
			&hotelInfoMasVisto.NombreHotel,
			&hotelInfoMasVisto.CiudadIdHotel,
			&hotelInfoMasVisto.DireccionHotel,
			&hotelInfoMasVisto.ValoracionHotel,
			&hotelInfoMasVisto.DescripcionHotel,
			&hotelInfoMasVisto.ServiciosHotel,
			&hotelInfoMasVisto.TelefonoHotel,
			&hotelInfoMasVisto.CorreoElectronico,
			&hotelInfoMasVisto.SitioWeb,
			&infoPaqueteMasVisto.RowNum,
		)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear los resultados", http.StatusInternalServerError)
			return
		}

		infoPaqueteMasVisto.HotelInfo = hotelInfoMasVisto
		paquetesMasVisto.InfoPaquete = infoPaqueteMasVisto

		paquetesMasVistos = append(paquetesMasVistos, paquetesMasVisto)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paquetesMasVistos); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}

}
