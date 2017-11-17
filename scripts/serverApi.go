package main

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/handlers"

	"github.com/CheYpe/Ingressou/api"
	"github.com/CheYpe/Ingressou/api/utils/database/postgres"
	"github.com/CheYpe/Ingressou/api/utils/logger"
	"github.com/CheYpe/Ingressou/api/utils/settings"
	"github.com/gorilla/mux"
)

var (
	app *mux.Router
	srv *http.Server
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := postgres.Open(); err != nil {
		if GoDetails, _ := strconv.ParseBool(os.Getenv("GO_DETAILS")); GoDetails {
			logger.Errorf("Conexão Postgres %v", err.Error())
		} else {
			logger.Errorln("Open Postgres")
		}
	}

	app = mux.NewRouter()

	srv = &http.Server{
		Handler:      app,
		Addr:         settings.GetSettings().Listen,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}

	allowedOrigins := handlers.AllowedOrigins([]string{
		"*",
	})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"PUT",
		"DELETE",
		"OPTIONS",
	})

	allowedHeaders := handlers.AllowedHeaders([]string{
		"Content-Type: application/json",
	})

	srv.Handler = handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(app)

}

func main() {

	api.Routes(app)

	logger.Infof("Inciando a aplicação, acesse: %v", "http://"+settings.GetSettings().Listen)

	err := srv.ListenAndServe()
	if err != nil {
		logger.Errorln("Fechando aplicação com erro:", err.Error())
	}
}
