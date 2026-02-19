:: SPDX-License-Identifier: Apache-2.0
:: Copyright (c) 2025 Qirashi
:: Project: packed_webp


echo off
chcp 65001

echo Модернизация кода...
go fix ./...

echo Начинаю сборку.
go build -o ./out/riff.exe -buildvcs=false -ldflags="-s -w -buildid=" -trimpath -buildmode=exe -tags=release -asmflags="-trimpath" -mod=readonly main_logic.go main.go
if %ERRORLEVEL% neq 0 (
    echo Ошибка: Сборка завершилась с ошибкой. Код ошибки: %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
echo Сборка выполнена успешно.

@REM where ResourceHacker >nul 2>nul
@REM if %errorlevel% == 0 (
@REM     echo Resource Hacker найден, выполняю команды...
@REM     ResourceHacker -open ./out/riff.exe -save ./out/riff.exe -action addoverwrite -res ".\res\ICO.ico" -mask ICONGROUP,MAINICON,
@REM ) else (
@REM     echo Ошибка: Resource Hacker не найден в PATH.
@REM     echo Иконка не установлена.
@REM )

@pause
