package domain

import (
	"errors"
	"math"
	"sort"
	"strings"
)

var ErrCriterionDuplicate = errors.New("критерий с таким названием уже есть в рубрикаторе")

type Criterion struct {
	name   string
	max    float64
	weight float64
}

func NewCriterion(name string, max, weight float64) (Criterion, error) {
	if strings.TrimSpace(name) == "" {
		return Criterion{}, errors.New("название критерия обязательно")
	}
	if math.IsNaN(max) || math.IsInf(max, 0) || max <= 0 {
		return Criterion{}, errors.New("максимальный балл должен быть положительным конечным числом")
	}
	if math.IsNaN(weight) || math.IsInf(weight, 0) || weight <= 0 {
		return Criterion{}, errors.New("вес должен быть положительным конечным числом")
	}
	return Criterion{name: strings.TrimSpace(name), max: max, weight: weight}, nil
}

func (c Criterion) Name() string    { return c.name }
func (c Criterion) Max() float64    { return c.max }
func (c Criterion) Weight() float64 { return c.weight }

type Rubricator struct {
	criteria map[string]Criterion
}

func NewRubricator() *Rubricator {
	return &Rubricator{criteria: make(map[string]Criterion)}
}

func (r *Rubricator) Add(criterion Criterion) error {
	if _, exists := r.criteria[criterion.Name()]; exists {
		return ErrCriterionDuplicate
	}
	r.criteria[criterion.Name()] = criterion
	return nil
}

func (r *Rubricator) ByName(name string) (Criterion, bool) {
	criterion, ok := r.criteria[name]
	return criterion, ok
}

func (r *Rubricator) All() []Criterion {
	result := make([]Criterion, 0, len(r.criteria))
	for _, criterion := range r.criteria {
		result = append(result, criterion)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].name < result[j].name })
	return result
}
