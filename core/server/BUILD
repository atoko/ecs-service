load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix goland/server
# gazelle:resolve go goland/protocol/gen/go //protocol/src:go
# gazelle:resolve go goland/protocol/gen/go/command //protocol/src/command:command
# gazelle:resolve go goland/protocol/gen/go/session //protocol/src/session
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=server/src/go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)
