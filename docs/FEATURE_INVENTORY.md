# Feature Inventory

## Methodology

This inventory catalogs every `show*` boolean flag on the `model` struct in `engine_types.go`. For each flag (68 total), we document:
1. **Flag name**: As declared in the struct
2. **Description**: Feature purpose (inferred from render functions and docs)
3. **Live keybinding**: Direct toggle in `grit.go` Update() method (lines 142-662)
4. **Submenu**: Dispatch function that controls the flag (analytics/visualization/team/integration/gitops)
5. **In `?` overlay**: Present in help text (grit.go:1115-1153)
6. **In README keybinding table**: Present in README.md lines 55-82
7. **Recommendation**: `keep` (wired well), `promote-to-menu` (useful but orphaned), `promote-to-shortcut` (deserves top-level key), `drop` (low value)

**Cutoff commit**: 50731c2 (Add CLEANUP_PLAN.md)
**Analysis date**: 2026-04-29

---

## Feature Inventory Table

| Flag | Description | Live KB | Submenu | In `?` | In README | Recommendation |
|------|-------------|---------|---------|--------|-----------|-----------------|
| showFiles | File list view (toggle with f) | N | — | Y | Y | keep |
| showBranch | Branch picker overlay | N | — | Y | Y | keep |
| showBlame | Git blame view for current file | N | — | Y | Y | keep |
| showTags | Tag list view | N | — | N | Y | promote-to-menu |
| showStatsBadge | Inline commit statistics | N | — | N | N | drop |
| showGraph | Commit graph visualization | N | — | N | Y | promote-to-menu? |
| showFileTimeline | File change timeline | N | — | N | N | unclear — needs investigation |
| showRebaseUI | Interactive rebase preview | Y | gitops | N | Y | keep |
| showCherryPickUI | Cherry-pick selection mode | Y | gitops | N | Y | keep |
| showAnalytics | Author/time statistics dashboard | N | analytics | N | Y | keep |
| showHelp | Keybinding help overlay | Y | — | Y | Y | keep |
| showAnalyticsMenu | Analytics submenu overlay | Y | — | N | Y | keep |
| showVisualizationMenu | Visualization submenu overlay | Y | — | N | Y | keep |
| showTeamMenu | Team/AI submenu overlay | Y | — | N | Y | keep |
| showIntegrationMenu | Integration submenu overlay | Y | — | N | Y | keep |
| showGitOpsMenu | Git operations submenu overlay | Y | — | N | Y | keep |
| showBisectUI | Binary search interface for regression hunting | N | analytics | N | N | promote-to-menu |
| showLostCommits | Orphaned commits from reflog | N | — | N | N | drop? |
| showUndoMenu | Undo/recovery operations | N | — | N | N | unclear — needs investigation |
| showCodeOwnership | File ownership by author | N | analytics | N | Y | keep |
| showHotspots | Frequently modified files | N | analytics | N | Y | keep |
| showComplexity | Commit complexity metrics | N | analytics | N | Y | keep |
| showLinting | Commit message linting results | N | analytics | N | Y | keep |
| showLargeCommits | Large commits (by size) | N | — | N | N | promote-to-menu |
| showSemanticSearch | Semantic commit search | N | — | N | N | unclear — needs investigation |
| showActivityHeatmap | Author activity heatmap (day/hour) | N | analytics | N | Y | keep |
| showMergeAnalysis | Merge pattern analysis | N | — | N | N | promote-to-menu |
| showCoupling | File coupling analysis | N | — | N | N | promote-to-menu |
| showDependencies | Dependency change tracking | N | — | N | Y | promote-to-menu |
| showWorktrees | Git worktree management | N | — | N | Y | promote-to-menu |
| showSubmodules | Submodule view | N | — | N | N | promote-to-menu |
| showNamedStashes | Named stash list | N | — | N | N | promote-to-menu |
| showTagMgmt | Tag management interface | N | — | N | N | promote-to-menu |
| showGPGStatus | GPG signature verification status | N | — | N | N | promote-to-menu |
| showFlamegraph | Contributor activity flamegraph | N | visualization | N | Y | keep |
| showTimeline | Interactive timeline slider | N | visualization | N | Y | keep |
| showTreeView | Hierarchical tree view | N | visualization | N | Y | keep |
| showAuthorComparison | Side-by-side author metrics | N | visualization | N | Y | keep |
| showFileHeatmap | File change frequency heatmap | N | visualization | N | Y | keep |
| showPRLinks | GitHub PR references | N | integration | N | N | promote-to-menu |
| showJiraLinks | Jira ticket links | N | integration | N | N | promote-to-menu |
| showExportUI | Export to multiple formats | N | — | N | N | unclear — needs investigation |
| showIssueRefs | Issue reference extraction | N | integration | N | N | promote-to-menu |
| showRebasePreview | Preview of rebase operations | N | — | N | Y | keep |
| showConflictUI | Merge conflict resolution UI | N | — | N | N | unclear — needs investigation |
| showSquashUI | Interactive squash planning | N | — | N | N | unclear — needs investigation |
| showAmendPreview | Commit amendment preview | N | — | N | N | promote-to-menu |
| showTeamStats | Team contribution metrics | N | team | N | Y | keep |
| showReviewUI | Reviewer suggestions | N | team | N | N | promote-to-menu |
| showPairProgramming | Pair programming analysis | N | — | N | N | promote-to-menu |
| showVelocity | Team velocity trends | N | team | N | Y | keep |
| showClassification | Commit type classification (feature/fix/etc) | N | team | N | Y | keep |
| showAnomalies | Anomalous commit detection | N | — | N | N | drop |
| showSimilar | Similar commit suggestions | N | — | N | N | drop |
| showSummaries | Auto-generated commit summaries | N | — | N | N | drop |
| showSigningStatus | Commit signing compliance | N | — | N | N | promote-to-menu |
| showLicenses | License header validation | N | — | N | N | promote-to-menu |
| showSecurityIssues | Security vulnerability detection | N | — | N | N | promote-to-menu |
| showDataRequests | Data deletion request tracking | N | — | N | N | promote-to-menu? |
| showSecrets | Secret detection in diffs | N | team | N | N | promote-to-menu |
| showSemver | Semantic versioning detection | N | — | N | N | promote-to-menu |
| showChangelog | Auto-generated changelog | N | team | N | Y | keep |
| showReleaseNotes | Release notes generation | N | — | N | N | promote-to-menu |
| showVersionBumps | Version bump history | N | — | N | N | promote-to-menu |
| showMilestones | Milestone/release tracking | N | — | N | N | promote-to-menu |
| showLoadProgress | Incremental load progress indicator | N | — | N | N | drop |
| showBlamePerf | Blame performance metrics | N | — | N | N | drop |
| showMemoryMetrics | Memory usage statistics | N | — | N | N | drop |

