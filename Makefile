.PHONY: run

run:
	@echo "Running server..."
	@{ \
		trap 'echo "Exit" >&2; exit 0;' INT; \
		go run app.go subscription.go validate.go database.go worker.go config.go json_processing.go cache.go http_server.go validation_schema.go tasks.go > app.log 2>&1; \
	} 2>/dev/null
