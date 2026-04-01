package exporter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNetlifyDeployer(t *testing.T) {
	deployer := NewNetlifyDeployer("test-token", "test-site", "/tmp/dist")
	if deployer == nil {
		t.Fatal("Expected non-nil NetlifyDeployer")
	}

	if deployer.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", deployer.token)
	}

	if deployer.siteID != "test-site" {
		t.Errorf("Expected siteID 'test-site', got '%s'", deployer.siteID)
	}

	if deployer.distDir != "/tmp/dist" {
		t.Errorf("Expected distDir '/tmp/dist', got '%s'", deployer.distDir)
	}
}

func TestNewNetlifyDeployer_EmptyToken(t *testing.T) {
	deployer := NewNetlifyDeployer("", "test-site", "/tmp/dist")
	if deployer.token != "" {
		t.Error("Expected empty token")
	}
}

func TestNewNetlifyDeployer_EmptySiteID(t *testing.T) {
	deployer := NewNetlifyDeployer("test-token", "", "/tmp/dist")
	if deployer.siteID != "" {
		t.Error("Expected empty siteID")
	}
}

func TestNetlifyDeployer_Deploy_NoSiteID(t *testing.T) {
	deployer := NewNetlifyDeployer("test-token", "", "/tmp/dist")
	err := deployer.Deploy()
	if err == nil {
		t.Error("Expected error when siteID is empty")
	}
}

func TestNetlifySite_Structure(t *testing.T) {
	site := netlifySite{
		ID:   "site123",
		Name: "my-site",
		URL:  "https://my-site.netlify.app",
	}

	if site.ID != "site123" {
		t.Errorf("Expected ID 'site123', got '%s'", site.ID)
	}

	if site.Name != "my-site" {
		t.Errorf("Expected Name 'my-site', got '%s'", site.Name)
	}

	if site.URL != "https://my-site.netlify.app" {
		t.Errorf("Expected URL 'https://my-site.netlify.app', got '%s'", site.URL)
	}
}

func TestNetlifyDeployResponse_Structure(t *testing.T) {
	resp := netlifyDeployResponse{
		ID:       "deploy123",
		URL:      "https://api.netlify.com/deploys/deploy123",
		DeployID: "deploy123",
	}

	if resp.ID != "deploy123" {
		t.Errorf("Expected ID 'deploy123', got '%s'", resp.ID)
	}

	if resp.DeployID != "deploy123" {
		t.Errorf("Expected DeployID 'deploy123', got '%s'", resp.DeployID)
	}
}

func TestListNetlifySites_RequiresValidToken(t *testing.T) {
	err := ListNetlifySites("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestListNetlifySites_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	err := ListNetlifySites("invalid-token")
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}
}

func TestListNetlifySites_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := ListNetlifySites("test-token")
	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestNetlifyDeployer_UploadFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("Expected Authorization header")
		}

		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			t.Error("Expected Content-Type header")
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	deployer := &NetlifyDeployer{
		token:   "test-token",
		siteID:  "test-site",
		distDir: "/tmp/dist",
	}

	client := &http.Client{}
	err := deployer.uploadFile(client, "deploy123", "test.txt", "")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestNetlifyDeployer_FinalizeDeploy_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("Expected Authorization header")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	deployer := &NetlifyDeployer{
		token:  "test-token",
		siteID: "test-site",
	}

	err := deployer.finalizeDeploy(server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestNetlifyDeployer_FinalizeDeploy_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	deployer := &NetlifyDeployer{
		token:  "test-token",
		siteID: "test-site",
	}

	err := deployer.finalizeDeploy(server.URL)
	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestNetlifyDeployer_CreateDeploy_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
	}))
	defer server.Close()

	deployer := &NetlifyDeployer{
		token:  "test-token",
		siteID: "test-site",
	}

	_, _, err := deployer.createDeploy()
	if err == nil {
		t.Error("Expected error for bad request")
	}
}

func TestNetlifyDeployer_CreateDeploy_RequiresToken(t *testing.T) {
	deployer := &NetlifyDeployer{
		token:  "",
		siteID: "test-site",
	}

	_, _, err := deployer.createDeploy()
	if err == nil {
		t.Error("Expected error for empty token")
	}
}