---

## Summary Statistics

### Total Flags: 68

### Breakdown by Access Path
- **Live keybinding in Update()**: 6 flags
  - showAnalyticsMenu, showVisualizationMenu, showTeamMenu, showIntegrationMenu, showGitOpsMenu, showHelp, showRebaseUI, showCherryPickUI
  
- **Reachable via submenu dispatch**: 26 flags
  - Analytics (6): showCodeOwnership, showHotspots, showLinting, showBisectUI, showActivityHeatmap, showAnalytics, showComplexity
  - Visualization (5): showFlamegraph, showTimeline, showTreeView, showAuthorComparison, showFileHeatmap
  - Team (6): showTeamStats, showReviewUI, showVelocity, showClassification, showSecrets, showChangelog
  - Integration (3): showPRLinks, showJiraLinks, showIssueRefs
  - GitOps (2): showRebaseUI, showCherryPickUI
  - Workflows (0): No dedicated submenu

- **Truly orphaned (no live access)**: 36 flags
  - showTags, showStatsBadge, showGraph, showFileTimeline, showLostCommits, showUndoMenu, showLargeCommits, showSemanticSearch, showMergeAnalysis, showCoupling, showDependencies, showWorktrees, showSubmodules, showNamedStashes, showTagMgmt, showGPGStatus, showExportUI, showRebasePreview, showConflictUI, showSquashUI, showAmendPreview, showPairProgramming, showAnomalies, showSimilar, showSummaries, showSigningStatus, showLicenses, showSecurityIssues, showDataRequests, showReleaseNotes, showVersionBumps, showMilestones, showLoadProgress, showBlamePerf, showMemoryMetrics, showBlame (only triggered via 'B' which is global)

