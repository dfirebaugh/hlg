
build-docs:
	bash ./scripts/build_docs.sh

serve-docs:
	bash ./scripts/serve_docs.sh

deploy-docs:
	bash scripts/deploy_docs.sh

clean:
	rm -rf .dist
	rm -rf docs/user_docs/.book
