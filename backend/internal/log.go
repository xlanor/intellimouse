package internal

// Backend structure and logging module with reference to
// Cantor by evercyan, released under the MIT license
// https://github.com/evercyan/cantor

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// NewLogger ...
func NewLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	})
	return logger
}
