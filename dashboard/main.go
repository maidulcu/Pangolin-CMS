package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
)

type PageData struct {
	SiteURL    string
	Exports    []ExportLog
	InProgress bool
	PagesCount int
}

type ExportLog struct {
	ID        int
	Timestamp string
	Pages     int
	Status    string
}

var exports = []ExportLog{}

func main() {
	engine := html.New("./dashboard/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Static("/static", "./dashboard/static")

	app.Get("/", indexHandler)
	app.Post("/export", exportHandler)
	app.Get("/exports", exportsHandler)
	app.Post("/deploy", deployHandler)

	log.Println("Dashboard starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func loadConfig() (string, string, string) {
	viper.SetConfigName("pangolin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.pangolin")

	viper.ReadInConfig()

	siteURL := viper.GetString("site_url")
	s3Bucket := viper.GetString("s3_bucket")
	s3Region := viper.GetString("s3_region")
	if s3Region == "" {
		s3Region = "us-east-1"
	}

	return siteURL, s3Bucket, s3Region
}

func indexHandler(c *fiber.Ctx) error {
	siteURL, _, _ := loadConfig()

	data := PageData{
		SiteURL:    siteURL,
		Exports:    exports,
		InProgress: false,
		PagesCount: countPages(),
	}

	return c.Render("dashboard/views/index.html", data)
}

func exportHandler(c *fiber.Ctx) error {
	siteURL, _, _ := loadConfig()

	if siteURL == "" {
		return c.JSON(fiber.Map{"error": "Not initialized. Run 'pangolin init' first"})
	}

	exports = append([]ExportLog{{
		ID:        len(exports) + 1,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Pages:     0,
		Status:    "running",
	}}, exports...)

	go func() {
		cmd := exec.Command("./pangolin", "export", "-d", "dist")
		cmd.Env = os.Environ()
		cmd.Run()

		pages := countPages()
		if len(exports) > 0 {
			exports[0].Pages = pages
			exports[0].Status = "completed"
		}
	}()

	return c.JSON(fiber.Map{"status": "started"})
}

func exportsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"exports": exports})
}

func deployHandler(c *fiber.Ctx) error {
	_, s3Bucket, s3Region := loadConfig()

	if s3Bucket == "" {
		return c.JSON(fiber.Map{"error": "S3 bucket not configured"})
	}

	go func() {
		cmd := exec.Command("./pangolin", "deploy", "-b", s3Bucket, "-r", s3Region)
		cmd.Env = os.Environ()
		cmd.Run()
	}()

	return c.JSON(fiber.Map{"status": "deploying"})
}

func countPages() int {
	info, err := os.Stat("dist")
	if err != nil || !info.IsDir() {
		return 0
	}

	count := 0
	entries, _ := os.ReadDir("dist")
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}
