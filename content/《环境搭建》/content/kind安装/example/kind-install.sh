#!/bin/sh
set -euo pipefail

default_controlplanes=1
default_workers=1
default_cluster_name=""
default_image=""

controlplanes="${1:-${CONTROLPLANES:=${default_controlplanes}}}"
workers="${2:-${WORKERS:=${default_workers}}}"
cluster_name="${3:-${CLUSTER_NAME:=${default_cluster_name}}}"
# IMAGE controls the K8s version as well (e.g. kindest/node:v1.11.10)
image="${4:-${IMAGE:=${default_image}}}"

have_kind() {
  [[ -n "$(command -v kind)" ]]
}

if ! have_kind; then
  echo "installing kind"
  [ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
  chmod +x kind
  sudo mv kind /usr/local/bin
else
  echo "kind detected with version: $(kind --version)"
fi

have_kubectl() {
    [[ -n "$(command -v kubectl)" ]]
}

if ! have_kubectl; then
    echo "Please install kubectl first:"
    echo "  https://kubernetes.io/docs/tasks/tools/#kubectl"
    exit 1
fi

#拼接kind命令
kind_cmd="kind create cluster"
if [[ -n "${cluster_name}" ]]; then
  kind_cmd+=" --name ${cluster_name}"
fi

if [[ -n "${image}" ]]; then
  kind_cmd+=" --image ${image}"
fi


