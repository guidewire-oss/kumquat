package renderer

import (
	"kumquat/repository"
)

// stripResourcesFromResults removes the Resource objects from the results, replacing them with the underlying content.
func StripResourcesFromResults(results []map[string]repository.Resource) []map[string]any {
	strippedResults := make([]map[string]any, len(results))
	for i, result := range results {
		stripped := make(map[string]any)
		for k, v := range result {
			stripped[k] = v.Content()
		}
		strippedResults[i] = stripped
	}

	return strippedResults
}
