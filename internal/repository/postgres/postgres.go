package postgres

import "fmt"

func makeQueryForSlug(segmentsIds []int, startIndex int) ([]interface{}, string) {
	slugsAsInterface := make([]interface{}, len(segmentsIds))
	slugsString := ""
	for i, slug := range segmentsIds {
		i, slug := i, slug

		slugsAsInterface[i] = slug

		if i > 0 {
			slugsString += ", "
		}
		slugsString += fmt.Sprintf("$%d", i+1+startIndex)
	}

	return slugsAsInterface, slugsString
}
