# Orange 部署打包脚本 - 只打包 Docker 构建需要的文件
# 使用方法: powershell -ExecutionPolicy Bypass -File deploy-pack.ps1

$ErrorActionPreference = "Stop"
$ProjectRoot = "C:\Users\24311\Documents\coding\orange"
$ServerIP = "123.56.99.154"
$OutputFile = "$ProjectRoot\orange-deploy.tar.gz"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Orange Docker 部署打包工具" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 排除规则（和 .dockerignore 保持一致）
$Excludes = @(
    "node_modules",
    "dist",
    ".git",
    ".gitignore",
    ".vscode",
    ".idea",
    ".trae",
    "*.md",
    "*.log",
    "ssl/",
    "data/",
    "uploads/",
    ".github/",
    "tests/",
    "test/",
    "e2e/",
    "playwright.config.*",
    "vitest.config.*",
    "apps/",
    "docs/",
    "scripts/",
    "deploy/",
    "*.sh",
    "*.bat",
    "*.exe*",
    "*.backup",
    "UTF8",
    "update.md",
    "genhash.js",
    "zeabur.json",
    "nginx-config.conf",
    "nginx-docker.conf",
    "docker-compose.local.yml",
    "check-backend.sh",
    "find-project.sh",
    "fix.sh"
)

Write-Host "[1/4] 清理旧的打包文件..." -ForegroundColor Yellow
if (Test-Path $OutputFile) { Remove-Item $OutputFile -Force }

Write-Host "[2/4] 打包必要文件..." -ForegroundColor Yellow

Push-Location $ProjectRoot

# 用 tar 打包，排除不需要的文件
$ExcludeArgs = $Excludes | ForEach-Object { "--exclude=$_" }
& tar czf $OutputFile @ExcludeArgs .

Pop-Location

$SizeMB = [math]::Round((Get-Item $OutputFile).Length / 1MB, 2)
Write-Host "[3/4] 打包完成: $OutputFile ($SizeMB MB)" -ForegroundColor Green

Write-Host ""
Write-Host "[4/4] 上传到服务器 $ServerIP ..." -ForegroundColor Yellow
Write-Host ""
Write-Host "请在服务器上执行以下命令:" -ForegroundColor White
Write-Host "----------------------------------------" -ForegroundColor Gray
Write-Host "  ssh root@$ServerIP" -ForegroundColor White
Write-Host "  mkdir -p /www/orange/data /www/orange/uploads" -ForegroundColor White
Write-Host "  cd /www && rm -rf orange && mkdir orange" -ForegroundColor White
Write-Host "" -ForegroundColor Gray
Write-Host "然后在本地执行上传命令:" -ForegroundColor White
Write-Host "----------------------------------------" -ForegroundColor Gray
Write-Host "  scp `"$OutputFile`" root@${ServerIP}:/www/" -ForegroundColor White
Write-Host "" -ForegroundColor Gray
Write-Host "再在服务器上解压并启动:" -ForegroundColor White
Write-Host "----------------------------------------" -ForegroundColor Gray
Write-Host "  cd /www/orange && tar xzf ../orange-deploy.tar.gz --strip-components=1" -ForegroundColor White
Write-Host "  echo `"JWT_SECRET=OrangeProd$(openssl rand -hex 16)`" > .env" -ForegroundColor White
Write-Host "  docker-compose up -d --build" -ForegroundColor White
Write-Host "----------------------------------------" -ForegroundColor Gray
Write-Host ""

Write-Host "是否立即执行上传？(Y/N)" -ForegroundColor Yellow
$Confirm = Read-Host
if ($Confirm -eq 'Y' -or $Confirm -eq 'y') {
    Write-Host "正在上传..." -ForegroundColor Yellow
    & scp $OutputFile "root@${ServerIP}:/www/"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "上传成功！" -ForegroundColor Green
        Write-Host ""
        Write-Host "现在 SSH 登录服务器执行解压和启动命令即可。" -ForegroundColor Cyan
    } else {
        Write-Host "上传失败，请检查网络连接和 SSH 配置。" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  打包大小: $SizeMB MB" -ForegroundColor Cyan
Write-Host "  打包文件: $OutputFile" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
