package csvreport

import (
	"fmt"
	"os"

	"pos/internal/application"
)

type Exporter struct {
	path string
}

func New(path string) *Exporter {
	return &Exporter{path: path}
}

func (e *Exporter) Export(rows []application.FinalScoreRow) (string, error) {
	file, err := os.Create(e.path)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл %s: %w", e.path, err)
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, "student;score"); err != nil {
		return "", err
	}
	for _, row := range rows {
		if row.HasScore {
			if _, err := fmt.Fprintf(file, "%s;%.2f\n", row.Name, row.Score); err != nil {
				return "", err
			}
		} else {
			if _, err := fmt.Fprintf(file, "%s;\n", row.Name); err != nil {
				return "", err
			}
		}
	}
	return e.path, nil
}
