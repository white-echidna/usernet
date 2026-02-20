_default:
    @just --list

alias r := run
alias b := build

run *args: build
    ./control_plane/bin/usernet {{ args }}

build: _build_cp

[working-directory('control_plane')]
_build_cp:
    go build -o ./bin/usernet ./cmd/usernet
