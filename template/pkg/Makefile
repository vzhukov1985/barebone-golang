RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

LAST_VER := $(shell git describe --tags  --match "$(RUN_ARGS)*" --abbrev=0)
LAST_MIN_VER := $(shell git describe --tags  --match "$(RUN_ARGS)*" --abbrev=0 | sed 's/$(RUN_ARGS)\/v1.0.//')
NEW_MIN_VER := $(shell expr $(LAST_MIN_VER) + 1)

v:
	@echo $(LAST_VER)

new:
	$(shell git add .)
	$(shell git commit -m 'upd')
	$(shell git push)
	$(shell git tag '$(RUN_ARGS)/v1.0.$(NEW_MIN_VER)')
	$(shell git push origin '$(RUN_ARGS)/v1.0.$(NEW_MIN_VER)')
