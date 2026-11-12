// SPDX-License-Identifier: Apache-2.0
// The script of every page: the theme switch, the search palette, the code tabs, the table of contents and the copy buttons. The pages work without it.
(function () {
  "use strict";
  var d = document, root = d.documentElement;

  // The theme switch goes through the system theme, light and dark.
  var themeBtn = d.querySelector("[data-theme-toggle]");
  if (themeBtn) {
    themeBtn.addEventListener("click", function () {
      var cur = root.getAttribute("data-theme");
      var sysDark = matchMedia("(prefers-color-scheme: dark)").matches;
      var next = !cur ? (sysDark ? "light" : "dark") : null;
      try {
        if (next) { root.setAttribute("data-theme", next); localStorage.setItem("al-theme", next); }
        else { root.removeAttribute("data-theme"); localStorage.removeItem("al-theme"); }
      } catch (e) {}
    });
  }

  // build is the commit of the deployment. A file that the script keeps in localStorage is kept with the build, so that a new deployment replaces it.
  var build = root.getAttribute("data-build") || "";
  function cached(key, url) {
    if (build) {
      try {
        var hit = JSON.parse(localStorage.getItem(key) || "null");
        if (hit && hit.b === build) return Promise.resolve(hit.d);
      } catch (e) {}
    }
    return fetch(url).then(function (r) { return r.json(); }).then(function (j) {
      if (build) try { localStorage.setItem(key, JSON.stringify({ b: build, d: j })); } catch (e) {}
      return j;
    });
  }

  // Search. The index is search.json. It loads when the palette opens for the first time, and the browser keeps it in localStorage until the next deployment.
  var index = null, loading = null;
  function load() {
    if (index) return Promise.resolve(index);
    if (!loading) loading = cached("prog-search", "/search.json").then(function (j) { index = j; return j; });
    return loading;
  }
  function norm(s) { return (s || "").toLowerCase().normalize("NFD").replace(/[̀-ͯ]/g, ""); }
  function score(e, q, words) {
    var n = norm(e.n), best = 0;
    if (n === q) best = 1000;
    else if (n.indexOf(q) === 0) best = 600;
    else if (n.indexOf(q) > 0) best = 400;
    (e.a || []).forEach(function (a) {
      var x = norm(a);
      if (x === q) best = Math.max(best, 900);
      else if (x.indexOf(q) === 0) best = Math.max(best, 500);
      else if (x.indexOf(q) >= 0) best = Math.max(best, 300);
    });
    if (!best) {
      var hay = n + " " + norm((e.a || []).join(" ")) + " " + norm(e.s) + " " + norm(e.k) + " " + norm(e.r);
      var all = words.every(function (w) { return hay.indexOf(w) >= 0; });
      if (all) best = 100;
    }
    if (e.t && q.toUpperCase() === e.t) best = 1000;
    return best;
  }
  function find(q) {
    q = norm(q).trim();
    if (!q || !index) return [];
    var words = q.split(/\s+/);
    var out = [];
    for (var i = 0; i < index.length; i++) {
      var s = score(index[i], q, words);
      if (s) out.push([s - i / 10000, index[i]]);
    }
    out.sort(function (a, b) { return b[0] - a[0]; });
    return out.map(function (x) { return x[1]; });
  }
  window.AlgoSearch = { load: load, find: find };

  var pal = d.getElementById("palette");
  if (pal && pal.showModal) {
    var input = pal.querySelector("input"), list = pal.querySelector("ul"), sel = 0, hits = [];
    function esc(s) { return String(s || "").replace(/[&<>"]/g, function (c) { return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]; }); }
    function draw() {
      hits = find(input.value).slice(0, 12);
      list.innerHTML = hits.map(function (e, i) {
        return '<li role="option" aria-selected="' + (i === sel) + '"><a href="' + esc(e.u) + '"><span class="k">' + esc(e.k) + "</span>" + esc(e.n) + '<span class="s">' + esc(e.s) + "</span></a></li>";
      }).join("") || (input.value.trim() ? '<li><a href="/search/?q=' + encodeURIComponent(input.value) + '">No page has this name. Search the full text.</a></li>' : "");
    }
    function open() { pal.showModal(); input.value = ""; sel = 0; list.innerHTML = ""; input.focus(); load().then(draw); }
    d.querySelectorAll("[data-palette]").forEach(function (a) {
      a.addEventListener("click", function (e) { if (e.metaKey || e.ctrlKey) return; e.preventDefault(); open(); });
    });
    d.addEventListener("keydown", function (e) {
      var t = e.target, typing = t && (t.tagName === "INPUT" || t.tagName === "TEXTAREA" || t.isContentEditable);
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || (e.key === "/" && !typing)) { e.preventDefault(); open(); }
    });
    input.addEventListener("input", function () { sel = 0; draw(); });
    input.addEventListener("keydown", function (e) {
      if (e.key === "ArrowDown") { sel = Math.min(sel + 1, hits.length - 1); draw(); e.preventDefault(); }
      else if (e.key === "ArrowUp") { sel = Math.max(sel - 1, 0); draw(); e.preventDefault(); }
      else if (e.key === "Enter") { if (hits[sel]) { e.preventDefault(); location.href = hits[sel].u; } }
    });
    pal.addEventListener("click", function (e) { if (e.target === pal) pal.close(); });
  }

  // Code tabs. A page may have more than one group of tabs. A choice of language applies to every group and is kept for the next page.
  var groups = Array.prototype.slice.call(d.querySelectorAll("section.code"));
  function pick(lang, save) {
    var hit = false;
    groups.forEach(function (g) {
      var b = g.querySelector('.tabs button[data-lang="' + lang + '"]');
      if (!b) return;
      hit = true;
      g.querySelectorAll(".tabs button").forEach(function (x) { x.setAttribute("aria-selected", String(x === b)); x.tabIndex = x === b ? 0 : -1; });
      g.querySelectorAll(".panel").forEach(function (p) { p.classList.toggle("on", p.dataset.lang === lang); });
    });
    if (hit && save) try { localStorage.setItem("prog-lang", lang); } catch (e) {}
    return hit;
  }
  if (groups.length) {
    var m = /^#(?:code|example)-(.+)$/.exec(location.hash);
    var saved = null;
    try { saved = localStorage.getItem("prog-lang"); } catch (e) {}
    if (!(m && pick(m[1], false)) && saved) pick(saved, false);
    groups.forEach(function (g) {
      var tabs = g.querySelector(".tabs");
      if (!tabs) return;
      tabs.addEventListener("click", function (e) { var b = e.target.closest("button[data-lang]"); if (b) pick(b.dataset.lang, true); });
      tabs.addEventListener("keydown", function (e) {
        if (e.key !== "ArrowRight" && e.key !== "ArrowLeft") return;
        var bs = Array.prototype.slice.call(tabs.querySelectorAll("button"));
        var i = bs.indexOf(d.activeElement);
        if (i < 0) return;
        var j = (i + (e.key === "ArrowRight" ? 1 : bs.length - 1)) % bs.length;
        bs[j].focus(); pick(bs[j].dataset.lang, true);
      });
    });
  }

  // Hover a line of the pseudocode to see the matching lines of the code.
  function codeLines(panel) { return panel ? panel.querySelectorAll("pre .line") : []; }
  function mark(id) {
    d.querySelectorAll(".panel").forEach(function (p) {
      var m = {};
      try { m = JSON.parse(p.dataset.markers || "{}"); } catch (e) {}
      var ls = codeLines(p);
      ls.forEach(function (l) { l.classList.remove("hl"); });
      (m[id] || []).forEach(function (n) { if (ls[n - 1]) ls[n - 1].classList.add("hl"); });
    });
  }
  window.AlgoMark = mark;
  d.querySelectorAll(".prose ol.pseudo li[data-line]").forEach(function (li) {
    li.addEventListener("mouseenter", function () { mark(li.dataset.line); });
    li.addEventListener("mouseleave", function () { mark(""); });
  });

  // Copy buttons.
  d.querySelectorAll("[data-copy]").forEach(function (b) {
    b.addEventListener("click", function () {
      var pre = b.closest(".panel").querySelector("pre");
      if (!pre || !navigator.clipboard) return;
      navigator.clipboard.writeText(pre.innerText).then(function () {
        var t = b.textContent; b.textContent = "Copied"; setTimeout(function () { b.textContent = t; }, 1400);
      });
    });
  });

  // The table of contents shows the current section.
  var toc = d.querySelector(".toc");
  if (toc && "IntersectionObserver" in window) {
    var links = {};
    toc.querySelectorAll("a").forEach(function (a) { links[a.getAttribute("href").slice(1)] = a; });
    var heads = Array.prototype.filter.call(d.querySelectorAll("main h2[id]"), function (h) { return links[h.id]; });
    var cur = null;
    var io = new IntersectionObserver(function () {
      var best = null;
      heads.forEach(function (h) { if (h.getBoundingClientRect().top < 120) best = h; });
      if (best !== cur) {
        if (cur && links[cur.id]) links[cur.id].classList.remove("on");
        cur = best;
        if (cur && links[cur.id]) links[cur.id].classList.add("on");
      }
    }, { rootMargin: "0px 0px -70% 0px", threshold: [0, 1] });
    heads.forEach(function (h) { io.observe(h); });
  }
  // The record page shows the JSON records of one collection. The fragment is the URL of the collection. The index gives the pack, the offset and the length of the records. The script asks for that range of the pack. A host that does not answer a range request sends the whole pack, and the script keeps it in the Cache API: the name of a pack has the hash of its bytes, so the copy stays good, and the next record of the same pack needs no request.
  if (d.body.classList.contains("page-record")) {
    var pre = d.querySelector("[data-record]"), status = d.querySelector("[data-record-status]"), links = d.querySelector("[data-record-links]");
    var coll = decodeURIComponent(location.hash.slice(1));
    var fail = function (msg) { status.textContent = msg; };
    var slice = function (buf, e) { return new TextDecoder().decode(new Uint8Array(buf, e[1], e[2])); };
    var show = function (text, pack) {
      try { text = JSON.stringify(JSON.parse(text), null, 2); } catch (e) {}
      pre.textContent = text; pre.hidden = false; status.hidden = true;
      links.querySelector("[data-record-raw]").href = pack;
      links.querySelector("[data-record-back]").href = coll;
      links.hidden = false;
    };
    var open = window.caches ? caches.open("prog-records") : Promise.resolve(null);
    fetch("/api/records.json").then(function (r) { return r.json(); }).then(function (idx) {
      var e = idx.records[coll];
      if (!e) return fail("No collection has the URL " + coll + ".");
      var pack = idx.packs[e[0]];
      return open.then(function (c) {
        // The packs of earlier deployments go out of the cache.
        if (c) c.keys().then(function (ks) { ks.forEach(function (k) { if (idx.packs.indexOf(new URL(k.url).pathname) < 0) c.delete(k); }); });
        return (c ? c.match(pack) : Promise.resolve(null)).then(function (hit) {
          if (hit) return hit.arrayBuffer().then(function (b) { show(slice(b, e), pack); });
          return fetch(pack, { headers: { Range: "bytes=" + e[1] + "-" + (e[1] + e[2] - 1) } }).then(function (r) {
            if (r.status === 206) return r.text().then(function (t) { show(t, pack); });
            if (!r.ok) return fail("The pack " + pack + " did not load.");
            if (c) c.put(pack, r.clone());
            return r.arrayBuffer().then(function (b) { show(slice(b, e), pack); });
          });
        });
      });
    }).catch(function () { fail("The records did not load."); });
  }

  // A link is fetched when the pointer rests on it or a finger touches it, so that the next page opens at once. The script does not do this when the reader asks to save data.
  var conn = navigator.connection;
  if (!(conn && conn.saveData)) {
    var done = {}, timer = null;
    var want = function (a) {
      if (!a || !a.href || a.origin !== location.origin || a.hasAttribute("download") || a.target) return false;
      var p = a.pathname;
      return !done[p] && p !== location.pathname && (p.slice(-1) === "/" || p.slice(-5) === ".html");
    };
    var prefetch = function (a) {
      done[a.pathname] = true;
      var l = d.createElement("link");
      l.rel = "prefetch"; l.href = a.pathname;
      d.head.appendChild(l);
    };
    d.addEventListener("mouseover", function (ev) {
      var a = ev.target.closest && ev.target.closest("a");
      if (!want(a)) return;
      clearTimeout(timer);
      timer = setTimeout(function () { prefetch(a); }, 65);
    }, { passive: true });
    d.addEventListener("mouseout", function () { clearTimeout(timer); }, { passive: true });
    d.addEventListener("touchstart", function (ev) {
      var a = ev.target.closest && ev.target.closest("a");
      if (want(a)) prefetch(a);
    }, { passive: true });
  }

  // The service worker (sw.js) keeps the assets and the pages that the reader opened, so that they open fast and with no network.
  if ("serviceWorker" in navigator && window.isSecureContext) {
    addEventListener("load", function () { navigator.serviceWorker.register("/sw.js").catch(function () {}); });
  }

  // The 404 page looks up the path in the map of moved paths, the tags and the aliases of the nodes, and sends the reader on. The paths of the cards are in 64 files; the name of the file is the FNV-1a hash of the path, as movedShard computes it in site.go.
  if (d.body.classList.contains("page-notfound")) {
    var path = location.pathname.replace(/\/+$/, "") || "/";
    try { path = decodeURIComponent(path); } catch (e) {}
    var shard = function (s) {
      var u = new TextEncoder().encode(s), h = 0x811c9dc5;
      for (var i = 0; i < u.length; i++) { h = Math.imul(h ^ u[i], 0x01000193) >>> 0; }
      return ("0" + (h % 64).toString(16)).slice(-2);
    };
    var get = function (u) { return fetch(u).then(function (r) { return r.ok ? r.json() : {}; }).catch(function () { return {}; }); };
    var go = function (to) { location.replace(to.indexOf("#") < 0 ? to + location.search + location.hash : to); };
    var recordPage = function (coll) { return "/api/record/#" + coll; };
    var low = path.toLowerCase();
    // The records of a collection were in <collection>/cards.json before the record packs.
    if (/\/cards\.json$/.test(path)) go(recordPage(path.replace(/cards\.json$/, "")));
    else
    get("/api/moved.json").then(function (m) {
      if (m[path] || m[low]) return go(m[path] || m[low]);
      return get("/api/moved/" + shard(path) + ".json").then(function (c) {
        if (c[path]) return go(c[path]);
        if (low !== path) return get("/api/moved/" + shard(low) + ".json").then(function (c2) { if (c2[low]) go(c2[low]); });
      });
    });
  }
})();
