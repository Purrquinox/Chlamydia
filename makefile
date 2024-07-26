BUILDTIME := $(shell date '+%Y-%m-%d | %H:%M:%S')
REPONAME := github.com/Purrquinox/Chlamydia
PROJECTNAME := chlamydia

COMBOS := linux/amd64 linux/arm64 windows/amd64 windows/arm64

all:
	CGO_ENABLED=0 go build -v $(GOFLAGS_DBG)
release:
	for combo in $(COMBOS); do \
		echo "$$combo"; \
		mkdir -p bin/$$combo; \
		CGO_ENABLED=0 GOOS=$${combo%/*} GOARCH=$${combo#*/} go build -o bin/$$combo/core $(GOFLAGS); \
		sha512sum bin/$$combo/core > bin/$$combo/core.sha512; \
	done

	for folder in bin/windows/*; do \
		mv -vf $$folder/core $$folder/core.exe; \
	done

	python build.py
fmt:
	go fmt ./...