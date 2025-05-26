package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	log "github.com/sirupsen/logrus"

	logger2 "github.com/zachmann/offa/internal/logger"
)

func addMiddlewares(s fiber.Router) {
	addRecoverMiddleware(s)
	addRequestIDMiddleware(s)
	addLoggerMiddleware(s)
	addFaviconMiddleware(s)
	addStaticFiles(s)
	addHelmetMiddleware(s)
	addCompressMiddleware(s)
}

func addLoggerMiddleware(s fiber.Router) {
	s.Use(
		logger.New(
			logger.Config{
				Format:     "${time} ${ip} ${ua} ${latency} - ${status} ${method} ${path} ${locals:requestid}\n",
				TimeFormat: "2006-01-02 15:04:05",
				Output:     logger2.MustGetAccessLogger(),
			},
		),
	)
}

func addCompressMiddleware(s fiber.Router) {
	s.Use(compress.New())
}

func addRecoverMiddleware(s fiber.Router) {
	s.Use(recover.New())
}

func addHelmetMiddleware(s fiber.Router) {
	s.Use(helmet.New())
}

func addRequestIDMiddleware(s fiber.Router) {
	s.Use(requestid.New())
}

//go:embed favicon.ico
var faviconFS embed.FS

//go:embed static
var _staticFS embed.FS
var staticFS fs.FS

func init() {
	var err error
	staticFS, err = fs.Sub(_staticFS, "static")
	if err != nil {
		log.WithError(err).Fatal()
	}
}

func addFaviconMiddleware(s fiber.Router) {
	s.Use(
		favicon.New(
			favicon.Config{
				File:       "favicon.ico",
				FileSystem: http.FS(faviconFS),
			},
		),
	)
}

func addStaticFiles(s fiber.Router) {
	s.Use(
		"/static", filesystem.New(
			filesystem.Config{
				Root: http.FS(staticFS),
			},
		),
	)
}
