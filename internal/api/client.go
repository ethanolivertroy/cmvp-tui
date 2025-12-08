package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethanolivertroy/nist-cmvp-cli/internal/model"
)

const (
	BaseURL            = "https://ethanolivertroy.github.io/NIST-CMVP-API/api"
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
		baseURL:    BaseURL,
	}
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

	var metadata MetadataJSON
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
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

	var response ModulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
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

	var response InProcessModulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
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
