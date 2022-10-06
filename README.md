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
git tag -a v0.0.27 -m "change cors handling, convert all headers to lowercase"

# push new tag
git push origin v0.0.27
```