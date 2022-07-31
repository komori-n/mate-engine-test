# mate-engine-test

[![ct](https://github.com/komori-n/mate-engine-test/actions/workflows/ct.yaml/badge.svg?branch=main)](https://github.com/komori-n/mate-engine-test/actions/workflows/ct.yaml)

A simple test application for Shogi mate engines (like [KomoringHeights](https://github.com/komori-n/KomoringHeights)).

## Usage

You can test mate engines like the following command:

```sh
$ mate-engine-test -e KomoringHeights-by-gcc -f tests/tests.yaml
- num_process:  8
- test_file:  ./tests/tests.yaml
- engine: KomoringHeightss-by-gcc(...)
basic 100% |██████████████████████████████| (4/4)
```

Test sets can be specified in `.yaml` with the following format:

```yaml
# Define test set
basic:
  # Time limit for each problem(ms)
  time_limit: 100

  # Engine options. These values are directly passed to the engine
  engine_opts:
    USI_Hash: 64
    PostSearchCount: 0
    RootIsAndNodeIfChecked: false
    PvInterval: 0
    YozumePrintLevel: 0

  # Test cases
  tests:
    - sfen: 4k4/9/4P4/9/9/9/9/9/9 b G2r2b3g4s4n4l17p 1
    # Non-mate positions are also testable
    - sfen: 6R+P1/8k/7sp/9/9/9/9/9/9 b r2b4g3s4n4l16p 1
      nomate: true
```

### Options

* `-e`, `--engine`: Path to a mate engine.
* `-f`, `--test-files`: Path to test sets(see below). It can be used multiple times.
* `-p`, `--process`: Number of processes for testing.
* `--exit-on-fail`: Exit immediately if one test fails.

## Build

```sh
bazel run :gazelle
bazel build cmd/mate-engine-test:mate-engine-test
```

Note that the command `:gazelle` generates `BUILD` for each directory.

You can start testing by

```sh
bazel-bin/cmd/mate-engine-test/mate-engine-test_/mate-engine-test #opts...
```

or

```sh
bazelisk run --run_under="cd $PWD &&" cmd/mate-engine-test:mate-engine-test -- #opts...
```

## Test

```sh
bazelisk test --test_timeout 5 //...
```

## License

MIT License
