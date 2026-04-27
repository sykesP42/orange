$SERVER = "root@123.56.99.154"
$CONTAINER = "orange-frontend"
$HTML_DIR = "/usr/share/nginx/html"
$TAR_FILE = "c:\Users\24311\Documents\coding\orange\tar.gz"
$TAR_REMOTE = "/tmp/orange-frontend.tar.gz"

Set-Location "c:\Users\24311\Documents\coding\orange"

Write-Host "Creating tar archive..." -ForegroundColor Green
tar -czf $TAR_FILE -C dist .

$size = [math]::Round((Get-Item $TAR_FILE).Length / 1MB, 2)
Write-Host "Uploading archive ($size MB)..." -ForegroundColor Yellow
scp $TAR_FILE "${SERVER}:${TAR_REMOTE}"

Write-Host "Extracting on server..." -ForegroundColor Cyan
ssh $SERVER "docker cp ${TAR_REMOTE} ${CONTAINER}:${TAR_REMOTE}"
ssh $SERVER "docker exec ${CONTAINER} sh -c 'rm -rf ${HTML_DIR}/* && cd ${HTML_DIR} && tar -xzf ${TAR_REMOTE} && rm ${TAR_REMOTE}'"

Write-Host "Reloading Nginx..." -ForegroundColor Green
ssh $SERVER "docker exec ${CONTAINER} nginx -s reload"

Remove-Item $TAR_FILE -Force
Write-Host "Deployment complete!" -ForegroundColor Green
