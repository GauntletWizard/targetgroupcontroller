load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

## Docker rules
# https://github.com/bazelbuild/rules_docker

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    # Head as of 5/10/2018
    commit = "e5ebe3d241775a220e37aceaa24c3e78700a4e0b",
)

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
    container_repositories = "repositories",
)

## Go rules
# https://github.com/bazelbuild/rules_go#generating-build-files

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.12.1/rules_go-0.12.1.tar.gz"],
    sha256 = "8b68d0630d63d95dacc0016c3bb4b76154fe34fca93efd65d1c366de3fcb4294",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

# Gazelle
# gazelle:prefix github.com/gauntletwizard/targetgroupcontroller
http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.12.0/bazel-gazelle-0.12.0.tar.gz"],
    sha256 = "ddedc7aaeb61f2654d7d7d4fd7940052ea992ccdb031b8f9797ed143ac7e8d43",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# Direct dependencies:
go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "d6a9817c4afc94d51115e4a30d449056a3fbf547",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "io_k8s_client_go",
    commit = "0ec73abb067faf46eeb80234f4df6e8413ef3e09",
    remote = "https://github.com/kubernetes/client-go.git",
    vcs = "git",
    importpath = "k8s.io/client-go",
)

go_repository(
    name = "io_k8s_api",
    commit = "8be2a0b24ed0dac9cfc1ac2d987ea16cfcdbecb6",
    remote = "https://github.com/kubernetes/api.git",
    vcs = "git",
    importpath = "k8s.io/api",
    # This magic comes from https://github.com/bazelbuild/rules_go/issues/1022
    build_file_proto_mode = "disable",
)

go_repository(
    name = "io_k8s_apimachinery",
    commit = "bce280dade67588308614434551c0b865433e44d",
    remote = "https://github.com/kubernetes/apimachinery.git",
    vcs = "git",
    importpath = "k8s.io/apimachinery",
    # This magic comes from https://github.com/bazelbuild/rules_go/issues/1022
    build_file_proto_mode = "disable",
)

go_repository(
    name = "com_github_aws_aws_sdk_go",
    commit = "bcf2dfef3f28c8ac3d245e6b0f684bbf34d55a05",
    importpath = "github.com/aws/aws-sdk-go",
)

go_repository(
    name = "com_github_stretchr_testify",
    commit = "f35b8ab0b5a2cef36673838d662e249dd9c94686",
    importpath = "github.com/stretchr/testify",
)

# Dependencies required for prometheus
go_repository(
    name = "com_github_beorn7_perks",
    commit = "3a771d992973f24aa725d07868b467d1ddfceafb",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "7600349dcfe1abd18d72d3a1770870d9800a7801",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "ae68e2d4c00fed4943b5f6698d504a5fe083da8a",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "99fa1f4be8e564e8a6b613da7fa6f46c9edafc6c",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "24818f796faf91cd76ec7bddd72458fbced7a6c1",
    importpath = "github.com/google/gofuzz",
)
