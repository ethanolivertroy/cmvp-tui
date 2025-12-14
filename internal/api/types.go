package api

// ModuleJSON represents the JSON structure from the API
// Used for active and historical modules
type ModuleJSON struct {
	CertificateNumber    string `json:"Certificate Number"`
	CertificateNumberURL string `json:"Certificate Number_url"`
	VendorName           string `json:"Vendor Name"`
	ModuleName           string `json:"Module Name"`
	ModuleType           string `json:"Module Type"`
	ValidationDate       string `json:"Validation Date"`

	// Extended fields from certificate detail extraction
	Standard           string   `json:"standard"`
	Status             string   `json:"status"`
	OverallLevel       interface{} `json:"overall_level"` // Can be int or string from API
	SunsetDate         string   `json:"sunset_date"`
	Caveat             string   `json:"caveat"`
	Embodiment         string   `json:"embodiment"`
	Description        string   `json:"description"`
	Lab                string   `json:"lab"`
	Algorithms         []string `json:"algorithms"`
	AlgorithmsDetailed []string `json:"algorithms_detailed"`
	SecurityPolicyURL  string   `json:"security_policy_url"`
	CertificateDetailURL string `json:"certificate_detail_url"`
}

// InProcessModuleJSON has slightly different structure for modules in process
type InProcessModuleJSON struct {
	ModuleName string `json:"Module Name"`
	VendorName string `json:"Vendor Name"`
	Standard   string `json:"Standard"`
	Status     string `json:"Status"`
}

// MetadataJSON represents the metadata endpoint response
type MetadataJSON struct {
	GeneratedAt            string `json:"generated_at"`
	TotalModules           int    `json:"total_modules"`
	TotalHistoricalModules int    `json:"total_historical_modules"`
	TotalModulesInProcess  int    `json:"total_modules_in_process"`
	Source                 string `json:"source"`
	Version                string `json:"version"`
}

// ModulesResponse is the wrapper for the modules.json endpoint
type ModulesResponse struct {
	Metadata MetadataJSON `json:"metadata"`
	Modules  []ModuleJSON `json:"modules"`
}

// HistoricalModulesResponse is the wrapper for historical-modules.json
type HistoricalModulesResponse struct {
	Metadata MetadataJSON `json:"metadata"`
	Modules  []ModuleJSON `json:"modules"`
}

// InProcessModulesResponse is the wrapper for modules-in-process.json
type InProcessModulesResponse struct {
	Metadata MetadataJSON          `json:"metadata"`
	Modules  []InProcessModuleJSON `json:"modules"`
}
