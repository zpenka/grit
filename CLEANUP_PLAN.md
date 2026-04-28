# grit Cleanup Plan

A three-workstream plan for tightening up the project. Each workstream is broken into self-contained tasks that can be farmed out to Haiku agents. Tasks are sized to be completable from a single prompt with the file paths, line numbers, and acceptance criteria already specified.

## Current State (baseline, captured on `main`)

- 19,920 lines of Go in the root package; 21,088 across the module
- Test coverage: root **71.1%**, `core/` 83.4%, `cmd/grit` 0%
- 38 `*_test.go` files at the repo root, many named after the work batch they came from rather than the feature they cover
- 68 `show*` boolean flags declared on `model`; only 37 are referenced from `grit.go`'s live key dispatcher
- 5 submenus (`a`/`v`/`t`/`i`/`g`) expose **26** features total, but `engine_rendering.go` defines **62** `render*UI` functions
- `?` help overlay documents the 5 submenus and basics — single-letter shortcuts and most submenu contents are not listed
- README claims "113+ integrated features" — reachable count is much lower

---

## Workstream A — Test coverage to 80% + reorg

### A1. Delete dead test files (mechanical, ~30 min)

Two files masquerade as tests but don't actually run, and the artifacts of running coverage are checked in.

- **`engine_test_coverage.go`** — has 60 `func Test*` declarations but the filename does not end in `_test.go`. The functions never execute, AND because the file imports `testing`, it pulls the test framework into the production binary. Delete the file (the *real* coverage of the symbols it pretends to test is established elsewhere; verify with `go test ./...` after deletion).
- **`coverage.html`**, **`coverage.out`** — generated artifacts checked into git. Delete and add to `.gitignore`.

**Acceptance**: `go build ./cmd/grit` still succeeds, `go test ./...` still passes, binary size shrinks.

### A2. Fold "process-named" test files into feature-named ones (per-file, ~15-30 min each)

Each of these files was created during a particular work batch rather than to cover a feature. Move their tests into the appropriate feature-named file (or create one). When a test duplicates one that already exists in the destination file, drop it.

| File to retire | Likely target file(s) |
|---|---|
| `coverage_gaps_test.go` | split by feature into `parsing_test.go`, `filtering_test.go`, `engine_rendering_test.go` |
| `engine_coverage_test.go` | same — split by what each test exercises |
| `docs_reorganization_test.go` | delete unless it tests live behavior; otherwise fold into a single `docs_test.go` |
| `engine_analytics_extended_test.go` | merge into `engine_analytics_test.go` (create if missing) |
| `engine_team_ai_extended_test.go` | merge into `engine_team_ai_test.go` |
| `grit_bulk_feature_wiring_test.go` | merge into `grit_feature_wiring_test.go` (or new `keybinding_dispatch_test.go`) |
| `grit_integration_improvements_test.go` | merge into a feature-scoped file (`integration_test.go` or `engine_integration_test.go` — create if missing) |
| `grit_menu_improvements_test.go` | merge into `submenus_test.go` (new) |
| `grit_remaining_features_test.go` | distribute by feature |
| `phase5_quality_ux_release_test.go` | distribute by feature |
| `render_consolidation_test.go` | merge into `engine_rendering_test.go` |
| `render_migration_test.go` | merge into `engine_rendering_test.go` |

**Per-file acceptance**: tests still pass; total test count does not decrease except via explicit dedup; coverage does not drop.

### A3. Target test layout (single Haiku task, after A2 is done)

End state:

```
core/                                # already feature-organized, no change
parsing_test.go                      # parseCommits, parseDiff, parseFileItems, parseBranches
filtering_test.go                    # filterCommits, visibleCommits, extension filters
navigation_test.go                   # cursor, panel switching, scrolling
engine_cache_test.go                 # diff cache, statistics cache, regex cache
engine_optimization_test.go          # incremental load, parallel
engine_rendering_test.go             # all render*UI tests (consolidated)
engine_analytics_test.go             # bisect, ownership, hotspots, complexity, etc.
engine_team_ai_test.go               # team stats, AI, compliance, signing
engine_visualization_test.go         # flamegraph, timeline, tree, heatmap
engine_workflows_test.go             # worktrees, stashes, tags, GPG
engine_integration_test.go           # PR/Jira/issue refs, export
engine_git_ops_test.go               # rebase, cherry-pick, reset, amend
keybinding_dispatch_test.go          # grit.go's Update() dispatch table
submenus_test.go                     # a/v/t/i/g menu navigation + dispatch
benchmarks_test.go                   # already focused, leave
integration_test.go                  # multi-feature scenarios only
```

### A4. Cover the zeroes to reach 80% (parallelizable, one Haiku per file)

The fastest path to 80% is to write tests for the functions reporting **0.0%** coverage in `go tool cover -func`. Each task below is a single file, sized for one Haiku.

