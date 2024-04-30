#!/usr/bin/env bash

set -x  # Debug: Print commands and their arguments as they are executed.

# Exiting and error handling
set -o errexit
set -o nounset
set -o pipefail

# Print environment debug info
echo "Debugging environment paths and variables"
echo "PATH: $PATH"
GOPATH=$(go env GOPATH)
echo "GOPATH: $GOPATH"
GOBIN="${GOBIN:-$GOPATH/bin}"  # Set GOBIN to $GOPATH/bin if not already set
echo "GOBIN: $GOBIN"

# Script content continues...

# Copyright 2017 The Kubernetes Authors.
# Licensed under the Apache License, Version 2.0

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[@]}")/..
echo "SCRIPT_ROOT: $SCRIPT_ROOT"

TOOLS_DIR=$(realpath ./hack/tools)
TOOLS_BIN_DIR="${TOOLS_DIR}/bin"
GO_INSTALL=$(realpath ./hack/go-install.sh)
CONTROLLER_GEN_VER=v0.11.1
CONTROLLER_GEN_BIN=controller-gen
CONTROLLER_GEN=${TOOLS_BIN_DIR}/${CONTROLLER_GEN_BIN}-${CONTROLLER_GEN_VER}
CRD_OPTIONS="crd:crdVersions=v1"

echo "Installing controller-gen at ${CONTROLLER_GEN}"
GOBIN=${TOOLS_BIN_DIR} ${GO_INSTALL} sigs.k8s.io/controller-tools/cmd/controller-gen ${CONTROLLER_GEN_BIN} ${CONTROLLER_GEN_VER}

CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
echo "CODEGEN_PKG: $CODEGEN_PKG"

echo "Running code generation for config API group"
bash "${CODEGEN_PKG}/kube_codegen.sh" \
  deepcopy,conversion,defaulter \
  sigs.k8s.io/scheduler-plugins/pkg/generated \
  sigs.k8s.io/scheduler-plugins/apis \
  sigs.k8s.io/scheduler-plugins/apis \
  config:v1 \
  --output-base "${SCRIPT_ROOT}" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate/boilerplate.generatego.txt"

echo "Running code generation for scheduling API group"
bash "${CODEGEN_PKG}/kube_codegen.sh" \
  deepcopy,client,informer,lister \
  sigs.k8s.io/scheduler-plugins/pkg/generated \
  sigs.k8s.io/scheduler-plugins/apis \
  scheduling:v1alpha1 \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate/boilerplate.generatego.txt"

echo "Generating controller objects"
${CONTROLLER_GEN} object:headerFile="${SCRIPT_ROOT}/hack/boilerplate/boilerplate.generatego.txt" \
  paths="./apis/scheduling/..."

echo "Generating CRDs, RBAC, and Webhook configurations"
${CONTROLLER_GEN} ${CRD_OPTIONS} rbac:roleName=manager webhook \
  paths="./apis/scheduling/..." \
  output:crd:artifacts:config=config/crd/bases
