package handlers

import (
	"backend/api/utils"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Abrir la conexión a la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	defer db.Close()

	// Intentar hacer ping a la base de datos
	err = db.Ping()
	if err != nil {
		// Si hay un error al hacer ping, devolver código 403 (Forbidden)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}

	// Si se hizo ping con éxito, devolver código 200 (OK)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
