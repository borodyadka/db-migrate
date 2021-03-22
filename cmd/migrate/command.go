package main

import (
	"errors"
	"math"
	"strconv"
)

type migrateCommand struct {
	Count int
}

func (*migrateCommand) IsCommand() {}

type createCommand struct {
	Name string
}

func (*createCommand) IsCommand() {}

type testCommand struct{}

func (*testCommand) IsCommand() {}

type commandType interface {
	IsCommand()
}

func parseCommand(cmd []string) (commandType, error) {
	if len(cmd) >= 2 {
		switch cmd[0] {
		case "create":
			return &createCommand{Name: cmd[1]}, nil
		case "up":
			n, err := strconv.Atoi(cmd[1])
			if err != nil {
				return nil, err
			}
			return &migrateCommand{Count: n}, nil
		case "down":
			n, err := strconv.Atoi(cmd[1])
			if err != nil {
				return nil, err
			}
			return &migrateCommand{Count: -n}, nil
		}
	} else if len(cmd) == 1 {
		switch cmd[0] {
		case "up":
			return &migrateCommand{Count: math.MaxInt32}, nil
		case "down":
			return &migrateCommand{Count: -1}, nil
		case "test":
			return &testCommand{}, nil
		}
	}

	return nil, errors.New("unknown command")
}
