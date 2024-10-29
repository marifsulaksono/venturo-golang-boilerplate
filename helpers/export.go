package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"simple-crud-rnd/structs"
	"time"
)

func ExportUsersToCSV(filePath string, users []structs.User, fields []string) error {
	// create CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(fields); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, user := range users {
		var record []string
		v := reflect.ValueOf(user)

		for _, field := range fields {
			f := v.FieldByName(field)
			if !f.IsValid() {
				record = append(record, "")
				continue
			}

			switch f.Kind() {
			case reflect.String:
				record = append(record, f.String())
			case reflect.Int, reflect.Int64:
				record = append(record, fmt.Sprintf("%d", f.Int()))
			case reflect.Struct:
				if t, ok := f.Interface().(time.Time); ok {
					record = append(record, t.Format(time.RFC3339))
				} else {
					record = append(record, fmt.Sprintf("%v", f.Interface()))
				}
			case reflect.Ptr:
				if f.IsNil() {
					record = append(record, "")
				} else {
					record = append(record, fmt.Sprintf("%v", f.Elem()))
				}
			default:
				record = append(record, fmt.Sprintf("%v", f.Interface()))
			}
		}

		// write the csv file
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
