# Releasing

1. Bump the version in `version.go`:

   ```go
   const Version = "1.2.2"
   ```

2. Commit and open a PR (main is branch-protected):

   ```sh
   git checkout -b bump/v1.2.2
   git add version.go
   git commit -m "chore: bump version to 1.2.2"
   git push origin bump/v1.2.2
   # open a PR, get it reviewed and merged
   ```

3. After the PR is merged, tag the merge commit on main:

   ```sh
   git checkout main && git pull
   git tag v1.2.2
   git push origin v1.2.2
   ```

4. Trigger the Go module proxy to index the new version:

   ```sh
   GOPROXY=proxy.golang.org go list -m github.com/linkvite/go@v1.2.2
   ```

   This pushes the release into the public proxy cache so it's immediately available via `go get`, without waiting for the proxy to discover it on its own.
