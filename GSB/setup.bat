@echo off

if not exist requirements.txt (
    echo requirements.txt not found!
    exit /b 1
)

for /f "delims=" %%i in (requirements.txt) do go get %%i

echo All imports have been installed.
pause
