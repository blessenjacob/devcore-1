APP_REPO := "github.com/kayako/service-location"

APP_VERSION := $$(git rev-list --abbrev-commit --abbrev=8 --max-count=1 HEAD)
APP_BUILD_DATE := $$(date -u +%Y-%m-%dT%H:%M)
APP_BUILD_BRANCH := $$(git rev-parse --abbrev-ref HEAD)

LD_FLAGS := "-X=$(APP_REPO)/geo.AppVersion=$(APP_VERSION)		\
	-X=$(APP_REPO)/geo.BuildTime=$(APP_BUILD_DATE)				\
	-X=$(APP_REPO)/geo.BuildBranch=$(APP_BUILD_BRANCH)"

BUILD_OUTPUT := ./build/location

MML := $(MM_LICENSE)

DB_ISP_URL  := "https://download.maxmind.com/app/geoip_download?edition_id=GeoIP2-ISP&suffix=tar.gz&date=20160517&license_key=$(MML)"
DB_CITY_URL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoIP2-City&suffix=tar.gz&date=20160517&license_key=$(MML)"
DB_CONN_URL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoIP2-Connection-Type&date=20160503&suffix=tar.gz&license_key=$(MML)"

export GO15VENDOREXPERIMENT=1

.PHONY : default
default: test \
	build

.PHONY : build
build: clean \
	geo

.PHONY : geo
geo:
	@echo "Building geo with $(APP_VERSION), $(APP_BUILD_DATE) and $(APP_BUILD_BRANCH), binary will be placed at $(BUILD_OUTPUT)"
	@go build -ldflags $(LD_FLAGS) -o $(BUILD_OUTPUT) main.go

.PHONY : test
test: deps      \
	lint		\
	testcmd     \
	testgeo     \
	testserver  \
	testmain    \

.PHONY : lint
lint:
	@echo "Performing lint operation"
	@golint ./
	@golint ./cmd
	@golint ./geo
	@golint ./server
	@echo "Performing vet operation"
	@go vet ./ ./cmd ./geo ./server

.PHONY : testmain
testmain:
	@echo "Running main test suite"
	@go test -race -v ./


.PHONY : testgeo
testgeo:
	@echo "Running geo test suite"
	@go test -race -v ./geo/

.PHONY : testserver
testserver:
	@echo "Running web test suite"
	@go test -race -v ./server/

.PHONY : testcmd
testcmd:
	@echo "Running commands test suite"
	@go test -race -v ./cmd/

.PHONY : release
release: test	\
	build

.PHONY : deps
deps:
	@echo "Installing required packages [golint]"
	@go get github.com/golang/lint/golint

.PHONY : dbfiles
dbfiles:
	mkdir -p db-files /tmp/dbs \
	&& curl -fsSL $(DB_ISP_URL) -o /tmp/isp.tar.gz \
	&& curl -fsSL $(DB_CITY_URL) -o /tmp/city.tar.gz \
	&& curl -fsSL $(DB_CONN_URL) -o /tmp/conn.tar.gz \
	&& tar -C /tmp/dbs -xzf /tmp/isp.tar.gz && mv /tmp/dbs/*-ISP* /tmp/isp && mv /tmp/isp/*ISP*.mmdb ./db-files/isp.mmdb \
	&& tar -C /tmp/dbs -xzf /tmp/city.tar.gz && mv /tmp/dbs/*-City* /tmp/city && mv /tmp/city/*City*.mmdb ./db-files/city.mmdb \
	&& tar -C /tmp/dbs -xzf /tmp/conn.tar.gz && mv /tmp/dbs/*-Conn* /tmp/conn && mv /tmp/conn/*Conn*.mmdb ./db-files/conn.mmdb

.PHONY : install
install:
	@cp $(BUILD_OUTPUT) /usr/local/bin/location

.PHONY : uninstall
uninstall:
	@rm -f /usr/local/bin/location

.PHONY : clean
clean:
	@echo "Removing build directory"
	-@rm -rf build/ 2> /dev/null
	@echo "Removing existing database files"
	-@rm -rf db-files
	@echo "Removing downloaded database files"
	-@rm -rf /tmp/isp /tmp/city /tmp/conn
