# Devingen API Core

Contains models and utilities for other API projects.

## Development

### Adding this as a local dependency

Add `replace` config to go.mod of the other project before require.

```
replace github.com/devingen/api-core => ../api-core
```

### Releasing a new version

Check the existing tags and releases on the repo to avoid conflicts and override.
```
git tag --list
```

Create a new release.

```
make release VERSION=0.1.4
```