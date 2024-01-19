#!/bin/bash

bbook build

GIT_REPO_URL=$(git config --get remote.origin.url)

cd docs/user_docs/.book

git init .
git remote add github $GIT_REPO_URL
git checkout -b gh-pages
git add .
git commit -am "Static site deploy"
git push github gh-pages --force

cd ../..
