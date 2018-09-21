load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push")

go_image(
    name = "dockerimage",
    srcs = ["main.go"],
    embed = [":go_default_library"],
    importpath = "github.com/lifeonairteam/houseparty-infrastructure/docker/targetgroupcontroller",
    visibility = ["//visibility:public"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
)

go_binary(
    name = "reaper",
    srcs = ["main.go"],
    embed = [":go_default_library"],
    importpath = "github.com/lifeonairteam/houseparty-infrastructure/docker/targetgroupcontroller",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = [
    ],
    importpath = "github.com/gauntletwizard/targetgroupcontroller",
    visibility = ["//visibility:private"],
    deps = [
        "//k8sclient:go_default_library",
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
    image = "dockerimage",
    format = "Docker",
    registry = "quay.io/lifeonair",
    repository = "targetgroupcontroller",
    tag = "{BUILD_EMBED_LABEL}",
    stamp = True,
    )
