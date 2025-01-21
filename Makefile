dev:
	@go run ./main.go

build.linux:
	@echo "Building for linux..."
	STAGE=prod GOOS=linux GOARCH=amd64 go build -o dist/cbr-linux ./main.go
	@echo "Build complete!"

include .env.deploy
deploy: .env.deploy build.linux
	@echo "Deploying..."
	@scp -rp ./dist/cbr-linux $(DEPLOY_USER)@$(DEPLOY_HOST):$(DEPLOY_PATH)
	@ssh $(DEPLOY_USER)@$(DEPLOY_HOST) "sh $(DEPLOY_PATH)/after_deploy.sh"
	@echo "Deploy complete!"
