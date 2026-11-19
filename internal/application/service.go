package application

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"pos/internal/domain"
)

type Service struct {
	students     StudentRepository
	rubricator   RubricatorRepository
	gradeBook    GradeBookRepository
	gradeFactory *domain.GradeFactory
	policies     map[string]domain.ScoringPolicy
	policy       string
	exporter     FinalReportExporter
	nextStudent  domain.StudentID
}

func NewService(students StudentRepository, rubricator RubricatorRepository, gradeBook GradeBookRepository, gradeFactory *domain.GradeFactory, policies map[string]domain.ScoringPolicy, defaultPolicy string, exporter FinalReportExporter) *Service {
	return &Service{
		students:     students,
		rubricator:   rubricator,
		gradeBook:    gradeBook,
		gradeFactory: gradeFactory,
		policies:     policies,
		policy:       defaultPolicy,
		exporter:     exporter,
		nextStudent:  1,
	}
}

func (s *Service) AddStudent(name string) (*domain.Student, error) {
	student, err := domain.NewStudent(s.nextStudent, name)
	if err != nil {
		return nil, err
	}
	if err := s.students.Add(student); err != nil {
		return nil, err
	}
	s.nextStudent++
	return student, nil
}

func (s *Service) Students() []*domain.Student {
	return s.students.All()
}

func (s *Service) AddCriterion(name string, max, weight float64) error {
	criterion, err := domain.NewCriterion(name, max, weight)
	if err != nil {
		return err
	}
	return s.rubricator.Get().Add(criterion)
}

func (s *Service) Criteria() []domain.Criterion {
	return s.rubricator.Get().All()
}

func (s *Service) AssignGrade(studentID domain.StudentID, criterionName string, score float64, date time.Time, comment string) (*domain.Grade, error) {
	if _, ok := s.students.ByID(studentID); !ok {
		return nil, domain.ErrStudentNotFound
	}
	grade, err := s.gradeFactory.Create(s.rubricator.Get(), studentID, criterionName, score, date, comment)
	if err != nil {
		return nil, err
	}
	s.gradeBook.Get().Put(grade)
	return grade, nil
}

func (s *Service) GradesOf(studentID domain.StudentID) ([]*domain.Grade, error) {
	if _, ok := s.students.ByID(studentID); !ok {
		return nil, domain.ErrStudentNotFound
	}
	return s.gradeBook.Get().ForStudent(studentID), nil
}

func (s *Service) CancelGrade(studentID domain.StudentID, criterionName string) error {
	if _, ok := s.students.ByID(studentID); !ok {
		return domain.ErrStudentNotFound
	}
	return s.gradeBook.Get().Cancel(studentID, criterionName)
}

func (s *Service) Policies() []string {
	names := make([]string, 0, len(s.policies))
	for name := range s.policies {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (s *Service) PolicyName() string {
	return s.policy
}

func (s *Service) SetPolicy(name string) error {
	if _, exists := s.policies[name]; !exists {
		return fmt.Errorf("неизвестное правило итога: %s", name)
	}
	s.policy = name
	return nil
}

func (s *Service) currentPolicy() (domain.ScoringPolicy, error) {
	policy, exists := s.policies[s.policy]
	if !exists {
		return nil, errors.New("правило итога не настроено")
	}
	return policy, nil
}

func (s *Service) FinalScores() ([]FinalScoreRow, error) {
	policy, err := s.currentPolicy()
	if err != nil {
		return nil, err
	}
	students := s.students.All()
	rows := make([]FinalScoreRow, 0, len(students))
	for _, student := range students {
		grades := s.gradeBook.Get().ForStudent(student.ID())
		score, err := policy.Final(grades, s.rubricator.Get())
		row := FinalScoreRow{StudentID: student.ID(), Name: student.Name()}
		if errors.Is(err, domain.ErrNoGrades) {
			row.HasScore = false
		} else if err != nil {
			return nil, err
		} else {
			row.HasScore = true
			row.Score = score
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (s *Service) ExportFinal() (string, error) {
	rows, err := s.FinalScores()
	if err != nil {
		return "", err
	}
	return s.exporter.Export(rows)
}
