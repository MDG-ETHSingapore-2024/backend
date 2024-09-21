package utils

import (
    "encoding/base64"
    "fmt"
)

func CreateSafeDbName(userId, dbName string) string {
    encodedUserId := base64.RawURLEncoding.EncodeToString([]byte(userId))
    return fmt.Sprintf("%s_%s", dbName, encodedUserId)
}

func PrepareColumnsAndValues(data map[string]string) ([]string, []string) {
	columns := make([]string, 0, len(data))
	values := make([]string, 0, len(data))

	for k, v := range data {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("'%s'", v))
	}

	return columns, values
}
