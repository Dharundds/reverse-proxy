package main

import (
	"context"
	"os"
	"os/signal"
	"reverse-proxy/internal/constants"
	"reverse-proxy/internal/helpers"
	"reverse-proxy/internal/logger"
	"reverse-proxy/internal/manager"
	"reverse-proxy/internal/server"
	"strconv"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

func serverInit(ctx context.Context) {
	var err error
	logger.NewLogger(
		logger.WithLogFilePath("logs"),
		logger.WithLevel("info"),
	)

	constants.Redis, err = manager.NewRedisManager()
	if err != nil {
		log.Error().Msgf("Error while initialising Redis -> %v", err)
	}

	constants.RPCtxManager = manager.NewRPManager()

}

func serverClose() {
	if constants.Redis != nil {
		constants.Redis.Close()
	}
}

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())
	serverInit(rootCtx)
	defer serverClose()

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	res := helpers.LoadRedisContext()
	if res == nil {
		log.Fatal().Msg("Failed to load Redis context")
	}

	// Get port from environment variable with fallback to 5000
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "80"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Error().Msgf("Invalid PORT value '%s', using default 5000: %v", portStr, err)
		port = 80
	}

	// Get host from environment variable with fallback to "0.0.0.0"
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	go func() {
		srv := server.NewServer(
			server.WithPort(port),
			server.WithAddr(host),
		)
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
