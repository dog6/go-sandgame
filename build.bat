@echo off
echo Building sand-game_win64-amd64..
go build -o ./builds/sand-game_win-amd64.exe
echo Done.
echo.

:: I think this enables the gpu. ¯\_(ツ)_/¯
echo Building sand-game_win64-amd64_gpu
set NvOptimusEnablement=1
set AmdPowerXpressRequestHighPerformance=1
go build -o ./builds/sand-game_win-amd64_gpu.exe
echo Done.
timeout /t 3 >NUL