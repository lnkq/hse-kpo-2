package domain

import (
	"errors"
	"math"
	"testing"
	"time"
)

func setupScores(t *testing.T) (*Rubricator, *GradeFactory) {
	t.Helper()
	rubricator := NewRubricator()
	c1, err := NewCriterion("тесты", 3, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := NewCriterion("код-ревью", 2, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if err := rubricator.Add(c1); err != nil {
		t.Fatal(err)
	}
	if err := rubricator.Add(c2); err != nil {
		t.Fatal(err)
	}
	return rubricator, NewGradeFactory()
}

func TestFinalScorePolicies(t *testing.T) {
	rubricator, factory := setupScores(t)
	g1, _ := factory.Create(rubricator, 1, "тесты", 1, time.Time{}, "")
	g2, _ := factory.Create(rubricator, 1, "код-ревью", 1, time.Time{}, "")
	grades := []*Grade{g1, g2}
	weighted, err := NewWeightedSumPolicy().Final(grades, rubricator)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(weighted-5.0/12.0) > 1e-9 {
		t.Fatalf("неверный взвешенный итог: %v", weighted)
	}
	average, err := NewAveragePolicy().Final(grades, rubricator)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(average-5.0/12.0) > 1e-9 {
		t.Fatalf("неверный средний итог: %v", average)
	}
	if _, err := NewWeightedSumPolicy().Final(nil, rubricator); !errors.Is(err, ErrNoGrades) {
		t.Fatalf("итог без оценок посчитан: %v", err)
	}
}

func TestRoundedPolicyDecorator(t *testing.T) {
	rubricator, factory := setupScores(t)
	grade, _ := factory.Create(rubricator, 1, "тесты", 1, time.Time{}, "")
	rounded, err := NewRoundedPolicy(NewWeightedSumPolicy(), 2).Final([]*Grade{grade}, rubricator)
	if err != nil {
		t.Fatal(err)
	}
	if rounded != 0.33 {
		t.Fatalf("декоратор не округлил итог: %v", rounded)
	}
	if _, err := NewRoundedPolicy(NewWeightedSumPolicy(), 2).Final(nil, rubricator); !errors.Is(err, ErrNoGrades) {
		t.Fatalf("декоратор должен пропускать ошибку стратегии: %v", err)
	}
}
