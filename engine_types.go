package grit

import (
	"regexp"
)

// Core data types

type commit struct {
	hash      string
	shortHash string
	author    string
	when      string
	subject   string
}

type panel int

const (
	panelList panel = iota
	panelDiff
)

type lineKind int

const (
	lineContext lineKind = iota
	lineAdded
	lineRemoved
	lineHunk
	lineMeta
)

type diffLine struct {
	kind lineKind
	text string
}

type fileItem struct {
	path    string
	diffIdx int
}

type blameLine struct {
	shortHash string
	author    string
	date      string
	lineNum   int
	text      string
}

// FilterCache caches filter results with metrics tracking.
type FilterCache struct {
	cache   map[string][]commit
	metrics CacheMetrics
}

// NewFilterCache creates a new filter cache with metrics.
func NewFilterCache() *FilterCache {
	return &FilterCache{
		cache: make(map[string][]commit),
		metrics: CacheMetrics{
			Size:    0,
			MaxSize: 100,
			Hits:    0,
			Misses:  0,
		},
	}
}

// model is the central state holder for the UI
type model struct {
	commits    []commit
	cursor     int
	focus      panel
	diffLines  []diffLine
	diffOffset int
	fileItems  []fileItem
	fileCursor int
	showFiles  bool
	searching  bool
	query      string
	flash      string
	// Branch picker
	showBranch    bool
	branches      []string
	branchCursor  int
	currentRef    string
	currentBranch string
	// Blame view
	showBlame   bool
	blameLines  []blameLine
	blameOffset int
	// Count prefix for j/k navigation
	countBuf string
	// Filtering
	authorFilter string
	sinceFilter  int // days; 0 = no filter
	// Breadcrumb trail
	navHistory    []int
	navHistoryIdx int
	// Bookmarks
	bookmarks []string // commit short hashes
	// Stats
	lastStats commitStatistics
	// Line comments
	comments map[int]string
	// Tag view
	showTags  bool
	tags      []string
	tagCursor int
	// Option 1: UI Integration
	showStatsBadge bool
	// Option 2: Commit Graph
	commitGraph []graphNode
	showGraph   bool
	// Option 3: File-Centric
	fileHistory      []commit
	currentFile      string
	showFileTimeline bool
	// Option 4: Stash & Reflog
	viewMode      string // "log", "stash", "reflog"
	stashes       []stashEntry
	reflogEntries []reflogEntry
	stashCursor   int
	reflogCursor  int
	// UI Integration
	inGoToCommitMode bool
	goToCommitInput  string
	inCommentMode    bool
	commentInput     string
	// Optimization: Caches
	dcache    *diffCache
	scache    *statCache
	recache   *regexCache
	// Performance tracking
	diffCacheHits   int
	statCacheHits   int
	regexCacheHits  int
	// Advanced Operations
	showRebaseUI      bool
	rebaseSequence    []rebaseOp
	showCherryPickUI  bool
	cherryPickList    []string
	resetMode         string // soft, mixed, hard
	amendMessage      string
	// Analytics
	showAnalytics     bool
	authorStats       map[string]int
	timeStats         map[string]int
	collaborators     map[string][]string // author -> co-authors
	reviewers         map[string][]string // commit hash -> reviewers
	productivity      map[string]interface{}
	repoPath          string
	width             int
	height            int
	loading           bool
	// Bisect & Recovery
	bisectState         bisectState
	showBisectUI        bool
	lostCommits         []lostCommit
	showLostCommits     bool
	undoStack           []string // commit hashes for undo
	undoStackIdx        int
	showUndoMenu        bool
	reflogRecoveryMode  bool
	recoveryCommits     []lostCommit
	// Code Patterns & Quality
	codeOwnership       map[string]codeOwnershipData
	showCodeOwnership   bool
	hotspots            []hotspotData
	showHotspots        bool
	commitMetrics       []commitMetrics
	showComplexity      bool
	lintingResults      []lintingResult
	showLinting         bool
	largeCommits        []commitMetrics
	showLargeCommits    bool
	// Commit Analysis & Search
	semanticSearchResults []semanticSearchResult
	showSemanticSearch    bool
	semanticQuery         string
	authorActivityHeatmap map[string]authorActivityData
	showActivityHeatmap   bool
	mergeAnalysisData     []mergeAnalysis
	showMergeAnalysis     bool
	commitCouplings       []commitCoupling
	showCoupling          bool
	// Performance & Filtering
	extensionFilters  []fileExtensionFilter
	currentExtFilter  string
	commitGroups      []commitGroup
	groupingMode      string // "", "pr", "branch", "date"
	dependencyChanges []dependencyChange
	showDependencies  bool
	// Advanced Workflows
	worktrees          []worktreeInfo
	showWorktrees      bool
	currentWorktree    string
	submodules         []submoduleInfo
	showSubmodules     bool
	namedStashes       []namedStash
	showNamedStashes   bool
	pendingTagOps      []tagOperation
	showTagMgmt        bool
	gpgStatuses        map[string]gpgSignatureStatus
	showGPGStatus      bool
	// Visualization
	contributorFlameData []contributorFlameData
	showFlamegraph       bool
	timelinePoints       []timelinePoint
	timelineSliderPos    int
	showTimeline         bool
	treeRoot             *treeNode
	showTreeView         bool
	authorComparisons    []authorComparison
	selectedAuthors      [2]string
	showAuthorComparison bool
	fileHeatmap          []fileHeatmapEntry
	showFileHeatmap      bool
	// Integration & Export
	prReferences      []githubPRReference
	showPRLinks       bool
	jiraLinks         []jiraTicketLink
	showJiraLinks     bool
	pendingExports    []exportData
	showExportUI      bool
	exportFormat      string
	issueReferences   []issueReference
	showIssueRefs     bool
	// Advanced Git Operations
	rebasePreview           rebasePreview
	showRebasePreview       bool
	conflictList            []conflictInfo
	showConflictUI          bool
	squashPlans             []squashPlan
	showSquashUI            bool
	cherryPickImprovements  []cherryPickImprovement
	amendPreview            amendPreview
	showAmendPreview        bool
	// Team & Collaboration
	teamStats               []teamStats
	showTeamStats           bool
	reviewWorkflows         []reviewWorkflow
	showReviewUI            bool
	reviewerSuggestions     []reviewerSuggestion
	pairProgrammingData     []pairProgrammingData
	showPairProgramming     bool
	velocityHistory         []velocityData
	showVelocity            bool
	// AI-Powered Insights
	messageCompletions      []messageCompletion
	commitClassifications   []commitClassification
	showClassification      bool
	anomalies               []anomalyData
	showAnomalies           bool
	similarCommits          []similarCommit
	showSimilar             bool
	autoSummaries           []autoSummary
	showSummaries           bool
	// Compliance & Security
	signingStatuses         map[string]signingStatus
	showSigningStatus       bool
	licenseHeaders          []licenseHeader
	showLicenses            bool
	securityIssues          []securityIssue
	showSecurityIssues      bool
	dataDeleteRequests      []dataDeleteRequest
	showDataRequests        bool
	secretDetections        []secretDetection
	showSecrets             bool
	// Release & Versioning
	semverVersions          []semverData
	showSemver              bool
	changelog               []changelogEntry
	showChangelog           bool
	releaseNotes            []releaseNote
	showReleaseNotes        bool
	versionBumps            []versionBump
	showVersionBumps        bool
	milestones              []milestone
	showMilestones          bool
	// Advanced Performance
	loadState               repoLoadState
	diffJobs                []diffProcessingJob
	indexData               indexData
	showLoadProgress        bool
	blameCache              map[string][]blameEntry
	showBlamePerf           bool
	memoryMetrics           memoryMetrics
	showMemoryMetrics       bool
}

