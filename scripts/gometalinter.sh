#!/bin/bash

clear

gometalinter                            \
    --cyclo-over=10                     \
    --disable-all                       \
    --enable=gofmt                      \
    --enable=goimports                  \
    --enable=dupl                       \
    --enable=golint                     \
    --enable=structcheck                \
    --enable=gocyclo                    \
    --enable=vet                        \
    --enable=errcheck                   \
    --enable=ineffassign                \
    --enable=vetshadow                  \
    --enable=varcheck                   \
    --enable=deadcode                   \
    --deadline=180s                     \
    $1
