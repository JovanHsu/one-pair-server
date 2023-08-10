#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs `make build` command.
# The command compiles all Makefile configs.
# Args:
#   WHAT: Directory names to build.  If any of these directories has a 'main'
#     package, the build will produce executable files under $(OUT_DIR)/bin/platforms OR $(OUT_DIR)/bin—tools/platforms.
#     If not specified, "everything" will be built.
# Usage: `scripts/build_all_server.sh`.
# Example: `hack/build-go.sh WHAT=cmd/kubelet`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# CPU core number
# Check the system type
system_type=$(uname)

if [[ "$system_type" == "Darwin" ]]; then
    # macOS (using sysctl)
    cpu_count=$(sysctl -n hw.ncpu)
elif [[ "$system_type" == "Linux" ]]; then
    # Linux (using lscpu)
    cpu_count=$(lscpu --parse | grep -E '^([^#].*,){3}[^#]' | sort -u | wc -l)
else
    echo "Unsupported operating system: $system_type"
    exit 1
fi
echo -e "${GREEN_PREFIX}======> cpu_count=$cpu_count${COLOR_SUFFIX}"

openim::log::status "Building OpenIM, Parallel compilation compile=$cpu_count"
compile_count=$((cpu_count / 2))

# For help output
ARGHELP=""
if [[ "$#" -gt 0 ]]; then
    ARGHELP="'$*'"
fi

openim::color::echo $COLOR_CYAN "NOTE: $0 has been replaced by 'make multiarch' or 'make build'"
echo
echo "The equivalent of this invocation is: "
echo "    make build ${ARGHELP}"
echo
echo " Example: "
echo "    Print a single binary:"
echo "    make build BINS=openim-api"
echo "    Print : Enable debugging and logging"
echo "    make build BINS=openim-api V=1 DEBUG=1"
echo
make --no-print-directory -C "${OPENIM_ROOT}" -j$compile_count build "$*"

if [ $? -eq 0 ]; then
    openim::log::success "all service build success, run 'make start' or './scripts/start_all.sh'"
else
    openim::log::error "make build Error, script exits"
fi
