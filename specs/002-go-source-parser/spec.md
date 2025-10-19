# Feature Specification: Go Source Parser

**Feature Branch**: `002-go-source-parser`
**Created**: 2025-10-19
**Status**: Draft
**Input**: User description: "Parse and load Go source code from a target directory to prepare it for graph analysis. Users need to analyze Go codebases by converting source code into structured data. The tool should discover all Go packages within a specified directory, parse the source files, and prepare them for the next phase of graph construction."

## Clarifications

### Session 2025-10-19

- Q: When the parse command completes, where should the parsed data and statistics be sent? → A: Terminal output for statistics, in-memory data for next phase
- Q: What should happen when the target directory contains no Go files? → A: Exit successfully with "0 packages found, 0 files parsed"
- Q: When files fail to parse due to syntax errors, how should the errors be presented in the terminal output? → A: List file paths with error messages (e.g., "path/file.go:15:3: expected ';', found 'EOF'")
- Q: When directory permissions prevent reading some subdirectories during parsing, what should happen? → A: packages.Load handles this automatically by skipping inaccessible directories and continuing
- Q: When the system encounters symbolic links within the project directory during discovery, how should they be handled? → A: Follow all symlinks and parse linked directories/files

## User Scenarios & Testing

### User Story 1 - Parse Entire Project (Priority: P1)

A developer runs the parse command on their Go project directory and receives a complete analysis of all packages and files, even if some files contain syntax errors.

**Why this priority**: This is the core value proposition - users must be able to load their codebase before any analysis can happen. Without this, the tool provides no value.

**Independent Test**: Can be fully tested by running the parse command on a Go project directory and verifying that all valid packages are discovered and parsed, with clear output showing statistics.

**Acceptance Scenarios**:

1. **Given** a valid Go project directory, **When** the user runs the parse command, **Then** all Go packages in the directory are discovered and successfully parsed
2. **Given** a Go project with some files containing syntax errors, **When** the user runs the parse command, **Then** valid files are parsed successfully and errors are listed in the terminal output showing file path, line, column, and error description for each failed file
3. **Given** a Go project with nested packages, **When** the user runs the parse command, **Then** all packages at all directory levels are discovered and parsed
4. **Given** the parse command completes, **When** the user reviews the output, **Then** they see clear statistics showing total packages found, files parsed, and any errors encountered
5. **Given** a directory with no Go files, **When** the user runs the parse command, **Then** the system exits successfully displaying "0 packages found, 0 files parsed"
6. **Given** a Go project with some subdirectories having restricted permissions, **When** the user runs the parse command, **Then** accessible directories are parsed and inaccessible directories are silently skipped
7. **Given** a Go project containing symbolic links to directories or files, **When** the user runs the parse command, **Then** packages.Load follows symlinks automatically according to Go's standard package loading behavior

---

### User Story 2 - Module-Aware Parsing (Priority: P2)

A developer working with Go modules expects the parser to understand module boundaries and package organization according to go.mod structure.

**Why this priority**: Modern Go projects use modules, and understanding module structure is essential for accurate dependency analysis. However, basic parsing can work without module awareness.

**Independent Test**: Can be tested by running the parse command on a project with go.mod and verifying that package paths respect module boundaries and import paths are correctly resolved.

**Acceptance Scenarios**:

1. **Given** a Go project with a go.mod file, **When** the user runs the parse command, **Then** package import paths are resolved according to the module definition
2. **Given** a multi-module project, **When** the user runs the parse command, **Then** each module's packages are correctly identified with their full import paths

---

### User Story 3 - Documentation Preservation (Priority: P3)

A developer analyzing code expects the parser to preserve comments and documentation so that documentation-based analysis features can be built later.

**Why this priority**: Preserving documentation enables future features like documentation coverage analysis and API documentation generation, but isn't required for basic dependency graphing.

**Independent Test**: Can be tested by parsing a file with comments and verifying that comment data is retained in the parsed output structure.

**Acceptance Scenarios**:

