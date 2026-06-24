# Windows-only production build. -clean=false is required to bypass
# the RemoveAll+Rename that fails under SearchIndexer/Defender file
# locks. wails3 build runs the full pipeline: tidy -> generate icons
# -> generate bindings -> build frontend -> generate syso -> go build,
# which embeds the Windows icon + version info (from
# build/windows/{icon.ico,wails.exe.manifest,info.json}) and the
# macOS bundle metadata (from build/darwin/Info.plist).

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

Write-Host "[1/1] wails3 build -clean=false..." -ForegroundColor Cyan
& wails3 build -clean=false
if ($LASTEXITCODE -ne 0) { throw "wails3 build failed" }

$exe = Join-Path $root "bin/NetstatUI.exe"
if (Test-Path $exe) {
    $size = [math]::Round((Get-Item $exe).Length / 1MB, 2)
    Write-Host "OK build succeeded: $exe ($size MB)" -ForegroundColor Green
} else {
    throw "exe not produced"
}
