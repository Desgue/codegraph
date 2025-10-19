# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CodeGraph is a Go-based code analysis and visualization tool managed using specification-driven development via the Specify framework. The project emphasizes clean architecture, idiomatic Go, and domain-driven design.

## Development Workflow

This project uses specification-driven development through `.specify/` and custom slash commands via `.claude/commands/`:

### Slash Commands

- `/commit` - Analyze unstaged changes, group logically, and commit following constitution guidelines
- `/speckit.specify` - Create/update feature specification from natural language description
- `/speckit.plan` - Execute implementation planning workflow to generate design artifacts
- `/speckit.tasks` - Generate actionable, dependency-ordered tasks.md from design artifacts
- `/speckit.implement` - Execute implementation plan by processing all tasks
- `/speckit.clarify` - Identify underspecified areas and ask clarification questions
- `/speckit.analyze` - Perform cross-artifact consistency analysis
- `/speckit.constitution` - Create/update project constitution
- `/speckit.checklist` - Generate custom checklist for current feature

### Feature Development Flow

1. Create specification: `/speckit.specify [feature description]`
2. Generate implementation plan: `/speckit.plan` (creates research, design artifacts in `specs/[###-feature]/`)
3. Generate tasks: `/speckit.tasks` (creates dependency-ordered tasks.md)
4. Implement: `/speckit.implement` (executes tasks systematically)
5. Commit changes: `/commit`

All feature specifications and plans live under `specs/[###-feature-name]/` directory.

## Project Constitution

The project follows strict principles defined in `.specify/memory/constitution.md` (v1.0.0):

### Code Standards

1. **Variable naming**: Descriptive names required, no abbreviations
2. **Function size**: Keep functions small and focused (single responsibility)
3. **Domain logic**: No anemic models - behavior belongs with data
4. **Go idioms**:
   - Use `any` over `interface{}`
   - Use modern `for range` syntax
   - Embrace composition over inheritance
   - Explicit error handling
5. **Simplicity (KISS & YAGNI)**: No speculative features or premature optimization
6. **Package names**: Short, single-word preferred (avoid `util`, `common`, `helper`)

### Commit Standards (NON-NEGOTIABLE)

Follow Conventional Commits format:
```
<type>(<scope>): <subject>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `style`, `build`, `ci`

**Rules**:
- One-line commits by default
- NO author attribution (no "authored by", "co-authored by", or LLM mentions)
- NO mentions of tools or generation method
- Focus on what and why, not who

Examples:
```
feat(parser): add support for multi-line strings
fix(api): handle null pointer in user lookup
refactor(storage): simplify cache invalidation logic
```

### Complexity Justification

When violating principles (e.g., performance-critical code), document:
```go
// COMPLEXITY JUSTIFIED: Profiling shows 40% CPU in this path.
// Simpler recursive approach causes stack overflow at 10k+ nodes.
```

## Architecture Approach

- **Domain-Driven Design**: Business logic in domain entities, not scattered in services
- **Idiomatic Go**: Follow Go conventions, not OOP/FP patterns from other languages
- **Constitution compliance**: All design decisions must align with constitution principles
- **Documentation-first**: Specifications drive implementation, not the reverse

## Git Workflow

- Main branch: `main`
- Feature branches: `[###-feature-name]` format
- All commits follow Conventional Commits
- Use `/commit` command to ensure proper grouping and formatting
