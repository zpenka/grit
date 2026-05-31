package grit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD93D"))
	hashStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#F39C12"))
	authorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
	whenStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	subjStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D")).Bold(true)
	addStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71"))
	removeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C"))
	hunkStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
	metaStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	msgStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	alertStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D")).Bold(true)
	divStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	focusIndStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D"))
	activePanelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D")).Bold(true)
	lineNumStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
	selectedBg       = lipgloss.Color("#3A3A4A")
)

type commitsLoadedMsg []commit
type diffLoadedMsg []diffLine
type branchesLoadedMsg struct {
	branches []string
	current  string
}
type blameLoadedMsg []blameLine
type flashClearMsg struct{}
type editorDoneMsg struct{}

func fetchCommits(repoPath, ref string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"-C", repoPath, "log",
			"--format=%H%x00%h%x00%an%x00%ar%x00%s", "-n", "200"}
		if ref != "" {
			args = append(args, ref)
		}
		out, err := exec.Command("git", args...).Output()
		if err != nil {
			return commitsLoadedMsg(nil)
		}
		return commitsLoadedMsg(parseCommits(string(out)))
	}
}

func fetchDiff(repoPath, hash string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("git", "-C", repoPath, "show",
			"--stat", "--patch", "--color=never", hash).Output()
		if err != nil {
			return diffLoadedMsg(nil)
		}
		return diffLoadedMsg(parseDiff(string(out)))
	}
}

func fetchBranches(repoPath string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("git", "-C", repoPath, "branch", "-a").Output()
		if err != nil {
			return branchesLoadedMsg{}
		}
		s := string(out)
		return branchesLoadedMsg{
			branches: parseBranches(s),
			current:  parseCurrentBranch(s),
		}
	}
}

func fetchBlame(repoPath, hash, file string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("git", "-C", repoPath, "blame",
			"--date=short", hash, "--", file).Output()
		if err != nil {
			return blameLoadedMsg(nil)
		}
		return blameLoadedMsg(parseBlame(string(out)))
	}
}

func flashCmd() tea.Cmd {
	return tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
		return flashClearMsg{}
	})
}

