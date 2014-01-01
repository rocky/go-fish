# Whatever it is you want to do, it should be forwarded to the
# to top-level directories
.PHONY: check check-quick test test-quick all

all:
	$(MAKE) -C ..

%:
	$(MAKE) -C ../.. $@