// Supporting types

type commitStatistics struct {
	filesChanged int
	insertions   int
	deletions    int
}

type graphNode struct {
	hash    string
	depth   int
	isMerge bool
	parents []string
}

type stashEntry struct {
	name   string
	branch string
	subject string
	hash   string
}

type reflogEntry struct {
	hash    string
	action  string
	message string
	date    string
}

type diffCache struct {
	data     map[string][]diffLine
	order    []string
	maxSize  int
	hitCount int
}

type statCache struct {
	data     map[string]commitStatistics
	order    []string
	maxSize  int
	hitCount int
}

type regexCache struct {
	data     map[string]*regexp.Regexp
	maxSize  int
	hitCount int
}

type rebaseOp struct {
	action  string // pick, squash, fixup, reword, drop
	hash    string
	subject string
}

type bisectOp struct {
	hash    string
	isBad   bool
	isGood  bool
	current bool
}

type bisectState struct {
	active       bool
	current      string
	good         []string
	bad          []string
	candidates   []string
	visualSteps  int
	totalSteps   int
}

type lostCommit struct {
	hash      string
	shortHash string
	author    string
	subject   string
	date      string
}

type codeOwnershipData struct {
	author        string
	files         map[string]int // file -> count of commits
	lines         int
	expertise     float64
	isOwner       bool
}

