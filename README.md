# Devingen API Core

Contains models and utilities for other API projects.

## Development

### Adding this as a local dependency

Add `replace` config to go.mod of the other project before require.

```
replace github.com/devingen/api-core => ../api-core
```

### Releasing a new version

Create a git tag with the desired version and push the tag.

```
# see tags
git tag --list

# create new tag
git tag -a v0.0.14 -m "get header function to request, add data transfer field to field types"

# push new tag
git push origin v0.0.14
```