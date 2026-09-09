MODULE  := github.com/dimkarp93/md-libs
VERSION := $(shell tr -d '[:space:]' < versions.txt)
TAG     := v$(VERSION)
DIST    := dist
PROXY   := $(DIST)/proxy/$(MODULE)/@v
STAGE   := $(DIST)/stage/$(MODULE)@$(TAG)

.PHONY: build test test-v test-run cover check test-all vet fmt pack clean bump-patch bump-minor bump-major

build:
	go build ./...

test:
	go test ./...

test-v:
	go test -v ./...

test-run:
	@test -n "$(T)" || { echo "укажите тест: make test-run T=TestParseInline"; exit 1; }
	go test -v -run '$(T)' ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	@echo "детальный отчёт: go tool cover -html=coverage.out"

check: vet test

test-all: check
	@for d in ../md-docx ../md-pdf; do \
		test -d $$d || continue; \
		echo; echo "==> $$d"; \
		$(MAKE) -C $$d test || exit 1; \
	done

vet:
	go vet ./...

fmt:
	gofmt -l -w .

pack: build test
	rm -rf $(DIST)
	mkdir -p $(PROXY) $(STAGE)
	tar -c --exclude=./$(DIST) --exclude=./.git --exclude=./go.work --exclude=./go.work.sum . | tar -x -C $(STAGE)
	cd $(DIST)/stage && zip -qr ../proxy/$(MODULE)/@v/$(TAG).zip $(MODULE)@$(TAG)
	cp go.mod $(PROXY)/$(TAG).mod
	printf '{"Version":"%s","Time":"%s"}\n' "$(TAG)" "$$(date -u +%Y-%m-%dT%H:%M:%SZ)" > $(PROXY)/$(TAG).info
	printf '%s\n' "$(TAG)" > $(PROXY)/list
	rm -rf $(DIST)/stage
	tar -czf $(DIST)/md-libs-$(VERSION)-proxy.tar.gz -C $(DIST) proxy
	@echo
	@echo "Пакет: $(DIST)/md-libs-$(VERSION)-proxy.tar.gz  ($(MODULE) $(TAG))"
	@echo
	@echo "Как использовать на другой машине (в репозитории, который зависит от md-libs):"
	@echo "  tar -xzf md-libs-$(VERSION)-proxy.tar.gz -C /opt"
	@echo
	@echo "  все прочие зависимости уже в кэше:"
	@echo "    GOFLAGS=-mod=mod GOPROXY=file:///opt/proxy GOSUMDB=off go build ./..."
	@echo
	@echo "  остальное тянется из сети:"
	@echo "    GOFLAGS=-mod=mod GOPROXY=file:///opt/proxy,https://proxy.golang.org,direct \\"
	@echo "    GONOSUMDB='github.com/dimkarp93/*' go build ./..."
	@echo
	@echo "Отключение sumdb обязательно: модуль ещё не в sum.golang.org. GOPRIVATE не подходит —"
	@echo "он выставляет GONOPROXY и уводит go мимо файлового прокси. Если есть go.work — GOWORK=off."

clean:
	rm -rf $(DIST)

bump-patch:
	@v=$$(tr -d '[:space:]' < versions.txt); MAJ=$${v%%.*}; rest=$${v#*.}; MIN=$${rest%%.*}; PAT=$${rest##*.}; printf '%s.%s.%s\n' "$$MAJ" "$$MIN" "$$((PAT + 1))" > versions.txt; cat versions.txt

bump-minor:
	@v=$$(tr -d '[:space:]' < versions.txt); MAJ=$${v%%.*}; rest=$${v#*.}; MIN=$${rest%%.*}; printf '%s.%s.0\n' "$$MAJ" "$$((MIN + 1))" > versions.txt; cat versions.txt

bump-major:
	@v=$$(tr -d '[:space:]' < versions.txt); MAJ=$${v%%.*}; printf '%s.0.0\n' "$$((MAJ + 1))" > versions.txt; cat versions.txt
