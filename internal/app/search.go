package app

import (
	"passman/internal/interfaces"
	"strings"
)

func (a *app) search(search string) {
	var filteredRecords []interfaces.Record

	for _, record := range a.records {
		if strings.Contains(record.Name, search) {
			filteredRecords = append(filteredRecords, record)
		}
	}

	a.renderer.Render(filteredRecords)
}
