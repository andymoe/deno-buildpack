## Basic Use

1. Checkout this buildpack repo to a `deno-buildpack` directory.
1. run `scripts/build.sh`
1. Create a `deno-sample-app` directory with at least a `main.ts` file
containing [simple deno program](https://deno.land/manual/examples/http_server)
1. Run `pack set-default-builder paketobuildpacks/builder:base`
1. Run `pack build test-deno-app --path ./deno-sample-app --buildpack ./deno-buildpack`
1. Run `docker run --rm --name deno-test -p 8080:8080 test-deno-app`
1. Finally, you can destroy the running container with `docker rm -f deno-test`

## Prerequisite

- [Pack](https://buildpacks.io/docs/tools/pack/cli/install/)
- [Docker](https://docs.docker.com/get-docker/)
- Some JavaScript or TypeScript you'd like to run with [deno](https://deno.land)
