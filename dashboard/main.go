package main

import (
	"encoding/json"
	"fmt"
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
	SiteURL     string
	S3Bucket    string
	S3Region    string
	Exports     []ExportLog
	Deploys     []DeployLog
	InProgress  bool
	PagesCount  int
	AssetsCount int
}

type ExportLog struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Pages     int    `json:"pages"`
	Assets    int    `json:"assets"`
	Status    string `json:"status"`
	Duration  string `json:"duration"`
}

type DeployLog struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Bucket    string `json:"bucket"`
	Files     int    `json:"files"`
	Status    string `json:"status"`
	Duration  string `json:"duration"`
}

var (
	exports          = []ExportLog{}
	deploys          = []DeployLog{}
	exportInProgress = false
	deployInProgress = false
	exportStartTime  time.Time
	deployStartTime  time.Time
)

func main() {
	loadHistory()

	engine := html.New("./dashboard/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Static("/static", "./dashboard/static")

	app.Get("/", indexHandler)
	app.Post("/export", exportHandler)
	app.Get("/export-status", exportStatusHandler)
	app.Post("/deploy", deployHandler)
	app.Get("/deploy-status", deployStatusHandler)
	app.Get("/history", historyHandler)
	app.Get("/settings", settingsHandler)
	app.Post("/settings", saveSettingsHandler)

	log.Println("Dashboard starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func loadHistory() {
	data, err := os.ReadFile(".pangolin-history.json")
	if err != nil {
		return
	}
	var history struct {
		Exports []ExportLog `json:"exports"`
		Deploys []DeployLog `json:"deploys"`
	}
	json.Unmarshal(data, &history)
	exports = history.Exports
	deploys = history.Deploys
}

func saveHistory() {
	data, _ := json.MarshalIndent(struct {
		Exports []ExportLog `json:"exports"`
		Deploys []DeployLog `json:"deploys"`
	}{exports, deploys}, "", "  ")
	os.WriteFile(".pangolin-history.json", data, 0644)
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
	siteURL, s3Bucket, s3Region := loadConfig()

	data := PageData{
		SiteURL:     siteURL,
		S3Bucket:    s3Bucket,
		S3Region:    s3Region,
		Exports:     exports,
		Deploys:     deploys,
		InProgress:  exportInProgress || deployInProgress,
		PagesCount:  countPages(),
		AssetsCount: countAssets(),
	}

	return c.Render("dashboard/views/index.html", data)
}

func exportStatusHandler(c *fiber.Ctx) error {
	pages := 0
	assets := 0
	status := "idle"
	duration := ""

	if exportInProgress {
		status = "running"
		pages = countPages()
		assets = countAssets()
		duration = time.Since(exportStartTime).Round(time.Second).String()
	}

	return c.JSON(fiber.Map{
		"in_progress": exportInProgress,
		"status":      status,
		"pages":       pages,
		"assets":      assets,
		"duration":    duration,
	})
}

func deployStatusHandler(c *fiber.Ctx) error {
	duration := ""

	if deployInProgress {
		duration = time.Since(deployStartTime).Round(time.Second).String()
	}

	return c.JSON(fiber.Map{
		"in_progress": deployInProgress,
		"status":      map[bool]string{true: "deploying", false: "idle"}[deployInProgress],
		"duration":    duration,
	})
}

func historyHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"exports": exports,
		"deploys": deploys,
	})
}

func exportHandler(c *fiber.Ctx) error {
	siteURL, _, _ := loadConfig()

	if siteURL == "" {
		return c.JSON(fiber.Map{"error": "Not initialized. Run 'pangolin init' first"})
	}

	if exportInProgress {
		return c.JSON(fiber.Map{"error": "Export already in progress"})
	}

	exportInProgress = true
	exportStartTime = time.Now()

	exports = append([]ExportLog{{
		ID:        len(exports) + 1,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Pages:     0,
		Assets:    0,
		Status:    "running",
		Duration:  "0s",
	}}, exports...)

	go func() {
		cmd := exec.Command("./pangolin", "export", "-d", "dist")
		cmd.Env = os.Environ()
		cmd.Run()

		exportInProgress = false
		duration := time.Since(exportStartTime).Round(time.Second).String()

		if len(exports) > 0 {
			exports[0].Pages = countPages()
			exports[0].Assets = countAssets()
			exports[0].Status = "completed"
			exports[0].Duration = duration
		}
		saveHistory()
	}()

	return c.JSON(fiber.Map{"status": "started"})
}

func deployHandler(c *fiber.Ctx) error {
	_, s3Bucket, s3Region := loadConfig()

	if s3Bucket == "" {
		return c.JSON(fiber.Map{"error": "S3 bucket not configured"})
	}

	if deployInProgress {
		return c.JSON(fiber.Map{"error": "Deploy already in progress"})
	}

	deployInProgress = true
	deployStartTime = time.Now()

	deploys = append([]DeployLog{{
		ID:        len(deploys) + 1,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Bucket:    s3Bucket,
		Files:     countPages() + countAssets(),
		Status:    "running",
		Duration:  "0s",
	}}, deploys...)

	go func() {
		cmd := exec.Command("./pangolin", "deploy", "-b", s3Bucket, "-r", s3Region)
		cmd.Env = os.Environ()
		cmd.Run()

		deployInProgress = false
		duration := time.Since(deployStartTime).Round(time.Second).String()

		if len(deploys) > 0 {
			deploys[0].Status = "completed"
			deploys[0].Duration = duration
		}
		saveHistory()
	}()

	return c.JSON(fiber.Map{"status": "deploying"})
}

func settingsHandler(c *fiber.Ctx) error {
	siteURL, s3Bucket, s3Region := loadConfig()

	data := fiber.Map{
		"SiteURL":  siteURL,
		"S3Bucket": s3Bucket,
		"S3Region": s3Region,
	}

	return c.Render("dashboard/views/settings.html", data)
}

func saveSettingsHandler(c *fiber.Ctx) error {
	siteURL := c.FormValue("site_url")
	apiKey := c.FormValue("api_key")
	s3Bucket := c.FormValue("s3_bucket")
	s3Region := c.FormValue("s3_region")

	viper.Set("site_url", siteURL)
	viper.Set("api_key", apiKey)
	viper.Set("s3_bucket", s3Bucket)
	viper.Set("s3_region", s3Region)

	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.pangolin"
	os.MkdirAll(configDir, 0755)

	viper.SetConfigName("pangolin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.WriteConfig()

	return c.Redirect("/settings?success=true")
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

func countAssets() int {
	total := 0

	paths := []string{"dist/images", "dist/assets"}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil || !info.IsDir() {
			continue
		}
		entries, _ := os.ReadDir(p)
		total += len(entries)
	}

	return total
}

func init() {
	fmt.Println("Pangolin Dashboard")
}
