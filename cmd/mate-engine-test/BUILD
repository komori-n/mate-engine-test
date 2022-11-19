load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "mate-engine-test_lib",
    srcs = ["main.go"],
    importpath = "github.com/komori-n/mate-engine-test/cmd/mate-engine-test",
    visibility = ["//visibility:private"],
    deps = [
        "//lib/engine",
        "//lib/test_cases",
        "@com_github_schollz_progressbar_v3//:progressbar",
        "@com_github_spf13_pflag//:pflag",
    ],
)

go_binary(
    name = "mate-engine-test",
    embed = [":mate-engine-test_lib"],
    visibility = ["//visibility:public"],
)
