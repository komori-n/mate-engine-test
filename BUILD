load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/komori-n/mate-engine-test
gazelle(
    name = "gazelle",
    args = [
        "-build_file_name=BUILD",
    ],
)

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)
