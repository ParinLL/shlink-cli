FROM alpine:3.20

RUN apk add --no-cache ca-certificates curl jq postgresql-client

COPY sync-short-codes.sh /usr/local/bin/sync-short-codes
RUN chmod +x /usr/local/bin/sync-short-codes

USER 65532:65532
ENTRYPOINT ["/usr/local/bin/sync-short-codes"]
