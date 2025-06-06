package handlers

import (
	"errors"
	"net/http"
	"workFileData/internal/domain"
	"workFileData/internal/service"
	"workFileData/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type FileDataHandlers struct {
	logger  *logger.Logger
	service *service.FileService
}

func NewFileDataHandlers(logger *logger.Logger, service *service.FileService) *FileDataHandlers {
	return &FileDataHandlers{logger: logger, service: service}
}

func (h *FileDataHandlers) GetFileByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileIDStr := vars["id"]

	content, err := h.service.ReadFile(fileIDStr)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			http.Error(w, domain.ErrFileNotFound.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("failed to read file", zap.Error(err))
		http.Error(w, domain.ErrFailedReadFile.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("file read successfully")

	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
