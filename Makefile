# Variable containing the name of the directory we want to clean
CACHE_DIR := ./cache

# List all subdirectories inside 'cache'
SUBDIRS := $(wildcard $(CACHE_DIR)/*)

# Default target
.PHONY: clean
clean:
	@echo "Removing directories in $(CACHE_DIR)"
	@for dir in $(SUBDIRS); do \
		if [ -d "$$dir" ]; then \
			echo "Removing $$dir"; \
			rm -rf "$$dir"; \
		fi; \
	done

