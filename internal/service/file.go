package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
	"workFileData/internal/domain"
	"workFileData/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type FileService struct {
	logger   *logger.Logger
	wg       sync.WaitGroup
	writerCh chan fileWriteTask
	cancel   context.CancelFunc
}

type fileWriteTask struct {
	fileIndex int
	data      domain.FileData
}

const (
	numWorkers      = 10
	countIterations = 12
)

func NewFileService(logger *logger.Logger) *FileService {
	if err := os.MkdirAll("files", 0755); err != nil {
		logger.Error("Failed to create files directory", zap.Error(err))
	}

	for i := 1; i <= 10; i++ {
		filename := fmt.Sprintf("files/file_%d.yaml", i)
		if _, err := os.Create(filename); err != nil {
			logger.Error("Failed to create file", zap.String("filename", filename), zap.Error(err))
		}
	}

	service := &FileService{
		logger:   logger,
		writerCh: make(chan fileWriteTask),
	}

	return service
}

func (s *FileService) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)

	go s.FileWriter(ctx)

	for i := 1; i <= numWorkers; i++ {
		s.wg.Add(1)
		go s.StartWorkerPool(ctx, numWorkers, countIterations, i)
	}

}

func (s *FileService) FileWriter(ctx context.Context) {
	defer func() {
		s.wg.Done()
		s.logger.Info("File writer fully stopped")
	}()

	s.wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("File writer stopped by context")
			return
		default:
			task, ok := <-s.writerCh
			if !ok {
				s.logger.Info("File writer channel closed")
				return
			}
			s.logger.Info("Received file write task", zap.String("data", task.data.Name))
			filename := fmt.Sprintf("files/file_%d.yaml", task.fileIndex)
			if err := s.saveToFile(filename, task.data); err != nil {
				s.logger.Error("Failed to save file",
					zap.Error(err),
					zap.String("filename", filename))
			}
		}
	}
}

func (s *FileService) generateRandomDataFile(workerID, iteration int) domain.FileData {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

	tags := make([]string, r.Intn(5)+1)
	for i := range tags {
		tags[i] = fmt.Sprintf("tag-%d-%d", workerID, rand.Intn(100))
	}

	metadata := make(map[string]interface{})
	metadata["worker_id"] = workerID
	metadata["iteration"] = iteration
	metadata["random_float"] = rand.Float64() * 100
	metadata["is_valid"] = rand.Intn(2) == 1
	metadata["nested"] = map[string]interface{}{
		"level": r.Intn(10),
		"type":  fmt.Sprintf("type-%d", r.Intn(5)),
	}

	return domain.FileData{
		ID:        uuid.New(),
		Name:      fmt.Sprintf("worker-%d-iter-%d", workerID, iteration),
		Active:    r.Intn(2) == 1,
		Temp:      r.Float64()*100 + 20,
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}
}

func (s *FileService) saveToFile(fileName string, data domain.FileData) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return nil
}

func (s *FileService) ReadFile(file string) ([]byte, error) {
	fileID, err := strconv.Atoi(file)
	if err != nil || fileID < 1 || fileID > 10 {
		return nil, fmt.Errorf("invalid file ID: %w", err)
	}

	filename := fmt.Sprintf("files/file_%d.yaml", fileID)

	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}
