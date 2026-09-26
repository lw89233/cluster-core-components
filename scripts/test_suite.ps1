Write-Host "Rozpoczynam ujednolicony potok testowy (Test Suite)..." -ForegroundColor Cyan

$packages = @(
    "./pkg/chaos",
    "./pkg/dispatcher",
    "./pkg/liveness",
    "./pkg/reconciler",
    "./pkg/recovery",
    "./pkg/scheduler",
    "./pkg/storage",
    "./pkg/graceful"
)

$failed = $false

foreach ($pkg in $packages) {
    Write-Host "Testowanie $pkg..." -ForegroundColor Yellow
    go test $pkg
    if ($LASTEXITCODE -ne 0) {
        $failed = $true
    }
}

Write-Host "Kompilacja narzedzi konsolowych..." -ForegroundColor Yellow
go build -o inspect.exe ./cmd/inspect
if ($LASTEXITCODE -ne 0) {
    $failed = $true
    Write-Host "Blad kompilacji ./cmd/inspect" -ForegroundColor Red
} else {
    Write-Host "Kompilacja inspect.exe przebiegla pomyslnie." -ForegroundColor Green
}

Write-Host "----------------------------------------"
if ($failed) {
    Write-Host "POTOK ZATRZYMANY: Wykryto bledy w testach lub kompilacji!" -ForegroundColor Red
    exit 1
} else {
    Write-Host "SUKCES: Wszystkie moduly dzialaja poprawnie." -ForegroundColor Green
    exit 0
}