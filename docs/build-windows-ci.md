# Windows Build in CI

This project builds Windows release artifacts on GitHub Actions using a dedicated
Windows job.

## Why this setup

- The app uses Fyne and OpenGL bindings, which require CGO for desktop builds.
- Building on `windows-latest` avoids fragile Linux-to-Windows CGO cross-toolchain setup in release CI.
- Toolchain installation is explicit in workflow history.

## CI behavior

- `pull_request` and `push` run:
  - Linux verify job
  - Linux build job
  - Windows build job
- Version tag push (`v*`) additionally runs release packaging and publishes:
  - `vrpoker-stats-<tag>-linux-wayland.tar.gz`
  - `vrpoker-stats-<tag>-windows-amd64.zip`

## Windows job toolchain

The workflow installs MinGW GCC via `msys2/setup-msys2` and builds with:

```powershell
$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CC = "C:\msys64\mingw64\bin\gcc.exe"
go build -o vrpoker-stats.exe .
```

## Local reproduction on Linux

Install MinGW and run:

```bash
mise run build-windows
```

This uses:

```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o vrpoker-stats.exe .
```

## Troubleshooting

- `x86_64-w64-mingw32-gcc not found`
  - Install MinGW (`sudo apt-get install -y mingw-w64`) or ensure PATH is correct.
- `build constraints exclude all Go files in .../go-gl/gl`
  - This usually indicates `CGO_ENABLED=0` or missing CGO compiler.
