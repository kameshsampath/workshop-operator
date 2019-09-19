package util

//WorkshopLabels attaches the common labels for all workshop controller created resources
func WorkshopLabels() map[string]string {
	return map[string]string{
		"category":   "rhd-workshop",
		"controller": "workshop-controller",
	}
}
