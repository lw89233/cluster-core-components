Write-Host "Rozpoczynam ujednolicony potok testowy (Test Suite)..." -ForegroundColor Cyan

$packages = @(
    "./pkg/chaos",
    "./pkg/dispatcher",
    "./pkg/liveness",
    "./pkg/reconciler",
    "./pkg/recovery",
    "./pkg/scheduler",
    "./pkg/storage"
)

$failed = $false

foreach ($pkg in $packages) {
    Write-Host "Testowanie $pkg..." -ForegroundColor Yellow
    go test $pkg
    if ($LASTEXITCODE -ne 0) {
        $failed = $true
    }
}

Write-Host "----------------------------------------"
if ($failed) {
    Write-Host "POTOK ZATRZYMANY: Wykryto bledy w testach!" -ForegroundColor Red
    exit 1
} else {
    Write-Host "SUKCES: Wszystkie moduly dzialaja poprawnie." -ForegroundColor Green
    exit 0
}