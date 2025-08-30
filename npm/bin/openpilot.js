#!/usr/bin/env node
const path = require('path');
const { spawn } = require('child_process');
const fs = require('fs');

const exe = path.join(__dirname, 'openpilot' + (process.platform === 'win32' ? '.exe' : ''));
if (!fs.existsSync(exe)) {
  console.error('[openpilot] executable missing; reinstall package');
  process.exit(1);
}
const args = process.argv.slice(2);
const child = spawn(exe, args, { stdio: 'inherit' });
child.on('exit', code => process.exit(code));
