package domain

import (
	"errors"
	"strings"
)

var ErrStudentNotFound = errors.New("студент не найден")

type StudentID int

type Student struct {
	id   StudentID
	name string
}

func NewStudent(id StudentID, name string) (*Student, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("имя студента обязательно")
	}
	return &Student{id: id, name: strings.TrimSpace(name)}, nil
}

func (s *Student) ID() StudentID { return s.id }
func (s *Student) Name() string  { return s.name }
