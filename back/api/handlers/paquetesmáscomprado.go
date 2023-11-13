package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"log"
	"net/http"
)

type ciudadPaqueteOferta struct {
	Ciudad string `json:"ciudad"`
}

func ObtenerPaqueteMasComprado(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data ciudadPaqueteOferta
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
	SELECT
    p.id AS id_paquete,
    COUNT(*) as count_por_paquete,
    p.nombre AS nombre_paquete,
    p.descripcion AS descripcion_paquete,
    p.detalles AS detalles_paquete,
    COALESCE(sq.total_personas, 0) AS total_personas,
    co.nombre AS nombre_ciudad_origen,
    cd.nombre AS nombre_ciudad_destino,
    MIN(p.precio_vuelo) AS precio_vuelo_min,
    MIN(fp.precio_oferta_vuelo) AS oferta_vuelo_min,
    MAX(fp.precio_oferta_vuelo) AS oferta_vuelo_max
FROM
    reserva r
    INNER JOIN fechapaquete fp ON r.id_fechapaquete = fp.id
    INNER JOIN paquete p ON fp.id_paquete = p.id
    INNER JOIN ciudad co ON p.id_origen = co.id
    INNER JOIN ciudad cd ON p.id_destino = cd.id
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
    ) sq ON p.id = sq.paquete_id
WHERE 
    r.estado = 'P' AND co.nombre = $1
GROUP BY 
    p.id, p.nombre, p.descripcion, p.detalles, co.nombre, cd.nombre, sq.total_personas
ORDER BY 
    count_por_paquete DESC;`, data.Ciudad)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	resultados := []models.PaqueteMasComprado{}

	for rows.Next() {
		var paqueteMasComprado models.PaqueteMasComprado

		err := rows.Scan(
			&paqueteMasComprado.ID,
			&paqueteMasComprado.CountPorPaquete,
			&paqueteMasComprado.NombrePaquete,
			&paqueteMasComprado.DescripcionPaquete,
			&paqueteMasComprado.DetallesPaquete,
			&paqueteMasComprado.TotalPersonas,
			&paqueteMasComprado.NombreCiudadOrigen,
			&paqueteMasComprado.NombreCiudadDestino,
			&paqueteMasComprado.PrecioVueloMin,
			&paqueteMasComprado.OfertaVueloMin,
			&paqueteMasComprado.OfertaVueloMax,
		)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error al escanear los resultados", http.StatusInternalServerError)
			return
		}

		// Hacer algo con paqueteMasComprado, por ejemplo, agregarlo a un slice de resultados
		resultados = append(resultados, paqueteMasComprado)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resultados); err != nil {
		log.Fatal(err)
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}
}
