// serviceWorker.js

importScripts("")

const serviceWorkerCacheName = "__ServiceWorkerCacheName_v3"

const appShellFiles = [
  "/",
  "/r/0.0.20/__main.css",
  "/r/0.0.20/__default_theme.js",
  "/r/0.0.20/__app.js",
]

self.addEventListener('install', function (e) {
  e.waitUntil(
    caches.open(serviceWorkerCacheName).then(function (cache) {
      console.log('[Service Worker] Caching all: app shell and content')
      return cache.addAll(appShellFiles);
    })
  );
  console.log('[Service Worker] Install')
})

// Fetching content using Service Worker
self.addEventListener('fetch', function (e) {
  e.respondWith(
    caches.match(e.request).then(function (r) {
      console.log('[Service Worker] Fetched resource ' + e.request.url)
      return fetch(e.request).then(function (response) {
        if (e.request.method.toLowerCase() != "get") return response
        return caches.open(serviceWorkerCacheName).then((cache) => {
          console.log('[Service Worker] Caching new resource: ' + e.request.url)
          cache.put(e.request, response.clone())
          return response
        })
      }).catch(function(err) {
        return r
      })
    })
  )
})

self.addEventListener('activate', function (event) {
  var cacheWhitelist = [serviceWorkerCacheName];

  event.waitUntil(
    caches.keys().then(function (keyList) {
      return Promise.all(keyList.map(function (key) {
        if (cacheWhitelist.indexOf(key) === -1) {
          return caches.delete(key);
        }
      }));
    })
  );
});