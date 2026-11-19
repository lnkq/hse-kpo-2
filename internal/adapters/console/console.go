package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"pos/internal/application"
	"pos/internal/domain"
)

type Console struct {
	in      *bufio.Scanner
	out     io.Writer
	service *application.Service
}

func New(in io.Reader, out io.Writer, service *application.Service) *Console {
	return &Console{in: bufio.NewScanner(in), out: out, service: service}
}

func (c *Console) Run() error {
	fmt.Fprintln(c.out, "пос - преподавательская оценка студентов")
	for {
		fmt.Fprintln(c.out, "\n1 - добавить студента\n2 - студенты\n3 - добавить критерий\n4 - рубрикатор\n5 - выставить оценку\n6 - оценки студента\n7 - исправить оценку\n8 - отменить оценку\n9 - итог\n10 - правило итога\n11 - экспорт итога (csv)\n0 - выход")
		choice, err := c.ask("выбор: ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		switch choice {
		case "1":
			err = c.addStudent()
		case "2":
			c.students()
		case "3":
			err = c.addCriterion()
		case "4":
			c.rubricator()
		case "5":
			err = c.assignGrade("оценка выставлена")
		case "6":
			err = c.gradesOf()
		case "7":
			err = c.assignGrade("оценка исправлена")
		case "8":
			err = c.cancelGrade()
		case "9":
			err = c.finalScores()
		case "10":
			err = c.choosePolicy()
		case "11":
			err = c.exportFinal()
		case "0":
			return nil
		default:
			fmt.Fprintln(c.out, "неизвестная команда")
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (c *Console) addStudent() error {
	name, err := c.ask("имя студента: ")
	if err != nil {
		return err
	}
	student, err := c.service.AddStudent(name)
	if err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintf(c.out, "студент добавлен: %d - %s\n", student.ID(), student.Name())
	return nil
}

func (c *Console) students() {
	students := c.service.Students()
	if len(students) == 0 {
		fmt.Fprintln(c.out, "студентов нет")
		return
	}
	for _, student := range students {
		fmt.Fprintf(c.out, "%d - %s\n", student.ID(), student.Name())
	}
}

func (c *Console) addCriterion() error {
	name, err := c.ask("название критерия: ")
	if err != nil {
		return err
	}
	max, err := c.askFloat("максимальный балл: ")
	if err != nil {
		return err
	}
	weight, err := c.askFloat("вес критерия: ")
	if err != nil {
		return err
	}
	if err := c.service.AddCriterion(name, max, weight); err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintln(c.out, "критерий добавлен в рубрикатор")
	return nil
}

func (c *Console) rubricator() {
	criteria := c.service.Criteria()
	if len(criteria) == 0 {
		fmt.Fprintln(c.out, "рубрикатор пуст")
		return
	}
	for _, criterion := range criteria {
		fmt.Fprintf(c.out, "%s: max %v, вес %v\n", criterion.Name(), criterion.Max(), criterion.Weight())
	}
}

func (c *Console) assignGrade(done string) error {
	id, err := c.askInt("номер студента: ")
	if err != nil {
		return err
	}
	criterion, err := c.ask("критерий: ")
	if err != nil {
		return err
	}
	score, err := c.askFloat("балл: ")
	if err != nil {
		return err
	}
	date, err := c.askDate("дата (ДД.ММ.ГГГГ, пусто - без даты): ")
	if err != nil {
		return err
	}
	comment, err := c.ask("комментарий (пусто - без комментария): ")
	if err != nil {
		return err
	}
	if _, err := c.service.AssignGrade(domain.StudentID(id), criterion, score, date, comment); err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintln(c.out, done)
	return nil
}

func (c *Console) gradesOf() error {
	id, err := c.askInt("номер студента: ")
	if err != nil {
		return err
	}
	grades, err := c.service.GradesOf(domain.StudentID(id))
	if err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	if len(grades) == 0 {
		fmt.Fprintln(c.out, "оценок нет")
		return nil
	}
	for _, grade := range grades {
		line := fmt.Sprintf("%s: %v", grade.Criterion(), grade.Score())
		if !grade.Date().IsZero() {
			line += fmt.Sprintf(" (дата %s)", grade.Date().Format("02.01.2006"))
		}
		if grade.Comment() != "" {
			line += " комментарий: " + grade.Comment()
		}
		fmt.Fprintln(c.out, line)
	}
	return nil
}

func (c *Console) cancelGrade() error {
	id, err := c.askInt("номер студента: ")
	if err != nil {
		return err
	}
	criterion, err := c.ask("критерий: ")
	if err != nil {
		return err
	}
	if err := c.service.CancelGrade(domain.StudentID(id), criterion); err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintln(c.out, "оценка отменена")
	return nil
}

func (c *Console) finalScores() error {
	rows, err := c.service.FinalScores()
	if err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintf(c.out, "правило итога: %s\n", c.service.PolicyName())
	if len(rows) == 0 {
		fmt.Fprintln(c.out, "студентов нет")
		return nil
	}
	for _, row := range rows {
		if row.HasScore {
			fmt.Fprintf(c.out, "%d - %s: итог %.2f\n", row.StudentID, row.Name, row.Score)
		} else {
			fmt.Fprintf(c.out, "%d - %s: нет оценок\n", row.StudentID, row.Name)
		}
	}
	return nil
}

func (c *Console) choosePolicy() error {
	names := c.service.Policies()
	for i, name := range names {
		fmt.Fprintf(c.out, "%d - %s\n", i+1, name)
	}
	choice, err := c.askInt("правило: ")
	if err != nil {
		return err
	}
	if choice < 1 || choice > len(names) {
		fmt.Fprintln(c.out, "нет такого правила")
		return nil
	}
	if err := c.service.SetPolicy(names[choice-1]); err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintf(c.out, "правило итога: %s\n", c.service.PolicyName())
	return nil
}

func (c *Console) exportFinal() error {
	path, err := c.service.ExportFinal()
	if err != nil {
		fmt.Fprintln(c.out, "ошибка:", err)
		return nil
	}
	fmt.Fprintln(c.out, "итог экспортирован:", path)
	return nil
}
