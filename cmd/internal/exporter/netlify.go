package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type NetlifyDeployer struct {
	token   string
	siteID  string
	distDir string
}

type netlifySite struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type netlifyDeployResponse struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	DeployID string `json:"deploy_id"`
}

func NewNetlifyDeployer(token, siteID, distDir string) *NetlifyDeployer {
	return &NetlifyDeployer{
		token:   token,
		siteID:  siteID,
		distDir: distDir,
	}
}

func (d *NetlifyDeployer) Deploy() error {
	if d.siteID == "" {
		return fmt.Errorf("Netlify site ID is required. Run 'pangolin init --netlify' or set 'netlify_site' in config")
	}

	deployID, deployURL, err := d.createDeploy()
	if err != nil {
		return fmt.Errorf("failed to create deploy: %w", err)
	}

	if err := d.uploadFiles(deployID); err != nil {
		return fmt.Errorf("failed to upload files: %w", err)
	}

	if err := d.finalizeDeploy(deployURL); err != nil {
		return fmt.Errorf("failed to finalize deploy: %w", err)
	}

	siteURL := fmt.Sprintf("https://%s.netlify.app", d.siteID)
	fmt.Printf("Successfully deployed to %s\n", siteURL)
	return nil
}

func (d *NetlifyDeployer) createDeploy() (string, string, error) {
	body := map[string]interface{}{
		"site_id": d.siteID,
		"title":   "Pangolin Static Export",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.netlify.com/api/v1/deploys", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var deployResp netlifyDeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&deployResp); err != nil {
		return "", "", err
	}

	return deployResp.DeployID, deployResp.URL, nil
}

func (d *NetlifyDeployer) uploadFiles(deployID string) error {
	entries, err := os.ReadDir(d.distDir)
	if err != nil {
		return err
	}

	client := &http.Client{}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := d.uploadDirectory(deployID, d.distDir, ""); err != nil {
				return err
			}
			continue
		}

		if err := d.uploadFile(client, deployID, entry.Name(), ""); err != nil {
			return err
		}
	}

	return nil
}

func (d *NetlifyDeployer) uploadDirectory(deployID, basePath, relPath string) error {
	dirPath := filepath.Join(basePath, relPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	client := &http.Client{}

	for _, entry := range entries {
		entryPath := filepath.Join(relPath, entry.Name())
		if entry.IsDir() {
			if err := d.uploadDirectory(deployID, basePath, entryPath); err != nil {
				return err
			}
			continue
		}

		if err := d.uploadFile(client, deployID, entry.Name(), filepath.ToSlash(relPath)); err != nil {
			return err
		}
	}

	return nil
}

func (d *NetlifyDeployer) uploadFile(client *http.Client, deployID, filename, subdir string) error {
	var filePath string
	if subdir == "" {
		filePath = filepath.Join(d.distDir, filename)
	} else {
		filePath = filepath.Join(d.distDir, subdir, filename)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Join(subdir, filename))
	if err != nil {
		return err
	}

	if _, err := part.Write(content); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	uploadPath := fmt.Sprintf("https://api.netlify.com/api/v1/deploys/%s/files", deployID)
	req, err := http.NewRequest("POST", uploadPath, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload %s: status %d: %s", filename, resp.StatusCode, string(respBody))
	}

	relativePath := filename
	if subdir != "" {
		relativePath = filepath.Join(subdir, filename)
	}
	fmt.Printf("Uploaded: %s\n", relativePath)

	return nil
}

func (d *NetlifyDeployer) finalizeDeploy(deployURL string) error {
	body := map[string]bool{"published": true}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", deployURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to finalize deploy: status %d", resp.StatusCode)
	}

	return nil
}

func ListNetlifySites(token string) error {
	req, err := http.NewRequest("GET", "https://api.netlify.com/api/v1/sites", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list sites: status %d", resp.StatusCode)
	}

	var sites []netlifySite
	if err := json.NewDecoder(resp.Body).Decode(&sites); err != nil {
		return err
	}

	fmt.Println("Available Netlify sites:")
	for _, site := range sites {
		fmt.Printf("  - %s (ID: %s) -> %s\n", site.Name, site.ID, site.URL)
	}

	return nil
}
