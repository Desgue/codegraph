# Tasks: Go Source Parser

**Input**: Design documents from `/specs/002-go-source-parser/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md, contracts/parser-api.md

**Tests**: No tests generated - validation-only implementation per user request.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- Single Go project at repository root
- Package structure: `parser/` for new parsing package
- CLI integration: `cli/parse.go` (existing)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency setup

- [X] T001 Add golang.org/x/tools/go/packages dependency via go get
- [X] T002 [P] Create parser package directory structure at parser/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core parsing infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Create packages.Config setup with LoadMode flags in parser/loader.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Parse Entire Project (Priority: P1) 🎯 MVP

**Goal**: Discover and parse all Go packages in a directory, reporting statistics and handling partial failures gracefully

**Independent Test**: Run parse command on a Go project directory with valid and invalid files, verify all accessible packages are discovered, statistics are accurate, and errors are clearly reported

### Implementation for User Story 1

- [X] T004 [US1] Implement Load() function with packages.Load call in parser/loader.go
- [X] T005 [US1] Implement error handling for catastrophic vs partial failures in parser/loader.go
- [X] T006 [US1] Implement packages.PrintErrors() call for stderr output in parser/loader.go
- [X] T007 [US1] Return []*packages.Package directly from Load() per YAGNI principle in parser/loader.go
- [X] T008 [US1] Integrate parser.Load() into ParseCommand.Execute() in cli/parse_command.go
- [X] T009 [US1] Add statistics output formatting in ParseCommand.Execute() in cli/parse_command.go
- [X] T010 [US1] Add error count reporting to stderr in ParseCommand.Execute() in cli/parse_command.go

**Checkpoint**: At this point, User Story 1 should be fully functional - can parse projects, report statistics, handle errors gracefully

---

## Phase 4: User Story 2 - Module-Aware Parsing (Priority: P2)

**Goal**: Respect Go module boundaries and resolve import paths according to go.mod structure

**Independent Test**: Run parse command on a project with go.mod and verify package import paths are correctly resolved according to module definition

### Implementation for User Story 2

- [X] T011 [US2] Verify packages.Load automatically detects go.mod by manual testing with module-based project
- [X] T012 [US2] Add module path validation to result output in cli/parse_command.go
- [X] T013 [US2] Update statistics output to show module information when present in cli/parse_command.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - module-aware parsing respects go.mod boundaries

**Note**: packages.Load with NeedImports automatically handles module boundary detection via go.mod files

---

## Phase 5: User Story 3 - Documentation Preservation (Priority: P3)

**Goal**: Preserve documentation comments during parsing for future documentation analysis features

**Independent Test**: Parse a file with package-level and function-level comments, verify comment data is retained in parsed AST structure

### Implementation for User Story 3

- [X] T014 [US3] Verify packages.Load with NeedSyntax preserves comment nodes in AST
- [X] T015 [US3] Validate comment preservation by manual inspection of AST comment nodes in parsed output
- [X] T016 [US3] Document comment access patterns for future graph features

**Checkpoint**: All user stories should now be independently functional - comments are preserved in AST for future use

**Note**: packages.Load with NeedSyntax automatically preserves comment nodes in AST structure

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T017 Verify performance target (<5s for 50k LOC) with benchmark test
- [X] T018 Add validation for 1000-file project support
- [X] T019 Run manual validation: (1) valid project, (2) project with syntax errors, (3) nested packages, (4) empty directory, (5) module-based project
- [X] T020 Code cleanup and ensure Go idiomatic practices

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3, 4, 5)**: All depend on Foundational phase completion
  - User Story 1 (P1): Can start after Foundational - No dependencies on other stories
  - User Story 2 (P2): Can start after US1 complete - Validates module-aware behavior on top of basic parsing
  - User Story 3 (P3): Can start after US1 complete - Validates comment preservation on top of basic parsing
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after US1 complete - Validates module-aware behavior on top of basic parsing
- **User Story 3 (P3)**: Can start after US1 complete - Validates comment preservation on top of basic parsing

### Within Each User Story

- Core implementation before integration with CLI
- CLI integration before statistics output
- Statistics output before error reporting

### Parallel Opportunities

- Setup tasks T001 and T002 can run in parallel
- Foundational phase T003 completes setup
- Once Foundational phase completes, User Story phases proceed sequentially (US1 → US2 → US3)

---

## Parallel Example: Setup Phase

```bash
# Launch setup tasks together:
Task: "Add golang.org/x/tools/go/packages dependency via go get"
Task: "Create parser package directory structure at parser/"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (tasks T001-T002)
2. Complete Phase 2: Foundational (task T003) - CRITICAL - blocks all stories
3. Complete Phase 3: User Story 1 (tasks T004-T010)
4. **STOP and VALIDATE**: Test User Story 1 independently with quickstart scenarios
5. Deploy/demo if ready - basic parsing works end-to-end

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (T004-T010) → Test independently → MVP complete (basic parsing works)
3. Add User Story 2 (T011-T013) → Test independently → Module-aware parsing works
4. Add User Story 3 (T014-T016) → Test independently → Comment preservation works
5. Each story adds value without breaking previous stories

### Sequential Strategy (Single Developer)

1. Complete Setup + Foundational together (T001-T003)
2. Complete User Story 1 fully before moving to US2 (T004-T010)
3. Complete User Story 2 fully before moving to US3 (T011-T013)
4. Complete User Story 3 fully (T014-T016)
5. Complete Polish phase (T017-T020)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- No tests generated per user request - validation-only implementation
- Return []*packages.Package directly from Load() per constitution (YAGNI/KISS)
- Integration point: ParseCommand.Execute() orchestrates parser package calls
- Statistics tracking happens in CLI layer, not parser layer
- Commit after each task or logical group following conventional commits
- Stop at any checkpoint to validate story independently
