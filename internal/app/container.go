package app

import (
	"io"

	"pos/internal/adapters/console"
	"pos/internal/adapters/csvreport"
	"pos/internal/adapters/memory"
	"pos/internal/application"
	"pos/internal/domain"
)

type Container struct {
	Service *application.Service
	Console *console.Console
}

func New(in io.Reader, out io.Writer, exportPath string) *Container {
	policies := map[string]domain.ScoringPolicy{
		"взвешенный по весам": domain.NewRoundedPolicy(domain.NewWeightedSumPolicy(), 2),
		"простое среднее":     domain.NewRoundedPolicy(domain.NewAveragePolicy(), 2),
	}
	service := application.NewService(
		memory.NewStudentRepository(),
		memory.NewRubricatorRepository(),
		memory.NewGradeBookRepository(),
		domain.NewGradeFactory(),
		policies,
		"взвешенный по весам",
		csvreport.New(exportPath),
	)
	return &Container{Service: service, Console: console.New(in, out, service)}
}
