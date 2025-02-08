APP_NAME := golang-url-shortener
GO := go
FUNCTIONS := generate redirect gateway

build:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})

build-%:
	cd internal/adapters/functions/$* && ${GO} build -o ${APP_NAME}

test:
	@echo "Запуск тестов..."
	go test ./... -cover
	

tidy:
	go mod tidy
