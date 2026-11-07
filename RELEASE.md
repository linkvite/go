# Releasing github.com/linkvite/go

There is no "publish" step for Go modules: **a release is a git tag** (`v2.3.1`) on `main`. Once the tag is on GitHub, `go get github.com/linkvite/go/v2@v2.3.1` works, and the Go module proxy and pkg.go.dev pick it up. `version.go` holds the same number so the SDK can send it in its User-Agent.

Tags are permanent: the module proxy and checksum database cache a version forever. **Never move or reuse a tag**; release a new version instead.

## Every release

### 1. If the API changed: refresh the spec

The SDK is generated from cultus's OpenAPI document. In cultus, `make openapi` writes `api/openapi.json`; then here:

```sh
make generate SPEC=../cultus/api/openapi.json   # copies the spec, regenerates *_gen.go
make check                                       # generated code current + lint + tests
```

Use the spec from cultus's `main`, the API that's deployed. Commit and merge it through a PR like any change.

### 2. Bump and tag

On a clean, up-to-date branch:

```sh
git checkout -b bump/v2.0.1
make release-patch   # 2.0.0 -> 2.0.1: fixes, new optional fields
make release-minor   # 2.0.1 -> 2.1.0: new methods or services
```

This bumps `version.go`, commits it and creates an annotated tag locally.

### 3. Merge, then push the tag

`main` is branch-protected, so the bump goes through a PR:

```sh
git push origin bump/v2.0.1
# open the PR and merge it
git checkout main && git pull
git push origin v2.0.1
```

If the PR was squash-merged, the tag points at a commit that isn't on `main`. Re-create it on `main` before pushing:

```sh
git tag -d v2.0.1
make tag             # tags version.go's version on the current commit
git push origin v2.0.1
```

### 4. Make it visible

```sh
make proxy           # asks proxy.golang.org for the new version, so it's installable right away
```

pkg.go.dev updates within a few minutes. Optionally, create a GitHub release from the tag with the highlights.

## Major versions (v2 → v3)

A breaking change needs a new major version, and Go puts the major in the import path:

1. Change the module path in `go.mod` to `github.com/linkvite/go/v3`, and update imports in the repo (`example/`, docs, README).
2. Set `version.go` to `3.0.0`.
3. Merge, then on `main`: `make tag` and `git push origin v3.0.0`.

`make release-major` only prints these steps: a `v3.0.0` tag on a `/v2` module path is rejected by the go command. `make tag`, `make release-minor`/`-patch` and `make proxy` refuse to run when `version.go`'s major doesn't match the module path.

### The first v2 release

The v2 rewrite already has `module github.com/linkvite/go/v2` and `version.go` at `2.0.0`, so don't bump; tag it once the v2 branch is merged:

```sh
git checkout main && git pull
make tag                 # v2.0.0
git push origin v2.0.0
make proxy
```

v1 users are unaffected: `github.com/linkvite/go` keeps resolving to v1.x until they change their import path.

## If something goes wrong

- **Bad release:** tag a fixed patch. To warn people off the broken one, add a `retract` directive to `go.mod` in that patch:

  ```go
  retract v2.0.1 // broken bookmark decoding
  ```

  `go get` then skips it and warns anyone who already depends on it.
- **Tag pushed by mistake** (before anyone fetched it): delete it locally and on GitHub. Once the proxy has seen it, it's cached forever, so retract it instead.
