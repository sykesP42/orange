$SERVER = "root@123.56.99.154"
$CONTAINER = "orange-frontend"
$HTML_DIR = "/usr/share/nginx/html"
$TEMP_ARCHIVE = "/tmp/orange-frontend.tar.gz"

Write-Host "Creating archive..." -ForegroundColor Green
Set-Location "c:\Users\24311\Documents\coding\orange"
tar -czf $TEMP_ARCHIVE -C dist .

Write-Host "Uploading archive ($([math]::Round((Get-Item $TEMP_ARCHIVE).Length / 1MB, 2)) MB)..." -ForegroundColor Yellow
scp $TEMP_ARCHIVE "${SERVER}:${TEMP_ARCHIVE}"

Write-Host "Extracting on server..." -ForegroundColor Cyan
ssh $SERVER "docker exec ${CONTAINER} sh -c 'rm -rf ${HTML_DIR}/* && tar -xzf ${TEMP_ARCHIVE} -C ${HTML_DIR} && rm ${TEMP_ARCHIVE}'"

Write-Host "Reloading Nginx..." -ForegroundColor Green
ssh $SERVER "docker exec ${CONTAINER} nginx -s reload"

Remove-Item $TEMP_ARCHIVE -Force
Write-Host "Deployment complete!" -ForegroundColor Green
