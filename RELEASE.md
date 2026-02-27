# Releasing

Make sure your working tree is clean before starting (`make check-clean`).

1. **Run the appropriate release target** on a new branch:

   ```sh
   git checkout -b bump/v1.2.2
   make release-patch   # 0.0.x
   make release-minor   # 0.x.0
   make release-major   # x.0.0
   ```

   This bumps `version.go`, commits the change, and creates an annotated tag locally.

2. **Push the branch and open a PR** (main is branch-protected):

   ```sh
   git push origin bump/v1.2.2
   # open a PR and get it merged
   ```

3. **After the PR is merged, push the tag**:

   ```sh
   git checkout main && git pull
   git push origin v1.2.2
   ```

   Tags are not subject to the branch protection rules, so this pushes directly.

4. **Trigger the Go module proxy** to index the new version immediately:

   ```sh
   make proxy
   ```
