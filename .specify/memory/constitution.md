<!--
Sync Impact Report:
Version: 1.0.0 → 1.1.0
Rationale: MINOR version bump - Added new principle (Dependency Management)

Modified principles:
- None (existing principles unchanged)

Added sections:
- NEW: Principle VIII - Dependency Management (Standard Library First)

Removed sections:
- None

Templates requiring updates:
✅ plan-template.md - No constitution-specific updates needed (Constitution Check section already present)
✅ spec-template.md - No constitution-specific updates needed
✅ tasks-template.md - No constitution-specific updates needed

Follow-up TODOs: None
-->

# CodeGraph Constitution

## Core Principles

### I. Clear Variable Naming

Variable names MUST be descriptive and convey intent without requiring additional context. Avoid abbreviations, single-letter names (except standard loop indices), and ambiguous terms.

**Rationale**: Code is read far more often than written. Clear names eliminate cognitive load and reduce bugs from misunderstood state.

### II. Function Size Discipline

Functions MUST be small to medium-sized. A function should do one thing well. If a function requires scrolling to understand, it should be refactored.

**Rationale**: Small functions are easier to test, understand, debug, and reuse. They enforce single responsibility and reduce cyclomatic complexity.

### III. Domain-Driven Design (No Anemic Model)

Domain logic MUST reside in domain entities, not scattered across services. Anemic models (data structures with no behavior) are prohibited unless justified.

**Rationale**: Behavior belongs with data. Anemic models lead to procedural code disguised as OOP, making business logic hard to locate and maintain.

### IV. Go Idiomatic Practices

Go code MUST follow Go paradigms, NOT object-oriented or functional programming patterns from other languages. Adapt concepts to Go's design philosophy.

**Rules**:
- Use `any` instead of `interface{}` for empty interfaces
- Use modern `for range` syntax where applicable
- Embrace composition over inheritance
- Prefer interfaces for behavior contracts, structs for data
- Keep error handling explicit; avoid exception-like patterns
- Use goroutines and channels idiomatically, not as thread/promise replacements

**Rationale**: Go has deliberate design decisions. Forcing other paradigms creates unidiomatic, hard-to-maintain code that fights the language.

### V. Simplicity Principles (KISS & YAGNI) (NON-NEGOTIABLE)

**KISS (Keep It Simple, Stupid)**: Solutions MUST be as simple as possible, but no simpler. Complexity requires explicit justification.

**YAGNI (You Aren't Gonna Need It)**: Do not implement features, abstractions, or infrastructure until they are actually needed. No speculative generalization.

**Enforcement**:
- All abstractions must solve a current, concrete problem
- Premature optimization is forbidden without profiling data
- Generic solutions allowed only when multiple concrete cases exist

**Rationale**: Complexity is the enemy of maintainability. Most "future-proofing" solves problems that never materialize while creating immediate maintenance burden.

### VI. Git Commit Standards (NON-NEGOTIABLE)

Commits MUST follow Conventional Commits format without author attribution.

**Format**: `<type>(<scope>): <subject>`

**Rules**:
- One-line commits (default); multi-line only when explicitly needed
- NO "authored by", "co-authored by", or LLM attribution
- NO mentions of tools, assistants, or generation method
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `style`, `build`, `ci`

**Example**:
```
feat(parser): add support for multi-line strings
fix(api): handle null pointer in user lookup
refactor(storage): simplify cache invalidation logic
```

**Rationale**: Commit history is a project artifact, not a credits roll. Consistent format enables tooling (changelogs, semver automation). Attribution lives in Git metadata (author field), not messages.

### VII. Package Naming Conventions

Package names MUST be short, meaningful, and designed to work with exported types to form readable identifiers.

**Rules**:
- Prefer single-word package names: `parser`, `graph`, `index`, `cache`
- Avoid `util`, `common`, `helper`, `lib`, `base` (unless specific purpose)
- Package name + type name should read naturally: `graph.Node`, `cache.Entry`, `parser.Token`
- No underscores or mixed caps: `httputil` (not `http_util` or `httpUtil`)

**Example**:
```go
// Good
import "codegraph/parser"
token := parser.Token{}

// Bad
import "codegraph/parser_utilities"
token := parser_utilities.ParserToken{}
```

**Rationale**: Go encourages package-qualified identifiers. Short, descriptive package names reduce stutter (`parser.Parser` is fine, `parser.ParserParser` is not) and improve readability across the codebase.

### VIII. Dependency Management (Standard Library First)

Favor the Go standard library for all functionality. External dependencies MUST be justified and used only when the standard library is insufficient or would require excessive implementation effort.

**Rules**:
- Prefer `net/http` over web frameworks unless complex routing/middleware required
- Prefer `encoding/json`, `encoding/xml` over third-party serialization libraries
- Prefer `go/parser`, `go/ast` over external code analysis tools
- Use standard `testing` package unless advanced features (mocking, BDD) are essential
- Document dependency justification in commit message or design docs

**When External Dependencies Are Acceptable**:
- Standard library lacks the functionality entirely (e.g., GraphQL, gRPC code generation)
- Implementing from scratch would create significant maintenance burden (e.g., JWT parsing, complex cryptography)
- Well-established, actively maintained libraries that solve complex domain problems (e.g., database drivers, cloud SDKs)

**Example Justification**:
```go
// Using standard library (preferred)
import "net/http"

// Using external dependency (must justify)
import "github.com/gorilla/mux"  // JUSTIFIED: Complex routing with middleware chaining and path parameters beyond net/http ServeMux
```

**Rationale**: External dependencies introduce maintenance burden, security risks, version conflicts, and upgrade complexity. The Go standard library is stable, well-documented, and maintained by the Go team. Minimizing dependencies reduces attack surface and simplifies long-term maintenance.

## Governance

This constitution supersedes all other coding practices and preferences. All code reviews, pull requests, and design decisions MUST verify compliance with these principles.

**Amendment Process**:
1. Proposed changes require documented rationale and impact analysis
2. Violations of principles must be explicitly justified in code/PR comments
3. Constitution updates require version bump per semantic versioning

**Complexity Justification**:
When a principle must be violated (e.g., performance-critical code requires complexity), document in code comments:
```go
// COMPLEXITY JUSTIFIED: Profiling shows 40% CPU in this path.
// Simpler recursive approach causes stack overflow at 10k+ nodes.
```

**Compliance Review**:
- PRs flagged for principle violations MUST include justification or remediation
- Regular audits to identify drift from principles
- Constitution alignment checked in planning phase (see plan-template.md)

**Version**: 1.1.0 | **Ratified**: 2025-10-19 | **Last Amended**: 2025-10-19
