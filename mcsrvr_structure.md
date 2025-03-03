mcsrvr
├── .gitignore
├── DOCUMENTATION.md
├── LICENSE
├── README.md
├── cmd
│   ├── backup.go
│   ├── cmd.go
│   ├── config.go
│   ├── console.go
│   ├── del.go
│   ├── init.go
│   ├── list.go
│   ├── log.go
│   ├── restart.go
│   ├── root.go
│   ├── start.go
│   └── stop.go
├── go.mod
├── go.sum
├── main.go
├── mcsrvr_structure.md
└── pkg
    ├── config
    │   └── config.go
    ├── downloader
    │   └── downloader.go
    └── server
        ├── backup
        │   └── backup.go
        ├── init
        │   └── init.go
        ├── process
        │   ├── process.go
        │   ├── sysprocattr_unix.go
        │   └── sysprocattr_windows.go
        ├── rcon
        │   └── rcon.go
        ├── server.go
        ├── sysprocattr_unix.go
        └── sysprocattr_windows.go
