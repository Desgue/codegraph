# Implementation Tasks: Parse Directory CLI

**Feature**: 001-parse-directory-cli | **Branch**: `001-parse-directory-cli` | **Date**: 2025-10-19

## Overview

This document provides a dependency-ordered task list for implementing the Parse Directory CLI feature. Tasks are organized by user story to enable independent implementation and testing.

**Tech Stack**:
- Language: Go 1.21+
- Dependencies: Go standard library only (`flag`, `os`, `path/filepath`)
- Testing: Go standard testing (`go test`)

**Key Entities**:
- `ParseCommand`: CLI command entity (in `cli/` package)
- `TargetDirectory`: Directory validation entity (in `path/` package)

---

## Task Legend

- `- [ ]` = Pending task
- `[P]` = Parallelizable (can run concurrently with other [P] tasks in same phase)
- `[US#]` = User Story number this task belongs to
- Task IDs are sequential (T001, T002, T003...)

---

## Phase 1: Setup & Project Initialization

**Goal**: Initialize Go module and create base project structure.

**Tasks**:

- [ ] T001 Initialize go.mod if not exists in repository root
- [ ] T002 Create cli/ package directory
- [ ] T003 Create path/ package directory

**Completion Criteria**: Directory structure matches plan.md, `go.mod` exists with correct module name.

---

## Phase 2: Foundational Components

**Goal**: Implement foundational validation logic that both user stories depend on.

**Note**: User input specifies "won't have any logic other than validation for flag inputs for now" - no tests are being generated for this validation-only implementation.

### 2.1: TargetDirectory Entity (Foundation)

- [ ] T004 [P] Implement TargetDirectory struct in path/validator.go with fields: InputPath, ResolvedPath, WasDefaulted
- [ ] T005 [P] Implement NewTargetDirectory constructor in path/validator.go (handles empty string → os.Getwd)
- [ ] T006 [P] Implement path resolution using filepath.Abs in path/validator.go
- [ ] T007 [P] Implement directory validation using os.Stat in path/validator.go (existence, type, permissions)
- [ ] T008 [P] Implement TargetDirectory.String() method in path/validator.go (returns ResolvedPath)
- [ ] T009 [P] Implement TargetDirectory.LogDefaultBehavior() method in path/validator.go (writes to stderr if defaulted)

**Completion Criteria**: `TargetDirectory` validates directories, resolves paths, handles empty input by defaulting to current directory.

---

## Phase 3: User Story 1 - Parse Explicit Directory (P1)

**User Story**: A developer wants to analyze a Go codebase in a specific directory by providing the path as an argument to the parse command.

**Independent Test**: Run `codegraph parse /path/to/project --output file.graphml` and verify:
- Valid directory → accepts and validates
- Relative path → resolves to absolute
- Invalid directory → displays clear error message

**Acceptance Scenarios**:
1. Valid directory path → system validates and outputs confirmation
2. Relative directory path → system resolves and validates
3. Invalid directory path → system displays error: "directory does not exist"

### 3.1: ParseCommand Entity (US1)

- [ ] T010 [P] [US1] Implement ParseCommand struct in cli/parse_command.go with fields: TargetDirectory, OutputFile, IncludeTests
- [ ] T011 [P] [US1] Implement flag parsing using flag.NewFlagSet in cli/parse_command.go (--output, --include-tests flags)
- [ ] T012 [US1] Implement NewParseCommand constructor in cli/parse_command.go (calls NewTargetDirectory, validates flags)
- [ ] T013 [P] [US1] Implement ParseCommand.Validate() method in cli/parse_command.go (validates OutputFile is non-empty)
- [ ] T014 [P] [US1] Implement ParseCommand.Execute() method in cli/parse_command.go (placeholder, returns nil for now)

### 3.2: Main CLI Entry Point (US1)

- [ ] T015 [US1] Implement main() function in main.go with subcommand routing (check len(os.Args) >= 2)
- [ ] T016 [US1] Implement "parse" subcommand case in main.go (calls NewParseCommand, handles errors, calls Execute)
- [ ] T017 [US1] Implement error handling in main.go (write to stderr, exit code 1 for user errors)
- [ ] T018 [US1] Implement unknown command handling in main.go (error message + usage, exit code 1)

**Completion Criteria**: Running `codegraph parse <directory> --output <file>` validates the directory, parses flags, and exits successfully if valid. Invalid inputs display specific error messages to stderr.

---

## Phase 4: User Story 2 - Parse Current Directory (P2)

**User Story**: A developer working inside a Go project wants to parse the current directory without explicitly typing the path.

**Independent Test**: Navigate to a directory, run `codegraph parse --output file.graphml` without directory argument, and verify:
- Current directory is used
- Log message is displayed: "No directory specified, using current directory: [path]"

**Acceptance Scenarios**:
1. No directory argument → system uses current directory and logs message
2. Non-Go directory → system still processes (may find no Go files) and logs directory

### 4.1: Current Directory Default Behavior (US2)

- [ ] T019 [US2] Add current directory default handling in cli/parse_command.go NewParseCommand (handle empty positional args)
- [ ] T020 [US2] Call TargetDirectory.LogDefaultBehavior() in cli/parse_command.go after successful creation when defaulted

**Completion Criteria**: Running `codegraph parse --output file.graphml` (without directory) uses current directory and displays log message to stderr.

---

## Phase 5: Edge Cases & Error Messages (Final Polish)

**Goal**: Ensure all edge cases from spec are handled with proper error messages.

### 5.1: Error Message Validation

