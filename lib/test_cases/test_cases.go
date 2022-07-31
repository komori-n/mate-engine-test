package test_cases

import (
	"gopkg.in/yaml.v3"
)

type TestCase struct {
	Sfen   string `yaml:"sfen"`
	NoMate bool   `yaml:"nomate"`
}

type TestSet struct {
	TimeLimit int               `yaml:"time_limit"`
	Opts      map[string]string `yaml:"engine_opts"`
	Tests     []TestCase        `yaml:"tests"`
}

func Decode(y string) (map[string]TestSet, error) {
	var ts map[string]TestSet
	err := yaml.Unmarshal([]byte(y), &ts)

	if err != nil {
		return nil, err
	}

	return ts, nil
}
