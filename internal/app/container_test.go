package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"pos/internal/domain"
)

func TestServiceScenario(t *testing.T) {
	container := New(strings.NewReader(""), os.Stdout, filepath.Join(t.TempDir(), "report.csv"))
	service := container.Service

	ivan, err := service.AddStudent("Иван")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddStudent("Мария"); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("тесты", 10, 0.3); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("тесты", 10, 0.3); !errors.Is(err, domain.ErrCriterionDuplicate) {
		t.Fatalf("повторный критерий принят: %v", err)
	}
	if err := service.AddCriterion("код-ревью", 5, 0.7); err != nil {
		t.Fatal(err)
	}

	if _, err := service.AssignGrade(ivan.ID(), "тесты", 12, zeroDate(), ""); !errors.Is(err, domain.ErrScoreOutOfRange) {
		t.Fatalf("балл 12 из 10 принят: %v", err)
	}
	if _, err := service.AssignGrade(ivan.ID(), "шахматы", 5, zeroDate(), ""); !errors.Is(err, domain.ErrUnknownCriterion) {
		t.Fatalf("оценка по критерию вне рубрикатора принята: %v", err)
	}
	if _, err := service.AssignGrade(42, "тесты", 5, zeroDate(), ""); !errors.Is(err, domain.ErrStudentNotFound) {
		t.Fatalf("оценка несуществующему студенту принята: %v", err)
	}

	if _, err := service.AssignGrade(ivan.ID(), "тесты", 7, zeroDate(), ""); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AssignGrade(ivan.ID(), "тесты", 9, zeroDate(), "после пересдачи"); err != nil {
		t.Fatal(err)
	}
	if grades, err := service.GradesOf(ivan.ID()); err != nil || len(grades) != 1 || grades[0].Score() != 9 {
		t.Fatalf("исправление оценки не заменило ее: %v, %v", grades, err)
	}

	rows, err := service.FinalScores()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || !rows[0].HasScore || mathAbs(rows[0].Score-0.9) > 1e-9 || rows[1].HasScore {
		t.Fatalf("неверный итог: %v", rows)
	}

	if err := service.SetPolicy("простое среднее"); err != nil {
		t.Fatal(err)
	}
	rows, err = service.FinalScores()
	if err != nil {
		t.Fatal(err)
	}
	if mathAbs(rows[0].Score-0.9) > 1e-9 {
		t.Fatalf("неверный итог по другой стратегии: %v", rows[0].Score)
	}
	if err := service.SetPolicy("нет такого"); err == nil {
		t.Fatal("неизвестное правило принято")
	}

	if err := service.CancelGrade(ivan.ID(), "тесты"); err != nil {
		t.Fatal(err)
	}
	rows, err = service.FinalScores()
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].HasScore {
		t.Fatalf("итог остался после отмены единственной оценки: %v", rows[0])
	}
	if err := service.CancelGrade(ivan.ID(), "тесты"); !errors.Is(err, domain.ErrGradeNotFound) {
		t.Fatalf("отмена несуществующей оценки прошла: %v", err)
	}
}

func TestExport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.csv")
	container := New(strings.NewReader(""), os.Stdout, path)
	service := container.Service
	if _, err := service.AddStudent("Иван"); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("тесты", 10, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AssignGrade(1, "тесты", 7.5, zeroDate(), ""); err != nil {
		t.Fatal(err)
	}
	got, err := service.ExportFinal()
	if err != nil {
		t.Fatal(err)
	}
	if got != path {
		t.Fatalf("экспорт ушел не в тот файл: %s", got)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "student;score") || !strings.Contains(string(content), "Иван;0.75") {
		t.Fatalf("неверное содержимое отчета: %q", string(content))
	}
}

func TestConsoleScenario(t *testing.T) {
	exportPath := filepath.Join(t.TempDir(), "report.csv")
	input := strings.NewReader(strings.Join([]string{
		"1", "Иван",
		"3", "тесты", "10", "0,3",
		"3", "код-ревью", "5", "0,7",
		"5", "1", "тесты", "12", "", "",
		"5", "1", "тесты", "7", "10.12.2025", "хорошо",
		"7", "1", "тесты", "9", "", "",
		"6", "1",
		"9",
		"10", "2",
		"9",
		"11",
		"0",
	}, "\n"))
	out := &strings.Builder{}
	container := New(input, out, exportPath)
	if err := container.Console.Run(); err != nil {
		t.Fatal(err)
	}
	report := out.String()
	for _, expected := range []string{
		"балл вне диапазона критерия: 12 из 10",
		"оценка исправлена",
		"тесты: 9",
		"правило итога: взвешенный по весам",
		"1 - Иван: итог 0.90",
		"правило итога: простое среднее",
		"итог экспортирован",
	} {
		if !strings.Contains(report, expected) {
			t.Fatalf("в выводе нет %q\nполный вывод:\n%s", expected, report)
		}
	}
	content, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "Иван;0.90") {
		t.Fatalf("в отчете нет итога Ивана: %q", string(content))
	}
}

func mathAbs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func zeroDate() time.Time { return time.Time{} }
