// SPDX-License-Identifier: Apache-2.0
// The service worker of the site. prog site build writes it to /sw.js.
// A page comes from the network first, so that a reader sees a new deployment at once. When the network fails, or is slower than 3 seconds, the copy that the reader opened last time is used. An asset whose URL has the hash of its bytes (?v=, the fonts and KaTeX) comes from the cache first. Each other asset comes from the cache and is fetched again in the background. The record packs are kept by the record page, and search.json by the script in localStorage, so this worker lets them through.
"use strict";
const PAGES = "prog-pages-1", ASSETS = "prog-assets-1", KEEP = [PAGES, ASSETS, "prog-records"], MAX_PAGES = 200;

self.addEventListener("install", () => self.skipWaiting());
self.addEventListener("activate", (ev) => {
  ev.waitUntil(caches.keys().then((ks) => Promise.all(ks.filter((k) => !KEEP.includes(k)).map((k) => caches.delete(k)))).then(() => self.clients.claim()));
});

async function trim(c) {
  const ks = await c.keys();
  for (let i = 0; i < ks.length - MAX_PAGES; i++) await c.delete(ks[i]);
}

async function page(ev) {
  const c = await caches.open(PAGES);
  const net = fetch(ev.request).then(async (r) => {
    if (r.ok && r.type === "basic") {
      await c.put(ev.request, r.clone());
      ev.waitUntil(trim(c));
    }
    return r;
  });
  ev.waitUntil(net.catch(() => {}));
  const slow = new Promise((ok) => setTimeout(ok, 3000)).then(() => c.match(ev.request));
  try {
    const r = await Promise.race([net, slow.then((h) => h || net)]);
    return r;
  } catch (e) {
    return (await c.match(ev.request)) || Response.error();
  }
}

async function asset(ev, fixed) {
  const c = await caches.open(ASSETS);
  const hit = await c.match(ev.request);
  const net = fetch(ev.request).then(async (r) => {
    if (r.ok && r.type === "basic") {
      // An older version of the same file goes out of the cache.
      const p = new URL(ev.request.url).pathname;
      for (const k of await c.keys()) if (k.url !== ev.request.url && new URL(k.url).pathname === p) await c.delete(k);
      await c.put(ev.request, r.clone());
    }
    return r;
  });
  if (hit) {
    if (!fixed) ev.waitUntil(net.catch(() => {}));
    return hit;
  }
  return net;
}

self.addEventListener("fetch", (ev) => {
  const req = ev.request;
  if (req.method !== "GET" || req.headers.has("range")) return;
  const u = new URL(req.url);
  if (u.origin !== location.origin) return;
  if (req.mode === "navigate") return ev.respondWith(page(ev));
  if (u.pathname.startsWith("/assets/")) {
    const fixed = u.searchParams.has("v") || u.pathname.startsWith("/assets/fonts/") || u.pathname.startsWith("/assets/katex/");
    return ev.respondWith(asset(ev, fixed));
  }
});
