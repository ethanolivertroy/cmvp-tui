package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethanolivertroy/cmvp-tui/internal/model"
)

const maxResponseSize = 10 * 1024 * 1024 // 10MB limit for API responses

const (
	BaseURL            = "https://hackidle.github.io/nist-cmvp-api/api"
	BaseURLEnvVar      = "CMVP_API_BASE_URL"
	ModulesEndpoint    = "/modules.json"
	HistoricalEndpoint = "/historical-modules.json"
	InProcessEndpoint  = "/modules-in-process.json"
	MetadataEndpoint   = "/metadata.json"
)

// Client is an HTTP client for the NIST CMVP API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    resolveBaseURL(),
	}
}

func resolveBaseURL() string {
	if override := strings.TrimRight(strings.TrimSpace(os.Getenv(BaseURLEnvVar)), "/"); override != "" {
		return override
	}

	return BaseURL
}

// FetchAllModules fetches all three datasets and combines them
func (c *Client) FetchAllModules() ([]model.Module, error) {
	var allModules []model.Module

	// Fetch active modules
	active, err := c.fetchModules(ModulesEndpoint, model.StatusActive)
	if err != nil {
		return nil, fmt.Errorf("fetching active modules: %w", err)
	}
	allModules = append(allModules, active...)

	// Fetch historical modules
	historical, err := c.fetchModules(HistoricalEndpoint, model.StatusHistorical)
	if err != nil {
		return nil, fmt.Errorf("fetching historical modules: %w", err)
	}
	allModules = append(allModules, historical...)

	// Fetch in-process modules
	inProcess, err := c.fetchInProcessModules()
	if err != nil {
		return nil, fmt.Errorf("fetching in-process modules: %w", err)
	}
	allModules = append(allModules, inProcess...)

	return allModules, nil
}

// FetchMetadata fetches the metadata from the API
func (c *Client) FetchMetadata() (*MetadataJSON, error) {
	resp, err := c.httpClient.Get(c.baseURL + MetadataEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for metadata", resp.StatusCode)
	}

	body := io.LimitReader(resp.Body, maxResponseSize)
	var metadata MetadataJSON
	if err := json.NewDecoder(body).Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (c *Client) fetchModules(endpoint string, status model.ModuleStatus) ([]model.Module, error) {
	resp, err := c.httpClient.Get(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for %s", resp.StatusCode, endpoint)
	}

	body := io.LimitReader(resp.Body, maxResponseSize)
	var response ModulesResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	modules := make([]model.Module, len(response.Modules))
	for i, jm := range response.Modules {
		modules[i] = model.Module{
			CertificateNumber: jm.CertificateNumber,
			CertificateURL:    jm.CertificateNumberURL,
			VendorName:        jm.VendorName,
			ModuleName:        jm.ModuleName,
			ModuleType:        jm.ModuleType,
			ValidationDate:    parseDate(jm.ValidationDate),
			Status:            status,

			// Extended fields
			Standard:           jm.Standard,
			OverallLevel:       parseOverallLevel(jm.OverallLevel),
			SunsetDate:         jm.SunsetDate,
			Caveat:             jm.Caveat,
			Embodiment:         jm.Embodiment,
			Description:        jm.Description,
			Lab:                jm.Lab,
			Algorithms:         jm.Algorithms,
			AlgorithmsDetailed: jm.AlgorithmsDetailed,
			SecurityPolicyURL:  jm.SecurityPolicyURL,
		}
	}
	return modules, nil
}

func (c *Client) fetchInProcessModules() ([]model.Module, error) {
	resp, err := c.httpClient.Get(c.baseURL + InProcessEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for in-process modules", resp.StatusCode)
	}

	body := io.LimitReader(resp.Body, maxResponseSize)
	var response InProcessModulesResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	modules := make([]model.Module, len(response.Modules))
	for i, jm := range response.Modules {
		modules[i] = model.Module{
			CertificateNumber: "", // In-process modules don't have certificates yet
			CertificateURL:    "",
			VendorName:        jm.VendorName,
			ModuleName:        jm.ModuleName,
			ModuleType:        jm.Standard, // Use Standard field as module type
			ValidationDate:    time.Time{}, // No validation date yet
			Status:            model.StatusInProcess,
		}
	}
	return modules, nil
}

// parseOverallLevel converts the API overall_level (can be int or string) to int
func parseOverallLevel(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	default:
		return 0 // String values like "Tested Configuration(s)" become 0
	}
}

// parseDate parses a date string in MM/DD/YYYY format
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("01/02/2006", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