- [ ] T021 [P] Verify file-vs-directory error message matches spec format in path/validator.go: "Error: '[path]' is a file, not a directory"
- [ ] T022 [P] Verify permission denied error message matches spec format in path/validator.go: "Error: permission denied accessing '[path]'"
- [ ] T023 [P] Verify directory not exist error message matches spec format in path/validator.go: "directory does not exist: [path]"
- [ ] T024 [P] Verify output flag validation error message in cli/parse_command.go: "--output flag requires a file path"

### 5.2: Cross-Platform & Special Characters

- [ ] T025 [P] Verify paths with spaces are handled correctly (Go's filepath handles this by default)
- [ ] T026 [P] Verify unicode paths are handled correctly (Go's UTF-8 support handles this by default)
- [ ] T027 [P] Verify symbolic links are followed (os.Stat does this by default, not os.Lstat)

**Completion Criteria**: All error messages match CLI contract specification. Edge cases (spaces, unicode, symlinks) work correctly per research.md decisions.

---

## Phase 6: Build & Final Validation

**Goal**: Ensure binary builds successfully and meets constitution standards.

- [ ] T028 Build binary using `go build -o codegraph .` from repository root
- [ ] T029 Verify binary runs with --help flag: `./codegraph parse --help` shows usage
- [ ] T030 Verify binary handles missing command: `./codegraph` shows usage message
- [ ] T031 Review code for constitution compliance (variable naming, function size, idiomatic Go)

**Completion Criteria**: Binary builds, runs, displays help text. Code follows all constitution principles.

---

## Dependency Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) - TargetDirectory entity
    ↓
    ├─→ Phase 3 (US1) - ParseCommand + main() - Core parsing functionality
    │       ↓
    └─→ Phase 4 (US2) - Current directory default - Extends US1
            ↓
        Phase 5 (Edge Cases) - Error message polish
            ↓
        Phase 6 (Build) - Final validation
```

**User Story Completion Order**:
1. US1 (Parse Explicit Directory) - **P1** - MUST complete first (core functionality)
2. US2 (Parse Current Directory) - **P2** - Extends US1 (convenience feature)

**Blocking Dependencies**:
- Phase 2 MUST complete before Phase 3 (US1 needs TargetDirectory)
- Phase 3 MUST complete before Phase 4 (US2 extends US1)

**Independent Stories**: US1 and US2 could theoretically be independent, but US2 reuses US1's infrastructure so sequential order is more efficient.

---

## Parallel Execution Opportunities

### Phase 2 (Foundational):
All TargetDirectory tasks (T004-T009) are parallelizable - different methods in same file.

### Phase 3 (US1):
- T010, T011, T013, T014 are parallelizable (different methods in ParseCommand)
- T012 must wait for T010, T011 (needs struct and flag parsing)
- T015-T018 must be sequential (main() function logic flow)

### Phase 5 (Edge Cases):
All tasks T021-T027 are parallelizable - verification only, no implementation dependencies.

**Example Parallel Batch (Phase 2)**:
```bash
# Implement all TargetDirectory methods in parallel
Task T004, T005, T006, T007, T008, T009 → can be done simultaneously
```

**Example Parallel Batch (Phase 3)**:
```bash
# Implement ParseCommand struct and methods in parallel
Task T010, T011, T013, T014 → can be done simultaneously
# Then Task T012 (constructor) → depends on T010, T011
```

---

## Implementation Strategy

### MVP Scope (Minimum Viable Product)
**Recommended MVP**: User Story 1 only (Phase 1-3)
- Delivers core value: parse explicit directory with validation
- ~200-300 LOC as estimated in plan.md
- Fully testable: `codegraph parse /path --output file.graphml`

### Full Feature Scope
**Complete Feature**: User Stories 1 + 2 (Phase 1-6)
- Adds convenience: parse current directory
- Adds polish: error message validation, edge cases
- Estimated ~350-400 LOC total

### Incremental Delivery Recommendation
1. **First PR**: Phase 1-3 (US1) - Core functionality
2. **Second PR**: Phase 4 (US2) - Current directory default
3. **Third PR**: Phase 5-6 - Polish and edge cases

This allows early user feedback on core functionality while maintaining small, reviewable PRs.

---

## Testing Strategy (Post-Validation Implementation)

**Current Scope**: Validation-only implementation (no tests generated per user input)

**Future Testing** (when actual parsing logic is added):
- Unit tests for `TargetDirectory` (data-model.md has test scenarios)
- Unit tests for `ParseCommand` (data-model.md has test scenarios)
- Integration tests for CLI contract (contracts/cli-interface.md has examples)
- Table-driven tests for all validation scenarios

---

## Task Summary

**Total Tasks**: 31
- Phase 1 (Setup): 3 tasks
- Phase 2 (Foundational): 6 tasks
- Phase 3 (US1): 9 tasks
- Phase 4 (US2): 2 tasks
- Phase 5 (Edge Cases): 7 tasks
- Phase 6 (Build): 4 tasks

**Parallelizable Tasks**: 17 tasks marked [P]
**User Story Breakdown**:
- US1: 9 tasks (T010-T018)
- US2: 2 tasks (T019-T020)
- Foundational: 6 tasks (T004-T009)
- Setup/Polish: 14 tasks (T001-T003, T021-T031)

**Estimated Effort**:
- Setup: 30 minutes
- Foundational: 2 hours
- US1: 3 hours
- US2: 30 minutes
- Edge Cases: 1 hour
- Build: 30 minutes
- **Total**: ~7.5 hours

---

## Notes

- All file paths are specified per task for clarity
- Tasks follow strict checkbox format: `- [ ] [ID] [P?] [Story?] Description with path`
- Constitution compliance is enforced throughout (descriptive names, small functions, idiomatic Go)
- Standard library only constraint is maintained (no external dependencies)
- Error messages match CLI contract specification exactly
- User input specified validation-only implementation - no actual parsing logic yet
