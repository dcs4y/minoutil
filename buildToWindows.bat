set appName=game.exe

set rootPath=%~dp0
set targetPath=%rootPath%bin\game

@REM ����Ŀ¼�ṹ
xcopy %rootPath%data %targetPath%\data /TEIY
@REM ������Դ�ļ�
xcopy %rootPath%resources %targetPath%\resources /EFIY

@REM ����
:: ����CGO
SET CGO_ENABLED=0
@REM "Ŀ��ƽ̨��linux��windows"
SET GOOS=windows
:: Ŀ�괦�����ܹ���amd64
SET GOARCH=amd64
go build

@REM �ƶ�ִ���ļ�
move %appName% %targetPath%
