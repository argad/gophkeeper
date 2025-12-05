@echo off
setlocal

:: Set version and build date
set "VERSION=1.0.0"
for /f "tokens=1-4 delims=/ " %%a in ('date /t') do (
    set "BUILD_DATE=%%c-%%a-%%b"
)
if "%BUILD_DATE%"=="unknown" (
    set "BUILD_DATE=%DATE%"
)
for /f "tokens=1-2 delims=: " %%a in ('time /t') do (
    set "BUILD_DATE=%BUILD_DATE%T%%a%%b"
)


echo Building GophKeeper Client...
cd client
go mod tidy
go build -ldflags "-X 'gophkeeper/client/internal/commands.Version=%VERSION%' -X 'gophkeeper/client/internal/commands.BuildDate=%BUILD_DATE%'" -o ../bin/gophkeeper-cli.exe ./cmd/gophkeeper-cli
if %errorlevel% neq 0 (
    echo Client build failed!
    exit /b %errorlevel%
)
cd ..
echo GophKeeper Client built successfully! (./bin/gophkeeper-cli.exe)

echo.
echo Building GophKeeper Server...
cd server
go mod tidy
go build -o ../bin/gophkeeper-server.exe ./cmd/gophkeeper-server
if %errorlevel% neq 0 (
    echo Server build failed!
    exit /b %errorlevel%
)
echo GophKeeper Server built successfully! (./bin/gophkeeper-server.exe)

endlocal
