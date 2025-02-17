APP_NAME := golang-url-shortener
GO := go
FUNCTIONS := generate redirect delete stats gateway


build:
	${MAKE} -j ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})

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

local:
	${MAKE} -j ${MAKEOPTS} $(foreach function, $(FUNCTIONS), start-${function})

local-%:
	cd internal/adapters/functions/$* && go run main.go
