# Releasing

1. Bump the version in `version.go`:

   ```go
   const Version = "1.2.2"
   ```

2. Commit and tag:

   ```sh
   git add version.go
   git commit -m "chore: bump version to 1.2.2"
   git tag v1.2.2
   git push origin main --tags
   ```

3. Trigger the Go module proxy to index the new version:

   ```sh
   GOPROXY=proxy.golang.org go list -m github.com/linkvite/go@v1.2.2
   ```

   This pushes the release into the public proxy cache so it's immediately available via `go get`, without waiting for the proxy to discover it on its own.
