@echo off
echo Downloading portable Tesseract...

mkdir tesseract-portable
cd tesseract-portable

REM Download portable Tesseract
curl -L -o tesseract.zip https://digi.bib.uni-mannheim.de/tesseract/tesseract-ocr-w64-setup-5.3.3.20231005.exe

echo.
echo Tesseract will be bundled with the app
echo Users won't need to install anything!
