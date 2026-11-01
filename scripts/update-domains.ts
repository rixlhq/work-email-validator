#!/usr/bin/env bun

/**
 * Domain Lists Update Script
 * Downloads and merges disposable and free email domain lists from multiple sources
 */

import { readFileSync, writeFileSync } from "fs";

// ============================================================================
// Types
// ============================================================================

interface Source {
  name: string;
  url: string;
}

interface Config {
  disposable_sources: Source[];
  free_source: Source;
  extra_free_domains: string[];
  exclude_domains: string[];
}

interface ConflictResolutionResult {
  disposable: Set<string>;
  free: Set<string>;
  removed: number;
}

interface ProcessingResult {
  disposableDomains: Set<string>;
  freeDomains: Set<string>;
  conflictsResolved: number;
}

// ============================================================================
// Constants
// ============================================================================

import { resolve, dirname } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const ROOT_DIR = resolve(__dirname, "..");

const CONFIG_FILE = resolve(ROOT_DIR, "config/repositories.json");
const OUTPUT_PATHS = {
  disposable: resolve(ROOT_DIR, "data/disposable_domains.txt"),
  free: resolve(ROOT_DIR, "data/free_domains.txt"),
} as const;

const DOMAIN_PATTERN = /^[a-z0-9-]+(\.[a-z0-9-]+)+$/;
const MAX_LIST_CHANGE_RATIO = 0.1;
const FORCE_FLAG = "--force";

// ANSI color codes
const COLORS = {
  reset: "\x1b[0m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  cyan: "\x1b[36m",
  red: "\x1b[31m",
} as const;

const SYMBOLS = {
  success: "✓",
  warning: "⚠",
  error: "✗",
  done: "✅",
} as const;

// ============================================================================
// Utilities
// ============================================================================

const log = {
  info: (msg: string): void => console.log(msg),
  success: (msg: string): void => 
    console.log(`${COLORS.green}${SYMBOLS.success}${COLORS.reset} ${msg}`),
  warn: (msg: string): void => 
    console.log(`${COLORS.yellow}${SYMBOLS.warning}${COLORS.reset}  ${msg}`),
  error: (msg: string): void => 
    console.log(`${COLORS.red}${SYMBOLS.error}${COLORS.reset} ${msg}`),
  header: (msg: string): void => 
    console.log(`\n${COLORS.cyan}${msg}${COLORS.reset}`),
};

function formatTimestamp(): string {
  return new Date().toISOString().replace("T", " ").split(".")[0] + " UTC";
}

function formatError(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  return String(error);
}

function logProgress(current: number, total: number, name: string, url: string): void {
  console.log(`  [${current}/${total}] ${name}`);
  console.log(`      URL: ${url}`);
}

function logDownloadResult(success: boolean, count: number): void {
  if (success) {
    console.log(`      ${COLORS.green}${SYMBOLS.success}${COLORS.reset} Downloaded ${count} domains`);
  } else {
    console.log(`      ${COLORS.yellow}${SYMBOLS.warning}${COLORS.reset}  Failed to download, skipping...`);
  }
  console.log();
}

// ============================================================================
// Core Functions
// ============================================================================

function parseDomainList(text: string): string[] {
  return text
    .split("\n")
    .map((line) => line.trim().toLowerCase())
    .filter((line) => DOMAIN_PATTERN.test(line));
}

async function downloadFile(url: string, name: string): Promise<string[]> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
  }
  const text = await response.text();
  return parseDomainList(text);
}

interface DownloadResult {
  source: Source;
  domains: string[];
  error?: string;
}

async function downloadSourceWithRetry(
  source: Source,
  index: number,
  total: number
): Promise<DownloadResult> {
  const num = index + 1;
  logProgress(num, total, source.name, source.url);
  
  try {
    const domains = await downloadFile(source.url, source.name);
    
    if (domains.length === 0) {
      throw new Error("Downloaded file is empty or contains no valid domains");
    }
    
    logDownloadResult(true, domains.length);
    return { source, domains };
  } catch (error: unknown) {
    const errorMsg = formatError(error);
    logDownloadResult(false, 0);
    return { source, domains: [], error: errorMsg };
  }
}

