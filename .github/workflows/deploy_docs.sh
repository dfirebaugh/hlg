name: Deploy to GitHub Pages

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '^1.19.x'

    - name: Download latest bbook release
      run: |
        curl -s https://api.github.com/repos/dfirebaugh/bbook/releases/latest \
        | grep "browser_download_url.*tar.gz" \
        | cut -d : -f 2,3 \
        | tr -d \" \
        | wget -qi -
        tar -xzf bbook*.tar.gz

    - name: Run bbook
      run: |
        cd docs/user_docs
        ../bbook build

    - name: Deploy
      uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./docs/user_docs/.book
