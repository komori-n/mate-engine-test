load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "test_cases",
    srcs = ["test_cases.go"],
    importpath = "github.com/komori-n/mate-engine-test/lib/test_cases",
    visibility = ["//visibility:public"],
    deps = ["@in_gopkg_yaml_v3//:yaml_v3"],
)

go_test(
    name = "test_cases_test",
    srcs = ["test_cases_test.go"],
    embed = [":test_cases"],
    deps = ["@com_github_stretchr_testify//assert"],
)
