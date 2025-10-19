# CLI Interface Contract: Parse Directory CLI

**Feature**: 001-parse-directory-cli | **Date**: 2025-10-19

## Overview

This document defines the command-line interface contract for the `codegraph parse` command.

---

## Command Signature

```bash
codegraph parse [directory] --output <file> [--include-tests]
```

### Positional Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `directory` | `string` | No | Path to directory to parse. Can be absolute or relative. If omitted, uses current working directory. |

### Flags

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `--output` | `string` | **Yes** | None | Output file path for parse results. Must be non-empty string. |
| `--include-tests` | `boolean` | No | `false` | Include test files (`*_test.go`) in parsing. |

---

## Exit Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| `0` | Success | Command completed successfully |
| `1` | User Error | Invalid arguments, directory doesn't exist, permission denied, validation failure |
| `2` | Internal Error | Unexpected error during execution (future implementation) |

---

## Output Channels

### stdout (Standard Output)
Reserved for actual command output. Currently unused (validation-only implementation).

**Future**: Will contain parse results or success confirmation messages.

### stderr (Standard Error)
Used for:
- Error messages
- Informational logs (e.g., "No directory specified, using current directory")
- Help text (when invalid usage)

---

## Usage Examples

### Example 1: Parse Specific Directory
```bash
$ codegraph parse /home/user/myproject --output graph.graphml
```

**Expected Output** (stderr):
```
(no output if successful - future implementation will add confirmation)
```

**Exit Code**: `0`

---

### Example 2: Parse Current Directory
```bash
$ cd /home/user/myproject
$ codegraph parse --output results.graphml
```

**Expected Output** (stderr):
```
No directory specified, using current directory: /home/user/myproject
```

**Exit Code**: `0`

---

### Example 3: Parse with Test Files Included
```bash
$ codegraph parse ./src --output graph.graphml --include-tests
```

**Expected Output** (stderr):
```
(no output if successful)
```

**Exit Code**: `0`

---

### Example 4: Parse Relative Directory
```bash
$ codegraph parse ../other-project --output graph.graphml
```

**Expected Behavior**: Resolves `../other-project` to absolute path, validates, proceeds.

**Exit Code**: `0`

---

## Error Scenarios

### Error 1: Missing Required Flag
```bash
$ codegraph parse /home/user/project
```

**Expected Output** (stderr):
```
Error: --output flag requires a file path
```

**Exit Code**: `1`

---

### Error 2: Directory Does Not Exist
```bash
$ codegraph parse /nonexistent/path --output graph.graphml
```

**Expected Output** (stderr):
```
Error: directory does not exist: /nonexistent/path
```

**Exit Code**: `1`

---

### Error 3: Path is File, Not Directory
```bash
$ codegraph parse /home/user/file.go --output graph.graphml
```

**Expected Output** (stderr):
```
Error: '/home/user/file.go' is a file, not a directory
```

**Exit Code**: `1`

---

### Error 4: Permission Denied
```bash
$ codegraph parse /root --output graph.graphml
```

**Expected Output** (stderr):
```
Error: permission denied accessing '/root'
```

**Exit Code**: `1`

---

### Error 5: Unknown Command
```bash
$ codegraph unknown-command
```

**Expected Output** (stderr):
```
Unknown command: unknown-command
Usage: codegraph <command> [options]
```

**Exit Code**: `1`

---

### Error 6: No Command Provided
```bash
$ codegraph
```

**Expected Output** (stderr):
```
Usage: codegraph <command> [options]
```

**Exit Code**: `1`

---

## Help Text

### Command Help
```bash
$ codegraph parse --help
```

**Expected Output** (stderr):
```
Usage of parse:
  -output string
        Output file path (required)
  -include-tests
        Include test files in parsing
```

**Exit Code**: `0`

---

## Validation Rules

### Rule 1: Output Flag Validation
- **Rule**: `--output` must be provided and non-empty
- **Checked**: Before directory validation
- **Error Format**: `Error: --output flag requires a file path`

### Rule 2: Directory Existence
- **Rule**: Directory path must exist on filesystem
- **Checked**: After path resolution
- **Error Format**: `Error: directory does not exist: [path]`

### Rule 3: Directory Type Check
- **Rule**: Path must be a directory, not a file
- **Checked**: After existence check
- **Error Format**: `Error: '[path]' is a file, not a directory`

### Rule 4: Directory Accessibility
- **Rule**: User must have read permission for directory
- **Checked**: During `os.Stat` call
- **Error Format**: `Error: permission denied accessing '[path]'`

---

## Special Behaviors

### Symbolic Link Resolution
When a symbolic link is provided as the directory argument:
- The symlink is followed to its target
- Validation is performed on the target directory
- Resolved absolute path is used internally
- No error message about symlink (transparent to user)

**Example**:
```bash
$ ln -s /real/path /tmp/link
$ codegraph parse /tmp/link --output graph.graphml
```

**Behavior**: Internally resolves to `/real/path`, validates `/real/path` is a directory.

---

### Path Resolution
All directory paths are resolved to absolute paths before validation:

| Input Path | Resolved Path | Notes |
|------------|---------------|-------|
| `/absolute/path` | `/absolute/path` | Already absolute |
| `./relative/path` | `[cwd]/relative/path` | Resolved relative to current directory |
| `../parent/path` | `[cwd]/../parent/path` → cleaned | Resolved and cleaned |
| ` ` (empty) | `[cwd]` | Default to current working directory |

---

## Future Compatibility

### Reserved Flags
The following flag names are reserved for future use:
- `--format` (output format selection)
- `--exclude` (file exclusion patterns)
- `--max-depth` (directory traversal depth limit)
- `--verbose` (verbose logging)
- `--quiet` (suppress non-error output)

### Reserved Commands
The following commands are reserved for future features:
- `codegraph export` (export to different formats)
- `codegraph analyze` (run analysis on graph)
- `codegraph version` (show version info)

---

## Cross-Platform Considerations

### Path Separators
- Unix/Linux/macOS: `/` (forward slash)
- Windows: `\` (backslash) or `/` (forward slash - also supported)
- **Contract**: CLI accepts both, Go's `filepath` package handles conversion

### Case Sensitivity
- Unix/Linux/macOS: Filesystem paths are case-sensitive
- Windows: Filesystem paths are case-insensitive
- **Contract**: CLI follows OS filesystem behavior (no case normalization)

### Line Endings
- Unix/Linux/macOS: `\n` (LF)
- Windows: `\r\n` (CRLF)
- **Contract**: Error messages use `\n`, Go runtime handles conversion on Windows

---

## Contract Versioning

**Version**: 1.0.0
**Status**: Draft (validation-only implementation)
**Breaking Changes**: Future addition of actual parsing logic will not change CLI interface