- **`engine_team_ai.go`** — `dispatchTeamFeature` (line 465), `renderTeamStatsUI` (541), `renderReviewerSuggestionsUI` (553), `renderVelocityUI` (565), `renderClassificationUI` (577), `renderSecretsUI` (589), `renderChangelogUI` (601), `detectSecretsInLine` (613), `renderCherryPickUI` (626), `createMilestone` (442). 10 functions, mostly thin renderers — each gets a fixture-driven test asserting non-empty output and key tokens.
- **`engine_integration.go`** — `exportToMarkdown` (84), `exportPatchSeries` (107). Round-trip-style tests with a small commit fixture.
- **`engine_analytics.go`** — `groupCommits` (886). Table-driven test for the grouping modes.
- **`core/filter.go`** — `filterByDateRange` (115), `parseDaysAgo` (140). Unit tests with synthetic dates.
- **`engine_render_consolidation.go`** — `BuildItemsList` (41). If it's truly unused (see B1), delete it instead of testing it.
- **`cmd/grit/main.go`** — out of scope for unit coverage; leave at 0% (it's a 7-line entry point).

**Acceptance per task**: `go test -cover ./...` shows the touched file's coverage above 80%; no existing tests regress.

### A5. Stretch coverage in big files (do after A4 — measure first)

After A4, re-run `go tool cover -func=cov.out | sort -k3 -n` and pick off the next-lowest functions in `engine_rendering.go` and `grit.go`. Stop once the package hits 80%.

---

## Workstream B — Refactor / dedup

### B1. Delete dead render functions in `engine_rendering.go`

The file is 2,734 lines. The following functions are defined there, take `commits []commit` (the older style), and have **zero callers** anywhere in the codebase. They predate the menu-based architecture and are superseded by the per-feature `render*UI(m model, width int)` versions in `engine_analytics.go`, `engine_team_ai.go`, etc.

```
renderFlameGraphUI            (line 926)   — superseded by renderFlamegraphUI in engine_visualization.go
renderHotspotUI               (line 1220)  — superseded by renderHotspotsUI in engine_analytics.go
renderTimelineUI              (line 970)   — superseded by renderTimelineSliderUI in engine_visualization.go
renderChurnAnalysisUI         (line 1086)
renderExpertiseMapUI          (line 1157)
renderRegressionAnalysisUI    (line 1285)
renderCoverageAnalysisUI      (line 1373)
renderDiffAnalysisUI          (line 1487)
renderAIInsightsUI            (line 1606)
renderPerformanceOptimizationUI (line 1751)
renderGitOperationsUI         (line 1906)
renderRepositoryManagementUI  (line 2060)
renderDeveloperExperienceUI   (line 2157)
renderIntegrationUI           (line 2242)
renderTeamAnalyticsUI         (line 2374)
renderComplianceUI            (line 2527)
renderCommitSigningUI         (line 828)
renderCollaborationUI         (line 890)
renderDependencyGraphUI       (line 933)
renderCommitComparisonUI      (line 1001)
renderSearchUI                (line 1014)
renderAdvancedFilterUI        (line 1026)
```

Plus `BuildItemsList` at `engine_render_consolidation.go:41` (0 callers).

**Task for Haiku**: For each function, run `grep -rn '<funcName>\b' .`, confirm it's defined exactly once and called zero times outside its own file, then delete the definition. After each deletion run `go build ./... && go test ./...`. Expected to remove ~1500 lines from `engine_rendering.go`.

### B2. Delete the dead `handleKeyBinding` dispatcher

`engine_rendering.go:335` defines `handleKeyBinding(m model, key string) model` — explicitly marked deprecated in its own doc comment ("This function is dead code. All keybinding dispatch happens in grit.go's Update() method. Kept for reference only. Do not use.")

It is ~200 lines and contains the rich keymap that was never re-wired (`O`, `H`, `M`, `S`, `X`, `N`, `E`, `Y`, `T`, `D`, `W`, `Z`, `1-9`, `0`, `p`, `j`, `e`). It also contains a `case "q":` that toggles `showIssueRefs`, conflicting with the live "quit" binding — a trap if anyone ever calls it.

**Decision needed before deleting**: Workstream C below proposes resurrecting the *intent* of these shortcuts (single-key access to features). If we proceed with C, port the case bodies into `grit.go`'s `Update()` and *then* delete `handleKeyBinding`. If we don't, just delete it.

**Task**: Once C1 is decided, delete the function and any imports it alone needs. Confirm `go test ./...` and `go vet ./...` clean.

### B3. Reduce the rendering file

`engine_rendering.go` will be ~1100 lines after B1+B2. Re-evaluate whether what remains belongs in dedicated files (e.g. main commit-list rendering can stay; submenu rendering can move to a `engine_submenus.go`). Defer this until B1/B2 land — don't pre-shuffle.

### B4. `model` struct flag audit (do after C2 lands)

After Workstream C inventories which features are reachable, delete `show*` flags + their associated `*Data` fields for any feature we explicitly choose not to expose. Goal is to cut the 200+ field model down to roughly the surface area we actually ship.

---

## Workstream C — Discoverability / making features usable

### C1. Inventory: every feature → reachability status (research only, no code)

Single Haiku task. Produce `docs/FEATURE_INVENTORY.md` with one row per `show*` flag (68 of them). Columns:

- flag name
- feature description (one line)
- live keybinding (if any) — read from `grit.go` `Update()` only, NOT `handleKeyBinding`
- submenu it appears in (if any) — read from `dispatch*Feature` in `engine_analytics.go:139`, `engine_visualization.go:255`, `engine_team_ai.go:465`, `engine_integration.go:184`, `engine_git_ops.go:190`
- documented in `?` overlay (`grit.go:1102`)?  Y/N
- documented in README keybinding table?  Y/N
- recommendation: `keep / promote-to-menu / promote-to-shortcut / drop`

This becomes the source of truth for C2/C3/B4.

### C2. Promote orphaned features into existing submenus

Findings from the survey: the following declared `show*` flags do not appear anywhere in `grit.go`'s live dispatcher:

```
showAmendPreview     showAnomalies        showBlamePerf
showDataRequests     showFileTimeline     showGPGStatus
showGraph            showLargeCommits     showLicenses
showLoadProgress     showLostCommits      showMemoryMetrics
showMergeAnalysis    showMilestones       showNamedStashes
showRebasePreview    showReleaseNotes     showSecurityIssues
showSemanticSearch   showSemver           showSigningStatus
showSimilar          showSquashUI         showStatsBadge
showSubmodules       showSummaries        showTagMgmt
showTags             showUndoMenu         showVersionBumps
showWorktrees
```

After C1 categorizes these, group the `keep`s into existing submenus and extend each menu's `dispatch*Feature` and `*MenuLen`:

- **Analytics** (`engine_analytics.go:126`): add `showLargeCommits`, `showLostCommits`, `showSemanticSearch`, `showMergeAnalysis`, `showAnomalies`, `showSimilar`
- **Visualization** (`engine_visualization.go:244`): add `showFileTimeline`, `showStatsBadge`
- **Team/AI** (`engine_team_ai.go:453`): add `showSecurityIssues`, `showLicenses`, `showSigningStatus`, `showSummaries`
- **Integration** (`engine_integration.go:155`): add `showReleaseNotes`, `showMilestones`, `showVersionBumps`
- **Git Ops** (`engine_git_ops.go:163`): add `showAmendPreview`, `showRebasePreview`, `showSquashUI`, `showUndoMenu`
- **New "Workflows" submenu** (key suggestion: `w`): `showWorktrees`, `showSubmodules`, `showNamedStashes`, `showTagMgmt`, `showGPGStatus`

Each addition is one Haiku task: extend the dispatch function, bump `MenuLen`, add a row to the menu render config, write a test in `submenus_test.go`.

**Acceptance per task**: launching the menu, navigating to the new entry, and pressing Enter toggles the feature; the menu's existing tests still pass.

### C3. Rebuild the `?` help overlay

`grit.go:1102` is hand-rolled and stale. Replace with a generator that walks the same data structure used by C2 (so they can never drift). Layout:

```
NAVIGATION   ...
VIEWING      ...
SUBMENUS     a Analytics → ownership, hotspots, bisect, ...
             v Visualizations → flamegraph, timeline, ...
             t Team/AI → ...
             i Integration → ...
             g Git Ops → ...
             w Workflows → worktrees, submodules, ... (new)
HELP & EXIT  ...
```

If the user is currently inside a submenu, the `?` overlay should show that submenu's contents instead of the global list.

**Acceptance**: snapshot test in `engine_rendering_test.go` confirms every menu's items appear in the help overlay output.

### C4. Sync README + DEVELOPER.md to reality

The README's keybinding table (`README.md:55-82`) only documents the 5 submenus and basics. After C2/C3 land, regenerate the table from the same source-of-truth data. Mention the new Workflows menu. Drop the "113+ features" claim in the header in favor of an accurate count derived from the inventory in C1.

---

## Suggested Ordering

The workstreams are mostly independent, but a few things should land in order:

1. **A1** (delete dead test files + checked-in coverage artifacts) — cheap, removes confusion for everyone else.
2. **C1** (feature inventory) — read-only, unblocks B4 and C2 with real data.
3. **B1** (delete dead render funcs) — independent; shrinks the search space.
4. **A2 + A4** in parallel (Haiku-friendly: each file/cluster is its own task).
5. **C2** (promote features), then **C3** (rebuild help), then **C4** (sync docs).
6. **B2** (delete `handleKeyBinding`) — only after C2 has made the decisions about which shortcuts to revive.
7. **A3** (final layout sweep), **B3** (split rendering file), **B4** (model audit) — cleanup passes.
8. **A5** (chase 80%) — measure after everything else; only do as much as needed.

## Out of Scope

- New features. This plan is strictly cleanup.
- Changing the public CLI surface (`grit -v`, `grit -h`).
- Touching the `core/` package's structure — it's already feature-organized.
- Any rewrite of `grit.go`'s Bubble Tea integration. We're cleaning up around it, not replacing it.
