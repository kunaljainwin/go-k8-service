@echo off
REM Loop from 1 to 5
for /l %%i in (1,1,50) do (
    echo Iteration %%i
    curl http://localhost/api/v1/readyz
)
pause

