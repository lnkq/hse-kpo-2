package application

import "pos/internal/domain"

type StudentRepository interface {
	Add(student *domain.Student) error
	ByID(id domain.StudentID) (*domain.Student, bool)
	All() []*domain.Student
}

type RubricatorRepository interface {
	Get() *domain.Rubricator
}

type GradeBookRepository interface {
	Get() *domain.GradeBook
}

type FinalScoreRow struct {
	StudentID domain.StudentID
	Name      string
	Score     float64
	HasScore  bool
}

type FinalReportExporter interface {
	Export(rows []FinalScoreRow) (string, error)
}
