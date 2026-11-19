package console

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func (c *Console) ask(prompt string) (string, error) {
	fmt.Fprint(c.out, prompt)
	if !c.in.Scan() {
		if err := c.in.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(c.in.Text()), nil
}

func (c *Console) askInt(prompt string) (int, error) {
	for {
		value, err := c.ask(prompt)
		if err != nil {
			return 0, err
		}
		result, err := strconv.Atoi(value)
		if err == nil {
			return result, nil
		}
		fmt.Fprintln(c.out, "введите целое число")
	}
}

func (c *Console) askFloat(prompt string) (float64, error) {
	for {
		value, err := c.ask(prompt)
		if err != nil {
			return 0, err
		}
		result, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
		if err == nil {
			return result, nil
		}
		fmt.Fprintln(c.out, "введите число")
	}
}

func (c *Console) askDate(prompt string) (time.Time, error) {
	for {
		value, err := c.ask(prompt)
		if err != nil {
			return time.Time{}, err
		}
		if value == "" {
			return time.Time{}, nil
		}
		for _, layout := range []string{"02.01.2006", "2006-01-02"} {
			result, err := time.Parse(layout, value)
			if err == nil {
				return result, nil
			}
		}
		fmt.Fprintln(c.out, "введите дату в формате ДД.ММ.ГГГГ или оставьте пустым")
	}
}
