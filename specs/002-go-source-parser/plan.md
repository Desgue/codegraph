# Implementation Plan: Go Source Parser

**Branch**: `002-go-source-parser` | **Date**: 2025-10-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-go-source-parser/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Parse and load Go source code from target directories using `golang.org/x/tools/go/packages` to discover packages, parse files, and prepare structured data for graph analysis. The system will handle partial failures gracefully, report clear statistics and errors to the terminal, and retain parsed data in memory for the next phase.

## Technical Context

**Language/Version**: Go 1.24.5
**Primary Dependencies**: `golang.org/x/tools/go/packages` (standard Go tooling library)
**Storage**: N/A (in-memory data structure)
**Testing**: Standard `testing` package
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)
**Project Type**: Single CLI project
**Performance Goals**: Parse 50,000 LOC in under 5 seconds on standard hardware
**Constraints**: <5s parse time for 50k LOC, support up to 1000 source files
**Scale/Scope**: Up to 1000 Go source files per project

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clear Variable Naming | ✅ PASS | Domain entities (Package, SourceFile, ParseResult) are descriptive |
| II. Function Size Discipline | ✅ PASS | Parser functions will be small, focused on single responsibilities |
| III. Domain-Driven Design | ✅ PASS | Using `parser` package with behavior (Load, Parse methods), not anemic models |
| IV. Go Idiomatic Practices | ✅ PASS | Using `any`, modern `for range`, `go/packages` standard tooling |
| V. Simplicity (KISS/YAGNI) | ✅ PASS | No premature abstractions; using `[]*packages.Package` directly, no wrapper types |
| VI. Git Commit Standards | ✅ PASS | Conventional commits without attribution |
| VII. Package Naming | ✅ PASS | `parser` package name is short, meaningful, readable with types (parser.Result) |
| VIII. Dependency Management | ✅ PASS | Using `golang.org/x/tools/go/packages` (official Go tools library, not standard lib but well-justified for AST/package loading) |

**Decision**: Proceed to Phase 0. `golang.org/x/tools/go/packages` is justified as the standard approach for Go tooling—it handles module resolution, build constraints, and AST parsing comprehensively, avoiding reimplementation of complex logic.

## Project Structure

### Documentation (this feature)

```
specs/002-go-source-parser/
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
├── main.go              # CLI entry point with subcommand routing
├── cli/
│   ├── parse.go         # ParseCommand implementation (existing)
│   └── parse_test.go    # ParseCommand tests (existing)
├── path/
│   ├── target.go        # TargetDirectory validation (existing)
│   └── target_test.go   # Path validation tests (existing)
├── parser/              # NEW: Go source parsing package
│   ├── loader.go        # Package loading using go/packages
│   └── loader_test.go   # Loader tests
├── go.mod
└── go.sum
```

**Structure Decision**: Single project structure. The existing CLI command (`cli/parse.go`) will integrate with the new `parser` package. The `parser` package encapsulates all parsing logic, keeping the CLI layer thin and focused on user interaction. This follows Go's idiomatic flat package structure with clear separation of concerns.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

N/A - No constitution violations.

## Post-Design Constitution Re-Check

*Re-evaluated after Phase 1 design completion*

| Principle | Status | Post-Design Validation |
|-----------|--------|------------------------|
| I. Clear Variable Naming | ✅ PASS | Result, Load, Packages, TotalPackages are all descriptive |
| II. Function Size Discipline | ✅ PASS | `Load()` is ~20 lines, focused on single responsibility (load and aggregate) |
| III. Domain-Driven Design | ✅ PASS | `Result` contains aggregation logic, `Load()` encapsulates loading behavior |
| IV. Go Idiomatic Practices | ✅ PASS | Uses `packages.Load()`, `packages.PrintErrors()`, no wrapper abstractions |
| V. Simplicity (KISS/YAGNI) | ✅ PASS | Minimal code: 2 files (loader.go, result.go), no abstractions, returns `[]*packages.Package` directly |
| VI. Git Commit Standards | ✅ PASS | Will use conventional commits (e.g., `feat(parser): add Go source parsing`) |
| VII. Package Naming | ✅ PASS | Package `parser`, types `parser.Result`, `parser.Load()` read naturally |
| VIII. Dependency Management | ✅ PASS | Single external dependency justified: `golang.org/x/tools/go/packages` is the standard for Go tooling |

**Final Decision**: Design approved. All constitution principles satisfied. Ready for implementation via `/speckit.tasks`.

