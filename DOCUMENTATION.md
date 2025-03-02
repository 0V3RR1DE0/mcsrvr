# MCSRVR Documentation

This document provides detailed information about the MCSRVR Minecraft server management tool.

## Table of Contents

1. [Command Reference](#command-reference)
2. [Server Types](#server-types)
3. [Configuration](#configuration)
4. [Server Management](#server-management)
5. [RCON](#rcon)
6. [Backups](#backups)
7. [Troubleshooting](#troubleshooting)
8. [Advanced Usage](#advanced-usage)

## Command Reference

### `init` - Initialize a new Minecraft server

```
mcsrvr init <path> -n <name> <server-type> -v <version> [options]
```

Parameters:
- `<path>`: Path where the server will be created
- `-n, --name <name>`: Name of the server (used for management)
- `<server-type>`: Type of server to create (vanilla, papermc, fabric)
- `-v, --version <version>`: Minecraft version (e.g., 1.21.4)

Options:
- `--memory <memory>`: Memory allocation (default: 2G)
- `--java-args <args>`: Additional Java arguments

Examples:
```bash
# Initialize a PaperMC server
mcsrvr init D:/MCServers/Paper/MyServer -n MyServer papermc -v 1.21.4

# Initialize a vanilla server with custom memory allocation
mcsrvr init D:/MCServers/Vanilla/MyServer -n MyServer vanilla -v 1.21.4 --memory 4G

# Initialize a Fabric server with custom Java arguments
mcsrvr init D:/MCServers/Fabric/MyServer -n MyServer fabric -v 1.21.4 --java-args "-XX:+UseG1GC"
```

### `list` - List all servers

```
mcsrvr list [options]
```

Options:
- `--online`: Show only online servers
- `--offline`: Show only offline servers

Examples:
```bash
# List all servers
mcsrvr list

# List only online servers
mcsrvr list --online
```

### `start` - Start a server

```
mcsrvr start <server-name>
```

Parameters:
- `<server-name>`: Name of the server to start

Example:
```bash
mcsrvr start MyServer
```

### `stop` - Stop a server

```
mcsrvr stop <server-name>
```

Parameters:
- `<server-name>`: Name of the server to stop

Example:
```bash
mcsrvr stop MyServer
```

### `restart` - Restart a server

```
mcsrvr restart <server-name>
```

Parameters:
- `<server-name>`: Name of the server to restart

Example:
```bash
mcsrvr restart MyServer
```

### `console` - Access server console

```
mcsrvr console <server-name>
```

Parameters:
- `<server-name>`: Name of the server to access

Example:
```bash
mcsrvr console MyServer
```

Use `exit` or `Ctrl+C` to exit the console.

### `cmd` - Execute a command on a server

```
mcsrvr cmd <server-name> "<command>"
```

Parameters:
- `<server-name>`: Name of the server
- `"<command>"`: Command to execute (without the leading slash)

Example:
```bash
mcsrvr cmd MyServer "say Hello, world!"
```

### `log` - View server logs

```
mcsrvr log <server-name> [options]
```

Parameters:
- `<server-name>`: Name of the server

Options:
- `--follow`: Follow log output in real-time
- `--lines <n>`: Number of lines to show (default: 20)

Examples:
```bash
# View the last 20 lines of the log
mcsrvr log MyServer

# View the last 100 lines of the log
mcsrvr log MyServer --lines 100

# Follow the log in real-time
mcsrvr log MyServer --follow
```

### `backup` - Create a server backup

```
mcsrvr backup <server-name> <backup-path>
```

Parameters:
- `<server-name>`: Name of the server to backup
- `<backup-path>`: Path where the backup will be stored

Example:
```bash
mcsrvr backup MyServer D:/MCBackups
```

### `backups` - List server backups

```
mcsrvr backups [server-name]
```

Parameters:
- `[server-name]`: (Optional) Name of the server to list backups for

Examples:
```bash
# List all backups
mcsrvr backups

# List backups for a specific server
mcsrvr backups MyServer
```

### `restore` - Restore a server backup

```
mcsrvr restore <backup-path> <restore-path>
```

Parameters:
- `<backup-path>`: Path to the backup
- `<restore-path>`: Path where the server will be restored

Example:
```bash
mcsrvr restore D:/MCBackups/MyServer_2025-03-02_12-34-56 D:/MCServers/Restored
```

### `del` - Delete a server

```
mcsrvr del <server-name> [options]
```

Parameters:
- `<server-name>`: Name of the server to delete

Options:
- `-y, --yes`: Skip confirmation prompt

Examples:
```bash
# Delete a server with confirmation prompt
mcsrvr del MyServer

# Delete a server without confirmation prompt
mcsrvr del MyServer -y
```

### `config` - Configure server settings

```
mcsrvr config [server-name] [config-type] [options]
```

Parameters:
- `[server-name]`: (Optional) Name of the server to configure
- `[config-type]`: (Optional) Type of configuration (start, properties, ops, rcon)

Options:
- `--default-memory <memory>`: Default memory allocation for new servers
- `--default-java-args <args>`: Default Java arguments for new servers
- `--port <port>`: RCON port (for rcon config-type)
- `--password <password>`: RCON password (for rcon config-type)

Examples:
```bash
# Set default memory allocation for new servers
mcsrvr config --default-memory 4G

# Set default Java arguments for new servers
mcsrvr config --default-java-args "-XX:+UseG1GC -XX:+ParallelRefProcEnabled"

# Configure RCON settings for a server
mcsrvr config MyServer rcon --port 25575 --password mypassword

# Edit server.properties
mcsrvr config MyServer properties

# Edit startup script
mcsrvr config MyServer start

# Edit ops.json
mcsrvr config MyServer ops
```

## Server Types

MCSRVR supports the following server types:

### Vanilla

The official Minecraft server provided by Mojang. This is the standard server without any modifications.

```bash
mcsrvr init D:/MCServers/Vanilla/MyServer -n MyServer vanilla -v 1.21.4
```

### PaperMC

A high-performance fork of Spigot that aims to fix gameplay and mechanics inconsistencies and improve performance.

```bash
mcsrvr init D:/MCServers/Paper/MyServer -n MyServer papermc -v 1.21.4
```

### Fabric

A lightweight, modular modding toolchain for Minecraft.

```bash
mcsrvr init D:/MCServers/Fabric/MyServer -n MyServer fabric -v 1.21.4
```

## Configuration

### Server Configuration

MCSRVR stores server configurations in `~/.mcsrvr/config.json`. Each server has the following configuration options:

- `name`: Server name
- `type`: Server type (vanilla, papermc, fabric)
- `version`: Minecraft version
- `path`: Path to the server directory
- `memory`: Memory allocation
- `javaArgs`: Additional Java arguments
- `lastStarted`: Timestamp of when the server was last started

### Default Configuration

You can set default configuration options for new servers:

```bash
# Set default memory allocation
mcsrvr config --default-memory 4G

# Set default Java arguments
mcsrvr config --default-java-args "-XX:+UseG1GC -XX:+ParallelRefProcEnabled"
```

### Server Properties

You can edit the server.properties file using the config command:

```bash
mcsrvr config MyServer properties
```

This will open the server.properties file in your default text editor.

### RCON Configuration

RCON (Remote Console) allows you to connect to a running Minecraft server and execute commands. You can configure RCON settings using the config command:

```bash
mcsrvr config MyServer rcon --port 25575 --password mypassword
```

## Server Management

### Process Management

MCSRVR runs Minecraft servers as hidden processes, similar to systemctl in Linux. This means that servers can run in the background without keeping a terminal window open.

Active server processes are tracked in `~/.mcsrvr/active_servers.json`. This file is updated whenever a server is started or stopped.

### PID Tracking

MCSRVR tracks both the command window process (cmd.exe on Windows) and the Java process. When stopping a server, MCSRVR will:

1. Try to stop the server gracefully using RCON
2. If RCON fails, kill the Java process
3. If the command window process is still running, kill it

### Server Status

You can check the status of all servers using the list command:

```bash
mcsrvr list
```

This will show:
- Server name
- Server type
- Minecraft version
- Server path
- Status (Online/Offline)
- PID (if online)
- Last started timestamp

## RCON

RCON (Remote Console) is a protocol that allows you to connect to a running Minecraft server and execute commands. MCSRVR uses RCON for:

- Executing commands on a server (`mcsrvr cmd`)
- Accessing the server console (`mcsrvr console`)
- Gracefully stopping a server (`mcsrvr stop`)

RCON is automatically configured when a server is initialized. The default RCON port is 25575 and the default password is "mcsrvr".

You can change the RCON settings using the config command:

```bash
mcsrvr config MyServer rcon --port 25575 --password mypassword
```

## Backups

### Creating Backups

You can create a backup of a server using the backup command:

```bash
mcsrvr backup MyServer D:/MCBackups
```

This will create a backup in the specified directory with a timestamp, e.g., `MyServer_2025-03-02_12-34-56`.

### Listing Backups

You can list all backups using the backups command:

```bash
# List all backups
mcsrvr backups

# List backups for a specific server
mcsrvr backups MyServer
```

### Restoring Backups

You can restore a backup using the restore command:

```bash
mcsrvr restore D:/MCBackups/MyServer_2025-03-02_12-34-56 D:/MCServers/Restored
```

## Troubleshooting

### Server Won't Start

If a server won't start, check the following:

1. Make sure Java is installed and in your PATH
2. Check the server log for errors: `mcsrvr log <server-name>`
3. Make sure the server directory exists and contains the necessary files
4. Check if the server is already running: `mcsrvr list`

### RCON Connection Failed

If RCON connection fails, check the following:

1. Make sure the server is running: `mcsrvr list`
2. Check if RCON is enabled in server.properties: `enable-rcon=true`
3. Check the RCON port and password in server.properties
4. Try configuring RCON again: `mcsrvr config <server-name> rcon --port 25575 --password mcsrvr`

### Server Crashes

If a server crashes, check the server log for errors:

```bash
mcsrvr log <server-name> --lines 100
```

Common causes of crashes:
- Insufficient memory allocation
- Incompatible plugins or mods
- Corrupted world data

### Process Tracking Issues

If MCSRVR is not correctly tracking server processes, try the following:

1. Stop all servers: `mcsrvr stop <server-name>` for each server
2. Check if any Java processes are still running: `tasklist | findstr java` (Windows) or `ps aux | grep java` (Linux/macOS)
3. Kill any remaining Java processes manually
4. Restart MCSRVR

## Advanced Usage

### Custom Java Installation

By default, MCSRVR uses the Java installation in your PATH. If you want to use a different Java installation, you can specify the full path to the Java executable in the startup script:

1. Initialize the server as usual
2. Edit the startup script: `mcsrvr config <server-name> start`
3. Replace `java` with the full path to the Java executable, e.g., `C:/Program Files/Java/jdk-17/bin/java`

### Multiple Server Instances

MCSRVR can manage multiple server instances. Each server must have a unique name and directory.

```bash
# Initialize multiple servers
mcsrvr init D:/MCServers/Paper/Server1 -n Server1 papermc -v 1.21.4
mcsrvr init D:/MCServers/Paper/Server2 -n Server2 papermc -v 1.21.4
mcsrvr init D:/MCServers/Fabric/Server3 -n Server3 fabric -v 1.21.4

# Start multiple servers
mcsrvr start Server1
mcsrvr start Server2
mcsrvr start Server3

# List all servers
mcsrvr list
```

### Automatic Backups

You can set up automatic backups using your system's task scheduler (Windows) or cron (Linux/macOS).

Windows Task Scheduler example:
1. Open Task Scheduler
2. Create a new task
3. Set the trigger (e.g., daily at 3:00 AM)
4. Set the action to run the following command:
   ```
   mcsrvr backup <server-name> D:/MCBackups
   ```

Cron example (Linux/macOS):
```
0 3 * * * /path/to/mcsrvr backup <server-name> /path/to/backups
```

### Server Migration

To migrate a server to a new machine:

1. Create a backup of the server: `mcsrvr backup <server-name> <backup-path>`
2. Copy the backup to the new machine
3. Install MCSRVR on the new machine
4. Restore the backup: `mcsrvr restore <backup-path> <restore-path>`
5. Start the server: `mcsrvr start <server-name>`
