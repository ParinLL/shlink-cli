#!/usr/bin/env sh
set -eu

: "${DATABASE_URL:?DATABASE_URL is required}"
: "${CLOUDFLARE_API_TOKEN:?CLOUDFLARE_API_TOKEN is required}"
: "${CF_ACCOUNT_ID:?CF_ACCOUNT_ID is required}"
: "${CF_KV_NAMESPACE_ID:?CF_KV_NAMESPACE_ID is required}"
: "${SHLINK_DEFAULT_DOMAIN:=parin.dev}"
: "${KV_KEY_PREFIX:=shlink}"

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

cut -f 1 "$ROWS_FILE" | sort -u > "$WORKDIR/domains.txt"

UPDATED=0
UNCHANGED=0
while IFS= read -r domain; do
  [ -n "$domain" ] || continue

  MANIFEST_FILE="$WORKDIR/manifest.json"
  MANIFEST_RAW_FILE="$WORKDIR/manifest.raw.json"
  CURRENT_FILE="$WORKDIR/current.json"
  KV_KEY="$KV_KEY_PREFIX:v2:$(printf '%s' "$domain" | tr '[:upper:]' '[:lower:]')"
  ENCODED_KEY="$(printf '%s' "$KV_KEY" | jq -sRr @uri)"
  VALUE_URL="https://api.cloudflare.com/client/v4/accounts/$CF_ACCOUNT_ID/storage/kv/namespaces/$CF_KV_NAMESPACE_ID/values/$ENCODED_KEY"

  jq -Rsc \
    --arg domain "$domain" \
    '
      split("\n")
      | map(select(length > 0) | split("\t"))
      | map(select((.[0] | ascii_downcase) == ($domain | ascii_downcase)) | .[1])
      | unique
      | {
          version: 2,
          codes: (reduce .[] as $code ({}; .[$code] = true))
        }
    ' "$ROWS_FILE" > "$MANIFEST_RAW_FILE"

  jq -eSc '
    if .version == 2 and (.codes | type) == "object"
    then .
    else error("invalid KV manifest")
    end
  ' "$MANIFEST_RAW_FILE" > "$MANIFEST_FILE"

  if curl -fsS \
    "$VALUE_URL" \
    -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
    -o "$CURRENT_FILE"; then
    jq -Sc . "$CURRENT_FILE" > "$CURRENT_FILE.normalized"
  else
    printf '{}\n' > "$CURRENT_FILE.normalized"
  fi

  if cmp -s "$MANIFEST_FILE" "$CURRENT_FILE.normalized"; then
    UNCHANGED=$((UNCHANGED + 1))
    echo "No KV change for $domain"
    continue
  fi

  curl -fsS \
    "$VALUE_URL" \
    -X PUT \
    -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
    -H "Content-Type: application/json" \
    --data-binary "@$MANIFEST_FILE" >/dev/null

  UPDATED=$((UPDATED + 1))
  CODE_COUNT="$(jq '.codes | length' "$MANIFEST_FILE")"
  echo "Updated $domain manifest with $CODE_COUNT short codes"
done < "$WORKDIR/domains.txt"

echo "Completed Cloudflare KV sync: updated=$UPDATED unchanged=$UNCHANGED"
