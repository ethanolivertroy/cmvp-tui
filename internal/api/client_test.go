package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethanolivertroy/cmvp-tui/internal/model"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.baseURL != BaseURL {
		t.Errorf("baseURL = %v, want %v", client.baseURL, BaseURL)
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestNewClient_UsesEnvOverride(t *testing.T) {
	t.Setenv(BaseURLEnvVar, "https://example.com/cmvp/api/")

	client := NewClient()
	if client.baseURL != "https://example.com/cmvp/api" {
		t.Errorf("baseURL = %v, want %v", client.baseURL, "https://example.com/cmvp/api")
	}
}

func TestClient_FetchAllModules(t *testing.T) {
	// Create mock responses
	modulesResp := ModulesResponse{
		Modules: []ModuleJSON{
			{
				CertificateNumber: "1234",
				VendorName:        "Test Vendor",
				ModuleName:        "Test Module",
				ModuleType:        "Hardware",
				ValidationDate:    "01/15/2024",
			},
		},
	}
	historicalResp := ModulesResponse{
		Modules: []ModuleJSON{
			{
				CertificateNumber: "5678",
				VendorName:        "Historical Vendor",
				ModuleName:        "Historical Module",
				ValidationDate:    "06/20/2020",
			},
		},
	}
	inProcessResp := InProcessModulesResponse{
		Modules: []InProcessModuleJSON{
			{
				ModuleName: "In Process Module",
				VendorName: "IP Vendor",
				Standard:   "FIPS 140-3",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp interface{}
		switch r.URL.Path {
		case "/api/modules.json":
			resp = modulesResp
		case "/api/historical-modules.json":
			resp = historicalResp
		case "/api/modules-in-process.json":
			resp = inProcessResp
		default:
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
	}

	modules, err := client.FetchAllModules()
	if err != nil {
		t.Fatalf("FetchAllModules() error = %v", err)
	}

	if len(modules) != 3 {
		t.Errorf("got %d modules, want 3", len(modules))
	}

	// Check active module
	if modules[0].Status != model.StatusActive {
		t.Errorf("first module status = %v, want StatusActive", modules[0].Status)
	}
	if modules[0].CertificateNumber != "1234" {
		t.Errorf("first module cert = %v, want 1234", modules[0].CertificateNumber)
	}

	// Check historical module
	if modules[1].Status != model.StatusHistorical {
		t.Errorf("second module status = %v, want StatusHistorical", modules[1].Status)
	}

	// Check in-process module
	if modules[2].Status != model.StatusInProcess {
		t.Errorf("third module status = %v, want StatusInProcess", modules[2].Status)
	}
}

func TestClient_FetchAllModules_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
	}

	_, err := client.FetchAllModules()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestClient_FetchMetadata(t *testing.T) {
	metadata := MetadataJSON{
		GeneratedAt:            "2024-01-15T10:00:00Z",
		TotalModules:           100,
		TotalHistoricalModules: 50,
		TotalModulesInProcess:  25,
		Source:                 "NIST CMVP",
		Version:                "1.0",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/metadata.json" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metadata)
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
	}

	result, err := client.FetchMetadata()
	if err != nil {
		t.Fatalf("FetchMetadata() error = %v", err)
	}

	if result.TotalModules != 100 {
		t.Errorf("TotalModules = %v, want 100", result.TotalModules)
	}
	if result.Source != "NIST CMVP" {
		t.Errorf("Source = %v, want NIST CMVP", result.Source)
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    string
	}{
		{
			name:    "valid date",
			input:   "01/15/2024",
			wantErr: false,
			want:    "2024-01-15",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
			want:    "0001-01-01",
		},
		{
			name:    "invalid format",
			input:   "2024-01-15",
			wantErr: false, // Returns zero time, not error
			want:    "0001-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.input)
			got := result.Format("2006-01-02")
			if got != tt.want {
				t.Errorf("parseDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestClient_FetchMetadata_Error(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    string
	}{
		{
			name:       "404 error",
			statusCode: http.StatusNotFound,
			wantErr:    "API returned status 404 for metadata",
		},
		{
			name:       "500 error",
			statusCode: http.StatusInternalServerError,
			wantErr:    "API returned status 500 for metadata",
		},
		{
			name:       "429 rate limit",
			statusCode: http.StatusTooManyRequests,
			wantErr:    "API returned status 429 for metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Error", tt.statusCode)
			}))
			defer server.Close()

			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
				baseURL:    server.URL + "/api",
			}

			_, err := client.FetchMetadata()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestClient_FetchAllModules_StatusErrors(t *testing.T) {
	tests := []struct {
		name        string
		failOn      string
		wantErrPart string
	}{
		{
			name:        "modules endpoint fails",
			failOn:      "/api/modules.json",
			wantErrPart: "fetching active modules",
		},
		{
			name:        "historical endpoint fails",
			failOn:      "/api/historical-modules.json",
			wantErrPart: "fetching historical modules",
		},
		{
			name:        "in-process endpoint fails",
			failOn:      "/api/modules-in-process.json",
			wantErrPart: "fetching in-process modules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == tt.failOn {
					http.Error(w, "Error", http.StatusInternalServerError)
					return
				}
				// Return empty valid responses for other endpoints
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/api/modules-in-process.json" {
					json.NewEncoder(w).Encode(InProcessModulesResponse{})
				} else {
					json.NewEncoder(w).Encode(ModulesResponse{})
				}
			}))
			defer server.Close()

			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
				baseURL:    server.URL + "/api",
			}

			_, err := client.FetchAllModules()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !containsString(err.Error(), tt.wantErrPart) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tt.wantErrPart)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParseOverallLevel(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int
	}{
		{
			name:  "float64 input",
			input: float64(3),
			want:  3,
		},
		{
			name:  "int input",
			input: 2,
			want:  2,
		},
		{
			name:  "string input",
			input: "Tested Configuration(s)",
			want:  0,
		},
		{
			name:  "nil input",
			input: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOverallLevel(tt.input); got != tt.want {
				t.Errorf("parseOverallLevel(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
