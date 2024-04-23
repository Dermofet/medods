#!/usr/bin/env bash

# Установите наименование микро-сервиса
APP_NAME=medods
PROTO_PATHS="internal/api/grpc internal/client/grpc"
UNIT_COVERAGE_MIN=33

####
# Далее менять по согласованию
####

GOARCH=
GOOS=

# Запуск bufbuild
run_prototool(){
  if [ $OSTYPE == "msys" ]; then
    MSYS_NO_PATHCONV=1 docker run --rm --platform=linux/x86_64 -v "$(pwd):/work" citilink/prototool:1.11.1 $@
  else
    docker run --rm --platform=linux/x86_64 -v "$(pwd):/work" citilink/prototool:1.11.1 $@
  fi
}

# Обрабатывает proto-файлы prototool
process_proto_files(){
  local COMMAND="$1"
  local PROTO_DIR="$2"

  if [ ! -d "$PROTO_DIR" ]; then
    return 0
  fi

  run_prototool prototool "$COMMAND" "$PROTO_DIR"
}

# Генерация proto-файлов
gen_proto(){
  for CURPATH in ${PROTO_PATHS}; do
    echo "start process $CURPATH..."
    rm -Rf "$CURPATH/gen/*"
    process_proto_files all "$CURPATH"
    if [ -d "$CURPATH/gen" ]; then
      run_prototool chown -R "$(id -u)":"$(id -g)" "/work/$CURPATH/gen"
    fi
    echo "finish process $CURPATH..."
  done
}

# Запуск линтера proto-файлов
lint_proto(){
  echo "run proto linter"
  for CURPATH in ${PROTO_PATHS}; do
    process_proto_files lint "$CURPATH"
  done
}

build_win32() {
  echo "Build for Win32"
  GOARCH=386
  GOOS=windows
  return
}

build_win64() {
  echo "Build for Win64"
  GOARCH=amd64
  GOOS=windows
  return
}

build_linux() {
  echo "Build for Linux"
  GOARCH=amd64
  GOOS=linux
  go tool dist install -v pkg/runtime
  go install -v -a std
  return
}

# Сборка приложения
build() {
  if [ -z "$2" ]; then
    echo "Не выбрана система для компиляции"
    return
  fi

  unit
  local APP_PATH="./cmd/"$APP_NAME

  case $2 in
  "win32")
    build-win32
    GOOS=$GOOS GOARCH=$GOARCH go build -o ./bin/$APP_NAME.exe $APP_PATH

    return
    ;;
  "win64")
    build-win64
    local APP_NAME64=$APP_NAME"64"
    GOOS=$GOOS GOARCH=$GOARCH go build -o ./bin/$APP_NAME64.exe $APP_PATH

    return
    ;;
  "linux")
    build-linux
    GOOS=$GOOS GOARCH=$GOARCH go build -o ./bin/$APP_NAME $APP_PATH

    return
    ;;
  esac

  echo "Неизвестная система для компиляции"
}

# Запуск unit-тестов
unit() {
  deps
  go test ./...
}

# Запуск coverage-тестов
unit_race() {
  deps
  go test -race ./...
}

# тест на покрытие
unit_coverage() {
  echo "run test coverage"
  go test -coverpkg=./... -coverprofile=cover_profile.out.tmp $(go list ./internal/...)
  # удаляем protobuf, validate и моки из тестов покрытия
  < cover_profile.out.tmp grep -v -e "mock" -e "\.pb\.go" -e "\.pb\.validate\.go" > cover_profile.out
  rm cover_profile.out.tmp
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE >= $UNIT_COVERAGE_MIN"
  fi
}

# генерация теста на покрытие в виде html
html_unit_coverage() {
  echo "run test coverage to html"
  go test -coverpkg=./... -coverprofile=cover_profile.out.tmp $(go list ./internal/...)
  # удаляем protobuf, validate и моки из тестов покрытия
  < cover_profile.out.tmp grep -v -e "mock" -e "\.pb\.go" -e "\.pb\.validate\.go" > cover_profile.out
  rm cover_profile.out.tmp
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  go tool cover -html=cover_profile.out -o cover.html
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE >= $UNIT_COVERAGE_MIN"
  fi
}

# Настроить прокси
set_private_repo() {
    git config --global credential.helper store
    echo "https://$GOPROXY_LOGIN:$GOPROXY_TOKEN@codebase.mos.ru" > ~/.git-credentials
	go env -w GOPRIVATE="codebase.mos.ru"
}

# Подтянуть зависимости
deps() {
  go get ./...
}

newmigrate() {
    local MIGRATENAME
    migrate create -ext sql -dir ./cmd/medods/migrations -seq $MIGRATENAME
}

# Добавьте сюда список командx 
using() {
  echo "Укажите команду при запуске: ./run.sh [command]"
  echo "Список команд:"
  echo "-build <СИСТЕМА> (win32/win64/linux) - сборка приложения с необходимой архитектурой"
  echo "-gen_proto - генерация protobuf-файлов"
  echo "-unit - запуск unit-тестов"
  echo "-unit_race - запуск тестов гонки"
  echo "-unit_coverage - запуск тестов покрытия"
  echo "-html_unit_coverage - запуск тестов покрытия с генерацией html файла"
  echo "-deps - загрузка зависимостей"
  echo "-newmigrate - создание миграции"
}

############### НЕ МЕНЯЙТЕ КОД НИЖЕ ЭТОЙ СТРОКИ #################

command="$1"
if [ -z "$command" ]; then
  using
  exit 0
else
  $command $@
fi
