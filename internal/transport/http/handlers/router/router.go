package router

import (
	"workFileData/internal/transport/http/handlers"
	"workFileData/pkg/logger"

	"github.com/gorilla/mux"
)

func NewRouter(log *logger.Logger, h *handlers.FileDataHandlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/files/{id}", h.GetFileByID).Methods("GET")

	return r
}
