<#
.SYNOPSIS
  Levanta el proyecto completo (API Go + UI React) con Docker.
.DESCRIPTION
  Construye las imágenes del backend y del frontend y las ejecuta a la vez.
  El frontend se hornea con VITE_API_BASE_URL apuntando al backend, porque
  el navegador llama a la API directamente (no hay proxy entre contenedores).
.EXAMPLE
  ./run.ps1
  Levanta todo: backend en http://localhost:8080 y frontend en http://localhost:3000.
.EXAMPLE
  ./run.ps1 -Action down
  Detiene y elimina ambos contenedores.
.EXAMPLE
  ./run.ps1 -FrontendPort 3001 -BackendPort 8081
  Usa otros puertos de host (la URL de la API se ajusta sola).
.EXAMPLE
  ./run.ps1 -SkipBuild
  Reutiliza las imágenes ya construidas (arranque rápido).
#>
[CmdletBinding()]
param(
  [ValidateSet('up', 'down', 'restart', 'logs', 'build')]
  [string]$Action = 'up',

  [int]$FrontendPort = 3000,
  [int]$BackendPort = 8080,

  [string]$ApiUrl = '',

  [switch]$SkipBuild
)

$ErrorActionPreference = 'Stop'

$BackendImage = 'seezle/calc-api'
$FrontendImage = 'seezle/front-calculator'
$BackendContainer = 'seezle-calc-api'
$FrontendContainer = 'seezle-front-calculator'
$BackendInternalPort = 8080
$Root = $PSScriptRoot

if ([string]::IsNullOrWhiteSpace($ApiUrl)) {
  $ApiUrl = "http://localhost:$BackendPort"
}

function Test-Docker {
  if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    throw 'Docker no está disponible en el PATH. Instálalo e inícialo antes de continuar.'
  }
  docker info 2>$null | Out-Null
  if ($LASTEXITCODE -ne 0) {
    throw 'El daemon de Docker no responde. Inícialo antes de continuar.'
  }
}

function Assert-ImagesExist {
  foreach ($image in @($BackendImage, $FrontendImage)) {
    docker image inspect $image 2>$null | Out-Null
    if ($LASTEXITCODE -ne 0) {
      throw "No existe la imagen '$image'. Ejecuta sin -SkipBuild para construirla."
    }
  }
}

function Build-Images {
  Write-Host "==> Construyendo $BackendImage ..." -ForegroundColor Cyan
  docker build -t $BackendImage (Join-Path $Root 'backend')

  Write-Host "==> Construyendo $FrontendImage (VITE_API_BASE_URL=$ApiUrl) ..." -ForegroundColor Cyan
  docker build --build-arg "VITE_API_BASE_URL=$ApiUrl" -t $FrontendImage (Join-Path $Root 'frontend/front-calculator')
}

function Remove-ContainerIfExists([string]$Name) {
  $existing = docker ps -a --filter ('name=^{0}$' -f $Name) --format '{{.Names}}'
  if ($existing -eq $Name) {
    Write-Host "Eliminando contenedor previo '$Name'..." -ForegroundColor Yellow
    docker rm -f $Name | Out-Null
  }
}

function Wait-Healthy([string]$Url, [string]$Label, [int]$TimeoutSec = 90) {
  $deadline = (Get-Date).AddSeconds($TimeoutSec)
  while ((Get-Date) -lt $deadline) {
    try {
      $res = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 5
      if ($res.StatusCode -ge 200 -and $res.StatusCode -lt 400) {
        Write-Host "$Label OK: $Url" -ForegroundColor Green
        return
      }
    }
    catch {
      Start-Sleep -Seconds 2
    }
  }
  throw "$Label no respondió en ${TimeoutSec}s: $Url"
}

function Invoke-Up {
  Test-Docker
  if ($SkipBuild) {
    Assert-ImagesExist
  }
  else {
    Build-Images
  }

  Remove-ContainerIfExists $BackendContainer
  docker run -d --name $BackendContainer -p "${BackendPort}:${BackendInternalPort}" $BackendImage | Out-Null

  Remove-ContainerIfExists $FrontendContainer
  docker run -d --name $FrontendContainer -p "${FrontendPort}:80" $FrontendImage | Out-Null

  Wait-Healthy "http://localhost:$BackendPort/healthz" 'Backend'
  Wait-Healthy "http://localhost:$FrontendPort/" 'Frontend'

  Write-Host ''
  Write-Host 'Proyecto arriba:' -ForegroundColor Green
  Write-Host "  Frontend: http://localhost:$FrontendPort/"
  Write-Host "  Backend:  http://localhost:$BackendPort/healthz"
  Write-Host "  API usada por la UI: $ApiUrl"
}

function Invoke-Down {
  Test-Docker
  foreach ($c in @($FrontendContainer, $BackendContainer)) {
    $existing = docker ps -a --filter ('name=^{0}$' -f $c) --format '{{.Names}}'
    if ($existing -eq $c) {
      docker rm -f $c | Out-Null
      Write-Host "Detenido y eliminado: $c"
    }
    else {
      Write-Host "No existe: $c"
    }
  }
}

function Invoke-Logs {
  Test-Docker
  Write-Host "----- $BackendContainer -----"
  docker logs --tail 100 $BackendContainer
  Write-Host "----- $FrontendContainer -----"
  docker logs --tail 100 $FrontendContainer
}

switch ($Action) {
  'up' { Invoke-Up }
  'down' { Invoke-Down }
  'restart' { Invoke-Down; Invoke-Up }
  'logs' { Invoke-Logs }
  'build' { Test-Docker; Build-Images }
}
