load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push")

# gazelle:prefix github.com/gauntletwizard/targetgroupcontroller
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")

go_image(
    name = "dockerimage",
    srcs = ["main.go"],
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    importpath = "github.com/gauntletwizard/targetgroupcontroller",
    pure = "on",
    visibility = ["//visibility:public"],
)

go_binary(
    name = "targetgroupcontroller",
    embed = [":go_default_library"],
    importpath = "github.com/gauntletwizard/targetgroupcontroller",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "k8sclient.go",
        "main.go",
        "targetGroup.go",
        "watcher.go",
    ],
    importpath = "github.com/gauntletwizard/targetgroupcontroller",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_aws_aws_sdk_go//aws:go_default_library",
        "@com_github_aws_aws_sdk_go//aws/session:go_default_library",
        "@com_github_aws_aws_sdk_go//service/elbv2:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
        "@com_github_prometheus_client_golang//prometheus/promhttp:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/watch:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//tools/clientcmd:go_default_library",
    ],
)

container_push(
    name = "push_dockerimage",
    format = "Docker",
    image = "dockerimage",
    registry = "index.docker.io",
    repository = "gauntletwizard/targetgroupcontroller",
    stamp = True,
    tag = "{BUILD_EMBED_LABEL}",
)
