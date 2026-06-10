#!/usr/bin/env sh
set -eu

: "${DATABASE_URL:?DATABASE_URL is required}"
: "${CLOUDFLARE_API_TOKEN:?CLOUDFLARE_API_TOKEN is required}"
: "${CF_ACCOUNT_ID:?CF_ACCOUNT_ID is required}"
: "${CF_KV_NAMESPACE_ID:?CF_KV_NAMESPACE_ID is required}"
: "${SHLINK_DEFAULT_DOMAIN:=parin.dev}"
: "${KV_KEY_PREFIX:=shlink}"
: "${KV_EXPIRATION_TTL:=86400}"
: "${CF_KV_BATCH_SIZE:=9000}"

DEFAULT_SQL='
SELECT
  COALESCE(d.authority, :'\''default_domain'\'') AS domain,
  s.short_code
FROM short_urls s
LEFT JOIN domains d ON d.id = s.domain_id
WHERE s.short_code IS NOT NULL
ORDER BY 1, 2;
'

SQL="${SHLINK_SHORT_CODES_SQL:-$DEFAULT_SQL}"
WORKDIR="$(mktemp -d)"
ROWS_FILE="$WORKDIR/short-codes.tsv"

cleanup() {
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

echo "Exporting Shlink short codes from Postgres"
psql "$DATABASE_URL" \
  -v ON_ERROR_STOP=1 \
  -v default_domain="$SHLINK_DEFAULT_DOMAIN" \
  -At \
  -F '	' > "$ROWS_FILE" <<EOF
$SQL
EOF

TOTAL_ROWS="$(wc -l < "$ROWS_FILE" | tr -d ' ')"
if [ "$TOTAL_ROWS" = "0" ]; then
  echo "No short codes returned by query; skipping KV update"
  exit 0
fi

echo "Uploading $TOTAL_ROWS short-code keys to Cloudflare KV"
split -l "$CF_KV_BATCH_SIZE" "$ROWS_FILE" "$WORKDIR/batch."

BATCH_COUNT=0
for batch in "$WORKDIR"/batch.*; do
  BATCH_COUNT=$((BATCH_COUNT + 1))
  BODY_FILE="$WORKDIR/body.$BATCH_COUNT.json"
  RESPONSE_FILE="$WORKDIR/response.$BATCH_COUNT.json"

  jq -Rsc \
    --arg prefix "$KV_KEY_PREFIX" \
    --argjson ttl "$KV_EXPIRATION_TTL" \
    '
      split("\n")
      | map(select(length > 0))
      | map(
          split("\t") as $row
          | {
              key: ($prefix + ":v1:" + ($row[0] | ascii_downcase) + ":" + $row[1]),
              value: "1"
            }
            + if $ttl > 0 then { expiration_ttl: $ttl } else {} end
        )
    ' "$batch" > "$BODY_FILE"

  curl -fsS \
    "https://api.cloudflare.com/client/v4/accounts/$CF_ACCOUNT_ID/storage/kv/namespaces/$CF_KV_NAMESPACE_ID/bulk" \
    -X PUT \
    -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
    -H "Content-Type: application/json" \
    --data-binary "@$BODY_FILE" > "$RESPONSE_FILE"

  jq -e '
    .success == true
    and ((.result.unsuccessful_keys // []) | length == 0)
  ' "$RESPONSE_FILE" >/dev/null

  KEYS_IN_BATCH="$(jq 'length' "$BODY_FILE")"
  echo "Uploaded batch $BATCH_COUNT with $KEYS_IN_BATCH keys"
done

echo "Completed Cloudflare KV short-code sync"
