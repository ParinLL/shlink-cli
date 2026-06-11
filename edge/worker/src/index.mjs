const DEFAULT_HOSTS = "parin.dev";
const DEFAULT_KV_PREFIX = "shlink";
const DEFAULT_RESERVED_PREFIXES = [
  ".well-known",
  "rest",
  "server",
  "admin",
  "health",
  "assets",
  "build",
  "css",
  "js",
  "img",
  "images",
  "static",
  "vendor",
];
const DEFAULT_RESERVED_PATHS = [
  "/favicon.ico",
  "/robots.txt",
  "/sitemap.xml",
];

export default {
  async fetch(request, env, ctx) {
    return handleRequest(request, env, ctx);
  },
};

export async function handleRequest(request, env) {
  const url = new URL(request.url);
  const host = url.hostname.toLowerCase();

  if (!isProtectedHost(host, env)) {
    return fetch(request);
  }

  const decision = classifyPath(url.pathname, env);
  if (decision.action === "pass") {
    return fetch(request);
  }

  try {
    if (await hasAnyShortCode(env.SHORT_CODES, host, decision.candidates, env)) {
      return fetch(request);
    }
  } catch (error) {
    console.error("shlink-edge-guard kv lookup failed", {
      message: error?.message || String(error),
      path: url.pathname,
      host,
    });

    if (boolEnv(env.EDGE_GUARD_FAIL_OPEN, true)) {
      return fetch(request);
    }

    return new Response("Edge guard unavailable\n", {
      status: 503,
      headers: {
        "cache-control": "no-store",
        "x-shlink-edge-guard": "kv-error",
      },
    });
  }

  console.warn("shlink-edge-guard blocked orphan candidate", {
    host,
    path: url.pathname,
    method: request.method,
    cfRay: request.headers.get("cf-ray"),
    clientIp: request.headers.get("cf-connecting-ip"),
  });

  return new Response("Not found\n", {
    status: intEnv(env.EDGE_GUARD_BLOCK_STATUS, 404),
    headers: {
      "cache-control": "no-store",
      "content-type": "text/plain; charset=utf-8",
      "x-shlink-edge-guard": "blocked",
    },
  });
}

export function classifyPath(pathname, env = {}) {
  if (pathname === "/" || pathname === "") {
    return { action: "pass", reason: "root" };
  }

  if (reservedPaths(env).has(pathname)) {
    return { action: "pass", reason: "reserved-path" };
  }

  const cleanPath = pathname.replace(/^\/+/, "");
  const firstSegment = cleanPath.split("/")[0] || "";
  if (reservedPrefixes(env).has(firstSegment.toLowerCase())) {
    return { action: "pass", reason: "reserved-prefix" };
  }

  const candidate = safeDecode(cleanPath).replace(/\/+$/, "");
  if (!candidate || candidate === "." || candidate === "..") {
    return { action: "block", candidates: [] };
  }

  const candidates = [candidate];
  if (boolEnv(env.EDGE_GUARD_ALLOW_EXTRA_PATH, true) && candidate.includes("/")) {
    candidates.push(candidate.split("/")[0]);
  }

  return { action: "check", candidates: unique(candidates) };
}

export function buildLookupKeys(host, candidates, env = {}) {
  const prefix = stringEnv(env.KV_KEY_PREFIX, DEFAULT_KV_PREFIX);
  return candidates.map((candidate) => `${prefix}:v1:${host}:${candidate}`);
}

export function buildManifestKey(host, env = {}) {
  const prefix = stringEnv(env.KV_KEY_PREFIX, DEFAULT_KV_PREFIX);
  return `${prefix}:v2:${host}`;
}

function isProtectedHost(host, env) {
  const hosts = csvEnv(env.EDGE_GUARD_HOSTS, DEFAULT_HOSTS).map((value) =>
    value.toLowerCase(),
  );
  return hosts.includes(host);
}

async function hasAnyShortCode(kv, host, candidates, env) {
  if (!kv || typeof kv.get !== "function") {
    throw new Error("SHORT_CODES KV binding is missing");
  }

  const ttl = cacheTtl(env);
  const manifest = await kv.get(buildManifestKey(host, env), {
    cacheTtl: ttl,
    type: "json",
  });
  if (manifest?.version === 2 && manifest.codes) {
    return candidates.some((candidate) => manifest.codes[candidate] === true);
  }

  for (const key of buildLookupKeys(host, candidates, env)) {
    const value = await kv.get(key, { cacheTtl: ttl });
    if (value !== null) {
      return true;
    }
  }
  return false;
}

function reservedPrefixes(env) {
  return new Set(csvEnvAllowExplicitEmpty(
    env.EDGE_GUARD_RESERVED_PREFIXES,
    DEFAULT_RESERVED_PREFIXES.join(","),
  ));
}

function reservedPaths(env) {
  return new Set(csvEnvAllowExplicitEmpty(
    env.EDGE_GUARD_RESERVED_PATHS,
    DEFAULT_RESERVED_PATHS.join(","),
  ));
}

function cacheTtl(env) {
  const ttl = intEnv(env.EDGE_GUARD_KV_CACHE_TTL, 60);
  return Math.max(30, ttl);
}

function csvEnv(value, fallback) {
  return stringEnv(value, fallback)
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

function csvEnvAllowExplicitEmpty(value, fallback) {
  if (value === "") {
    return [];
  }
  return csvEnv(value, fallback);
}

function stringEnv(value, fallback) {
  return value === undefined || value === null || value === "" ? fallback : String(value);
}

function intEnv(value, fallback) {
  const parsed = Number.parseInt(stringEnv(value, fallback), 10);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function boolEnv(value, fallback) {
  const normalized = stringEnv(value, fallback ? "true" : "false").toLowerCase();
  return ["1", "true", "yes", "on"].includes(normalized);
}

function safeDecode(value) {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function unique(values) {
  return [...new Set(values)];
}
