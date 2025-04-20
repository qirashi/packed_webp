:: SPDX-License-Identifier: Apache-2.0
:: Copyright (c) 2025 Qirashi
:: Project: packed_webp

echo off
chcp 65001

echo Начинаю сборку .exe
go build -o ./riff.exe -ldflags="-s -w" -trimpath -buildmode=exe -tags=release -asmflags="-trimpath" -mod=readonly ^
main_logic.go main.go
if %ERRORLEVEL% neq 0 (
    echo Ошибка: Сборка завершилась с ошибкой. Код ошибки: %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
echo .exe Успешно собран.

set TOOL1=%cd%\res\upx.exe
if defined UPX (
    set "UPX="
)
if exist "%TOOL1%" (
    echo UPX найден, выполняю команды...
    "%TOOL1%" -9 "%cd%\riff.exe"
) else (
    echo Ошибка: UPX не найден по пути "%TOOL1%".
	echo Exe не сжат UPX.
)

@pause
