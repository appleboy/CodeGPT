#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  let platformName;
  let archName;
  let armVersion = '';

  // Map platform
  switch (platform) {
    case 'darwin':
      platformName = 'darwin';
      break;
    case 'linux':
      platformName = 'linux';
      break;
    case 'win32':
      platformName = 'windows';
      break;
    case 'freebsd':
      platformName = 'freebsd';
      break;
    default:
      throw new Error(`Unsupported platform: ${platform}`);
  }

  // Map architecture
  switch (arch) {
    case 'x64':
      archName = 'amd64';
      break;
    case 'arm':
      archName = 'arm';
      // Try to detect ARM version
      try {
        if (platform === 'linux') {
          const cpuinfo = fs.readFileSync('/proc/cpuinfo', 'utf8');
          if (cpuinfo.includes('ARMv7')) {
            armVersion = '-7';
          } else if (cpuinfo.includes('ARMv6')) {
            armVersion = '-6';
          } else {
            armVersion = '-5';
          }
        }
      } catch (e) {
        armVersion = '-7'; // Default to ARMv7
      }
      break;
    case 'arm64':
      archName = 'arm64';
      break;
    default:
      throw new Error(`Unsupported architecture: ${arch}`);
  }

  return { platformName, archName, armVersion };
}

function getVersion() {
  const packageJson = require('./package.json');
  return packageJson.version;
}

function getBinaryName() {
  const { platformName, archName, armVersion } = getPlatform();
  const version = getVersion();
  const ext = platformName === 'windows' ? '.exe' : '';

  return `CodeGPT-${version}-${platformName}-${archName}${armVersion}${ext}`;
}

function install() {
  try {
    const binaryName = getBinaryName();
    const binaryPath = path.join(__dirname, 'binaries', binaryName);
    const targetPath = path.join(__dirname, 'bin', 'codegpt-bin');

    if (!fs.existsSync(binaryPath)) {
      console.error(`Error: Binary not found for your platform: ${binaryName}`);
      console.error('Supported platforms:');
      console.error('  - darwin (macOS): x64, arm64');
      console.error('  - linux: x64, arm, arm64');
      console.error('  - windows: x64');
      console.error('  - freebsd: x64');
      process.exit(1);
    }

    // Copy binary to bin directory
    fs.copyFileSync(binaryPath, targetPath);

    // Make it executable on Unix-like systems
    if (process.platform !== 'win32') {
      fs.chmodSync(targetPath, 0o755);
    }

    console.log(`✓ CodeGPT installed successfully for ${process.platform}-${process.arch}`);
  } catch (error) {
    console.error('Installation failed:', error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  install();
}

module.exports = { install };