1. **Given** Go source files with package-level and function-level documentation comments, **When** the parser processes these files, **Then** all documentation comments are preserved in the parsed data structure
2. **Given** source files with inline comments, **When** the parser processes these files, **Then** inline comments are accessible in the parsed output

---

### Edge Cases

- **Empty directory**: See FR-001a
- **Permission denied**: See FR-001b
- **Symbolic links**: See FR-001c

**Open edge cases requiring future consideration:**
- How does the system handle extremely large projects with thousands of files?
- What happens when a file is valid Go syntax but imports non-existent packages?
- How does the system handle build-tag-specific files (e.g., _linux.go, _windows.go)?

## Requirements

### Functional Requirements

- **FR-001**: System MUST discover all directories containing Go source files (*.go) within the target directory using packages.Load
- **FR-001a**: System MUST exit successfully with zero statistics when no Go files are found in the target directory
- **FR-001b**: System MUST handle directories with insufficient read permissions by relying on packages.Load's built-in behavior to skip inaccessible directories
- **FR-001c**: System MUST handle symbolic links using packages.Load's standard symlink behavior
- **FR-002**: System MUST identify Go packages by analyzing package declarations in source files
- **FR-003**: System MUST parse Go source files using the standard parsing approach that preserves AST structure
- **FR-004**: System MUST continue processing remaining files when individual files fail to parse due to syntax errors
- **FR-005**: System MUST collect and report all parse errors to the terminal, listing each failed file path with its specific error message (including line, column, and error description)
- **FR-006**: System MUST recognize Go module boundaries when a go.mod file is present
- **FR-007**: System MUST preserve documentation comments (package-level, type-level, and function-level) during parsing
- **FR-008**: System MUST output summary statistics to the terminal including: total packages discovered, total files processed, and parse error count
- **FR-009**: System MUST retain parsed data in memory for immediate use by subsequent graph construction phase
- **FR-010**: System MUST complete parsing within reasonable time limits for typical projects (see SC-004: <5s for 50k LOC)

### Key Entities

- **Package**: Represents a Go package from packages.Package with name, import path, directory location, parsed AST, and collection of source files
- **Parse Error**: Embedded in packages.Package.Errors; records file path, line number, column number, and error description for files that failed to parse; formatted for terminal output as "path/file.go:line:column: error message"

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can successfully parse any valid Go project directory with at least one Go package
- **SC-002**: Parse errors in individual files do not prevent processing of other files in the same package or other packages
- **SC-003**: Users receive clear terminal output showing: number of packages found, number of files parsed, number of parse failures, and detailed list of parse errors with file paths, line numbers, columns, and error descriptions printed to stderr
- **SC-004**: Parsing completes in under 5 seconds for projects with up to 50,000 lines of code on standard hardware
- **SC-005**: The parsed data structure enables subsequent graph construction without requiring re-parsing
- **SC-006**: Users can parse projects containing up to 1000 source files without memory or performance issues

## Assumptions

- The target directory contains valid Go source code organized in standard Go project structure
- Users have read permissions for all directories and files they want to parse
- Standard Go file naming conventions are followed (*.go files)
- The parsing phase focuses on discovery and loading; dependency resolution and graph construction happen in subsequent phases
- Performance target (5 seconds for 50k LOC) assumes modern hardware (multi-core CPU, SSD storage)
- "Preserving comments" means retaining them in the parsed AST structure, not displaying them in output
- Parsed data is retained in memory for immediate consumption by graph construction phase within the same execution session
- Symbolic links are handled by packages.Load using Go's standard behavior

## Dependencies

- Builds on the existing directory validation functionality from feature 001-parse-directory-cli
- Requires the validated target directory path as input
- Prepares data for the future graph construction phase (not yet implemented)

## Out of Scope

- Dependency graph construction (handled in future feature)
- Visualization of the parsed code structure
- Code quality analysis or linting
- Modification or transformation of source files
- Cross-module dependency resolution beyond basic import path recognition
- Build constraint evaluation (build tags)
- Generated code detection or special handling
