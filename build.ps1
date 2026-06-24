# Windows-only production build. Runs the Wails 3 pipeline manually
# rather than `wails3 build` so we can:
#   - pass -clean=false to `wails3 generate bindings` to bypass the
#     RemoveAll+Rename of frontend/bindings that fails under
#     SearchIndexer/Defender file locks on local Windows
#   - generate the wails_windows_amd64.syso (icon + version info)
#     explicitly so the binary carries proper Windows metadata
#   - clean up the .syso file after linking
#
# CI uses `wails3 build` directly (no lock contention on the runner).

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

Write-Host "[1/4] wails3 generate bindings (-clean=false to bypass rename)..." -ForegroundColor Cyan
& wails3 generate bindings -ts -clean=false
if ($LASTEXITCODE -ne 0) { throw "bindings generation failed" }

Write-Host "[2/4] npm run build (frontend production)..." -ForegroundColor Cyan
Push-Location frontend
& npm run build -q
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "frontend build failed" }
Pop-Location

Write-Host "[3/4] wails3 generate syso (Windows icon + version info)..." -ForegroundColor Cyan
Push-Location build
& wails3 generate syso -arch amd64 -icon windows/icon.ico -manifest windows/wails.exe.manifest -info windows/info.json -out ../wails_windows_amd64.syso
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "syso generation failed" }
Pop-Location

Write-Host "[4/4] go build (backend production, links .syso)..." -ForegroundColor Cyan
& go build -tags production -trimpath -buildvcs=false -ldflags="-w -s -H windowsgui" -o bin/NetstatUI.exe .
if ($LASTEXITCODE -ne 0) { throw "backend build failed" }

Remove-Item -Path wails_windows_amd64.syso -ErrorAction SilentlyContinue

$exe = Join-Path $root "bin/NetstatUI.exe"
if (Test-Path $exe) {
    $size = [math]::Round((Get-Item $exe).Length / 1MB, 2)
    Write-Host "OK build succeeded: $exe ($size MB)" -ForegroundColor Green
} else {
    throw "exe not produced"
}
