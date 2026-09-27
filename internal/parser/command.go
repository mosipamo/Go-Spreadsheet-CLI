package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(line string) (Command, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}
	return Command{
		Name: strings.ToLower(fields[0]),
		Args: fields[1:],
	}, nil
}
