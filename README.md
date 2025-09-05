# Enapter CLI
![Build Status](https://github.com/enapter/enapter-cli/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/enapter/enapter-cli)](/LICENSE)
[![Release](https://img.shields.io/github/release/enapter/enapter-cli.svg)](https://github.com/enapter/enapter-cli/releases/latest)


This tool helps Enapter customers to work with devices it is alternative for [Enapter IDE for EMS Toolkit 3.0](https://marketplace.visualstudio.com/items?itemName=Enapter.enapter-ems-toolkit-ide). 
It helpful in the following cases:

1. Managing all your EMS setup as a code with Git and Ansible / Puppet
2. Establishing CI/CD workflow
3. Development and debugging of Enapter Blueprints
4. Development and debugging of Enapter Gateway Rules

## How to install

###  macOS - recommended

Version 1:

```bash
brew tap enapter/tap && brew install enapter
```

Version 3:

```bash
brew tap enapter/tap && brew install enapter@3
```

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

## How to use Version 1:

---
**NOTE**

Version 1 works only with Enapter Cloud connection.
---

### API token

Enapter CLI requires access token for authentication. Obtaining of the token is easy and can be done by following few steps.

1. Ensure you have registed [Enapter Cloud](https://cloud.enapter.com) account. If not, sign up [here](https://sso.enapter.com/users/new).
2. Log in to your Enapter Cloud account, click on your profile name in top right corner and choose `Account Settings`
3. Select `API Tokens` menu and click `New Token` button
4. Follow the instructions on the screen
<img src="./.assets/token.png">

5. Set environment variable `ENAPTER_API_TOKEN` with new token. To make it permanent don't forget to add it to configuration files of your shell.

  ```bash
  export ENAPTER_API_TOKEN="your token"
  ```

Please note that if you don't save your token, it is not possible to reveal it anymore. You need generate new token.

## How to use Version 3:

### API token

Enapter CLI requires access token for authentication. Obtaining of the token is easy and can be done by following few steps.

1. Navigate to your Enapter Gateway 3.0 Web Interface `Settings` page by using IP address or mDNS name [http://enapter-gateway.local/settings](https://enapter-gateway.local/settings)
2. Enapter your Enapter Gateway password
3. Click `API Token` and copy token to clipboard
4. Set environment variables `ENAPTER3_API_TOKEN`, `ENAPTER3_API_URL` and `ENAPTER3_API_ALLOW_INSECURE`. To make it permanent don't forget to add it to configuration files of your shell.

  ```bash
  export ENAPTER3_API_TOKEN="your token"
  export ENAPTER3_API_URL="http://ip_address/api"
  export ENAPTER3_API_ALLOW_INSECURE=true
  ```

5. Check connection works by running

  ```bash
  enapter3 device list
  ```

### Autocompletion in your favourite terminal app

---
**NOTE**

Available for Version 1 now.
---

In order to make life easier with command line interface, you may use [Amazon Q](https://aws.amazon.com/q/). This autocompletion tool has native support for the Enapter CLI for Mac OS X and Linux.

<img src="./.assets/enapter-cli-fig-integration.gif">

