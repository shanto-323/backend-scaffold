build: confirm
	echo "Running [go mod tidy]"
	@go mod tidy; 
	echo "Running [go mod vendor]"
	@go mod vendor; 
	echo "Running [docker-compose build --no-cache]"
	docker-compose build --no-cache
	echo "Build successful"

run: 
	echo "Running [docker-compose up -d $(service)]"
	docker-compose up -d $(service)

stop: confirm 
	echo "Running [docker-compose down -v]"
	docker-compose down -v

re_build: confirm
	@if [ -z "$(service)" ]; then \
		echo "Invalid service name"; \
		exit 1; \
	fi
	@echo "Rebuilding service: $(service)"
	@docker-compose rm -sf $(service)
	@docker-compose build --no-cache $(service)
	@make run service=$(service)

confirm:
	@read -p "Press [y/Y] to continue: " value; \
	value=$$(echo $$value | tr '[:upper:]' '[:lower:]'); \
	if [ "$$value" = "y" ]; then \
		echo "Running target..."; \
	else \
		exit 1; \
	fi

