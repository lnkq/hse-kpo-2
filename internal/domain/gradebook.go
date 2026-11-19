package domain

import (
	"errors"
	"fmt"
	"sort"
)

var ErrGradeNotFound = errors.New("оценка не найдена")

type gradeKey struct {
	student   StudentID
	criterion string
}

type GradeBook struct {
	grades map[gradeKey]*Grade
}

func NewGradeBook() *GradeBook {
	return &GradeBook{grades: make(map[gradeKey]*Grade)}
}

func (b *GradeBook) Put(grade *Grade) {
	b.grades[gradeKey{grade.student, grade.criterion}] = grade
}

func (b *GradeBook) Cancel(student StudentID, criterionName string) error {
	key := gradeKey{student, criterionName}
	if _, exists := b.grades[key]; !exists {
		return fmt.Errorf("%w: %s", ErrGradeNotFound, criterionName)
	}
	delete(b.grades, key)
	return nil
}

func (b *GradeBook) ForStudent(student StudentID) []*Grade {
	var result []*Grade
	for _, grade := range b.grades {
		if grade.student == student {
			result = append(result, grade)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].criterion < result[j].criterion })
	return result
}
