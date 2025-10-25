@echo off
REM PRISM Windows Installer (Command Prompt version)
REM Usage: install.cmd

setlocal EnableDelayedExpansion

set REPO=JohanBellander/prism
set BINARY_NAME=prism.exe

echo Installing PRISM...

REM Determine install location
set INSTALL_DIR=%LOCALAPPDATA%\prism\bin
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo Installing to %INSTALL_DIR%...

REM Create temporary directory
set TMP_DIR=%TEMP%\prism-install-%RANDOM%
mkdir "%TMP_DIR%"
cd /d "%TMP_DIR%"

REM Clone and build
echo Cloning repository...
git clone --depth 1 "https://github.com/%REPO%.git" prism 2>nul
if errorlevel 1 (
    echo Error: Failed to clone repository. Is git installed?
    goto :cleanup
)

cd prism

echo Building PRISM...
where go >nul 2>nul
if errorlevel 1 (
    echo Error: Go is required to build PRISM
    echo Install Go from https://go.dev/doc/install
    goto :cleanup
)

go build -o %BINARY_NAME% ./cmd/prism
if errorlevel 1 (
    echo Error: Build failed
    goto :cleanup
)

move /Y %BINARY_NAME% "%INSTALL_DIR%\"

REM Add to PATH if not already there
echo Checking PATH...
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if errorlevel 1 (
    echo Adding to PATH...
    REM Use PowerShell to add to user PATH to avoid setx 1024 char limit
    powershell -Command "$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if ($currentPath -notlike '*%INSTALL_DIR%*') { [Environment]::SetEnvironmentVariable('Path', $currentPath + ';%INSTALL_DIR%', 'User') }" >nul 2>&1
    if errorlevel 1 (
        echo.
        echo [33mWARNING: Could not automatically add to PATH.[0m
        echo Please manually add this to your PATH:
        echo   %INSTALL_DIR%
        echo.
        echo Or run this command:
        echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
        echo.
    ) else (
        REM Update current session PATH
        set PATH=%PATH%;%INSTALL_DIR%
    )
)

echo.
echo [32mPRISM installed successfully to %INSTALL_DIR%\%BINARY_NAME%[0m
echo.
echo To use PRISM in this session, run:
echo   set PATH=%%PATH%%;%INSTALL_DIR%
echo.
echo Or restart your command prompt.
echo.
echo Run 'prism --help' to get started!

goto :end

:cleanup
cd /d %TEMP%
if exist "%TMP_DIR%" rmdir /S /Q "%TMP_DIR%"
exit /b 1

:end
cd /d %TEMP%
if exist "%TMP_DIR%" rmdir /S /Q "%TMP_DIR%"
