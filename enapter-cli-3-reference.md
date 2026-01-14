# Enapter CLI 3 Reference

## Overview

The Enapter CLI is a command-line interface tool for managing Enapter energy management services, including sites, devices, blueprints, and the rule engine. It provides a comprehensive set of commands for interacting with the Enapter Cloud platform and Gateway devices.

## Getting Started

### Authentication

The Enapter CLI requires an access token for authentication. You can obtain your access token from your Enapter Cloud account settings at [Enapter Cloud](https://cloud.enapter.com).

### Setting Up Your First Enapter Cloud Connection

The recommended way to use the Enapter CLI is by setting up named connections. This approach allows you to:
- Manage multiple environments (production, staging, development)
- Switch between Enapter Cloud and Gateway connections easily
- Associate connections with specific sites
- Store configuration securely

**Step 1: Add a Enapter Cloud connection**

```bash
enapter3 connection add --name my-cloud --token YOUR_ACCESS_TOKEN
```

**Step 2: Set it as default (optional)**

```bash
enapter3 connection set-default --name my-cloud
```

**Step 3: Verify the connection**

```bash
enapter3 connection list
```

### Setting Up Your First Enapter Gateway Connection

**Step 1: Navigate to your Enapter Gateway 3.0 Web Interface `System Settings` page by using Gateway IP address or mDNS name http://enapter-gateway.local/settings** 

**Step 2: Enter your Enapter Gateway password**

**Step 3: Click `API Token` and copy token to clipboard**

**Step 4: Add a named connection**

  ```bash
  enapter3 connection add --gateway \
    --name my-gateway \
    --url http://GATEWAY_IP/api \
    --token GATEWAY_API_TOKEN \
    --allow-insecure
  ```
**Step 5: Set it as default (optional)**

  ```bash
  enapter3 connection set-default --name my-gateway
  ```

### Quick Start Examples

Once your connection is set up, you can start managing your Enapter resources:

**For Enapter Cloud connections:**

```bash
# List all sites
enapter3 site list

# List all devices for a specific site
enapter3 device list --site-id SITE_ID

# Get device information
enapter3 device get --site-id SITE_ID --device-id DEVICE_ID

# Upload a blueprint (from file or directory)
enapter3 blueprint upload --path ./my-blueprint.enbp
# or
enapter3 blueprint upload --path ./my-blueprint/

# Create a new Lua device
enapter3 device create lua-device \
  --site-id SITE_ID \
  --runtime-id UCM_DEVICE_ID \
  --device-name "My Device" \
  --device-slug my-device \
  --blueprint-path ./blueprint/  # or ./blueprint.enbp
```

**For Gateway connections (local-first):**

```bash
# List all devices (no site-id needed)
enapter3 device list

# Get device information
enapter3 device get --device-id DEVICE_ID

# Create a new Lua device
enapter3 device create lua-device \
  --runtime-id UCM_DEVICE_ID \
  --device-name "My Device" \
  --device-slug my-device \
  --blueprint-path ./blueprint/  # or ./blueprint.enbp
```

## Connection Management

The connection commands allow you to manage multiple connections to Enapter Cloud and Gateway devices.

::: tip Cloud vs. Gateway
**Important distinction:**
- **Enapter Cloud connections** require `--site-id` parameter for most device, rule-engine, and site commands
- **Gateway connections** work in local-first mode and do NOT require `--site-id` parameter

You can create site-scoped Cloud connections using `--site-id` flag in `connection add` to avoid specifying it in every command.
:::

### connection add

Add a new connection to Enapter Cloud or a Gateway.

**Usage:**
```bash
enapter3 connection add [options]
```

**Options:**

| Option | Description | Required |
|--------|-------------|----------|
| `--name` | Connection name | Yes |
| `--token` | Enapter API access token | Yes |
| `--url` | API base URL | No (default: https://api.enapter.com) |
| `--gateway` | Connection is to a Gateway | No (default: false) |
| `--site-id` | Limit connection to specific site (Cloud only) | No |
| `--allow-insecure` | Allow insecure connections | No (default: false) |

**Example:**
```bash
enapter3 connection add \
  --name production \
  --token abc123... \
  --site-id site-456
```

### connection list

List all configured connections.

**Usage:**
```bash
enapter3 connection list
```

### connection remove

Remove a connection.

**Usage:**
```bash
enapter3 connection remove --name CONNECTION_NAME
```

### connection set-default

Set the default connection for CLI operations.

**Usage:**
```bash
enapter3 connection set-default --name CONNECTION_NAME
```

## Site Management

Manage your Enapter sites.

### Common Options

Most site commands support these options:

| Option | Description |
|--------|-------------|
| `--connection, -c` | Name of the connection to use |
| `--api-allow-insecure` | Allow insecure connections |
| `--verbose` | Log extra details about the operation |

### site list

List all user sites.

**Usage:**
```bash
enapter3 site list [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--my-sites` | Show only sites where user is owner or installer |
| `--limit` | Maximum number of sites to retrieve |

**Example:**
```bash
# List all sites
enapter3 site list

# List only your sites
enapter3 site list --my-sites

# Limit results
enapter3 site list --limit 10
```

### site get

Retrieve detailed information about a specific site.

**Usage:**
```bash
enapter3 site get --site-id SITE_ID
```

**Example:**
```bash
enapter3 site get --site-id 12345
```

## Device Management

Comprehensive commands for managing Enapter devices.

### Common Options

Most device commands support these options:

| Option | Description |
|--------|-------------|
| `--connection, -c` | Name of the connection to use |
| `--site-id` | Site ID (auto-detected from connection if available) |
| `--device-id, -d` | Device ID |
| `--api-allow-insecure` | Allow insecure connections |
| `--verbose` | Log extra details |

### device list

List all devices ordered by device ID.

**Usage:**
```bash
enapter3 device list [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--expand` | Expand device information (connectivity, manifest, properties, communication, site) |
| `--limit` | Maximum number of devices to retrieve |

**Example:**
```bash
# List all devices
enapter3 device list

# List with expanded information
enapter3 device list --expand connectivity --expand manifest

# List devices for specific site
enapter3 device list --site-id 12345
```

### device get

Retrieve detailed information about a specific device.

**Usage:**
```bash
enapter3 device get --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--expand` | Expand device information (connectivity, manifest, properties, communication, site) |

**Example:**
```bash
enapter3 device get --device-id abc123 --expand properties --expand manifest
```

### device create standalone

Create a new standalone device.

**Usage:**
```bash
enapter3 device create standalone [options]
```

**Options:**

| Option | Description | Required |
|--------|-------------|----------|
| `--site-id, -s` | Site ID where device will be created | Yes |
| `--device-name, -n` | Name for the new device | Yes |
| `--device-slug` | Slug for the device | Yes |

**Example:**
```bash
enapter3 device create standalone \
  --site-id 12345 \
  --device-name "My Device" \
  --device-slug my-device
```

### device create lua-device

Create a new Lua device.

**Usage:**
```bash
enapter3 device create lua-device [options]
```

**Options:**

| Option | Description | Required |
|--------|-------------|----------|
| `--site-id` | Site ID | Yes |
| `--runtime-id, -r` | UCM device ID where Lua device will run | Yes |
| `--device-name, -n` | Name for the new device | Yes |
| `--device-slug` | Slug for the device | Yes |
| `--blueprint-id, -b` | Blueprint ID to use | Yes* |
| `--blueprint-path` | Blueprint path (.enbp file or directory) | Yes* |

*Either `--blueprint-id` or `--blueprint-path` is required.

**Example:**
```bash
# Using blueprint file
enapter3 device create lua-device \
  --site-id 12345 \
  --runtime-id ucm-789 \
  --device-name "My Lua Device" \
  --device-slug my-lua-device \
  --blueprint-path ./my-blueprint.enbp

# Using blueprint directory
enapter3 device create lua-device \
  --site-id 12345 \
  --runtime-id ucm-789 \
  --device-name "My Lua Device" \
  --device-slug my-lua-device \
  --blueprint-path ./my-blueprint/
```

### device update

Update device properties.

**Usage:**
```bash
enapter3 device update --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--name` | New device name |
| `--slug` | New device slug |

**Example:**
```bash
enapter3 device update --device-id abc123 --name "Updated Device Name"
```

### device delete

Delete a device.

**Usage:**
```bash
enapter3 device delete --device-id DEVICE_ID
```

**Example:**
```bash
enapter3 device delete --device-id abc123 --site-id 12345
```

### device change-blueprint

Change the blueprint associated with a device.

**Usage:**
```bash
enapter3 device change-blueprint --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description | Required |
|--------|-------------|----------|
| `--blueprint-id, -b` | Blueprint ID | Yes* |
| `--blueprint-path` | Blueprint path (.enbp file or directory) | Yes* |

*Either `--blueprint-id` or `--blueprint-path` is required.

**Example:**
```bash
# Using blueprint ID
enapter3 device change-blueprint --device-id abc123 --blueprint-id bp-456

# Using local blueprint file
enapter3 device change-blueprint --device-id abc123 --blueprint-path ./new-blueprint.enbp

# Using local blueprint directory
enapter3 device change-blueprint --device-id abc123 --blueprint-path ./new-blueprint/
```

### device logs

Show device logs with filtering and streaming options.

**Usage:**
```bash
enapter3 device logs --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--follow, -f` | Follow log output in real-time |
| `--from` | From timestamp (RFC 3339 format) |
| `--to` | To timestamp (RFC 3339 format) |
| `--limit, -l` | Maximum number of logs to retrieve |
| `--offset, -o` | Number of logs to skip |
| `--severity, -s` | Filter by severity |
| `--order` | Sort order (RECEIVED_AT_ASC, RECEIVED_AT_DESC) |
| `--show` | Filter criteria (ALL, PERSISTED_ONLY, TEMPORARY_ONLY) |

**Example:**
```bash
# Stream logs in real-time
enapter3 device logs --device-id abc123 --follow

# Get last 100 error logs
enapter3 device logs --device-id abc123 --limit 100 --severity error

# Get logs from specific time range
enapter3 device logs --device-id abc123 \
  --from 2024-01-01T00:00:00Z \
  --to 2024-01-31T23:59:59Z
```

### device telemetry

Show device telemetry data.

**Usage:**
```bash
enapter3 device telemetry --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--follow, -f` | Follow telemetry output in real-time |

**Example:**
```bash
# Stream telemetry in real-time
enapter3 device telemetry --device-id abc123 --follow
```

### device monitor

Monitor device traffic in real-time. This command allows you to observe all incoming and outgoing device communications.

**Usage:**
```bash
enapter3 device monitor --device-id DEVICE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--include-runtime` | Monitor device's runtime traffic too (default: false) |

**Example:**
```bash
# Monitor device traffic
enapter3 device monitor --device-id abc123

# Monitor device traffic including runtime
enapter3 device monitor --device-id abc123 --include-runtime
```

### device run-terminal

Open a remote terminal session to a Gateway device.

::: warning
Remote terminal feature must be enabled in gateway settings. Use `Ctrl+]` to force connection close.
:::

**Usage:**
```bash
enapter3 device run-terminal --device-id GATEWAY_ID
```

**Example:**
```bash
enapter3 device run-terminal --device-id gateway-123
```

### device command execute

Execute a command on a device.

**Usage:**
```bash
enapter3 device command execute --device-id DEVICE_ID --name COMMAND_NAME [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--name` | Command name |
| `--arguments` | Command arguments (JSON string) |

**Example:**
```bash
# Execute command without arguments
enapter3 device command execute --device-id abc123 --name start

# Execute command with arguments
enapter3 device command execute --device-id abc123 \
  --name set_temperature \
  --arguments '{"value": 25.5}'
```

### device command list

List command executions for a device.

**Usage:**
```bash
enapter3 device command list --device-id DEVICE_ID
```

### device command get

Retrieve information about a specific command execution.

**Usage:**
```bash
enapter3 device command get --device-id DEVICE_ID --execution-id EXECUTION_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--expand` | Expand execution information (log) |

**Example:**
```bash
enapter3 device command get \
  --device-id abc123 \
  --execution-id exec-456 \
  --expand log
```

### device communication-config generate

Generate a new communication configuration for a device.

**Usage:**
```bash
enapter3 device communication-config generate --device-id DEVICE_ID --protocol PROTOCOL
```

**Options:**

| Option | Description |
|--------|-------------|
| `--protocol` | Connection protocol (MQTT, MQTTS) |

**Example:**
```bash
enapter3 device communication-config generate \
  --device-id abc123 \
  --protocol MQTTS
```

## Blueprint Management

Manage device blueprints for your Enapter devices.

### Common Options

Blueprint commands support these options:

| Option | Description |
|--------|-------------|
| `--connection, -c` | Name of the connection to use |
| `--api-allow-insecure` | Allow insecure connections |
| `--verbose` | Log extra details |

### blueprint get

Retrieve blueprint metadata.

**Usage:**
```bash
enapter3 blueprint get --blueprint-id BLUEPRINT_ID
```

**Example:**
```bash
enapter3 blueprint get --blueprint-id my-blueprint
```

### blueprint download

Download a blueprint from the platform.

**Usage:**
```bash
enapter3 blueprint download --blueprint-id BLUEPRINT_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--output, -o` | Output file name |

**Example:**
```bash
enapter3 blueprint download \
  --blueprint-id my-blueprint \
  --output my-blueprint.enbp
```

### blueprint upload

Upload a blueprint to the platform.

**Usage:**
```bash
enapter3 blueprint upload --path PATH
```

**Options:**

| Option | Description |
|--------|-------------|
| `--path, -p` | Blueprint path (.enbp file or directory) |

**Example:**
```bash
# Upload from enbp file
enapter3 blueprint upload --path ./my-blueprint.enbp

# Upload from directory
enapter3 blueprint upload --path ./my-blueprint/
```

### blueprint profiles download

Download blueprint profiles from the platform.

**Usage:**
```bash
enapter3 blueprint profiles download [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--output, -o` | Output file name |

**Example:**
```bash
enapter3 blueprint profiles download --output profiles.zip
```

### blueprint profiles upload

Upload blueprint profiles to the platform.

**Usage:**
```bash
enapter3 blueprint profiles upload --path PATH
```

**Options:**

| Option | Description |
|--------|-------------|
| `--path, -p` | Profiles zip file path |

**Example:**
```bash
enapter3 blueprint profiles upload --path ./profiles.zip
```

## Rule Engine Management

Manage automation rules for your Enapter sites.

### Common Options

Rule engine commands support these options:

| Option | Description |
|--------|-------------|
| `--connection, -c` | Name of the connection to use |
| `--site-id` | Site ID |
| `--api-allow-insecure` | Allow insecure connections |
| `--verbose` | Log extra details |

### rule-engine get

Retrieve rule engine information.

**Usage:**
```bash
enapter3 rule-engine get --site-id SITE_ID
```

**Example:**
```bash
enapter3 rule-engine get --site-id 12345
```

### rule-engine suspend

Suspend execution of all rules on a site.

**Usage:**
```bash
enapter3 rule-engine suspend --site-id SITE_ID
```

**Example:**
```bash
enapter3 rule-engine suspend --site-id 12345
```

### rule-engine resume

Resume execution of rules on a site.

**Usage:**
```bash
enapter3 rule-engine resume --site-id SITE_ID
```

**Example:**
```bash
enapter3 rule-engine resume --site-id 12345
```

### rule-engine rule create

Create a new automation rule.

**Usage:**
```bash
enapter3 rule-engine rule create --site-id SITE_ID [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--slug` | Unique slug for the rule | Required |
| `--script` | Path to script file | Required |
| `--runtime-version` | Runtime version (V1, V3) | V3 |
| `--exec-interval` | Execution interval (V1 only, e.g., 5s, 2m) | - |
| `--disable` | Create rule in disabled state | false |

**Example:**
```bash
enapter3 rule-engine rule create \
  --site-id 12345 \
  --slug my-automation-rule \
  --script ./rule-script.lua \
  --runtime-version V3
```

### rule-engine rule list

List all rules for a site.

**Usage:**
```bash
enapter3 rule-engine rule list --site-id SITE_ID
```

**Example:**
```bash
enapter3 rule-engine rule list --site-id 12345
```

### rule-engine rule get

Retrieve information about a specific rule.

**Usage:**
```bash
enapter3 rule-engine rule get --site-id SITE_ID --rule-id RULE_ID
```

**Example:**
```bash
enapter3 rule-engine rule get --site-id 12345 --rule-id my-rule
```

### rule-engine rule update

Update rule metadata.

**Usage:**
```bash
enapter3 rule-engine rule update --site-id SITE_ID --rule-id RULE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--slug` | New slug for the rule |

**Example:**
```bash
enapter3 rule-engine rule update \
  --site-id 12345 \
  --rule-id old-slug \
  --slug new-slug
```

### rule-engine rule update-script

Update the script of an existing rule.

**Usage:**
```bash
enapter3 rule-engine rule update-script --site-id SITE_ID --rule-id RULE_ID [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--script` | Path to new script file | Required |
| `--runtime-version` | Runtime version (V1, V3) | V3 |
| `--exec-interval` | Execution interval (V1 only) | - |

**Example:**
```bash
enapter3 rule-engine rule update-script \
  --site-id 12345 \
  --rule-id my-rule \
  --script ./electrolyser-controller.lua
```

### rule-engine rule delete

Delete a rule.

**Usage:**
```bash
enapter3 rule-engine rule delete --site-id SITE_ID --rule-id RULE_ID
```

**Example:**
```bash
enapter3 rule-engine rule delete --site-id 12345 --rule-id my-rule
```

### rule-engine rule enable

Enable one or more rules.

**Usage:**
```bash
enapter3 rule-engine rule enable --site-id SITE_ID --rule-id RULE_ID [--rule-id RULE_ID ...]
```

**Example:**
```bash
# Enable single rule
enapter3 rule-engine rule enable --site-id 12345 --rule-id rule1

# Enable multiple rules
enapter3 rule-engine rule enable \
  --site-id 12345 \
  --rule-id rule1 \
  --rule-id rule2 \
  --rule-id rule3
```

### rule-engine rule disable

Disable one or more rules.

**Usage:**
```bash
enapter3 rule-engine rule disable --site-id SITE_ID --rule-id RULE_ID [--rule-id RULE_ID ...]
```

**Example:**
```bash
# Disable single rule
enapter3 rule-engine rule disable --site-id 12345 --rule-id rule1

# Disable multiple rules
enapter3 rule-engine rule disable \
  --site-id 12345 \
  --rule-id rule1 \
  --rule-id rule2
```

### rule-engine rule logs

Show logs for a specific rule.

**Usage:**
```bash
enapter3 rule-engine rule logs --site-id SITE_ID --rule-id RULE_ID [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--follow, -f` | Follow log output in real-time |

**Example:**
```bash
# Stream rule logs in real-time
enapter3 rule-engine rule logs --site-id 12345 --rule-id my-rule --follow
```

## Advanced Usage

### Working with Multiple Connections

You can manage multiple connections and switch between them:

```bash
# Add production connection
enapter3 connection add \
  --name production \
  --token prod-token

# Add staging connection
enapter3 connection add \
  --name staging \
  --token staging-token

# Use specific connection
enapter3 device list --connection production
enapter3 device list --connection staging

# Set default
enapter3 connection set-default --name production
```

### Gateway Connections

Connect directly to an Enapter Gateway:

```bash
enapter3 connection add \
  --name my-gateway \
  --gateway \
  --url https://gateway.local \
  --token gateway-token \
  --allow-insecure
```

### Site-Scoped Connections

When working with Enapter Cloud, you can create site-scoped connections to avoid specifying `--site-id` for every command:

```bash
enapter3 connection add \
  --name site-specific \
  --token your-token \
  --site-id 12345
```

When using this connection, the `--site-id` flag is automatically set for all commands:

```bash
# Without site-scoped connection (Cloud)
enapter3 device list --site-id 12345

# With site-scoped connection (Cloud)
enapter3 device list --connection site-specific
```

### Cloud vs. Gateway Connections

**Enapter Cloud connections** require `--site-id` for most device, rule-engine, and site-related commands:

```bash
# Cloud connection setup
enapter3 connection add --name cloud --token YOUR_TOKEN

# Commands require site-id
enapter3 device list --site-id 12345
enapter3 device get --site-id 12345 --device-id abc123
enapter3 rule-engine get --site-id 12345

# Or use site-scoped connection
enapter3 connection add --name cloud-site --token YOUR_TOKEN --site-id 12345
enapter3 device list --connection cloud-site  # site-id is automatic
```

**Gateway connections** work in local-first mode and do not require `--site-id`:

```bash
# Gateway connection setup
enapter3 connection add \
  --name my-gateway \
  --gateway \
  --url https://gateway.local \
  --token gateway-token

# Commands work without site-id
enapter3 device list
enapter3 device get --device-id abc123
enapter3 rule-engine get
```

### Expanding Data

Many commands support the `--expand` flag to retrieve additional information:

```bash
# Get device with all available expansions
enapter3 device get --device-id abc123 \
  --expand connectivity \
  --expand manifest \
  --expand properties \
  --expand communication \
  --expand site

# Get command execution with logs
enapter3 device command get \
  --device-id abc123 \
  --execution-id exec-456 \
  --expand log
```

### Streaming Data

Commands that support real-time data streaming:

```bash
# Stream device logs
enapter3 device logs --device-id abc123 --follow

# Stream device telemetry
enapter3 device telemetry --device-id abc123 --follow

# Monitor device traffic
enapter3 device monitor --device-id abc123

# Stream rule logs
enapter3 rule-engine rule logs --site-id 12345 --rule-id my-rule --follow
```

Use `Ctrl+C` to stop streaming.

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:

1. Verify your token is correct:
```bash
echo $ENAPTER3_API_TOKEN
```

2. Check your connection configuration:
```bash
enapter3 connection list
```

3. Test with verbose logging:
```bash
enapter3 device list --verbose
```

### Insecure Connections

For development or local Gateway connections, you may need to allow insecure connections:

```bash
# Global setting
export ENAPTER3_API_ALLOW_INSECURE=true

# Per-connection setting
enapter3 connection add \
  --name dev-gateway \
  --gateway \
  --url https://192.168.1.100 \
  --token token \
  --allow-insecure

# Per-command setting
enapter3 device list --api-allow-insecure
```

### Verbose Logging

Enable verbose logging for debugging:

```bash
enapter3 device get --device-id abc123 --verbose
```

## Command Reference Summary

### Connection Commands
- `connection add` - Add a new connection
- `connection list` - List all connections
- `connection remove` - Remove a connection
- `connection set-default` - Set default connection

### Site Commands
- `site list` - List user sites
- `site get` - Get site details

### Device Commands
- `device create standalone` - Create standalone device
- `device create lua-device` - Create Lua device
- `device list` - List devices
- `device get` - Get device details
- `device update` - Update device
- `device delete` - Delete device
- `device change-blueprint` - Change device blueprint
- `device logs` - Show device logs
- `device telemetry` - Show device telemetry
- `device monitor` - Monitor device traffic
- `device run-terminal` - Open remote terminal
- `device command execute` - Execute device command
- `device command list` - List command executions
- `device command get` - Get command execution details
- `device communication-config generate` - Generate communication config

### Blueprint Commands
- `blueprint get` - Get blueprint metadata
- `blueprint download` - Download blueprint
- `blueprint upload` - Upload blueprint
- `blueprint profiles download` - Download blueprint profiles
- `blueprint profiles upload` - Upload blueprint profiles

### Rule Engine Commands
- `rule-engine get` - Get rule engine info
- `rule-engine suspend` - Suspend rule execution
- `rule-engine resume` - Resume rule execution
- `rule-engine rule create` - Create new rule
- `rule-engine rule list` - List rules
- `rule-engine rule get` - Get rule details
- `rule-engine rule update` - Update rule metadata
- `rule-engine rule update-script` - Update rule script
- `rule-engine rule delete` - Delete rule
- `rule-engine rule enable` - Enable rule(s)
- `rule-engine rule disable` - Disable rule(s)
- `rule-engine rule logs` - Show rule logs

## Best Practices

### 1. Use Connections

Set up named connections instead of relying solely on environment variables:

```bash
# Good
enapter3 connection add --name prod --token TOKEN
enapter3 device list --connection prod

# Less flexible
export ENAPTER3_API_TOKEN=TOKEN
enapter3 device list
```

### 2. Use Site-Scoped Connections

For multi-site setups, create separate connections for each site:

```bash
enapter3 connection add --name site-a --token TOKEN --site-id SITE_A_ID
enapter3 connection add --name site-b --token TOKEN --site-id SITE_B_ID
```

### 3. Leverage Verbose Mode for Debugging

When troubleshooting issues, always use verbose mode:

```bash
enapter3 device get --device-id abc123 --verbose
```

### 4. Use Expand Flags Wisely

Only request expanded data when needed to minimize API load:

```bash
# Only expand what you need
enapter3 device get --device-id abc123 --expand properties
```

### 5. Follow Logs for Real-Time Monitoring

Use follow mode for development and debugging:

```bash
enapter3 device logs --device-id abc123 --follow --severity error
```

## Examples and Use Cases

### Device Lifecycle Management

**Cloud example:**

```bash
# Set site ID for all commands
SITE_ID="12345"

# 1. Create a new Lua device
enapter3 device create lua-device \
  --site-id "${SITE_ID}" \
  --runtime-id ucm-001 \
  --device-name "Temperature Sensor" \
  --device-slug temp-sensor-01 \
  --blueprint-path ./temp-sensor-blueprint/  # or .enbp file

# 2. Monitor device startup
enapter3 device logs --site-id "${SITE_ID}" --device-id temp-sensor-01 --follow

# 3. Check device telemetry
enapter3 device telemetry --site-id "${SITE_ID}" --device-id temp-sensor-01

# 4. Update device if needed
enapter3 device change-blueprint \
  --site-id "${SITE_ID}" \
  --device-id temp-sensor-01 \
  --blueprint-path ./temp-sensor-v2/  # or .enbp file

# 5. Execute commands on device
enapter3 device command execute \
  --site-id "${SITE_ID}" \
  --device-id temp-sensor-01 \
  --name calibrate \
  --arguments '{"offset": 1.5}'
```

**Gateway example:**

```bash
# 1. Create a new Lua device (no site-id needed)
enapter3 device create lua-device \
  --runtime-id ucm-001 \
  --device-name "Temperature Sensor" \
  --device-slug temp-sensor-01 \
  --blueprint-path ./temp-sensor-blueprint/  # or .enbp file

# 2. Monitor device startup
enapter3 device logs --device-id temp-sensor-01 --follow

# 3. Check device telemetry
enapter3 device telemetry --device-id temp-sensor-01

# 4. Update device if needed
enapter3 device change-blueprint \
  --device-id temp-sensor-01 \
  --blueprint-path ./temp-sensor-v2/  # or .enbp file

# 5. Execute commands on device
enapter3 device command execute \
  --device-id temp-sensor-01 \
  --name calibrate \
  --arguments '{"offset": 1.5}'
```

### Blueprint Development Workflow

```bash
# 1. Download existing blueprint
enapter3 blueprint download \
  --blueprint-id existing-blueprint \
  --output current-version.enbp

# 2. Make modifications locally
# ... edit files ...

# 3. Upload updated blueprint
enapter3 blueprint upload --path ./modified-blueprint/

# 4. Update device to use new blueprint
enapter3 device change-blueprint \
  --device-id test-device \
  --blueprint-path ./modified-blueprint/

# 5. Test and monitor
enapter3 device logs --device-id test-device --follow
```

### Rule Engine Automation

```bash
# 1. Create automation rule
enapter3 rule-engine rule create \
  --site-id 12345 \
  --slug temperature-alert \
  --script ./temperature-alert.lua \
  --runtime-version V3

# 2. Monitor rule execution
enapter3 rule-engine rule logs \
  --site-id 12345 \
  --rule-id temperature-alert \
  --follow

# 3. Disable rule for maintenance
enapter3 rule-engine rule disable \
  --site-id 12345 \
  --rule-id temperature-alert

# 4. Update rule script
enapter3 rule-engine rule update-script \
  --site-id 12345 \
  --rule-id temperature-alert \
  --script ./temperature-alert-v2.lua

# 5. Re-enable rule
enapter3 rule-engine rule enable \
  --site-id 12345 \
  --rule-id temperature-alert
```

### Multi-Site Management

```bash
# Set up connections for each site
enapter3 connection add --name site-factory --token TOKEN --site-id FACTORY_ID
enapter3 connection add --name site-warehouse --token TOKEN --site-id WAREHOUSE_ID

# List devices per site
enapter3 device list --connection site-factory
enapter3 device list --connection site-warehouse

# Deploy same rule to multiple sites
for site in site-factory site-warehouse; do
  enapter3 rule-engine rule create \
    --connection $site \
    --slug monitoring-rule \
    --script ./monitoring.lua
done
```

## Getting Help

### Command-Line Help

Get help for any command using the `--help` flag:

```bash
# General help
enapter3 --help

# Command group help
enapter3 device --help

# Specific command help
enapter3 device create lua-device --help
```

### Additional Resources

- [Enapter Cloud Platform](https://cloud3.enapter.com)
- [Enapter Documentation](https://developers.enapter.com)
- [Blueprint Marketplace](https://marketplace.enapter.com)

---

## Appendix: Environment Variables

For unattended setups, CI/CD pipelines, and shell scripting, you can use environment variables instead of managing connections. This approach is recommended for automation scenarios where interactive connection management is not practical.

### Supported Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENAPTER3_API_TOKEN` | Enapter API access token | - |
| `ENAPTER3_API_URL` | Enapter API base URL | `https://api.enapter.com` |
| `ENAPTER3_API_ALLOW_INSECURE` | Allow insecure connections | `false` |

### Usage in Scripts

**Gateway (local-first) script example:**

```bash
#!/bin/bash

# Set authentication for Gateway
export ENAPTER3_API_TOKEN="your-gateway-token"
export ENAPTER3_API_URL="http://enapter-gateway.local/api"
export ENAPTER3_API_ALLOW_INSECURE=true

# Run commands - no site-id needed for Gateway
enapter3 device list
enapter3 device get --device-id abc123
enapter3 rule-engine get
```

**Cloud script example:**

```bash
#!/bin/bash

# Set authentication for Cloud
export ENAPTER3_API_TOKEN="your-cloud-token"
export ENAPTER3_API_URL="https://api.enapter.com"

# Site ID required for Cloud commands
SITE_ID="12345"

enapter3 device list --site-id "${SITE_ID}"
enapter3 device get --site-id "${SITE_ID}" --device-id abc123
enapter3 rule-engine get --site-id "${SITE_ID}"
```

**CI/CD pipeline example (Cloud):**

```bash
#!/bin/bash

# Use secrets from CI/CD environment
export ENAPTER3_API_TOKEN="${ENAPTER_TOKEN}"
export ENAPTER3_API_URL="https://api.enapter.com"

# Site ID is required for Cloud deployments
SITE_ID="${ENAPTER_SITE_ID}"

# Deploy blueprint
enapter3 blueprint upload --path ./my-blueprint.enbp

# Create device (site-id required for Cloud)
enapter3 device create lua-device \
  --site-id "${SITE_ID}" \
  --runtime-id "${RUNTIME_ID}" \
  --device-name "Automated Device" \
  --device-slug "auto-device-${CI_BUILD_ID}" \
  --blueprint-id "my-blueprint"
```

**CI/CD pipeline example (Gateway):**

```bash
#!/bin/bash

# Use secrets from CI/CD environment for Gateway
export ENAPTER3_API_TOKEN="${GATEWAY_TOKEN}"
export ENAPTER3_API_URL="http://${GATEWAY_ADDRESS}/api"
export ENAPTER3_API_ALLOW_INSECURE=true

# Deploy blueprint
enapter3 blueprint upload --path ./my-blueprint.enbp

# Create device (no site-id needed for Gateway)
enapter3 device create lua-device \
  --runtime-id "${RUNTIME_ID}" \
  --device-name "Automated Device" \
  --device-slug "auto-device-${CI_BUILD_ID}" \
  --blueprint-id "my-blueprint"
```

### Precedence Rules

When both environment variables and connection configuration are present:

1. Command-line flags (e.g., `--connection`) have highest priority
2. Connection configuration (from `connection add`) is used if no command-line flags or environment variables specified
3. Environment variables are used only if no connection is configured or specified

---

*This documentation is for Enapter CLI version 3.x*
