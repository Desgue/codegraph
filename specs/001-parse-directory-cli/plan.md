# Implementation Plan: Parse Directory CLI

**Branch**: `001-parse-directory-cli` | **Date**: 2025-10-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-parse-directory-cli/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a basic CLI using Go's standard library `flag` package that accepts a `parse` command with a directory argument and optional flags (`--output <file>`, `--include-tests`). For now, the implementation will only validate flag inputs without executing parsing logic.

## Technical Context

**Language/Version**: Go 1.21+ (using modern for range syntax as per constitution)
**Primary Dependencies**: Go standard library only (`flag`, `os`, `path/filepath`)
**Storage**: N/A (CLI tool, no persistent storage)
**Testing**: Go standard testing (`go test`)
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: Single project (CLI application)
**Performance Goals**: Command validation <10ms (immediate feedback)
**Constraints**: Standard library only, no external dependencies
**Scale/Scope**: Single command (`parse`) with 2 flags, ~200-300 LOC for initial validation

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Clear Variable Naming
**Status**: ✅ PASS
**Evidence**: Feature requires descriptive names for directory paths, validation states. No abbreviations needed.

### Principle II: Function Size Discipline
**Status**: ✅ PASS
**Evidence**: CLI validation logic naturally decomposes into small functions (validate path, check permissions, resolve symlinks). Each validation rule is a separate responsibility.

### Principle III: Domain-Driven Design
**Status**: ✅ PASS
**Evidence**: Domain entities identified in spec (ParseCommand, TargetDirectory) will encapsulate validation behavior with their data. No anemic models.

### Principle IV: Go Idiomatic Practices
**Status**: ✅ PASS
**Evidence**: Using standard library `flag` package, `os.Stat` for path validation, `filepath` for resolution. Modern `for range` will be used where applicable. Constitution explicitly requires `any` over `interface{}`.

### Principle V: Simplicity (KISS & YAGNI)
**Status**: ✅ PASS
**Evidence**: User explicitly stated "won't have any logic other than validation for flag inputs for now". No premature features. Standard library only, no external dependencies.

### Principle VI: Git Commit Standards
**Status**: ✅ PASS
**Evidence**: Will follow Conventional Commits format without attribution. One-line commits by default.

### Principle VII: Package Naming Conventions
**Status**: ✅ PASS
**Evidence**: Will use short package names like `cli`, `command`, `parser` that combine naturally with exported types (e.g., `cli.ParseCommand`).

**Overall Gate Status**: ✅ ALL CHECKS PASS - Proceed to Phase 0

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
codegraph/
├── main.go                  # CLI entry point, flag parsing
├── cli/
│   ├── parse_command.go     # ParseCommand domain entity
│   └── parse_command_test.go
├── path/
│   ├── validator.go         # TargetDirectory validation logic
│   └── validator_test.go
└── go.mod
```

**Structure Decision**: Single project structure. CLI application with domain-driven packages: `cli` for command entities, `path` for directory validation logic. Follows Go convention of flat, small packages with clear responsibilities.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

**Status**: ✅ No violations - No complexity tracking required

All constitution principles are followed without exceptions. This feature uses simple, idiomatic Go with standard library only, small focused functions, and domain-driven entities.

---

## Phase 0: Research (Completed)

**Artifacts Generated**:
- [research.md](./research.md) - Technical decisions, stdlib patterns, best practices

**Key Decisions**:
1. Use `flag.NewFlagSet` for subcommand pattern
2. Use `os.Stat` + `FileInfo.IsDir()` for validation
3. Use `filepath.Abs` for path resolution
4. Use `os.Getwd()` for current directory default

All "NEEDS CLARIFICATION" items from Technical Context resolved.

---

## Phase 1: Design (Completed)

**Artifacts Generated**:
- [data-model.md](./data-model.md) - Domain entities (ParseCommand, TargetDirectory)
- [contracts/cli-interface.md](./contracts/cli-interface.md) - CLI contract specification
- [quickstart.md](./quickstart.md) - Developer and user guides

**Design Summary**:
- Two domain entities with validation behavior
- CLI contract defines exact input/output/error behavior
- Immutable value objects, fail-fast validation
- Clear separation: `cli` package for commands, `path` package for validation

---

## Post-Design Constitution Re-check

### Principle I: Clear Variable Naming
**Status**: ✅ PASS
**Evidence**: Data model uses descriptive names (`TargetDirectory`, `ResolvedPath`, `WasDefaulted`). No abbreviations in design.

### Principle II: Function Size Discipline
**Status**: ✅ PASS
**Evidence**: Each validation rule is separate method. Constructors are focused. Data model shows clear single responsibilities.

### Principle III: Domain-Driven Design
**Status**: ✅ PASS
**Evidence**: Both `ParseCommand` and `TargetDirectory` encapsulate validation logic with data. Methods like `Validate()`, `LogDefaultBehavior()` live on entities.

### Principle IV: Go Idiomatic Practices
**Status**: ✅ PASS
**Evidence**: Design uses constructor pattern (`NewTargetDirectory`), error returns, composition. Data model shows no inheritance patterns.

### Principle V: Simplicity (KISS & YAGNI)
**Status**: ✅ PASS
**Evidence**: Minimal design - only validates inputs, no parsing logic yet. Standard library only. No speculative features.

### Principle VI: Git Commit Standards
**Status**: ✅ PASS
**Evidence**: Implementation will follow Conventional Commits. No attribution in design docs.

### Principle VII: Package Naming Conventions
**Status**: ✅ PASS
**Evidence**: Packages are `cli` and `path` (short, single-word). Types combine naturally: `cli.ParseCommand`, `path.TargetDirectory`.

**Overall Re-check Status**: ✅ ALL CHECKS PASS - Design adheres to constitution

---

## Next Steps

**Phase 2**: Generate implementation tasks using `/speckit.tasks` command
- Task order will follow dependency graph
- Each task implements one focused piece of functionality
- Tests are written alongside implementation

**Implementation**: Execute tasks using `/speckit.implement` command

