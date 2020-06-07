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
git tag -a v0.0.1 -m "initial version"
git push origin v0.0.1
```