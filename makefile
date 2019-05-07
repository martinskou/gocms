version=0.0.1
date=$(shell date -j "+(%b %Y)")
base=$(shell pwd)
src="$(base)/src"
exec="$(base)/bin/gocms"

.PHONY: all

all:
	@echo " make <cmd>"
	@echo ""
	@echo "commands:"
	@echo " build          - runs go build"
	@echo " build_version  - runs go build with ldflags version=${version} & date=${date}"
	@echo ""

build: clean	
	cd $(src) && go build -v -o ${exec} -ldflags '-X "main.version=${version}" -X "main.date=${date}"'

build_version: check_version
	@go build -v -ldflags '-X "main.version=${version}" -X "main.date=${date}"' -o ${exec}_${version}

clean:
	@rm -f ${exec}

check_version:
	@if [ -a "${exec}_${version}" ]; then \
		echo "${exec}_${version} already exists"; \
		exit 1; \
fi;
