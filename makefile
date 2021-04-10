  ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif
MAINDIR = src
FRONTDIR = src/gramui

.PHONY: node_modules test tidy

backend:
	cd $(MAINDIR) && go run main.go

dev: node_modules
	cd $(MAINDIR) && go run main.go &
	cd $(FRONTDIR) && npm run serve

install: node_modules

node_modules: $(FRONTDIR)/package.json
	cd $(FRONTDIR)
	npm install

test:
	cd $(MAINDIR) && go test -v ./...

tidy:
	cd $(MAINDIR) && go fmt ./...

# TODO(Jovan): Build npm
build:
	cd $(MAINDIR) && go build -v
