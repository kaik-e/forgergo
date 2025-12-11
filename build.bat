@echo off
echo ========================================
echo   FORGER COMPANION GO - BUILD
echo ========================================
echo.

echo [*] Installing dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to download dependencies
    pause
    exit /b 1
)

echo.
echo [*] Building executable...
go build -ldflags="-s -w -H windowsgui" -o ForgerCompanion.exe .
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Build failed
    pause
    exit /b 1
)

echo.
echo ========================================
echo   BUILD SUCCESSFUL!
echo ========================================
echo.

for %%F in (ForgerCompanion.exe) do (
    set size=%%~zF
    set /a sizeMB=!size! / 1048576
    echo Output: ForgerCompanion.exe
    echo Size: !sizeMB! MB
)

echo.
echo To further compress, install UPX and run:
echo   upx --best --lzma ForgerCompanion.exe
echo.

pause
