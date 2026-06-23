# wails3 build on Windows often fails with "Access is denied" during
# frontend/bindings RemoveAll+Rename (Windows SearchIndexer / Defender locks).
# This script uses -clean=false to overwrite in-place and bypass the rename.
# Flags must match wails3 build's internal generate call exactly,
# except -clean=true -> -clean=false.

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

Write-Host "[1/3] wails3 generate bindings (-clean=false to bypass rename)..." -ForegroundColor Cyan
& wails3 generate bindings -ts -clean=false
if ($LASTEXITCODE -ne 0) { throw "bindings generation failed" }

Write-Host "[2/3] npm run build (frontend production)..." -ForegroundColor Cyan
Push-Location frontend
& npm run build -q
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "frontend build failed" }
Pop-Location

Write-Host "[3/3] go build (backend production)..." -ForegroundColor Cyan
& go build -tags production -ldflags="-w -s -H windowsgui" -o bin/NetstatUI.exe .
if ($LASTEXITCODE -ne 0) { throw "backend build failed" }

$exe = Join-Path $root "bin/NetstatUI.exe"
if (Test-Path $exe) {
    $size = [math]::Round((Get-Item $exe).Length / 1MB, 2)
    Write-Host "OK build succeeded: $exe ($size MB)" -ForegroundColor Green
} else {
    throw "exe not produced"
}
