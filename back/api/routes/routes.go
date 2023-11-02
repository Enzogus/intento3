package routes

import (
	"backend/api/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router) {
	//allowedOrigins := []string{"http://www.lumonidy.studio", "http://localhost:3000", "http://lumonidy.studio"}
	//c := middleware.CorsMiddleware(allowedOrigins)
	//r.Use(c)

	//Ruta estandar
	r.Handle("/", http.HandlerFunc(handlers.HomeHandler))

	//Ruta para los aeropuertos
	r.Handle("/aeropuertos", http.HandlerFunc(handlers.ListarAeropuertos)).Methods("GET")

	//Ruta para los paquetes
	r.Handle("/paquetes", http.HandlerFunc(handlers.ObtenerPaquetes)).Methods("GET")

	//Ruta para los paquetes por mes
	r.Handle("/paquetes-mes", http.HandlerFunc(handlers.ObtenerPaquetesMes)).Methods("GET")

	//Ruta para los paquetes destacados
	r.Handle("/paquetes-destacados", http.HandlerFunc(handlers.ObtenerPaquetesDestacados)).Methods("GET")

	//Ruta para los paquetes en oferta
	r.Handle("/paquetes-ofertas", http.HandlerFunc(handlers.ObtenerPaquetesOfertas)).Methods("GET")

	//Ruta para los paquetes mas vistos
	r.Handle("/paquetes-mas-vistos", http.HandlerFunc(handlers.ObtenerMasVistos)).Methods("GET")

	//Ruta para añadir una visita
	r.Handle("/anadir-vista", http.HandlerFunc(handlers.AgregarVista)).Methods("POST")

	//Ruta para los aeropuertos
	r.Handle("/aeropuerto", http.HandlerFunc(handlers.ObtenerAeropuertos))
}
