@echo off
setlocal EnableDelayedExpansion

echo Ejecutando pruebas unitarias para Stock Analyzer...

set GO_TEST_FLAGS=-v
set TOTAL_TESTS=0
set PASSED_TESTS=0
set FAILED_TESTS=0
set FAILED_PACKAGES=

echo.
echo ================ PRUEBAS DE CLIENTE API ================
go test %GO_TEST_FLAGS% ./tests/client/
if %errorlevel% equ 0 (
    set /a PASSED_TESTS+=1
) else (
    set /a FAILED_TESTS+=1
    set FAILED_PACKAGES=!FAILED_PACKAGES! client
)
set /a TOTAL_TESTS+=1

echo.
echo ================ PRUEBAS DE SERVICIOS ================
go test %GO_TEST_FLAGS% ./tests/services/
if %errorlevel% equ 0 (
    set /a PASSED_TESTS+=1
) else (
    set /a FAILED_TESTS+=1
    set FAILED_PACKAGES=!FAILED_PACKAGES! services
)
set /a TOTAL_TESTS+=1

echo.
echo ================ PRUEBAS DE REPOSITORIOS ================
REM Ejecutar con patrón para hacer matching flexible con expresiones regulares
go test %GO_TEST_FLAGS% -count=1 ./tests/repositories/
if %errorlevel% equ 0 (
    set /a PASSED_TESTS+=1
) else (
    set /a FAILED_TESTS+=1
    set FAILED_PACKAGES=!FAILED_PACKAGES! repositories
)
set /a TOTAL_TESTS+=1

echo.
echo ================ PRUEBAS DE CONTROLADORES ================
go test %GO_TEST_FLAGS% ./tests/routes/
if %errorlevel% equ 0 (
    set /a PASSED_TESTS+=1
) else (
    set /a FAILED_TESTS+=1
    set FAILED_PACKAGES=!FAILED_PACKAGES! routes
)
set /a TOTAL_TESTS+=1

echo.
echo ================ RESUMEN DE PRUEBAS ================
echo Total de paquetes probados: %TOTAL_TESTS%
echo Paquetes exitosos: %PASSED_TESTS%
echo Paquetes fallidos: %FAILED_TESTS%

if %FAILED_TESTS% gtr 0 (
    echo.
    echo ================ PAQUETES CON FALLOS ================
    echo Los siguientes paquetes tienen pruebas fallidas:!FAILED_PACKAGES!
    
    echo.
    echo ================ SUGERENCIAS DE RESOLUCIÓN ================
    if "!FAILED_PACKAGES!" == " repositories" (
        echo Para repositorios: Verifica que las consultas SQL incluyan la condición "stocks"."deleted_at" IS NULL
        echo Modifica TestSearchStocks para usar la consulta SQL completa con paréntesis y soft delete.
    )
    
    if "!FAILED_PACKAGES!" == " routes" (
        echo Para controladores: Verifica que las pruebas inicialicen correctamente los mocks y que
        echo los métodos mock retornen valores del tipo correcto.
    )
    
    exit /b 1
) else (
    echo.
    echo ================ ÉXITO ================
    echo Todas las pruebas han pasado correctamente.
)

echo.
echo Ejecución de pruebas completada.