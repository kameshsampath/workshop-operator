package util

//PatchAsString specifies a patch operation, path and value as String .
type PatchAsString struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

//WorkshopLabels attaches the common labels for all workshop controller created resources
func WorkshopLabels() map[string]string {
	return map[string]string{
		"category":   "rhd-workshop",
		"controller": "workshop-controller",
	}
}