async function downloadDisposableSources(
  sources: Source[]
): Promise<Map<string, Set<string>>> {
  log.header("Downloading disposable domains from multiple sources...");
  
  // Download all sources in parallel
  const downloadPromises = sources.map((source, index) =>
    downloadSourceWithRetry(source, index, sources.length)
  );
  
  const results = await Promise.all(downloadPromises);
  
  // Check for failures
  const failures = results.filter((r) => r.error !== undefined);
  
  if (failures.length > 0) {
    log.error(`\nFailed to download ${failures.length} source(s):`);
    for (const failure of failures) {
      log.error(`  - ${failure.source.name}: ${failure.error}`);
    }
    throw new Error(`Failed to download ${failures.length} of ${sources.length} disposable domain sources`);
  }
  
  // Build results map
  const resultsMap = new Map<string, Set<string>>();
  for (const result of results) {
    resultsMap.set(result.source.name, new Set(result.domains));
  }
  
  return resultsMap;
}

function normalizeExcludeDomains(excludeDomains: string[]): Set<string> {
  return new Set(excludeDomains.map((d) => d.toLowerCase()));
}

function mergeDomains(
  sourcesMap: Map<string, Set<string>>,
  excludeDomains: string[]
): Set<string> {
  log.header("Merging and deduplicating disposable domains...");
  
  const excludeSet = normalizeExcludeDomains(excludeDomains);
  log.info(`  Excluding ${excludeDomains.length} domains: ${excludeDomains.join(", ")}`);
  
  const merged = new Set<string>();
  
  for (const [_name, domains] of sourcesMap) {
    for (const domain of domains) {
      if (!excludeSet.has(domain)) {
        merged.add(domain);
      }
    }
  }
  
  return merged;
}

async function downloadFreeDomains(source: Source): Promise<Set<string>> {
  log.header("Downloading free email domains...");
  console.log(`  Source: ${source.name}`);
  console.log(`  URL: ${source.url}`);
  
  try {
    const domains = await downloadFile(source.url, source.name);
    
    if (domains.length === 0) {
      throw new Error("Downloaded file is empty or contains no valid domains");
    }
    
    log.success(`Downloaded ${domains.length} free domains`);
    return new Set(domains);
  } catch (error: unknown) {
    log.error(`Failed to download free domains: ${formatError(error)}`);
    throw new Error(`Critical: Free domains download failed - ${formatError(error)}`);
  }
}

function findConflicts(disposable: Set<string>, free: Set<string>): string[] {
  const conflicts: string[] = [];
  for (const domain of free) {
    if (disposable.has(domain)) {
      conflicts.push(domain);
    }
  }
  return conflicts;
}

function removeConflicts(disposable: Set<string>, conflicts: string[]): void {
  for (const domain of conflicts) {
    disposable.delete(domain);
  }
}

function resolveConflicts(
  disposable: Set<string>,
  free: Set<string>
): ConflictResolutionResult {
  log.header("Filtering free domains and removing conflicts...");
  log.info("  Removing legitimate free providers from disposable list...");
  
  const conflicts = findConflicts(disposable, free);
  removeConflicts(disposable, conflicts);
  
  log.success(`Removed ${conflicts.length} legitimate free providers from disposable list`);
  
  return {
    disposable,
    free,
    removed: conflicts.length,
  };
}

function createFileHeader(sources: string[], lastUpdated: string): string[] {
  return [
    ...sources,
    `# Last updated: ${lastUpdated}`,
  ];
}

function createDisposableHeader(sources: Source[]): string[] {
  return createFileHeader(
    [
      "# Disposable Email Domains",
      "# Sources:",
      ...sources.map((s) => `#   - ${s.name}`),
    ],
    formatTimestamp()
  );
}

function createFreeHeader(source: Source): string[] {
  return createFileHeader(
    [
      "# Free Email Providers",
      `# Source: ${source.name}`,
    ],
    formatTimestamp()
  );
}

