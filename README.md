# Enapter CLI

This tool helps Enapter customers to work with devices. It useful in the following cases:
1. Develop devices via blueprints.
2. Update and monitor devices.

## How to install

### Get prebuilt binaries

Choose your platform and required release on the [Releases page](https://github.com/golangci/golangci-lint/releases).

### Build from source

The following command builds enapter binary:
```
./build.sh
```

Or you can pass custom output path:
```
./build.sh /usr/local/bin/enapter
```

## How to use

### API token

First of all you need an API token. Please contact with customer support to get one.

Store token into environment variable `ENAPTER_API_TOKEN` to use it with enapter cli tool.
