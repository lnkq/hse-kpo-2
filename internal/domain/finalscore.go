package domain

import (
	"errors"
	"fmt"
	"math"
)

var ErrNoGrades = errors.New("по студенту нет оценок")

type ScoringPolicy interface {
	Final(grades []*Grade, rubricator *Rubricator) (float64, error)
	Name() string
}

type WeightedSumPolicy struct{}

func NewWeightedSumPolicy() *WeightedSumPolicy { return &WeightedSumPolicy{} }

func (p *WeightedSumPolicy) Name() string { return "взвешенный по весам" }

func (p *WeightedSumPolicy) Final(grades []*Grade, rubricator *Rubricator) (float64, error) {
	if len(grades) == 0 {
		return 0, ErrNoGrades
	}
	var total, weightSum float64
	for _, grade := range grades {
		criterion, ok := rubricator.ByName(grade.Criterion())
		if !ok {
			return 0, fmt.Errorf("%w: %s", ErrUnknownCriterion, grade.Criterion())
		}
		total += grade.Score() / criterion.Max() * criterion.Weight()
		weightSum += criterion.Weight()
	}
	return total / weightSum, nil
}

type AveragePolicy struct{}

func NewAveragePolicy() *AveragePolicy { return &AveragePolicy{} }

func (p *AveragePolicy) Name() string { return "простое среднее" }

func (p *AveragePolicy) Final(grades []*Grade, rubricator *Rubricator) (float64, error) {
	if len(grades) == 0 {
		return 0, ErrNoGrades
	}
	var sum float64
	for _, grade := range grades {
		criterion, ok := rubricator.ByName(grade.Criterion())
		if !ok {
			return 0, fmt.Errorf("%w: %s", ErrUnknownCriterion, grade.Criterion())
		}
		sum += grade.Score() / criterion.Max()
	}
	return sum / float64(len(grades)), nil
}

type RoundedPolicy struct {
	inner  ScoringPolicy
	digits int
}

func NewRoundedPolicy(inner ScoringPolicy, digits int) *RoundedPolicy {
	return &RoundedPolicy{inner: inner, digits: digits}
}

func (p *RoundedPolicy) Name() string { return p.inner.Name() }

func (p *RoundedPolicy) Final(grades []*Grade, rubricator *Rubricator) (float64, error) {
	result, err := p.inner.Final(grades, rubricator)
	if err != nil {
		return 0, err
	}
	factor := math.Pow10(p.digits)
	return math.Round(result*factor) / factor, nil
}
