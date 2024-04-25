package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	DownloadHandler "jupyter-download/helper"
)

func main() {
	isolateJupyter := DownloadHandler.InitDownload()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(
		helmet.New(),
	)

	app.Use(cors.New())

	app.Post("/dl-folder", isolateJupyter.FolderDownload)
	app.Post("/dl-files", isolateJupyter.FilesDownload)
	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":9090"))
}