type hotspotData struct {
	path             string
	changeFrequency  int
	recentChanges    int
	collaborators    int
	avgCommitSize    int
	riskLevel        string // low, medium, high
}

type commitMetrics struct {
	hash          string
	linesChanged  int
	filesChanged  int
	complexity    int // estimated
	isLarge       bool
	isComplex     bool
	messageQuality int // 0-100
}

type lintingResult struct {
	hash    string
	subject string
	issues  []string
	score   int // 0-100
}

// Commit Analysis & Search
type semanticSearchResult struct {
	hash      string
	subject   string
	matches   []string // matched items (function names, variables)
	relevance int      // 0-100
}

type authorActivityData struct {
	author      string
	hourOfDay   map[int]int // hour -> count
	dayOfWeek   map[int]int // day -> count
	peakHour    int
	peakDay     string
	avgPerDay   float64
}

type mergeAnalysis struct {
	hash          string
	isMerge       bool
	isFastForward bool
	parentCount   int
	conflictRisk  int // 0-100
}

type commitCoupling struct {
	file1       string
	file2       string
	coChangeCount int
	correlation float64 // 0-1
}

// Performance & Filtering
type fileExtensionFilter struct {
	extension string
	enabled   bool
}

type commitGroup struct {
	name     string
	commits  []string // hashes
	label    string   // PR, branch, or time period
	groupBy  string   // "pr", "branch", "date"
}

type dependencyChange struct {
	hash    string
	dep     string
	oldVer  string
	newVer  string
	reason  string
}

// Advanced Workflows
type worktreeInfo struct {
	path   string
	branch string
	hash   string
}

type submoduleInfo struct {
	path   string
	url    string
	hash   string
	branch string
}

type namedStash struct {
	index       int
	name        string
	description string
	hash        string
}

type tagOperation struct {
	name    string
	hash    string
	action  string // create, delete, push
	message string
}

type gpgSignatureStatus struct {
	hash      string
	signed    bool
	signer    string
	verified  bool
	algorithm string
}

// Visualization
type contributorFlameData struct {
	author     string
	commits    int
	lines      int
	percentage float64
	timeline   map[string]int // date -> commit count
}

type timelinePoint struct {
	date    string
	commits int
	hash    string
}

type treeNode struct {
	hash     string
	subject  string
	children []*treeNode
	depth    int
}

type authorComparison struct {
	author1      string
	author2      string
	commits1     int
	commits2     int
	files1       int
	files2       int
	additions1   int
	additions2   int
	deletions1   int
	deletions2   int
	similarity   float64
}

type fileHeatmapEntry struct {
	path      string
	frequency int
	recent    int
	risk      string // low, medium, high
}

// Integration & Export
type githubPRReference struct {
	hash    string
	prNumber int
	status  string // open, merged, closed
	title   string
}

type jiraTicketLink struct {
	hash   string
	ticket string
	status string
}

