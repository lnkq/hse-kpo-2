package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnknownCriterion = errors.New("критерий отсутствует в рубрикаторе")
	ErrScoreOutOfRange  = errors.New("балл вне диапазона критерия")
)

type Grade struct {
	student   StudentID
	criterion string
	score     float64
	date      time.Time
	comment   string
}

func (g *Grade) Student() StudentID { return g.student }
func (g *Grade) Criterion() string  { return g.criterion }
func (g *Grade) Score() float64     { return g.score }
func (g *Grade) Date() time.Time    { return g.date }
func (g *Grade) Comment() string    { return g.comment }

type GradeFactory struct{}

func NewGradeFactory() *GradeFactory { return &GradeFactory{} }

func (f *GradeFactory) Create(rubricator *Rubricator, student StudentID, criterionName string, score float64, date time.Time, comment string) (*Grade, error) {
	criterion, ok := rubricator.ByName(criterionName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownCriterion, criterionName)
	}
	if score < 0 || score > criterion.Max() {
		return nil, fmt.Errorf("%w: %v из %v", ErrScoreOutOfRange, score, criterion.Max())
	}
	return &Grade{student: student, criterion: criterion.Name(), score: score, date: date, comment: comment}, nil
}
