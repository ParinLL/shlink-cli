import assert from "node:assert/strict";
import {
  handleRequest,
  classifyPath,
  buildLookupKeys,
  buildManifestKey,
} from "./src/index.mjs";

const originFetch = globalThis.fetch;
globalThis.fetch = async () => new Response("origin", { status: 200 });

function envWithCodes(codes) {
  const values = new Set(codes);
  return {
    EDGE_GUARD_HOSTS: "parin.dev",
    KV_KEY_PREFIX: "shlink",
    SHORT_CODES: {
      async get(key, options) {
        if (options?.type === "json") {
          return null;
        }
        return values.has(key) ? "1" : null;
      },
    },
  };
}

function envWithManifest(codes) {
  const manifest = {
    version: 2,
    codes: Object.fromEntries(codes.map((code) => [code, true])),
  };
  return {
    EDGE_GUARD_HOSTS: "parin.dev",
    KV_KEY_PREFIX: "shlink",
    SHORT_CODES: {
      async get(key, options) {
        if (key === "shlink:v2:parin.dev" && options?.type === "json") {
          return manifest;
        }
        return null;
      },
    },
  };
}

try {
  assert.equal(classifyPath("/").action, "pass");
  assert.equal(classifyPath("/rest/v3/short-urls").action, "pass");
  assert.equal(classifyPath("/rest/v3/short-urls", { EDGE_GUARD_RESERVED_PREFIXES: "" }).action, "check");
  assert.equal(classifyPath("/favicon.ico", { EDGE_GUARD_RESERVED_PATHS: "" }).action, "check");
  assert.deepEqual(classifyPath("/abc").candidates, ["abc"]);
  assert.deepEqual(classifyPath("/abc/extra").candidates, ["abc/extra", "abc"]);
  assert.deepEqual(buildLookupKeys("parin.dev", ["abc"], {}).at(0), "shlink:v1:parin.dev:abc");
  assert.equal(buildManifestKey("parin.dev", {}), "shlink:v2:parin.dev");

  let response = await handleRequest(
    new Request("https://parin.dev/abc"),
    envWithManifest(["abc"]),
  );
  assert.equal(response.status, 200);
  assert.equal(await response.text(), "origin");

  response = await handleRequest(
    new Request("https://parin.dev/abc/anything"),
    envWithManifest(["abc"]),
  );
  assert.equal(response.status, 200);

  response = await handleRequest(
    new Request("https://parin.dev/legacy"),
    envWithCodes(["shlink:v1:parin.dev:legacy"]),
  );
  assert.equal(response.status, 200);

  response = await handleRequest(
    new Request("https://parin.dev/api"),
    envWithCodes([]),
  );
  assert.equal(response.status, 404);
  assert.equal(response.headers.get("x-shlink-edge-guard"), "blocked");

  response = await handleRequest(
    new Request("https://shlink-dashboard.x300-local.parinll.com/server"),
    envWithCodes([]),
  );
  assert.equal(response.status, 200);
} finally {
  globalThis.fetch = originFetch;
}

console.log("ok");
