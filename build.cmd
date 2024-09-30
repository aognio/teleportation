@echo off
REM Navigate to the project root directory
cd /d "%~dp0"

REM Build the Go project, specifying the output binary
go build -o tlp2-bin.exe ./cmd/tlp2-bin

REM Check if build succeeded and provide feedback
if %errorlevel% == 0 (
    echo Build successful: tlp2-bin.exe created.
) else (
    echo Build failed.
)