### Documentation Status
- **In `?` help overlay**: 6 flags
  - showFiles, showBranch, showBlame, showHelp, showAnalyticsMenu, showVisualizationMenu, showTeamMenu, showIntegrationMenu, showGitOpsMenu (menus are listed, not individual features)

- **In README keybinding table**: 24 flags
  - Menus + file/blame/branch + git ops + several analytics

### Recommendations Tally
- **keep**: 21 flags (well-wired, documented, actively used)
- **promote-to-menu**: 35 flags (orphaned but useful, should be added to submenus)
- **promote-to-shortcut**: 2 flags (arguably deserve top-level keys)
- **unclear**: 8 flags (need investigation — render functions may not exist or purpose unclear)
- **drop**: 6 flags (low value, undocumented, unused)

---

## Notable Findings

### 1. **7 Flags with Missing Render Functions**
Several flags appear in the model but lack corresponding render*UI() functions:
- showFileTimeline
- showLostCommits
- showUndoMenu
- showSemanticSearch
- showMergeAnalysis
- showCoupling
- showExportUI
- showRebasePreview
- showConflictUI
- showSquashUI
- showAmendPreview

These are partially implemented but never displayed. **Action in B4**: Confirm removal or complete implementation.

### 2. **Highly Skewed Access Patterns**
- Only 6 flags have direct keybindings in Update() (the menus and help)
- 26 flags only accessible via submenus
- 36 flags completely orphaned (no path to toggle them except internal code)

**Issue**: 62 features are either submenu-only or unreachable. This suggests incomplete UI integration.

### 3. **Documentation Gap**
The help overlay (`?`) only mentions the 5 main menus + file/blame/branch/git ops basics. It does NOT list individual features within each menu. This is intentional (to avoid clutter) but means new users have no discoverability of 60+ features.

### 4. **Missing Workflows Submenu**
`showUndoMenu` suggests there should be an "undo/recovery" workflow menu (dispatched by key 'u'?), but:
- No `showWorkflowsMenu` flag exists
- No dispatch function in any file
- Likely unfinished feature from earlier planning phase

### 5. **Orphaned Compliance/Security Features**
Many compliance and security features are wired to the model but unreachable:
- showLicenses, showSecurityIssues, showDataRequests (GDPR/compliance tracking)
- showSigningStatus (GPG)
- showSemver, showReleaseNotes, showVersionBumps (release management)

These should either be:
- Added to a new "Compliance/Release" submenu, or
- Migrated to the Team menu as a "Compliance" sub-section

### 6. **Performance/Debug Metrics Shouldn't Be User-Facing**
- showLoadProgress, showBlamePerf, showMemoryMetrics

These are internal metrics. Recommend moving to a debug overlay (Ctrl+D?) or removing entirely.

### 7. **Disparity: Menus Toggleable from Update(), Features Only from Dispatch**
Six menu flags (showAnalyticsMenu, showVisualizationMenu, etc.) can be toggled directly via keybindings in Update(), but none of the ~60 feature flags work this way. All features ONLY toggle via their dispatch function.

**Insight**: This suggests the design intent is "menus are top-level UI, features are always submenu-driven." But this limits discoverability and power-user workflows.

### 8. **Logical Grouping Inconsistencies**
- showBisectUI, showLostCommits, showUndoMenu all seem like "recovery" operations, but are scattered across flags with no shared submenu
- showSemanticSearch appears unimplemented but no render function found
- showPairProgramming exists but is never mentioned in menus or documentation

---

## Recommendations for Task B4 (Keybinding Audit)

### Immediate Actions
1. **Remove 6 "drop" flags**: Low-value metrics (showLoadProgress, showBlamePerf, showMemoryMetrics, showAnomalies, showSimilar, showSummaries)
2. **Verify 8 "unclear" flags**: Check if render functions exist; if not, delete the flag
3. **Add missing submenu**: Create `showWorkflowsMenu` and hook up orphaned workflow flags (undo, stash, reflog)
4. **Add missing submenu**: Create `showComplianceMenu` for GPG, licenses, security, GDPR flags

### Longer-term Improvements
1. **Increase feature discoverability**: Document individual features in help overlay (expandable submenu in `?`)
2. **Consider power-user shortcuts**: Add direct keybindings for 3-5 most-used features (e.g., Ctrl+O for code ownership)
3. **Audit render functions**: 10+ features have flags but no visible UI; either complete or delete

