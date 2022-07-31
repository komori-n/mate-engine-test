package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/komori-n/mate-engine-test/lib/engine"
	"github.com/komori-n/mate-engine-test/lib/test_cases"
	"github.com/schollz/progressbar/v3"
	flag "github.com/spf13/pflag"
)

const (
	readyWaitTime time.Duration = 10 * time.Second
)

func errorExit(message string) {
	fmt.Fprintln(os.Stderr, "error:")
	fmt.Fprintf(os.Stderr, "  %s\n", message)
	fmt.Fprintln(os.Stderr, "")
	flag.Usage()
	os.Exit(1)
}

func main() {
	cpu_num := runtime.NumCPU()
	num_process := flag.IntP("process", "p", cpu_num, "the number of process")
	exit_on_fail := flag.Bool("exit-on-fail", true, "exit immediately after mate failure")
	test_files := flag.StringArrayP("test-files", "f", []string{}, "test specification files")
	engine_path := flag.StringP("engine", "e", "", "engine path")

	flag.Parse()
	if len(*test_files) == 0 {
		errorExit("test files(--test-files, -f) are needed")
	}

	if len(*engine_path) == 0 {
		errorExit("engine(--engine, -e) is needed")
	}

	abs_engine_path, err := filepath.Abs(*engine_path)
	if err != nil {
		message := fmt.Sprintf("failed to initialize engine '%s': %s",
			*engine_path, err.Error())
		errorExit(message)
	}

	fmt.Println("## Parameters")
	fmt.Println("- num_process: ", *num_process)
	fmt.Println("- test_file: ", strings.Join(*test_files, ", "))
	fmt.Printf("- engine: %s(%s)\n", *engine_path, abs_engine_path)

	var engines [](*engine.Engine)
	for i := 0; i < *num_process; i++ {
		en, err := engine.New(abs_engine_path)
		if err != nil {
			message := fmt.Sprintf("failed to start engine: %s", err.Error())
			errorExit(message)
		}
		engines = append(engines, en)
	}

	var fs []string
	for _, test_file_glob := range *test_files {
		glob_files, err := filepath.Glob(test_file_glob)
		if err != nil {
			message := fmt.Sprintf("failed to read '%s': %s",
				test_file_glob, err)
			errorExit(message)
		} else if len(glob_files) == 0 {
			message := fmt.Sprintf("failed to read '%s'", test_file_glob)
			errorExit(message)
		}

		fs = append(fs, glob_files...)
	}

	for _, test_file := range fs {
		txt, err := os.ReadFile(test_file)
		if err != nil {
			message := fmt.Sprintf("failed to read '%s': %s",
				test_file, err.Error())
			errorExit(message)
		}

		ts_map, err := test_cases.Decode(string(txt))
		if err != nil {
			message := fmt.Sprintf("failed to read '%s': %s",
				test_file, err.Error())
			errorExit(message)
		}

		for key, ts := range ts_map {
			err = testMain(engines, key, &ts, *exit_on_fail)
		}
	}
}

func testMain(engines [](*engine.Engine),
	tc_name string,
	ts *test_cases.TestSet,
	exit_on_fail bool) error {
	desc := fmt.Sprintf("[cyan]%s[reset]", tc_name)
	bar := progressbar.NewOptions(len(ts.Tests),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(30),
		progressbar.OptionThrottle(200*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}))

	var wg sync.WaitGroup
	wg.Add(len(engines))

	tc_chan := make(chan test_cases.TestCase)
	for _, en := range engines {
		en.Set(ts.Opts)
		go workerMain(en, &wg, tc_chan, ts.TimeLimit, exit_on_fail)
	}

	for _, tc := range ts.Tests {
		tc_chan <- tc
		bar.Add(1)
	}

	close(tc_chan)
	wg.Wait()
	return nil
}

func workerMain(en *engine.Engine,
	wg *sync.WaitGroup,
	tc_chan chan test_cases.TestCase,
	time_limit int,
	exit_on_fail bool) {
	defer wg.Done()

	read_chan := make(chan error)
	timer := time.NewTicker(readyWaitTime)
	go func() {
		read_chan <- en.Ready()
	}()

	var err error
	select {
	case err = <-read_chan:
		// do nothing

	case <-timer.C:
		err = fmt.Errorf("Ready() didn't finished in limit time")
	}
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	for tc := range tc_chan {
		solve_result := make(chan struct {
			*engine.MateInfo
			error
		})
		timer := time.NewTicker(time.Millisecond * time.Duration(time_limit))
		go func() {
			mi, err2 := en.Solve(tc.Sfen)
			solve_result <- struct {
				*engine.MateInfo
				error
			}{mi, err2}
		}()

		select {
		case res := <-solve_result:
			if res.error != nil {
				err = res.error
			} else {
				if res.MateInfo.Mate != !tc.NoMate {
					err = fmt.Errorf("Expected '%s' but got '%s'",
						mateString(!tc.NoMate), mateString(res.MateInfo.Mate))
				}
			}

		case <-timer.C:
			err = fmt.Errorf("Solve() didn't finished in limit time")
		}

		if err != nil {
			fmt.Fprintln(os.Stderr)
			fmt.Printf("%s: sfen %s\n", err, tc.Sfen)
			if exit_on_fail {
				os.Exit(1)
			}
		}
	}
}

func mateString(mate bool) string {
	if mate {
		return "mate"
	} else {
		return "nomate"
	}
}
