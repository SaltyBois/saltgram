  ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif
MAINDIR = src
FRONTDIR = src/frontend
WEBSERVERDIR = src/webserver
APIDIR = src/api
AUTHDIR = src/auth
USERSDIR = src/users
EMAILDIR = src/email
CONTENTDIR = src/content

.PHONY: node_modules test tidy protos api auth users content

protos:
	cd $(MAINDIR) && protoc -I protos/ protos/*/*.proto --go_out=protos/ --go-grpc_out=protos/

backend:
	make -j 5 api auth users email content

api:
	cd $(APIDIR) && go run main.go

auth:
	cd $(AUTHDIR) && go run main.go

users:
	cd $(USERSDIR) && go run main.go

email:
	cd $(EMAILDIR) && go run main.go

content:
	cd $(CONTENTDIR) && go run main.go

frontend:
	cd $(FRONTDIR) && npm run build
	cd $(WEBSERVERDIR) && go run main.go

front_build: 
	cd $(FRONTDIR) && npm run build

# NOTE(Jovan): Should be deprecated?
dev: node_modules
	cd $(MAINDIR) && go run main.
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

kill:
	npx kill-port $(SALT_API_PORT) &
	npx kill-port $(SALT_AUTH_PORT) &
	npx kill-port $(SALT_USERS_PORT) &
	npx kill-port $(SALT_EMAIL_PORT) &
	exit
	

