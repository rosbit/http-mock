SHELL=/bin/bash

EXES = http-mock

all:
	@for exe in $(EXES); do \
		echo "building $$exe ..."; \
		$(MAKE) -s -f make.inc s=static; \
	done

clean:
	rm -f $(EXES)
