# Enapter CLI
![Build Status](https://github.com/enapter/enapter-cli/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/enapter/enapter-cli)](/LICENSE)


This tool helps Enapter customers to work with devices. It useful in the following cases:
1. Develop devices via blueprints.
2. Update and monitor devices.

## How to install

### Get prebuilt binaries

Choose your platform and required release on the [Releases page](https://github.com/Enapter/enapter-cli/releases).

### Build from source

You should have [installed Go tools](https://golang.org/doc/install). Then you can build CLI via the following command:
```
./build.sh
```

Also you can pass custom output path:
```
./build.sh /usr/local/bin/enapter
```

## How to use

### API token

Enapter CLI requires access token for authentication. At the moment we provide it only to selected partners. Contact us at support@enapter.com to get your token.

Store token into environment variable `ENAPTER_API_TOKEN` to use it with enapter cli tool.
