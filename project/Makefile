GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m


.PHONY: help
help: ## Показати допомогу
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

.PHONY: dev, stop
dev: ## Запустити в режимі розробки з live reload
	@echo "$(GREEN)Запуск в режимі розробки...$(RESET)"
	docker compose up -d
	air -c .air.toml
stop:
	@echo "$(RESET)Stop docker...$(RESET)"
	docker compose down -v

.PHONY: swagger
swagger:
	swag init -g cmd/app/main.go -d .