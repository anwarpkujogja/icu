Write-Output "--- VERIFYING ICU DATABASE CONTENT ---"
Write-Output "Container: icu-db-container"
Write-Output "Database:  icu_db"
Write-Output "Table:     hasil_lab"
Write-Output "----------------------------------------"

# 1. Check if table exists and count rows
Write-Output "`n[1] Checking Row Count..."
docker exec -i icu-db-container psql -U postgres -d icu_db -c "SELECT count(*) as total_rows FROM hasil_lab;"

# 2. Dump Content
Write-Output "`n[2] Dumping First 5 Rows..."
docker exec -i icu-db-container psql -U postgres -d icu_db -c "SELECT id, created_at, result_data FROM hasil_lab LIMIT 5;"

Write-Output "`n[3] Checking Tables List (to insure we are in right DB)..."
docker exec -i icu-db-container psql -U postgres -d icu_db -c "\dt"

Write-Output "`nDone."
