package service

import (
	"context"
	"go.uber.org/zap"
)

func (s *FileService) StartWorkerPool(ctx context.Context, numWorkers, iterations, workerID int) {
	defer s.wg.Done()
	s.logger.Info("Worker started", zap.Int("id", workerID))

	for iter := 1; iter <= iterations; iter++ {
		select {
		case <-ctx.Done():
			s.logger.Info("Worker stopped by context", zap.Int("id", workerID))
			return
		default:
			data := s.generateRandomDataFile(workerID, iter)
			fileIndex := (iter-1)%10 + 1
			s.writerCh <- fileWriteTask{fileIndex, data}

		}
	}
	s.logger.Info("Worker completed", zap.Int("id", workerID))

}
