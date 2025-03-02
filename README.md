# MCSRVR - Minecraft Server Manager

MCSRVR is a command-line tool for easily setting up and managing Minecraft servers. It automates the process of downloading server files, configuring servers, and provides a systemctl-like interface for managing server processes.

## Features

- **Easy Server Setup**: Initialize vanilla, PaperMC, and Fabric servers with a single command
- **Server Management**: Start, stop, restart, and monitor your Minecraft servers
- **Console Access**: Access server console and execute commands remotely
- **Process Management**: Servers run as hidden processes, similar to systemctl in Linux
- **Backup & Restore**: Create and restore server backups
- **Configuration Management**: Easily configure server properties and settings
- **Multi-Server Support**: Manage multiple Minecraft servers from one interface

## Installation

### Prerequisites

- Windows, macOS, or Linux
- Go 1.16 or higher (for building from source or using Go installation)
- Java 17 or higher (for running Minecraft servers)

### From Binary

1. Download the latest release for your platform from the [Releases](https://github.com/0v3rr1de0/mcsrvr/releases) page.
2. Extract the archive.
3. Add the extracted directory to your PATH (optional).

### From Go (Recommended)

If you have Go installed, you can install MCSRVR directly from the source. This method downloads and compiles the application on your machine:

```bash
go install github.com/0v3rr1de0/mcsrvr@latest
```

Alternatively, to build from source:

```bash
# Clone the repository
git clone https://github.com/0v3rr1de0/mcsrvr.git
cd mcsrvr

# Build the project
go build

# (Optional) Install the binary to your Go bin path
go install
```

### Planned Package Manager Support

Future releases will include precompiled packages for easier installation through popular package managers:

- **Chocolatey** for Windows
- **Winget** for Windows
- **APT/Snapcraft** for Debian/Ubuntu-based Linux distributions
- **packman** for Arch-based Linux distributions
- **Homebrew** for macOS

Stay tuned for updates!

## Quick Start

### Initialize a New Server

```bash
# Initialize a PaperMC server
mcsrvr init D:/MCServers/Paper/MyServer -n MyServer papermc -v 1.21.4

# Initialize a vanilla server
mcsrvr init D:/MCServers/Vanilla/MyServer -n MyServer vanilla -v 1.21.4

# Initialize a Fabric server
mcsrvr init D:/MCServers/Fabric/MyServer -n MyServer fabric -v 1.21.4
```

### Manage Servers

```bash
# Start a server
mcsrvr start MyServer

# Stop a server
mcsrvr stop MyServer

# Restart a server
mcsrvr restart MyServer

# List all servers
mcsrvr list

# List only online servers
mcsrvr list --online

# Show server logs
mcsrvr log MyServer --lines <number of wanted lines eg. 50>

# Show server logs continously
mcsrvr log MyServer --follow

# Access server console
mcsrvr console MyServer

# Execute a command on a server
mcsrvr cmd MyServer "say Hello, world!"
```

### Backup and Restore (Not working currently)

```bash
# Create a backup
mcsrvr backup MyServer C:/MCBackups

# List backups
mcsrvr backups MyServer

# Restore a backup
mcsrvr restore C:/MCBackups/MyServer_2025-03-02_12-34-56 D:/MCServers/Restored
```

## Supported Server Types

- **Vanilla**: Official Minecraft server
- **PaperMC**: High-performance fork of Spigot
- **Fabric**: Lightweight, modular modding toolchain

Coming soon:
- Spigot
- Bukkit
- Forge
- Velocity
- BungeeCord
- Cuberite

## Configuration

MCSRVR stores server configurations in `~/.mcsrvr/config.json`. You can edit this file directly or use the `config` command:

```bash
# Set default memory allocation for new servers
mcsrvr config --default-memory 4G

# Set default Java arguments for new servers
mcsrvr config --default-java-args "-XX:+UseG1GC -XX:+ParallelRefProcEnabled"

# Configure RCON settings for a server
mcsrvr config MyServer rcon --port 25575 --password mypassword
```

## Documentation

- For detailed documentation, see [DOCUMENTATION.md](DOCUMENTATION.md).
- For detailed file structure, see [mcsrvr_structure.md](mcsrvr_structure.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgements

- [PaperMC](https://papermc.io/) for their high-performance Minecraft server
- [Fabric](https://fabricmc.net/) for their modding toolchain
- [Mojang](https://www.mojang.com/) for Minecraft