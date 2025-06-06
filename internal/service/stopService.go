package service

import (
	"fmt"
	"time"
)

func (s *FileService) Stop() {
    if s.cancel != nil {
        s.cancel()
    }

    select {
    case <-s.writerCh: 
    default:
        close(s.writerCh)
    }

    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        fmt.Println("All workers completed")
    case <-time.After(5 * time.Second):
        fmt.Println("Timeout waiting for workers")
    }

    fmt.Println("FileService stopped.")
}
