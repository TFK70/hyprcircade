package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
)

func CreateDaemonLogsCmd(bin string) *cli.Command {
	run := func(context context.Context, cmd *cli.Command) error {
		f, err := os.Open(fmt.Sprintf("/tmp/%s.log", bin)) // path to your colored log file
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		escSeq := regexp.MustCompile(`\\(033|x1b)`)

		for scanner.Scan() {
			line := scanner.Text()

			line = escSeq.ReplaceAllStringFunc(line, func(s string) string {
				return "\x1b"
			})
			line = strings.ReplaceAll(line, `\n`, "\n")
			line = strings.ReplaceAll(line, `\\`, `\`)

			fmt.Println(line)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		return nil
	}

	return &cli.Command{
		Name:   "logs",
		Usage:  "view logs of hyprcircade daemon",
		Action: run,
	}
}
