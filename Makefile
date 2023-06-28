.PHONY: build package

export GO111MODULE=on
export GOARCH=amd64
export GOOS=linux
export GOSUMDB=off
export GOPRIVATE=git.100tal.com

BuildPkg=websocket-pool

build-server:
	go build -o $(BuildPkg) -v
	echo "build success"

run-server:stop-server
	nohup ./$(BuildPkg) > log/nohup.log 2>&1 &
	echo "server start"

stop-server:
	@echo Stopping $(BuildPkg)
	@for PID in $$(ps -ef | grep -v grep | grep ./$(BuildPkg) | awk '{ print $$2 }'); do \
    		echo stopping $(BuildPkg) $$PID; \
    		kill $$PID; \
    done
