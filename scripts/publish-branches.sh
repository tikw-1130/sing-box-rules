#!/usr/bin/env bash

set -euo pipefail

: "${GITHUB_TOKEN:?GITHUB_TOKEN is required}"
: "${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}"

publish_branch() {
  local source_dir="$1"
  local branch="$2"
  local temp_dir

  temp_dir="$(mktemp -d)"
  cp -R "${source_dir}/." "${temp_dir}/"

  (
    cd "${temp_dir}"
    git init -b "${branch}"
    git config user.name "github-actions[bot]"
    git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
    git add --all
    git commit -m "Update ${branch}"
    git remote add origin "https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
    git push --force origin "HEAD:${branch}"
  )

  rm -rf "${temp_dir}"
}

publish_branch "dist/rule-set-geosite" "rule-set-geosite"
publish_branch "dist/rule-set-geoip" "rule-set-geoip"