function writeDomainsFile(
  filename: string,
  domains: Set<string>,
  header: string[]
): void {
  const sortedDomains = Array.from(domains).sort();
  const content = [...header, "", ...sortedDomains].join("\n") + "\n";
  writeFileSync(filename, content, "utf-8");
}

function assertListSizeChange(filename: string, domains: Set<string>): void {
  const previous = parseDomainList(readFileSync(filename, "utf-8")).length;
  const changeRatio = Math.abs(domains.size - previous) / previous;

  if (changeRatio > MAX_LIST_CHANGE_RATIO) {
    throw new Error(
      `${filename} would change from ${previous} to ${domains.size} domains, rerun with ${FORCE_FLAG} if this is expected`
    );
  }
}

function loadConfig(): Config {
  try {
    const configContent = readFileSync(CONFIG_FILE, "utf-8");
    return JSON.parse(configContent) as Config;
  } catch (error: unknown) {
    log.error(`Failed to load configuration: ${formatError(error)}`);
    process.exit(1);
  }
}

function logConfigSummary(config: Config): void {
  log.success(`Loaded ${config.disposable_sources.length} disposable sources`);
  log.success(`Loaded 1 free email source`);
  log.success(`Loaded ${config.exclude_domains.length} domains to exclude`);
}

async function processDomainsData(config: Config): Promise<ProcessingResult> {
  // Download disposable domains
  const disposableSourcesMap = await downloadDisposableSources(
    config.disposable_sources
  );
  
  // Merge and deduplicate
  let disposableDomains = mergeDomains(
    disposableSourcesMap,
    config.exclude_domains
  );
  
  log.success(
    `Disposable domains (before filtering): ${disposableDomains.size} unique domains`
  );
  
  // Download free domains
  let freeDomains = await downloadFreeDomains(config.free_source);
  config.extra_free_domains.forEach((domain) => freeDomains.add(domain));
  
  // Resolve conflicts (free takes priority)
  const result = resolveConflicts(disposableDomains, freeDomains);
  disposableDomains = result.disposable;
  freeDomains = result.free;
  
  log.success(`Free domains: ${freeDomains.size} domains (after filtering)`);
  
  return {
    disposableDomains,
    freeDomains,
    conflictsResolved: result.removed,
  };
}

function writeOutputFiles(
  result: ProcessingResult,
  config: Config
): void {
  if (!process.argv.includes(FORCE_FLAG)) {
    assertListSizeChange(OUTPUT_PATHS.disposable, result.disposableDomains);
    assertListSizeChange(OUTPUT_PATHS.free, result.freeDomains);
  }

  const disposableHeader = createDisposableHeader(config.disposable_sources);
  writeDomainsFile(OUTPUT_PATHS.disposable, result.disposableDomains, disposableHeader);
  
  const freeHeader = createFreeHeader(config.free_source);
  writeDomainsFile(OUTPUT_PATHS.free, result.freeDomains, freeHeader);
}

function printSummary(result: ProcessingResult, duration: string): void {
  log.header(`${SYMBOLS.done} Done! Lists updated successfully.`);
  console.log(`   Total disposable: ${result.disposableDomains.size}`);
  console.log(`   Total free: ${result.freeDomains.size}`);
  console.log(`   Conflicts resolved: ${result.conflictsResolved}`);
  console.log(`   Time taken: ${duration}s\n`);
}

async function main(): Promise<void> {
  const startTime = Date.now();
  
  console.log("Updating domain lists...");
  log.info(`Reading configuration from: ${CONFIG_FILE}\n`);
  
  const config = loadConfig();
  logConfigSummary(config);
  
  const result = await processDomainsData(config);
  writeOutputFiles(result, config);
  
  const duration = ((Date.now() - startTime) / 1000).toFixed(2);
  printSummary(result, duration);
}

main().catch((error: unknown) => {
  log.error(`Fatal error: ${formatError(error)}`);
  process.exit(1);
});
