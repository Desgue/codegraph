# Feature Specification: Parse Directory CLI

**Feature Branch**: `001-parse-directory-cli`
**Created**: 2025-10-19
**Status**: Draft
**Input**: User description: "initi cli commands to receive the folder to parse. when user run the program with the parse command it should output the directory passed as argument. when user does not specify a directory as argument the system must use the root of where the call was made and exibit a log specifying this behaviour"

## Clarifications

### Session 2025-10-19

- Q: How should the system respond when a user provides a file path instead of a directory path? → A: Detect it's a file and show error: "Error: '/path/to/file.go' is a file, not a directory"
- Q: How should the system behave when a directory exists but the user lacks read permissions? → A: Immediately fail with error: "Error: Permission denied accessing '[path]'"
- Q: When a user provides a path containing symbolic links, how should the system process it? → A: Follow symlinks and resolve to the actual directory, display resolved path

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Parse Explicit Directory (Priority: P1)

A developer wants to analyze a Go codebase in a specific directory by providing the path as an argument to the parse command.

**Why this priority**: This is the core functionality - parsing a specified directory is the primary use case that delivers immediate value.

**Independent Test**: Can be fully tested by running `codegraph parse /path/to/project` and verifying the directory is processed correctly. Delivers the core value of directory-based code analysis.

**Acceptance Scenarios**:

1. **Given** a valid directory path, **When** user runs `codegraph parse /path/to/project`, **Then** the system parses the specified directory and outputs confirmation showing the path being analyzed
2. **Given** a relative directory path, **When** user runs `codegraph parse ./subdirectory`, **Then** the system resolves the relative path and parses the correct directory
3. **Given** an invalid directory path, **When** user runs `codegraph parse /nonexistent/path`, **Then** the system displays a clear error message indicating the directory does not exist

---

### User Story 2 - Parse Current Directory (Priority: P2)

A developer working inside a Go project wants to parse the current directory without explicitly typing the path.

**Why this priority**: This is a convenience feature that improves user experience for the common case of analyzing the current working directory. Still delivers value independently.

**Independent Test**: Can be fully tested by navigating to a directory containing Go code, running `codegraph parse` without arguments, and verifying the current directory is parsed with appropriate logging.

**Acceptance Scenarios**:

1. **Given** no directory argument provided, **When** user runs `codegraph parse` from /home/user/myproject, **Then** the system parses the current directory and displays a log message: "No directory specified, using current directory: /home/user/myproject"
2. **Given** user is in a non-Go directory, **When** user runs `codegraph parse` without arguments, **Then** the system still processes the directory (may find no Go files) and logs the directory being used

---

### Edge Cases

- **File path instead of directory**: System detects the path is a file and displays error message: "Error: '[path]' is a file, not a directory"
- What happens when the directory path contains special characters or spaces?
- **Permission denied**: When user lacks read permissions, system immediately fails with error: "Error: Permission denied accessing '[path]'"
- **Symbolic links**: System follows symlinks to resolve the actual directory and displays the resolved absolute path
- What happens when the directory path is empty string vs null/undefined?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a `parse` command as the primary CLI entry point
- **FR-002**: System MUST accept an optional directory path argument for the `parse` command
- **FR-003**: System MUST validate that the provided directory path exists and is accessible
- **FR-004**: System MUST detect when the provided path is a file (not a directory) and display error: "Error: '[path]' is a file, not a directory"
- **FR-005**: System MUST output the directory path being analyzed before starting the parse operation
- **FR-006**: System MUST use the current working directory when no directory argument is provided
- **FR-007**: System MUST log a clear message when defaulting to current directory: "No directory specified, using current directory: [path]"
- **FR-008**: System MUST resolve relative paths to absolute paths before processing
- **FR-009**: System MUST follow symbolic links to resolve the actual directory and display the resolved absolute path
- **FR-010**: System MUST immediately fail with error "Error: Permission denied accessing '[path]'" when user lacks read permissions for the specified directory
- **FR-011**: System MUST display appropriate error messages for invalid paths or inaccessible directories
- **FR-012**: System MUST handle directory paths with spaces, special characters, and unicode characters correctly

### Key Entities

- **ParseCommand**: Represents the CLI parse command with its arguments and options
- **TargetDirectory**: Represents the directory to be parsed, including validation state and resolved absolute path

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can parse any valid directory by providing the path as an argument
- **SC-002**: Users can parse the current directory by running the command without arguments
- **SC-003**: All directory parsing operations display the target directory path before execution
- **SC-004**: 100% of cases where no directory is specified result in a clear log message indicating the current directory is being used
- **SC-005**: Invalid directory paths result in clear, actionable error messages within 100ms

## Assumptions

- Users have basic CLI experience and understand directory paths
- The parse command is the first/primary command of the CLI tool
- Standard output (stdout) is the appropriate channel for directory confirmation and logs
- Error messages should go to standard error (stderr)
- The tool follows standard POSIX conventions for path handling
- Users typically work in terminal environments where the current working directory is meaningful

## Dependencies

- Go standard library `os` package for working directory and path operations
- Go standard library `path/filepath` package for path resolution
- Existing or to-be-built CLI framework for command parsing
- Existing or to-be-built directory parsing logic that accepts a directory path

## Out of Scope

- Recursive directory traversal logic (assumed to be part of parsing implementation)
- Output format of parsing results (GraphML, JSON, etc.) - handled by separate feature
- Performance optimization for large codebases
- Watch mode or continuous parsing
- Multiple directory arguments in a single invocation
- Configuration file for default directories
