#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="${SCRIPT_DIR}/../bridge_ctypes/types"
OUT_DIR="${SCRIPT_DIR}/cffi/proto"

mkdir -p "${OUT_DIR}"

# Resolve protoc-gen-go
if [ -n "${PROTOC_GEN_GO_PATH:-}" ]; then
    export PATH="$(dirname "${PROTOC_GEN_GO_PATH}"):${PATH}"
elif command -v mise &>/dev/null && mise which protoc-gen-go &>/dev/null; then
    export PATH="$(dirname "$(mise which protoc-gen-go)")":${PATH}
fi

# Verify protoc-gen-go is available
if ! command -v protoc-gen-go &>/dev/null; then
    echo "ERROR: protoc-gen-go not found. Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" >&2
    exit 1
fi

protoc \
    --proto_path="${PROTO_DIR}" \
    --go_out="${OUT_DIR}" \
    --go_opt=paths=source_relative \
    "${PROTO_DIR}/baml/cffi/v1/baml_inbound.proto" \
    "${PROTO_DIR}/baml/cffi/v1/baml_outbound.proto"

echo "Generated Go proto files in ${OUT_DIR}/baml/cffi/v1/"
