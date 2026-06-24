#Requires -Version 5.1
<#
.SYNOPSIS
    Download ip2region xdb data files into the data/ directory.

.DESCRIPTION
    Fetches ip2region_v4.xdb and ip2region_v6.xdb from the official
    lionsoul2014/ip2region repository on GitHub. These files are required
    by the services/geo package for IP -> geographic location lookups.

.NOTES
    Re-run this script whenever you want to refresh the offline database.
#>

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$dataDir = Join-Path $root "..\data"
$dataDir = [System.IO.Path]::GetFullPath($dataDir)

if (-not (Test-Path $dataDir)) {
    New-Item -ItemType Directory -Path $dataDir -Force | Out-Null
}

$base = "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data"
$files = @(
    @{ Name = "ip2region_v4.xdb"; SizeMB = 11 },
    @{ Name = "ip2region_v6.xdb"; SizeMB = 25 }
)

foreach ($f in $files) {
    $name = $f.Name
    $dest = Join-Path $dataDir $name
    if (Test-Path $dest) {
        $sizeMB = [math]::Round((Get-Item $dest).Length / 1MB, 1)
        Write-Host "[skip] $name already present ($sizeMB MB)" -ForegroundColor DarkGray
        continue
    }
    $url = "$base/$name"
    Write-Host "[get ] $url" -ForegroundColor Cyan
    try {
        $ProgressPreference = "SilentlyContinue"
        Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
        $sizeMB = [math]::Round((Get-Item $dest).Length / 1MB, 1)
        Write-Host "[done] $name ($sizeMB MB)" -ForegroundColor Green
    } catch {
        Write-Host "[fail] $name : $_" -ForegroundColor Red
        if (Test-Path $dest) { Remove-Item $dest -Force }
        throw
    }
}

Write-Host ""
Write-Host "All ip2region data files are ready in $dataDir" -ForegroundColor Green
