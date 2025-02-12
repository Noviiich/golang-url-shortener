APP_NAME := golang-url-shortener
GO := go
FUNCTIONS := generate redirect gateway


build:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})

build-%:
	cd internal/adapters/functions/$* && ${GO} build -o ${APP_NAME}

unit-test:
	@echo "Запуск тестов..."
	cd internal/tests/unit/$* && ${GO} test -v .

benchmark-test:
	cd internal/tests/benchmark/$* && ${GO} test -v -bench=.

run:
	${MAKE} ${MAKEOPTS} $(foreach function, $(FUNCTIONS), run-${function})

run-%:
	cd internal/adapters/functions/$* && ./${APP_NAME}

tidy:
	go mod tidy