func copyToClipboard(s string) error {
	for _, args := range [][]string{
		{"pbcopy"},
		{"wl-copy"},
		{"xclip", "-selection", "clipboard"},
		{"xsel", "--clipboard", "--input"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = strings.NewReader(s)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no clipboard tool found")
}

func openInEditor(repoPath, hash string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	out, _ := exec.Command("git", "-C", repoPath, "show",
		"--color=always", "--no-pager", hash).Output()
	f, err := os.CreateTemp("", "grit-*.diff")
	if err != nil {
		return func() tea.Msg { return editorDoneMsg{} }
	}
	_, _ = f.Write(out)
	_ = f.Close()
	name := f.Name()
	cmd := exec.Command(editor, name)
	return tea.ExecProcess(cmd, func(error) tea.Msg {
		_ = os.Remove(name)
		return editorDoneMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return fetchCommits(m.repoPath, "")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case commitsLoadedMsg:
		m.commits = []commit(msg)
		m.cursor = 0
		vc := visibleCommits(m)
		if len(vc) > 0 {
			m.loading = true
			return m, fetchDiff(m.repoPath, vc[0].hash)
		}
		return m, nil

	case diffLoadedMsg:
		m.loading = false
		m.diffLines = []diffLine(msg)
		m.diffOffset = 0
		m.fileItems = parseFileItemsFromDiff(m.diffLines)
		m.fileCursor = 0
		m.showBlame = false
		return m, nil

	case branchesLoadedMsg:
		m.branches = msg.branches
		m.currentBranch = msg.current
		// position cursor on current branch
		for i, b := range m.branches {
			if b == msg.current {
				m.branchCursor = i
				break
			}
		}
		m.showBranch = true
		return m, nil

	case blameLoadedMsg:
		m.loading = false
		m.blameLines = []blameLine(msg)
		m.blameOffset = 0
		m.showBlame = true
		return m, nil

	case flashClearMsg:
		m.flash = ""
		return m, nil

	case editorDoneMsg:
		return m, nil
	}

	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if km.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Esc clears all feature panels (except structural ones)
	if km.String() == "esc" && !m.showBlame && !m.showBranch && !m.showFiles && !m.searching && !m.showAnalyticsMenu && !m.showVisualizationMenu && !m.showTeamMenu && !m.showIntegrationMenu && !m.showGitOpsMenu {
		m.showHelp = false
		m.showCodeOwnership = false
		m.showHotspots = false
		m.showLinting = false
		m.showComplexity = false
		m.showActivityHeatmap = false
		m.showAnalytics = false
		m.showFlamegraph = false
		m.showTimeline = false
		m.showTreeView = false
		m.showAuthorComparison = false
		m.showFileHeatmap = false
		m.showTeamStats = false
		m.showReviewUI = false
		m.showPairProgramming = false
		m.showVelocity = false
		m.showClassification = false
		m.showSecrets = false
		m.showChangelog = false
		m.showPRLinks = false
		m.showJiraLinks = false
		m.showIssueRefs = false
		return m, nil
	}

	// Blame view: all navigation goes to blame scrolling
	if m.showBlame {
		panelH := diffPanelHeight(m)
		switch km.String() {
		case "q":
			return m, tea.Quit
		case "j", "down":
			maxOff := len(m.blameLines) - panelH
			if maxOff < 0 {
				maxOff = 0
			}
			if m.blameOffset < maxOff {
				m.blameOffset++
			}
		case "k", "up":
			if m.blameOffset > 0 {
				m.blameOffset--
			}
		case "d", " ":
			maxOff := len(m.blameLines) - panelH
			if maxOff < 0 {
				maxOff = 0
			}
			m.blameOffset += panelH / 2
			if m.blameOffset > maxOff {
				m.blameOffset = maxOff
			}
		case "u":
			m.blameOffset -= panelH / 2
			if m.blameOffset < 0 {
				m.blameOffset = 0
			}
		case "g":
			m.blameOffset = 0
		case "G":
			m.blameOffset = len(m.blameLines)
		case "B", "esc":
			m.showBlame = false
			// restore diff to the blamed file's position
			for _, fi := range m.fileItems {
				if fi.path == currentFile(m) {
					m = scrollToDiffLine(m, fi.diffIdx)
					break
				}
			}
		}
		return m, nil
	}

	// Search mode
	if m.searching {
		switch km.String() {
		case "esc":
			m.searching = false
		case "enter":
			m.searching = false
		case "ctrl+/":
			m.query = ""
			m.cursor = 0
			vc := visibleCommits(m)
			if len(vc) > 0 {
				m.loading = true
				return m, fetchDiff(m.repoPath, vc[0].hash)
			}
		case "backspace":
			runes := []rune(m.query)
			if len(runes) > 0 {
				prevFirst := firstVisibleHash(m)
				m.query = string(runes[:len(runes)-1])
				m.cursor = 0
				if h := firstVisibleHash(m); h != prevFirst && h != "" {
					m.loading = true
					return m, fetchDiff(m.repoPath, h)
				}
			}
		default:
			if len(km.Runes) > 0 {
				prevFirst := firstVisibleHash(m)
				m.query += string(km.Runes)
				m.cursor = 0
				if h := firstVisibleHash(m); h != prevFirst && h != "" {
					m.loading = true
					return m, fetchDiff(m.repoPath, h)
				}
			}
		}
		return m, nil
	}

	vc := visibleCommits(m)

	// Global bindings (work in any panel mode)
	switch km.String() {
	case "q":
		return m, tea.Quit
	case "tab":
		switch {
		case m.showBranch:
			m = toggleBranchView(m)
		case m.showFiles:
			m = toggleFileView(m)
		default:
			m = switchPanel(m)
		}
		m.countBuf = ""
		return m, nil
	case "/":
		m.searching = true
		m.countBuf = ""
		return m, nil
	case "f":
		if len(m.fileItems) > 0 || m.showFiles {
			m = toggleFileView(m)
			m.showBranch = false
		}
		m.countBuf = ""
		return m, nil
	case "b":
		if !m.showBranch {
			m.showFiles = false
			m.countBuf = ""
			return m, fetchBranches(m.repoPath)
		}
		m = toggleBranchView(m)
		m.countBuf = ""
		return m, nil
	case "B":
		if file := currentFile(m); file != "" && len(vc) > 0 {
			m.loading = true
			m.countBuf = ""
			return m, fetchBlame(m.repoPath, vc[m.cursor].hash, file)
		}
		return m, nil
	case "y":
		if m.cursor < len(vc) {
			if err := copyToClipboard(vc[m.cursor].hash); err != nil {
				m.flash = "clipboard unavailable"
			} else {
				m.flash = "copied " + vc[m.cursor].shortHash
			}
			m.countBuf = ""
			return m, flashCmd()
		}
		return m, nil
	case "e":
		if m.cursor < len(vc) {
			m.countBuf = ""
			return m, openInEditor(m.repoPath, vc[m.cursor].hash)
		}
		return m, nil
	case "?":
		m.showHelp = !m.showHelp
		m.countBuf = ""
		return m, nil
	case "a":
		m.showAnalyticsMenu = !m.showAnalyticsMenu
		m.analyticsMenuIdx = 0
		m.countBuf = ""
		return m, nil
	case "v":
		m.showVisualizationMenu = !m.showVisualizationMenu
		m.vizMenuIdx = 0
		m.countBuf = ""
		return m, nil
	case "t":
		m.showTeamMenu = !m.showTeamMenu
		m.teamMenuIdx = 0
		m.countBuf = ""
		return m, nil
	case "i":
		m.showIntegrationMenu = !m.showIntegrationMenu
		m.integrationMenuIdx = 0
		m.countBuf = ""
		return m, nil
	case "r":
		m.showRebaseUI = !m.showRebaseUI
		if m.showRebaseUI && len(m.rebaseSequence) == 0 {
			m.rebaseSequence = parseRebaseSequence(m.commits)
		}
		m.countBuf = ""
		return m, nil
	case "c":
		if m.cursor < len(vc) {
			m = toggleCherryPick(m, vc[m.cursor].hash)
			m.showCherryPickUI = len(m.cherryPickList) > 0
		}
		m.countBuf = ""
		return m, nil
	case "x":
		// Reset mode cycles: "" → soft → mixed → hard → ""
		switch m.resetMode {
		case "":
			m.resetMode = "soft"
		case "soft":
			m.resetMode = "mixed"
		case "mixed":
			m.resetMode = "hard"
		default:
			m.resetMode = ""
		}
		m.countBuf = ""
		return m, nil
	case "g":
		m.showGitOpsMenu = !m.showGitOpsMenu
		m.gitOpsMenuIdx = 0
		m.countBuf = ""
		return m, nil
	case "w":
		m.showWorkflowsMenu = !m.showWorkflowsMenu
		m.workflowsMenuIdx = 0
		m.countBuf = ""
		return m, nil
	}

	// Analytics menu navigation
	if m.showAnalyticsMenu {
		switch km.String() {
		case "j", "down":
			if m.analyticsMenuIdx < analyticsMenuLen-1 {
				m.analyticsMenuIdx++
			}
		case "k", "up":
			if m.analyticsMenuIdx > 0 {
				m.analyticsMenuIdx--
			}
		case "enter", " ":
			m = dispatchAnalyticsFeature(m, m.analyticsMenuIdx)
			m.showAnalyticsMenu = false
		case "esc", "a":
			m.showAnalyticsMenu = false
		}
		return m, nil
	}

	// Visualization menu navigation
	if m.showVisualizationMenu {
		switch km.String() {
		case "j", "down":
			if m.vizMenuIdx < vizMenuLen-1 {
				m.vizMenuIdx++
			}
		case "k", "up":
			if m.vizMenuIdx > 0 {
				m.vizMenuIdx--
			}
		case "enter", " ":
			m = dispatchVizFeature(m, m.vizMenuIdx)
			m.showVisualizationMenu = false
		case "esc", "v":
			m.showVisualizationMenu = false
		}
		return m, nil
	}

	// Team menu navigation
	if m.showTeamMenu {
		switch km.String() {
		case "j", "down":
			if m.teamMenuIdx < teamMenuLen-1 {
				m.teamMenuIdx++
			}
		case "k", "up":
			if m.teamMenuIdx > 0 {
				m.teamMenuIdx--
			}
		case "enter", " ":
			m = dispatchTeamFeature(m, m.teamMenuIdx)
			m.showTeamMenu = false
		case "esc", "t":
			m.showTeamMenu = false
		}
		return m, nil
	}

	// Integration menu navigation
	if m.showIntegrationMenu {
		switch km.String() {
		case "j", "down":
			if m.integrationMenuIdx < integrationMenuLen-1 {
				m.integrationMenuIdx++
			}
		case "k", "up":
			if m.integrationMenuIdx > 0 {
				m.integrationMenuIdx--
			}
		case "enter", " ":
			m = dispatchIntegrationFeature(m, m.integrationMenuIdx)
			m.showIntegrationMenu = false
		case "esc", "i":
			m.showIntegrationMenu = false
		}
		return m, nil
	}

	// Git ops menu navigation
	if m.showGitOpsMenu {
		switch km.String() {
		case "j", "down":
			if m.gitOpsMenuIdx < gitOpsMenuLen-1 {
				m.gitOpsMenuIdx++
			}
		case "k", "up":
			if m.gitOpsMenuIdx > 0 {
				m.gitOpsMenuIdx--
			}
		case "enter", " ":
			m = dispatchGitOpsFeature(m, m.gitOpsMenuIdx)
			m.showGitOpsMenu = false
		case "esc", "g":
			m.showGitOpsMenu = false
		}
		return m, nil
	}

	// Workflows menu navigation
	if m.showWorkflowsMenu {
		switch km.String() {
		case "j", "down":
			if m.workflowsMenuIdx < workflowsMenuLen-1 {
				m.workflowsMenuIdx++
			}
		case "k", "up":
			if m.workflowsMenuIdx > 0 {
				m.workflowsMenuIdx--
			}
		case "enter", " ":
			m = dispatchWorkflowsFeature(m, m.workflowsMenuIdx)
			m.showWorkflowsMenu = false
		case "esc", "w":
			m.showWorkflowsMenu = false
		}
		return m, nil
	}

	panelH := diffPanelHeight(m)

	// Branch picker navigation
	if m.showBranch {
		switch km.String() {
		case "j", "down":
			if m.branchCursor < len(m.branches)-1 {
				m.branchCursor++
			}
		case "k", "up":
			if m.branchCursor > 0 {
				m.branchCursor--
			}
		case "enter", " ":
			if m.branchCursor < len(m.branches) {
				selected := m.branches[m.branchCursor]
				m.currentRef = selected
				m.showBranch = false
				m.cursor = 0
				m.query = ""
				m.loading = true
				return m, fetchCommits(m.repoPath, selected)
			}
		case "esc":
			m = toggleBranchView(m)
		}
		return m, nil
	}

	// File list navigation
	if m.showFiles {
		switch km.String() {
		case "j", "down":
			if m.fileCursor < len(m.fileItems)-1 {
				m.fileCursor++
			}
		case "k", "up":
			if m.fileCursor > 0 {
				m.fileCursor--
			}
		case "enter", " ":
			if m.fileCursor < len(m.fileItems) {
				m = scrollToDiffLine(m, m.fileItems[m.fileCursor].diffIdx)
				m.showFiles = false
				m.focus = panelDiff
			}
		case "esc":
			m = toggleFileView(m)
		}
		return m, nil
	}

	if m.focus == panelList {
		// Digit accumulation for count prefix
		if len(km.String()) == 1 && km.String() >= "0" && km.String() <= "9" {
			m.countBuf += km.String()
			return m, nil
		}
		n := parseCount(m.countBuf)
		m.countBuf = ""

		switch km.String() {
		case "j", "down":
			newCursor := m.cursor + n
			if newCursor >= len(vc) {
				newCursor = len(vc) - 1
			}
			if newCursor >= 0 && newCursor != m.cursor {
				m.cursor = newCursor
				m.diffOffset = 0
				m.loading = true
				return m, fetchDiff(m.repoPath, vc[m.cursor].hash)
			}
		case "k", "up":
			newCursor := m.cursor - n
			if newCursor < 0 {
				newCursor = 0
			}
			if newCursor != m.cursor {
				m.cursor = newCursor
				m.diffOffset = 0
				m.loading = true
				return m, fetchDiff(m.repoPath, vc[m.cursor].hash)
			}
		case "g":
			if m.cursor != 0 && len(vc) > 0 {
				m.cursor = 0
				m.diffOffset = 0
				m.loading = true
				return m, fetchDiff(m.repoPath, vc[0].hash)
			}
		case "G":
			if last := len(vc) - 1; last >= 0 && m.cursor != last {
				m.cursor = last
				m.diffOffset = 0
				m.loading = true
				return m, fetchDiff(m.repoPath, vc[last].hash)
			}
		case "l":
			m = switchPanel(m)
		}
	} else {
		m.countBuf = ""
		switch km.String() {
		case "j", "down":
			return scrollDiffDown(m, 1), nil
		case "k", "up":
			return scrollDiffUp(m, 1), nil
		case "d", " ":
			return scrollDiffDown(m, panelH/2), nil
		case "u":
			return scrollDiffUp(m, panelH/2), nil
		case "g":
			m.diffOffset = 0
		case "G":
			return scrollDiffDown(m, len(m.diffLines)), nil
		case "h":
			m = switchPanel(m)
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "\n  loading…\n"
	}

	// Show help overlay if requested
	if m.showHelp {
		return renderHelpOverlay(m)
	}

	// Show analytics menu overlay if requested
	if m.showAnalyticsMenu {
		return renderAnalyticsMenuOverlay(m, m.width)
	}

	// Show visualization menu overlay if requested
	if m.showVisualizationMenu {
		return renderVisualizationMenuOverlay(m, m.width)
	}

	// Show team menu overlay if requested
	if m.showTeamMenu {
		return renderTeamMenuOverlay(m, m.width)
	}

	// Show integration menu overlay if requested
	if m.showIntegrationMenu {
		return renderIntegrationMenuOverlay(m, m.width)
	}

	// Show git ops menu overlay if requested
	if m.showGitOpsMenu {
		return renderGitOpsMenuOverlay(m, m.width)
	}

	// Show workflows menu overlay if requested
	if m.showWorkflowsMenu {
		return renderWorkflowsMenuOverlay(m, m.width)
	}

	// Show analytics feature panels
	if m.showCodeOwnership {
		return renderCodeOwnershipUI(m, m.width)
	}
	if m.showHotspots {
		return renderHotspotsUI(m, m.width)
	}
	if m.showLinting {
		return renderLintingUI(m, m.width)
	}
	if m.showBisectUI {
		return renderBisectUI(m, m.width)
	}
	if m.showActivityHeatmap {
		return renderActivityHeatmapUI(m, m.width)
	}
	if m.showAnalytics {
		return renderAnalyticsPanel(m, m.width)
	}
	if m.showComplexity {
		return renderComplexityUI(m, m.width)
	}
	if m.showCoupling {
		return renderCouplingAnalysisUI(m, m.width)
	}
	if m.showDependencies {
		return renderDependenciesUI(m, m.width)
	}

	// Show visualization feature panels
	if m.showFlamegraph {
		return renderFlamegraphUI(m, m.width)
	}
	if m.showTimeline {
		return renderTimelineSliderUI(m, m.width)
	}
	if m.showTreeView {
		return renderTreeViewUI(m, m.width)
	}
	if m.showAuthorComparison {
		return renderAuthorComparisonUI(m, m.width)
	}
	if m.showFileHeatmap {
		return renderFileHeatmapUI(m, m.width)
	}

	// Show team/AI feature panels (will be added in team functions)
	if m.showTeamStats {
		return renderTeamStatsUI(m, m.width)
	}
	if m.showReviewUI {
		return renderReviewerSuggestionsUI(m, m.width)
	}
	if m.showVelocity {
		return renderVelocityUI(m, m.width)
	}
	if m.showClassification {
		return renderClassificationUI(m, m.width)
	}
	if m.showSecrets {
		return renderSecretsUI(m, m.width)
	}
	if m.showChangelog {
		return renderChangelogUI(m, m.width)
	}

	// Show integration feature panels
	if m.showPRLinks {
		return renderPRLinksUI(m, m.width)
	}
	if m.showJiraLinks {
		return renderJiraLinksUI(m, m.width)
	}
	if m.showIssueRefs {
		return renderIssueRefsUI(m, m.width)
	}
	if m.showConflictUI {
		return renderConflictUI(m, m.width)
	}
	if m.showExportUI {
		return renderExportUI(m, m.width)
	}

	// Show git ops panels
	if m.showRebaseUI {
		return renderRebaseUI(m, m.width)
	}
	if m.showCherryPickUI {
		return renderCherryPickUI(m, m.width)
	}

	listW := listPanelWidth(m.width)
	diffW := diffPanelWidth(m.width)
	panelH := diffPanelHeight(m)
	vc := visibleCommits(m)

	var sb strings.Builder

	// Title bar
	sb.WriteString(" ")
	sb.WriteString(titleStyle.Render("git log"))
	if m.currentRef != "" {
		sb.WriteString("  ")
		sb.WriteString(hashStyle.Render(m.currentRef))
	}
	sb.WriteString("  ")
	sb.WriteString(msgStyle.Render(m.repoPath))
	sb.WriteString("\n\n")

	// Left panel header
	var listHdr string
	switch {
	case m.showBranch:
		hdrText := fmt.Sprintf("Branches (%d)", len(m.branches))
		if m.focus == panelList {
			listHdr = "▌ " + activePanelStyle.Render(hdrText)
		} else {
			listHdr = "  " + msgStyle.Render(hdrText)
		}
	case m.showFiles:
		hdrText := fmt.Sprintf("Files (%d)", len(m.fileItems))
		if m.focus == panelList {
			listHdr = "▌ " + activePanelStyle.Render(hdrText)
		} else {
			listHdr = "  " + msgStyle.Render(hdrText)
		}
	case m.query != "":
		hdrText := fmt.Sprintf("Commits [/%s] %d", m.query, len(vc))
		if m.focus == panelList {
			listHdr = "▌ " + activePanelStyle.Render(hdrText)
		} else {
			listHdr = "  " + msgStyle.Render(hdrText)
		}
	default:
		if m.focus == panelList {
			listHdr = "▌ " + activePanelStyle.Render("Commits")
		} else {
			listHdr = "  " + msgStyle.Render("Commits")
		}
	}

	// Right panel header
	var diffHdr string
	switch {
	case m.showBlame:
		file := currentFile(m)
		hdrText := hashStyle.Render("BLAME") + " " + msgStyle.Render(file)
		if m.focus == panelDiff {
			diffHdr = "▌ " + activePanelStyle.Render(hdrText)
		} else {
			diffHdr = "  " + hdrText
		}
	case m.loading:
		diffHdr = "  " + msgStyle.Render("loading…")
	case len(vc) == 0:
		diffHdr = "  " + msgStyle.Render("no commits")
	default:
		c := vc[m.cursor]
		hdrText := hashStyle.Render(c.shortHash) + "  " +
			msgStyle.Render(truncate(c.subject, diffW-12))
		if m.focus == panelDiff {
			diffHdr = "▌ " + activePanelStyle.Render(hdrText)
		} else {
			diffHdr = "  " + hdrText
		}
	}

	sb.WriteString(lipgloss.NewStyle().Width(listW).Render(listHdr))
	sb.WriteString(divStyle.Render("│"))
	sb.WriteString(diffHdr)
	sb.WriteString("\n")

	// Windowed commit list
	commitStart := m.cursor - panelH/2
	if commitStart < 0 {
		commitStart = 0
	}
	if commitStart+panelH > len(vc) {
		commitStart = len(vc) - panelH
		if commitStart < 0 {
			commitStart = 0
		}
	}

	// Windowed file list
	fileStart := m.fileCursor - panelH/2
	if fileStart < 0 {
		fileStart = 0
	}
	if fileStart+panelH > len(m.fileItems) {
		fileStart = len(m.fileItems) - panelH
		if fileStart < 0 {
			fileStart = 0
		}
	}

	// Windowed branch list
	branchStart := m.branchCursor - panelH/2
	if branchStart < 0 {
		branchStart = 0
	}
	if branchStart+panelH > len(m.branches) {
		branchStart = len(m.branches) - panelH
		if branchStart < 0 {
			branchStart = 0
		}
	}

	for i := 0; i < panelH; i++ {
		var leftLine string
		switch {
		case m.showBranch:
			leftLine = renderBranchRow(m, branchStart+i, listW)
		case m.showFiles:
			leftLine = renderFileRow(m, fileStart+i, listW)
		default:
			leftLine = renderCommitRow(m, vc, commitStart+i, listW)
		}

		var rightLine string
		if m.showBlame {
			rightLine = renderBlameRow(m, m.blameOffset+i, diffW)
		} else {
			rightLine = renderDiffRow(m, m.diffOffset+i, diffW)
		}

		sb.WriteString(leftLine)
		sb.WriteString(divStyle.Render("│"))
		sb.WriteString(rightLine)
		sb.WriteString("\n")
	}

	// Footer
	sb.WriteString("\n")
	switch {
	case m.flash != "":
		sb.WriteString(alertStyle.Render(m.flash))
	case m.searching:
		sb.WriteString(msgStyle.Render("[/] " + m.query + "█   Esc exit   Ctrl+/ clear"))
	case m.countBuf != "":
		sb.WriteString(msgStyle.Render(fmt.Sprintf("[%s] j/k  jump   other key  cancel", m.countBuf)))
	case m.showBlame:
		sb.WriteString(msgStyle.Render("j/k  scroll   d/u  half-page   g/G  top/bottom   B/Esc  back to diff"))
	case m.showBranch:
		sb.WriteString(msgStyle.Render("j/k  navigate   Enter  switch branch   b/Esc  back"))
	case m.showFiles:
		sb.WriteString(msgStyle.Render("j/k  navigate   Enter  jump to file   f/Esc  back"))
	case m.focus == panelDiff:
		sb.WriteString(msgStyle.Render("j/k  scroll   d/u  half-page   g/G  top/bottom   h/Tab  commits   ?  help   q  quit"))
	default:
		sb.WriteString(msgStyle.Render("j/k  navigate   5j  jump 5   l/Tab  diff   /  search   f  files   b  branches   ?  help   q  quit"))
	}
	sb.WriteString("\n\n")

	return sb.String()
}

func renderCommitRow(m model, vc []commit, idx int, w int) string {
	const (
		cursorW = 2
		hashW   = 7
		whenW   = 8
		authW   = 9
	)
	subjectW := w - cursorW - hashW - 1 - authW - 1 - whenW - 1
	if subjectW < 4 {
		subjectW = 4
	}
	if idx < 0 || idx >= len(vc) {
		return lipgloss.NewStyle().Width(w).Render("")
	}
	c := vc[idx]
	selected := idx == m.cursor

	bg := func(st lipgloss.Style) lipgloss.Style {
		if selected && m.focus == panelList && !m.showBranch && !m.showFiles {
			return st.Background(selectedBg)
		}
		return st
	}
	var cur string
	if idx == m.cursor {
		cur = bg(cursorStyle).Width(cursorW).Render("▶")
	} else {
		cur = bg(lipgloss.NewStyle()).Width(cursorW).Render("")
	}
	hash := bg(hashStyle).Width(hashW + 1).Render(c.shortHash)
	subj := bg(subjStyle).Width(subjectW + 1).Render(truncate(c.subject, subjectW))
	auth := bg(authorStyle).Width(authW + 1).Render(truncate(firstWord(c.author), authW))
	when := bg(whenStyle).Width(whenW).Render(truncate(c.when, whenW))
	return cur + hash + subj + auth + when
}

func renderFileRow(m model, idx int, w int) string {
	const cursorW = 2
	if idx < 0 || idx >= len(m.fileItems) {
		return lipgloss.NewStyle().Width(w).Render("")
	}
	fi := m.fileItems[idx]
	selected := idx == m.fileCursor
	bg := func(st lipgloss.Style) lipgloss.Style {
		if selected {
			return st.Background(selectedBg)
		}
		return st
	}
	var cur string
	if selected {
		cur = bg(cursorStyle).Width(cursorW).Render("▶")
	} else {
		cur = bg(lipgloss.NewStyle()).Width(cursorW).Render("")
	}
	path := bg(subjStyle).Width(w - cursorW).Render(truncate(fi.path, w-cursorW))
	return cur + path
}

func renderBranchRow(m model, idx int, w int) string {
	const cursorW = 2
	if idx < 0 || idx >= len(m.branches) {
		return lipgloss.NewStyle().Width(w).Render("")
	}
	name := m.branches[idx]
	selected := idx == m.branchCursor
	isCurrent := name == m.currentBranch

	bg := func(st lipgloss.Style) lipgloss.Style {
		if selected {
			return st.Background(selectedBg)
		}
		return st
	}
	var cur string
	if selected {
		cur = bg(cursorStyle).Width(cursorW).Render("▶")
	} else {
		cur = bg(lipgloss.NewStyle()).Width(cursorW).Render("")
	}

	nameW := w - cursorW - 2
	if nameW < 4 {
		nameW = 4
	}
	label := truncate(name, nameW)
	var marker string
	if isCurrent {
		marker = bg(alertStyle).Width(2).Render("●")
	} else {
		marker = bg(lipgloss.NewStyle()).Width(2).Render("")
	}
	text := bg(subjStyle).Width(nameW).Render(label)
	return cur + text + marker
}

func renderBlameRow(m model, idx int, w int) string {
	if idx < 0 || idx >= len(m.blameLines) {
		return ""
	}
	bl := m.blameLines[idx]
	const hashW = 7
	const authW = 12
	const dateW = 10
	const numW = 5
	contentW := w - hashW - 1 - authW - 1 - dateW - 1 - numW - 1
	if contentW < 4 {
		contentW = 4
	}
	hash := hashStyle.Width(hashW + 1).Render(bl.shortHash)
	auth := authorStyle.Width(authW + 1).Render(truncate(bl.author, authW))
	date := whenStyle.Width(dateW + 1).Render(bl.date)
	num := lineNumStyle.Width(numW + 1).Render(fmt.Sprintf("%d", bl.lineNum))
	text := truncate(bl.text, contentW)
	return hash + auth + date + num + text
}

func renderDiffRow(m model, idx int, w int) string {
	if idx < 0 || idx >= len(m.diffLines) {
		return ""
	}
	dl := m.diffLines[idx]
	text := truncate(dl.text, w)
	switch dl.kind {
	case lineAdded:
		return addStyle.Render(text)
	case lineRemoved:
		return removeStyle.Render(text)
	case lineHunk:
		return hunkStyle.Render(text)
	case lineMeta:
		return metaStyle.Render(text)
	default:
		return text
	}
}

func firstVisibleHash(m model) string {
	vc := visibleCommits(m)
	if len(vc) == 0 {
		return ""
	}
	return vc[0].hash
}

func renderHelpOverlay(m model) string {
	var sb strings.Builder

	w := m.width
	if w < 40 {
		w = 40
	}
	h := m.height
	if h < 10 {
		h = 10
	}

	// Build help text with dynamic menu items
	var help strings.Builder
	help.WriteString("KEYBINDING REFERENCE\n\n")
	help.WriteString("NAVIGATION\n")
	help.WriteString("  j/k          Scroll commit list / diff panel\n")
	help.WriteString("  g/G          Jump to first / last commit (in list) / git ops menu (g key)\n")
	help.WriteString("  5j           Jump 5 commits (count prefix works with j/k)\n")
	help.WriteString("  l/h          Switch between commit and diff panels\n")
	help.WriteString("  Tab          Toggle between commit and diff\n\n")

	help.WriteString("VIEWING\n")
	help.WriteString("  f            Toggle file list\n")
	help.WriteString("  b            Show branch picker\n")
	help.WriteString("  B            Show blame for current file\n")
	help.WriteString("  /            Search commits by subject, author, or hash\n")
	help.WriteString("  Ctrl+/       Clear search\n\n")

	help.WriteString("CLIPBOARD & EDITOR\n")
	help.WriteString("  y            Copy current commit hash\n")
	help.WriteString("  e            Open current commit in $EDITOR\n\n")

	// Dynamic SUBMENUS section built from *MenuItems slices
	help.WriteString("SUBMENUS (j/k to navigate, Enter to select, key again or Esc to close)\n")

	// Analytics menu
	help.WriteString("  a  Analytics\n")
	for i := 0; i < len(analyticsMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(analyticsMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + analyticsMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	// Visualization menu
	help.WriteString("  v  Visualizations\n")
	for i := 0; i < len(vizMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(vizMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + vizMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	// Team/AI menu
	help.WriteString("  t  Team / AI\n")
	for i := 0; i < len(teamMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(teamMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + teamMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	// Integration menu
	help.WriteString("  i  Integration\n")
	for i := 0; i < len(integrationMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(integrationMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + integrationMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	// Git Ops menu
	help.WriteString("  g  Git Ops\n")
	for i := 0; i < len(gitOpsMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(gitOpsMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + gitOpsMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	// Workflows menu
	help.WriteString("  w  Workflows\n")
	for i := 0; i < len(workflowsMenuItems); i += 3 {
		help.WriteString("     ")
		for j := 0; j < 3 && i+j < len(workflowsMenuItems); j++ {
			if j > 0 {
				help.WriteString("  ")
			}
			help.WriteString("- " + workflowsMenuItems[i+j])
		}
		help.WriteString("\n")
	}
	help.WriteString("\n")

	help.WriteString("HELP & EXIT\n")
	help.WriteString("  ?            This help overlay\n")
	help.WriteString("  Esc          Close any panel / clear search\n")
	help.WriteString("  q            Quit application\n")

	// Split into lines and render
	lines := strings.Split(strings.TrimSpace(help.String()), "\n")
	contentHeight := h - 4

	// Ensure HELP & EXIT section is always visible by reserving space at the bottom
	// "HELP & EXIT\n  ?\n  Esc\n  q\n" is at least 4 lines
	helpSectionStartIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "HELP & EXIT" {
			helpSectionStartIdx = i
			break
		}
	}

	if helpSectionStartIdx > 0 && len(lines) > contentHeight {
		// If HELP & EXIT section would be cut off, trim from the menus section instead
		helpSectionLines := len(lines) - helpSectionStartIdx
		if helpSectionLines <= 4 && contentHeight > helpSectionLines {
			lines = append(lines[:contentHeight-helpSectionLines], lines[helpSectionStartIdx:]...)
		} else {
			lines = lines[:contentHeight]
		}
	} else if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	for _, line := range lines {
		sb.WriteString(msgStyle.Render(line))
		sb.WriteString("\n")
	}

	// Add footer
	footerText := " Press ? to close  │  q to quit "
	padding := w - len(footerText)
	if padding < 0 {
		padding = 0
	}
	footer := strings.Repeat(" ", padding/2) + footerText
	sb.WriteString(alertStyle.Render(footer))

	return sb.String()
}

const Version = "1.0.1"

func Run() {
	// Handle CLI flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("grit v%s\n", Version)
			return
		case "-h", "--help":
			fmt.Println("grit - Terminal UI for exploring git history")
			fmt.Println("Usage: grit [options]")
			fmt.Println("Options:")
			fmt.Println("  -v, --version  Show version")
			fmt.Println("  -h, --help     Show this help message")
			return
		}
	}

	out, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	repoPath := strings.TrimSpace(string(out))
	if repoPath == "" {
		repoPath = "."
	}
	p := tea.NewProgram(newModel(repoPath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
