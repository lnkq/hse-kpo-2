package domain

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestCriterionValidation(t *testing.T) {
	if _, err := NewCriterion("", 10, 0.5); err == nil {
		t.Fatal("пустое название принято")
	}
	if _, err := NewCriterion("тесты", 0, 0.5); err == nil {
		t.Fatal("нулевой максимум принят")
	}
	if _, err := NewCriterion("тесты", math.NaN(), 0.5); err == nil {
		t.Fatal("nan принят как максимум")
	}
	if _, err := NewCriterion("тесты", 10, -1); err == nil {
		t.Fatal("отрицательный вес принят")
	}
	criterion, err := NewCriterion(" тесты ", 10, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if criterion.Name() != "тесты" || criterion.Max() != 10 || criterion.Weight() != 0.5 {
		t.Fatalf("неверный критерий: %v", criterion)
	}
}

func TestRubricatorRejectsDuplicate(t *testing.T) {
	rubricator := NewRubricator()
	criterion, _ := NewCriterion("тесты", 10, 0.5)
	if err := rubricator.Add(criterion); err != nil {
		t.Fatal(err)
	}
	duplicate, _ := NewCriterion("тесты", 5, 0.7)
	if err := rubricator.Add(duplicate); !errors.Is(err, ErrCriterionDuplicate) {
		t.Fatalf("ожидался отказ для повторного критерия, получено %v", err)
	}
	if _, ok := rubricator.ByName("тесты"); !ok {
		t.Fatal("критерий не найден в рубрикаторе")
	}
	if _, ok := rubricator.ByName("код-ревью"); ok {
		t.Fatal("найден несуществующий критерий")
	}
}

func TestGradeFactoryInvariants(t *testing.T) {
	rubricator := NewRubricator()
	criterion, _ := NewCriterion("тесты", 10, 0.5)
	if err := rubricator.Add(criterion); err != nil {
		t.Fatal(err)
	}
	factory := NewGradeFactory()
	if _, err := factory.Create(rubricator, 1, "код-ревью", 5, time.Time{}, ""); !errors.Is(err, ErrUnknownCriterion) {
		t.Fatalf("оценка по неизвестному критерию принята: %v", err)
	}
	if _, err := factory.Create(rubricator, 1, "тесты", -1, time.Time{}, ""); !errors.Is(err, ErrScoreOutOfRange) {
		t.Fatalf("отрицательный балл принят: %v", err)
	}
	if _, err := factory.Create(rubricator, 1, "тесты", 12, time.Time{}, ""); !errors.Is(err, ErrScoreOutOfRange) {
		t.Fatalf("балл 12 из 10 принят: %v", err)
	}
	date := time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC)
	grade, err := factory.Create(rubricator, 1, "тесты", 10, date, "отлично")
	if err != nil {
		t.Fatal(err)
	}
	if grade.Score() != 10 || grade.Criterion() != "тесты" || grade.Comment() != "отлично" || !grade.Date().Equal(date) {
		t.Fatalf("неверная оценка: %v", grade)
	}
	if _, err := factory.Create(rubricator, 1, "тесты", 0, time.Time{}, ""); err != nil {
		t.Fatalf("граничный балл 0 отклонен: %v", err)
	}
}

func TestGradeBookReplaceAndCancel(t *testing.T) {
	rubricator := NewRubricator()
	criterion, _ := NewCriterion("тесты", 10, 0.5)
	if err := rubricator.Add(criterion); err != nil {
		t.Fatal(err)
	}
	factory := NewGradeFactory()
	book := NewGradeBook()
	first, _ := factory.Create(rubricator, 1, "тесты", 3, time.Time{}, "")
	book.Put(first)
	second, _ := factory.Create(rubricator, 1, "тесты", 8, time.Time{}, "")
	book.Put(second)
	grades := book.ForStudent(1)
	if len(grades) != 1 || grades[0].Score() != 8 {
		t.Fatalf("исправление оценки должно заменять ее, получено %v", grades)
	}
	if err := book.Cancel(1, "тесты"); err != nil {
		t.Fatal(err)
	}
	if grades := book.ForStudent(1); len(grades) != 0 {
		t.Fatalf("отмененная оценка осталась: %v", grades)
	}
	if err := book.Cancel(1, "тесты"); !errors.Is(err, ErrGradeNotFound) {
		t.Fatalf("отмена несуществующей оценки прошла: %v", err)
	}
}
