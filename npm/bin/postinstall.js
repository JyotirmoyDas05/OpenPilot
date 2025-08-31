#!/usr/bin/env node
const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { spawnSync } = require('child_process');

// Requested version ("latest" means query the latest GitHub release). We also
// fall back to the wrapper's own package.json version if latest assets are not yet published.
const version = process.env.OPENPILOT_VERSION || 'latest';
const repo = process.env.OPENPILOT_REPO || 'JyotirmoyDas05/openpilot'; // owner/repo
const debug = !!process.env.OPENPILOT_DEBUG;
const platformMap = { win32: 'windows', darwin: 'darwin', linux: 'linux' };
// Common Go arch names to release arch names mapping; add more as needed
const archMap = { x64: 'x86_64', amd64: 'x86_64', arm64: 'arm64', arm: 'arm', ia32: 'i386' };

function log(msg) { console.log('[openpilot]', msg); }
function debugLog(msg) { if (debug) log('[debug] ' + msg); }
function fail(msg) { console.error('[openpilot] ' + msg); process.exit(1); }

const plat = platformMap[process.platform];
if (!plat) fail('Unsupported platform: ' + process.platform);
const arch = archMap[process.arch];
if (!arch) fail('Unsupported architecture: ' + process.arch);

// Generate potential asset names in decreasing preference order.
function candidateAssets(v) {
  const baseVersionTag = v === 'latest' ? 'latest' : 'v' + v;
  const namePatterns = [];
  const titlePlat = plat.charAt(0).toUpperCase() + plat.slice(1); // GoReleaser uses Title case for OS
  const versPrefix = v === 'latest' ? '' : v + '_';

  // Helper to push both lowercase and TitleCase OS variants (only adds distinct names)
  function pushPattern(template) {
    const variants = new Set();
    variants.add(template(plat));
    variants.add(template(titlePlat));
    for (const p of variants) namePatterns.push(p);
  }

  // Canonical tar.gz
  pushPattern(osName => v === 'latest' ? `openpilot_${osName}_${arch}.tar.gz` : `openpilot_${v}_${osName}_${arch}.tar.gz`);
  // Alternate dash style
  pushPattern(osName => `openpilot-${osName}-${arch}.tar.gz`);
  // Legacy tgz
  pushPattern(osName => v === 'latest' ? `openpilot_${osName}_${arch}.tgz` : `openpilot_${v}_${osName}_${arch}.tgz`);
  // Zip variants (primarily for Windows, but try generically)
  pushPattern(osName => v === 'latest' ? `openpilot_${osName}_${arch}.zip` : `openpilot_${v}_${osName}_${arch}.zip`);
  pushPattern(osName => `openpilot-${osName}-${arch}.zip`);

  return namePatterns.map(n => ({
    url: `https://github.com/${repo}/releases/${baseVersionTag === 'latest' ? 'latest/download' : 'download/' + baseVersionTag}/${n}`,
    name: n
  }));
}

const candidates = candidateAssets(version);
const destDir = path.join(__dirname, '..', 'dist');
const binDir = path.join(__dirname);
fs.mkdirSync(destDir, { recursive: true });

log('detect platform=' + plat + ' arch=' + arch + ' version=' + version + ' repo=' + repo);
debugLog('candidates: ' + candidates.map(c => c.name).join(', '));

function download(url, file) {
  return new Promise((resolve, reject) => {
    const f = fs.createWriteStream(file);
    https.get(url, res => {
      if (res.statusCode !== 200) {
        reject(new Error('HTTP ' + res.statusCode));
        return;
      }
      res.pipe(f);
      f.on('finish', () => f.close(resolve));
    }).on('error', reject);
  });
}

(async () => {
  const archivePath = path.join(destDir, 'openpilot_download');
  let lastErr = null;
  for (const c of candidates) {
    log('trying ' + c.url);
    try {
      await download(c.url, archivePath);
      log('downloaded ' + c.name);
      lastErr = null;
      break;
    } catch (e) {
      debugLog('failed ' + c.name + ': ' + e.message);
      lastErr = e;
    }
  }
  if (lastErr) {
    // Fallback: if we asked for 'latest', derive a concrete version from our package.json
    // and try versioned asset naming patterns (with and without v prefix) before giving up.
    if (version === 'latest') {
      try {
        const pkg = require(path.join(__dirname, '..', 'package.json'));
        const pv = pkg.version.replace(/^v/, '');
        if (pv) {
          log('latest assets missing; retrying with package version ' + pv);
          const versioned = candidateAssets(pv);
          for (const c of versioned) {
            log('trying ' + c.url);
            try {
              await download(c.url, archivePath);
              log('downloaded ' + c.name);
              lastErr = null;
              break;
            } catch (e) {
              debugLog('failed ' + c.name + ': ' + e.message);
              lastErr = e;
            }
          }
        }
      } catch (e) {
        debugLog('failed to load package version fallback: ' + e.message);
      }
    }
    if (lastErr) {
      fail('All download attempts failed. Tried: ' + candidates.map(c => c.url).join(', '));
    }
  }
  // Determine archive type by extension and extract accordingly.
  let extracted = false;
  const lower = candidates.concat([]).find(c => fs.existsSync(archivePath)) ? archivePath.toLowerCase() : archivePath.toLowerCase();
  if (archivePath.endsWith('.zip') || lower.endsWith('.zip')) {
    if (process.platform === 'win32') {
      // Use PowerShell Expand-Archive (available by default) for zip
      const ps = spawnSync('powershell', ['-NoProfile', '-Command', `Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force`]);
      if (ps.status !== 0) {
        fail('ZIP extraction failed');
      }
      extracted = true;
    } else {
      const unzip = spawnSync('unzip', ['-o', archivePath, '-d', destDir]);
      if (unzip.status !== 0) fail('ZIP extraction failed');
      extracted = true;
    }
  } else {
    const tar = process.platform === 'win32' ? 'tar.exe' : 'tar';
    const result = spawnSync(tar, ['-xzf', archivePath, '-C', destDir]);
    if (result.status !== 0) {
      fail('Tar extraction failed');
    }
    extracted = true;
  }
  if (!extracted) fail('Unknown archive type');
  // Find binary
  function findBin(dir) {
    const entries = fs.readdirSync(dir);
    for (const e of entries) {
      const p = path.join(dir, e);
      const stat = fs.statSync(p);
      if (stat.isFile() && e === 'openpilot' + (process.platform === 'win32' ? '.exe' : '')) {
        return p;
      }
      if (stat.isDirectory()) {
        const f = findBin(p);
        if (f) return f;
      }
    }
    return null;
  }
  const binary = findBin(destDir);
  if (!binary) fail('Binary not found after extraction');
  const finalPath = path.join(binDir, 'openpilot' + (process.platform === 'win32' ? '.exe' : ''));
  fs.copyFileSync(binary, finalPath);
  fs.chmodSync(finalPath, 0o755);
  log('installed binary to ' + finalPath);
})();
