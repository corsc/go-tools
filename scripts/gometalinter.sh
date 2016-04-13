#!/bin/bash

clear

if [ -z "$LINTER_EXCLUDE" ]; then
    LINTER_EXCLUDE="/vendor/|_gonerated.go"
fi

gometalinter                            						\
    --disable-all                       						\
	--exclude="$LINTER_EXCLUDE"                         		\
    --enable=vet                                        		\
    --enable=vetshadow                                  		\
    --enable=deadcode                                   		\
    --enable=golint                                     		\
    --enable=varcheck                                   		\
    --enable=structcheck                                		\
    --enable=errcheck                                   		\
    --enable=ineffassign                                		\
    --enable=unconvert                                  		\
    --enable=goconst                                    		\
    --min-occurrences=3                                 		\
    --enable=gofmt                                      		\
    --linter='gofmt:gofmt -l -d -s ./*.go:^(?P<path>[^\n]+)$' 	\
    --enable=goimports                                  		\
    --enable=gocyclo                                    		\
    --cyclo-over=15                                     		\
    --enable=dupl                                       		\
    --enable=lll                                        		\
    --line-length=150                                   		\
    --deadline=1000s                                    		\
    $1
