.PHONY: build test e2e lint

TESTABLE = apps/eigendb apps/controller libs/faissgo libs/hnswgo
LINTABLE = apps/eigendb apps/controller libs/faissgo libs/hnswgo

all: build

build:
	$(MAKE) -C apps/eigendb build
	@echo "[*] Compiled binary is in 'apps/eigendb/dist'. Make sure to be in 'apps/eigendb/' when executing it."

run:
	$(MAKE) -C apps/eigendb run

test:
	@for target in $(TESTABLE); do \
		echo "[*] Testing $$target . . ."; \
		$(MAKE) -C $$target -j4 test || exit 1; \
	done

e2e:
	$(MAKE) -C apps/eigendb e2e

lint:
	@for target in $(LINTABLE); do \
		echo "[*] Linting $$target . . ."; \
		$(MAKE) -C $$target -j4 lint || exit 1; \
	done

