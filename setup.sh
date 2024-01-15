#!/bin/bash

## Environment files required by the services
# env_files=(
# 	"alerts"
# 	"assets"
# 	"assist"
# 	"chalice"
# 	"common"
# 	"db"
# 	"ender"
# 	"heuristics"
# 	"http"
# 	"imagestorage"
# 	"integrations"
# 	"peers"
# 	"sink"
# 	"sourcemapreader"
# 	"storage"
# 	"videostorage"
# )
#
# ## Other files required by the services
# other_files=(
# 	"Caddyfile"
# 	"nginx.conf"
# )
#
# for i in "${env_files[@]}"; do
# 	:
# 	if ! [ -f "./$i.env" ]; then
# 		curl fsSL "https://raw.githubusercontent.com/openreplay/openreplay/main/scripts/docker-compose/$i.env" >"$i".env
# 	else
# 		echo "Skipping $i.env; already exists"
# 	fi
# done
#
# for i in "${other_files[@]}"; do
# 	:
# 	if ! [ -f "./$i" ]; then
# 		touch "$i"
# 		curl fsSL "https://raw.githubusercontent.com/openreplay/openreplay/main/scripts/docker-compose/$i" >"$i"
# 	else
# 		echo "Skipping $i; already exists"
# 	fi
# done

git clone https://github.com/openreplay/openreplay.git

cd ./openreplay/scripts/docker-compose || exit &&
	rm install.sh docker-compose.yaml &&
	cp ../../../install.sh . &&
	cp ../../../docker-compose.yml . &&
	chmod +x ./install.sh
# ./install.sh