type exportData struct {
	format   string // "markdown", "patch", "json"
	commits  []commit
	content  string
	filename string
}

type issueReference struct {
	hash      string
	references []string // "#123", "#456"
	keywords  []string  // "fixes", "closes", "resolves"
}

// Advanced Git Operations
type rebasePreview struct {
	operations []rebaseOp
	conflicts  []string
	willApply  bool
	message    string
}

type conflictInfo struct {
	file     string
	hash     string
	markers  []string
	resolved bool
}

type squashPlan struct {
	targetHash  string
	toSquash    []string
	resultMsg   string
	lineCount   int
}

type cherryPickImprovement struct {
	hash       string
	autoConflict bool
	suggestions []string
}

type amendPreview struct {
	originalMsg string
	newMsg      string
	changes     map[string]int // file -> change count
}

// Team & Collaboration
type teamStats struct {
	author           string
	commits          int
	additions        int
	deletions        int
	avgCommitSize    int
	specialization   string
	collaborators    []string
}

type reviewWorkflow struct {
	prNumber      int
	author        string
	reviewers     []string
	approved      bool
	commentCount  int
	status        string
}

type reviewerSuggestion struct {
	reviewer   string
	expertise  float64
	availability float64
	score      float64
}

type pairProgrammingData struct {
	pair1      string
	pair2      string
	commits    int
	files      int
	coChangeRate float64
}

type velocityData struct {
	week    string
	commits int
	files   int
	additions int
	deletions int
}

// AI-Powered Insights
type messageCompletion struct {
	prefix      string
	suggestions []string
	confidence  []float64
}

type commitClassification struct {
	hash       string
	category   string // "feature", "fix", "refactor", "docs", "test"
	confidence float64
	reason     string
}

type anomalyData struct {
	hash      string
	type_     string // "large", "unusual-pattern", "unusual-time"
	severity  int    // 1-10
	description string
}

type similarCommit struct {
	hash1   string
	hash2   string
	subject1 string
	subject2 string
	similarity float64
}

type autoSummary struct {
	hash    string
	summary string
	length  int
	tokens  int
}

// Compliance & Security
type signingStatus struct {
	hash      string
	isSigned  bool
	enforced  bool
	compliant bool
}

type licenseHeader struct {
	file      string
	hasHeader bool
	license   string
	hash      string
}

type securityIssue struct {
	hash     string
	severity string // "low", "medium", "high", "critical"
	type_    string // "hardcoded-secret", "sql-injection", etc.
	location string
}

type dataDeleteRequest struct {
	hash    string
	date    string
	reason  string
	status  string // "pending", "executed"
	email   string
}

type secretDetection struct {
	hash      string
	type_     string // "api-key", "password", "token"
	location  string
	severity  string
}

// Release & Versioning
type semverData struct {
	hash       string
	version    string
	versionType string // "major", "minor", "patch"
	isRelease  bool
}

type changelogEntry struct {
	version   string
	date      string
	commits   []string
	features  []string
	bugfixes  []string
	breaking  []string
}

type releaseNote struct {
	version     string
	summary     string
	highlights  []string
	contributors []string
	date        string
}

type versionBump struct {
	hash    string
	from    string
	to      string
	date    string
	message string
}

type milestone struct {
	name    string
	version string
	commits []string
	date    string
	status  string
}

// Advanced Performance
type repoLoadState struct {
	totalCommits   int
	loadedCommits  int
	percentage     int
	isComplete     bool
	estimatedTime  int // seconds
}

type diffProcessingJob struct {
	hash       string
	status     string // "pending", "processing", "done"
	result     []diffLine
	error      string
}

type indexData struct {
	lastIndexed string
	entries     int
	isUpToDate  bool
	nextUpdate  string
}

type blameEntry struct {
	hash   string
	author string
	date   string
	line   int
	text   string
}

type memoryMetrics struct {
	usageBytes    int64
	cacheSize     int
	percentUsed   float64
	estimatedMax  int64
}
