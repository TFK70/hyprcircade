package config

import (
	"fmt"
	"strconv"

	hp "github.com/anotherhadi/hyprlang-parser"
)

type General struct {
	Anchor  string
	DarkAt  int
	LightAt int
}

type File struct {
	Path         string
	DayValue     string
	NightValue   string
	IgnoreAnchor bool
}

type Command struct {
	DayExec   string
	NightExec string
}

type Config struct {
	General  *General
	Files    []*File
	Commands []*Command
}

func NewConfig(cfgPath string) (*Config, error) {
	cfg, err := parseConfig(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	return cfg, nil
}

func parseConfig(cfgPath string) (*Config, error) {
	c, err := hp.LoadConfig(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	anchor := c.GetFirst("general", "anchor")

	darkAt, err := strconv.Atoi(c.GetFirst("general", "dark-at"))
	if err != nil {
		return &Config{}, err
	}
	lightAt, err := strconv.Atoi(c.GetFirst("general", "light-at"))
	if err != nil {
		return &Config{}, err
	}

	filePaths := c.GetAll("file", "path")
	fileDayValues := c.GetAll("file", "day-value")
	fileNightValues := c.GetAll("file", "night-value")
	fileIgnoreAnchors := c.GetAll("file", "ignore-anchor")

	if len(fileIgnoreAnchors) < len(filePaths) {
		return &Config{}, fmt.Errorf("ignoreAnchor property is missing in %d files", len(filePaths)-len(fileIgnoreAnchors))
	}

	files := []*File{}

	for idx, path := range filePaths {
		ignoreAnchor, err := strconv.ParseBool(fileIgnoreAnchors[idx])
		if err != nil {
			return &Config{}, err
		}

		file := &File{
			Path:         path,
			DayValue:     fileDayValues[idx],
			NightValue:   fileNightValues[idx],
			IgnoreAnchor: ignoreAnchor,
		}

		files = append(files, file)
	}

	dayCommandExecs := c.GetAll("command", "day-exec")
	nightCommandExecs := c.GetAll("command", "night-exec")

	commands := []*Command{}

	for _, exec := range dayCommandExecs {
		command := &Command{
			DayExec: exec,
		}

		commands = append(commands, command)
	}

	for _, exec := range nightCommandExecs {
		command := &Command{
			NightExec: exec,
		}

		commands = append(commands, command)
	}

	cfg := &Config{
		General: &General{
			Anchor:  anchor,
			DarkAt:  darkAt,
			LightAt: lightAt,
		},
		Files:    files,
		Commands: commands,
	}

	return cfg, nil
}
