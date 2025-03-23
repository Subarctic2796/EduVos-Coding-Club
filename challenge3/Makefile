# Simple Makefile for a Go project

TWCSS_IN = cmd/web/styles/input.css
TWCSS_OUT = cmd/web/assets/css/output.css

# Build the application
all: build
templ-install:
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if [ ! -x "$$(command -v templ)" ]; then \
				echo "templ installation failed. Exiting..."; \
				exit 1; \
			fi; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi

tailwind-install:
	@if [ ! -f tailwindcss ]; then curl -L https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss; fi
	@chmod +x tailwindcss

build: tailwind-install templ-install
	@echo "Building..."
	@echo "templ..."
	@templ generate --lazy
	@echo "tailwind..."
	@./tailwindcss -i $(TWCSS_IN) -o $(TWCSS_OUT)
	@echo "go..."
	@go build -o main cmd/api/main.go

build-prod: tailwind-install templ-install
	@echo "Building for prod..."
	@templ generate
	@./tailwindcss -i $(TWCSS_IN) -o $(TWCSS_OUT) --minify
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run clean watch tailwind-install templ-install build-prod
