<!--toc:start-->

- [Start the docker registry](#start-the-docker-registry)
- [The Go program](#the-go-program)
- [The (other) shell scripts/files](#the-other-shell-scriptsfiles)
- [Kind of important](#kind-of-important)
<!--toc:end-->

> [!NOTE]
> Everything here is kind of janky, just to get it to work first but in the future, it would be cleaned up to make it significantly easier to run

## Start the docker registry

```bash
docker compose -f docker-compose.registry.yml up
```

This starts a local docker registry and binds the volume to `~/registry-data/` for persistence across restarts.

## The Go program

The Go program essentially exists for the purpose of parsing YAML, we need this to extract the image name from the compose file.

> This part is only necessary to beat the rate-limiting issue and well, fix broken images
>
> I have included a pre-built version of this binary so you would not have to install Go but if you want to, you can just run `go build -o bin .`

This binary takes those services (image name and service name) and uses it to generate a shell script (commands.sh) that will pull the images, re-tag them and then push to the local registry so we don't have to reach out to the remote registry everytime.

The program generates two types; a version you can just pipe into sh (`cat commands.sh | sh`) and a version you can copy and paste individually (useful for seeing what pull failed due to rate-limits); by default it generates the former and you can change that by changing the second argument to `generateFile` to false (there is a makefile to rebuild and run the generate command)

```go
switch cmd {
case "generate":
	generateFile(content, true)
//...                      ^
}

```

After starting the registry and using the Go program to generate the shell script, you can proceed to run that shell script - you would want to add some sort of delay, but if it fails, run the script again, docker will skip the layers you have pulled by default (I intend to automate this in the future)

## The (other) shell scripts/files

I have made a modified version of the docker-compose file (mainly to fix the database issues and point the images at the local registry) and install script (I have commented out a few things like OpenSSL that I already have installed, you may need to uncomment it if it applies to your system).

`setup.sh` clones the repo and copies the modified files for you. You need to go into `openreplay/scripts/docker-compose` and run the `install.sh` script which will run the migrations for you and setup OpenReplay.

> [!NOTE]
> Before setting up Openreplay, you need to setup your domain name to point to the proper DNS server, SSL is required and Caddy will generate it automatically for you as long as it can verify the domain. Also, any change you will make to any of the files should be done before you run any of the shell scripts.

> [!WARNING]
> The `install.sh` script generates random passwords for you, replace the `create_passwords` function's use of randomPass with the values you want. For example;

```sh
function create_passwords() {
	POSTGRES_PASSWORD=$(randomPass)
	info "Creating dynamic passwords..."
	sed -i "s/change_me_domain/${DOMAIN_NAME}/g" common.env
	sed -i "s/change_me_domain/${DOMAIN_NAME}/g" docker-compose.yml
	sed -i "s/change_me_jwt/$(randomPass)/g" common.env
	sed -i "s/change_me_s3_key/$(randomPass)/g" common.env
	sed -i "s/change_me_s3_secret/$(randomPass)/g" common.env
	sed -i "s/change_me_pg_password/$POSTGRES_PASSWORD/g" common.env
	sed -i "s/change_me_pg_password/$POSTGRES_PASSWORD/g" docker-compose.yml
	info "Passwords created and updated in common.env file."
}
```

will become

```sh
function create_passwords() {
	POSTGRES_PASSWORD="mypassword"
	info "Creating dynamic passwords..."
	sed -i "s/change_me_domain/${DOMAIN_NAME}/g" common.env
	sed -i "s/change_me_domain/${DOMAIN_NAME}/g" docker-compose.yml
	sed -i "s/change_me_jwt/super_jwt_secret/g" common.env
	sed -i "s/change_me_s3_key/s3_key_meeh/g" common.env
    # ...
}
```

## Kind of important

If your version of docker somehow doesn't use namespaces, you probably want to uncomment these parts (in the docker-compose.yml file):

```yaml
heuristics-openreplay:
  image: localhost:5000/p1t3u8a3/heuristics:v1.16.0
  # domainname: app.svc.cluster.local
  container_name: heuristics
  networks:
    opereplay-net:
      # aliases:
      #   - heuristics-openreplay.app.svc.cluster.local
```
