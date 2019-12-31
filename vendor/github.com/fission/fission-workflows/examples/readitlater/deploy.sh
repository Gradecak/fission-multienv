#!/usr/bin/env bash

# Setup external services
kubectl apply -f redis/redis.yaml

# Setup environments
fission env create --name python3 --version 2 --image fission/python-env --builder fission/python-build-env
fission env create --name binary --image fission/binary-env

# Prepare functions
zip -jr notify-pushbullet.zip notify-pushbullet/
zip -jr ogp.zip ogp/
zip -jr parse-article-body.zip parse-article-body/
zip -jr redis.zip redis/

# Setup functions
fission fn create --env python3 --name notify-pushbullet --src notify-pushbullet.zip --entrypoint "notify.main" --buildcmd "./build.sh"
fission fn create --env python3 --name ogp-extract --src ogp.zip --entrypoint "ogp.extract" --buildcmd "./build.sh"
fission fn create --env python3 --name parse-article-body --src parse-article-body.zip --entrypoint "article.main" --buildcmd "./build.sh"
fission fn create --env binary  --name http --deploy http/http.sh

# Setup redis api functions
fission fn create --env python3 --name redis-list --src redis.zip --entrypoint "user.list" --buildcmd "./build.sh"
fission fn create --env python3 --name redis-append  --src redis.zip --entrypoint "user.append" --buildcmd "./build.sh"
fission fn create --env python3 --name redis-get --src redis.zip --entrypoint "user.get" --buildcmd "./build.sh"
fission fn create --env python3 --name redis-set  --src redis.zip --entrypoint "user.set" --buildcmd "./build.sh"

# Setup workflows
fission fn create --name parse-article --env workflow --src ./parse-article.wf.yaml
fission fn create --name save-article --env workflow --src ./save-article.wf.yaml