#!/bin/sh
# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` currently does not work
# with multiple packages, see https://code.google.com/p/go/issues/detail?id=6909
#
# Usage: script/coverage [--html|--coveralls]
#
#     --html      Additionally create HTML report and open it in browser
#     --coveralls Push coverage statistics to coveralls.io
#
# Source: https://github.com/mlafeldt/chef-runner/blob/v0.7.0/script/coverage

set -e
set -u  # Exit on undefined variables

# Cleanup on script exit
trap 'rm -rf "$workdir"' EXIT

workdir=.cover
profile="$workdir/cover.out"
mode=count
generate_cover_data() {
    rm -rf "${workdir:?}"
    mkdir -p "$workdir" || { echo "Failed to create $workdir"; exit 1; }

    for pkg in "$@"; do
        f="$workdir/$(echo "$pkg" | tr '/' '-').cover"
        go test -v -covermode="$mode" -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    find "$workdir" -name '*.cover' -exec grep -h -v "^mode:" {} + >>"$profile"
}

show_cover_report() {
    local report_type="$1"
    if [ "$report_type" != "func" ] && [ "$report_type" != "html" ]; then
        echo "Invalid report type: $report_type" >&2
        exit 1
    fi
    go tool cover "-${report_type}=${profile}" || { echo "Failed to generate coverage report"; exit 1; }
}

push_to_coveralls() {
    command -v goveralls >/dev/null 2>&1 || { echo "goveralls is required but not installed" >&2; exit 1; }
    echo "Pushing coverage statistics to coveralls.io"
    goveralls -coverprofile="$profile" || { echo "Failed to push to Coveralls.io" >&2; exit 1; }
}

mapfile -t packages < <(go list ./...)
generate_cover_data "${packages[@]}"

# Improve argument handling
case "$1" in
"")
    show_cover_report func
    ;;
--html)
    show_cover_report html
    ;;
--coveralls)
    push_to_coveralls
    ;;
*)
    echo "Usage: $0 [--html|--coveralls]" >&2
    exit 1
    ;;
esac