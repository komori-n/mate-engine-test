package engine

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type MateInfo struct {
	Mate bool
	Len  int
}

type Engine struct {
	cmd     *exec.Cmd
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func New(command string, args ...string) (*Engine, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(stdin)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Engine{cmd, writer, scanner}, nil
}

func (en *Engine) Set(opts map[string]string) error {
	fmt.Fprintf(en.writer, "setoption name Threads value 1\n")
	fmt.Fprintf(en.writer, "setoption name PostSearchLevel value None\n")
	for k, v := range opts {
		fmt.Fprintf(en.writer, "setoption name %s value %s\n", k, v)
	}

	return en.writer.Flush()
}

func (en *Engine) Wait(target string) error {
	for en.scanner.Scan() {
		text := en.scanner.Text()
		if text == target {
			return nil
		}
	}

	if err := en.scanner.Err(); err != nil {
		return err
	}

	return fmt.Errorf("got no '%s'", target)
}

func (en *Engine) Ready() error {
	fmt.Fprintln(en.writer, "isready")
	en.writer.Flush()

	return en.Wait("readyok")
}

func (en *Engine) Solve(sfen string) (*MateInfo, error) {
	fmt.Fprintf(en.writer, "position sfen %s\n", sfen)
	fmt.Fprintln(en.writer, "go mate infinite")
	if err := en.writer.Flush(); err != nil {
		return nil, err
	}

	for en.scanner.Scan() {
		text := en.scanner.Text()
		switch {
		case strings.Contains(text, "nomate"):
			return &MateInfo{Mate: false, Len: 0}, nil

		case strings.Contains(text, "checkmate "):
			s := strings.Fields(text)
			if len(s) <= 1 {
				return nil, fmt.Errorf("got checkmate without mate moves")
			} else if s[1] == "timeout" {
				return nil, fmt.Errorf("timeout")
			} else {
				return &MateInfo{Mate: true, Len: len(s) - 1}, nil
			}

		case strings.Contains(text, "Failed to detect PV"):
			return nil, fmt.Errorf("Failed to detect PV")
		}
	}
	if err := en.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("unexpected EOF")
}
