set appName=game.exe

set rootPath=%~dp0
set targetPath=%rootPath%bin\game

@REM 创建目录结构
xcopy %rootPath%data %targetPath%\data /TEIY
@REM 复制资源文件
xcopy %rootPath%resources %targetPath%\resources /EFIY

@REM 编译
:: 禁用CGO
SET CGO_ENABLED=0
@REM "目标平台是linux、windows"
SET GOOS=windows
:: 目标处理器架构是amd64
SET GOARCH=amd64
go build

@REM 移动执行文件
move %appName% %targetPath%
