package main

import (
	"context"
	"os"
	"os/signal"
	"reverse-proxy/internal/constants"
	"reverse-proxy/internal/logger"
	"reverse-proxy/internal/manager"
	"reverse-proxy/internal/server"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

func serverInit() {
	var err error
	logger.NewLogger(
		logger.WithLogFilePath("logs"),
		logger.WithLevel("info"),
	)

	constants.Redis, err = manager.NewRedisManager()
	if err != nil {
		log.Error().Msgf("Error while initialising Redis -> %v", err)
	}
}

func serverClose() {
	if constants.Redis != nil {
		constants.Redis.Close()
	}
}

func main() {
	serverInit()
	defer serverClose()

	var wg sync.WaitGroup
	_, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		srv := server.NewServer()
		if err := srv.StartBackend(); err != nil {
			log.Error().Msgf("Error while starting backend -> %v", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigs
		log.Info().Msgf("Received signal: %v", sig)
		cancel()
	}()

	wg.Wait()
}
