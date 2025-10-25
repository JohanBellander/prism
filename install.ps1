# PRISM Windows Installer
# Usage: iwr -useb https://raw.githubusercontent.com/JohanBellander/prism/master/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "JohanBellander/prism"
$BINARY_NAME = "prism.exe"

# Determine install location
$INSTALL_DIR = "$env:LOCALAPPDATA\prism\bin"
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

Write-Host "Installing PRISM to $INSTALL_DIR..." -ForegroundColor Cyan

# Create temporary directory
$TMP_DIR = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
Set-Location $TMP_DIR

try {
    # Clone and build from source
    Write-Host "Cloning repository..." -ForegroundColor Yellow
    git clone --depth 1 "https://github.com/$REPO.git" prism
    Set-Location prism

    # Get version info
    $VERSION = (git describe --tags --always 2>$null) ?? "dev"
    $COMMIT = (git rev-parse --short HEAD 2>$null) ?? "none"
    $DATE = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

    Write-Host "Building PRISM version $VERSION..." -ForegroundColor Yellow
    if (Get-Command go -ErrorAction SilentlyContinue) {
        go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" -o $BINARY_NAME ./cmd/prism
        Move-Item $BINARY_NAME $INSTALL_DIR -Force
    } else {
        Write-Host "Error: Go is required to build PRISM" -ForegroundColor Red
        Write-Host "Install Go from https://go.dev/doc/install" -ForegroundColor Red
        exit 1
    }
} finally {
    # Cleanup
    Set-Location ~
    Remove-Item $TMP_DIR -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "✅ PRISM installed successfully to $INSTALL_DIR\$BINARY_NAME" -ForegroundColor Green
Write-Host ""

# Check if install dir is in PATH
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$INSTALL_DIR*") {
    Write-Host "⚠️  Adding $INSTALL_DIR to your PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$INSTALL_DIR", "User")
    $env:Path = "$env:Path;$INSTALL_DIR"
    Write-Host "✅ PATH updated. Restart your terminal for changes to take effect." -ForegroundColor Green
}

Write-Host ""
Write-Host "Run 'prism --help' to get started!" -ForegroundColor Cyan
