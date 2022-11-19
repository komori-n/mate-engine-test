load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "engine",
    srcs = ["engine.go"],
    importpath = "github.com/komori-n/mate-engine-test/lib/engine",
    visibility = ["//visibility:public"],
)

go_test(
    name = "engine_test",
    srcs = ["engine_test.go"],
    embed = [":engine"],
    deps = ["@com_github_stretchr_testify//assert"],
)
