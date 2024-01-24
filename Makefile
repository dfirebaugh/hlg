
build-docs:
	bash ./scripts/build_docs.sh

serve-docs:
	cd docs/user_docs/ && bbook serve

deploy-docs:
	bash scripts/deploy_docs.sh

clean:
	rm -rf .dist
	rm -rf docs/user_docs/.book
