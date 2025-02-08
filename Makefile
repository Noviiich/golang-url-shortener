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

run:
	${MAKE} ${MAKEOPTS} $(foreach function, $(FUNCTIONS), run-${function})

run-%:
	cd internal/adapters/functions/$* && ./${APP_NAME}

tidy:
	go mod tidy
