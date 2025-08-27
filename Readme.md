# 1) Start Postgres and create db `otpdb`.
# 2) Apply migration migrations/001_create_otps.sql.
# 3) cp .env.example .env and set DATABASE_URL.
# 4) go mod tidy && go run ./cmd/api.
# 5) Try:
# curl -s -X POST http://localhost:8098/api/v1/otp \
# -H 'Content-Type: application/json' \
# -d '{"tenant_id":"itdyaingenieria","channel":"sms","destination":"+15551234567"}' | jq
# # take the returned id and code from logs, then validate:
# curl -s -X POST http://localhost:8098/api/v1/otp/validate \
# -H 'Content-Type: application/json' \
# -d '{"tenant_id":"itdyaingenieria","id":"<id-here>","code":"<code-from-log>"}' | jq