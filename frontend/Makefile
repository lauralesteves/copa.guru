BUCKET := copa.guru

CF_ID := $(shell aws cloudfront list-distributions --query "DistributionList.Items[?Aliases.Items[?contains(@, '$(BUCKET)')]] | [0].Id" --output text)

export AWS_PAGER :=

SYNC_IMMUTABLE = aws s3 sync $(1) s3://$(2) \
	--delete \
	--cache-control "public, max-age=31536000, immutable" \
	--exclude "*.html" \
	--exclude "robots.txt" \
	--exclude "site.webmanifest" \
	--exclude "ai.txt" \
	--exclude "humans.txt"

SYNC_METADATA = aws s3 sync $(1) s3://$(2) \
	--delete \
	--cache-control "public, max-age=3600" \
	--exclude "*" \
	--include "*.html" \
	--include "robots.txt" \
	--include "site.webmanifest" \
	--include "ai.txt" \
	--include "humans.txt"

.PHONY: server build deploy clean

## Dev server
server:
	bun run dev --port 5006

## Build
build:
	bun run build

## Deploy para S3 + invalidacao do CloudFront
deploy: build
	@echo "Deploying to s3://$(BUCKET)..."
	$(call SYNC_IMMUTABLE,dist/,$(BUCKET))
	$(call SYNC_METADATA,dist/,$(BUCKET))
	@echo "Invalidating CloudFront cache..."
	aws cloudfront create-invalidation \
		--distribution-id $(CF_ID) \
		--paths "/*"
	@echo "Deploy complete!"

## Limpar diretorio de build
clean:
	rm -rf dist/
