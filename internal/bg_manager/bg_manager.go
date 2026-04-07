package bgManager

import (
	"fmt"
	"log/slog"
	"sync"
)

type BgManager struct {
	Logger    *slog.Logger
	WaitGroup sync.WaitGroup
}

func New(logger *slog.Logger) *BgManager {
	return &BgManager{
		Logger: logger,
	}
}

func (bm *BgManager) Run(logger *slog.Logger, fn func()) {
	bm.WaitGroup.Add(1)

	go func() {
		defer bm.WaitGroup.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("recovered from panic: %v", err))
			}
		}()

		fn()
	}()
}
