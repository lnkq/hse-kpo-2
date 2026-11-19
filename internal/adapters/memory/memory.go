package memory

import (
	"sort"

	"pos/internal/domain"
)

type StudentRepository struct {
	students map[domain.StudentID]*domain.Student
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{students: make(map[domain.StudentID]*domain.Student)}
}

func (r *StudentRepository) Add(student *domain.Student) error {
	r.students[student.ID()] = student
	return nil
}

func (r *StudentRepository) ByID(id domain.StudentID) (*domain.Student, bool) {
	student, ok := r.students[id]
	return student, ok
}

func (r *StudentRepository) All() []*domain.Student {
	result := make([]*domain.Student, 0, len(r.students))
	for _, student := range r.students {
		result = append(result, student)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID() < result[j].ID() })
	return result
}

type RubricatorRepository struct {
	rubricator *domain.Rubricator
}

func NewRubricatorRepository() *RubricatorRepository {
	return &RubricatorRepository{rubricator: domain.NewRubricator()}
}

func (r *RubricatorRepository) Get() *domain.Rubricator {
	return r.rubricator
}

type GradeBookRepository struct {
	gradeBook *domain.GradeBook
}

func NewGradeBookRepository() *GradeBookRepository {
	return &GradeBookRepository{gradeBook: domain.NewGradeBook()}
}

func (r *GradeBookRepository) Get() *domain.GradeBook {
	return r.gradeBook
}
