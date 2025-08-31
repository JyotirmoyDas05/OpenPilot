// Package logo renders an OpenPilot wordmark in a stylized way.
package logo

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/JyotirmoyDas05/openpilot/internal/tui/styles"
	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/slice"
)

// letterform represents a letterform. It can be stretched horizontally by
// a given amount via the boolean argument.
type letterform func(bool) string

const diag = `╱`

// Opts are the options for rendering the OpenPilot title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	BrandColor   color.Color // Surya™ text color
	VersionColor color.Color // Version text color
	Width        int         // width of the rendered logo, used for truncation
}

// Render renders the OpenPilot logo. Set the argument to true to render the narrow
// version, intended for use in a sidebar.
//
// The compact argument determines whether it renders compact for the sidebar
// or wider for the main pane.
func Render(version string, compact bool, o Opts) string {
	const charm = " Surya™"

	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title.
	const spacing = 1
	letterforms := []letterform{letterO, letterP, letterE, letterN, letterP, letterI, letterL, letterO, letterT}
	// Preliminary word without stretch to measure base width.
	baseWord := renderWord(spacing, -1, letterforms...)
	baseWidth := lipgloss.Width(baseWord)

	// Decide layout.
	dec := decideLayout(o.Width, baseWidth, 6, 15, len(letterforms))
	if dec.Compact {
		compact = true
	}

	// Re-render with stretch if applicable.
	word := renderWord(spacing, dec.StretchIndex, letterforms...)
	brandWidth := lipgloss.Width(word)

	// Apply gradient per line if enabled.
	if dec.ApplyGradient {
		b := new(strings.Builder)
		for r := range strings.SplitSeq(word, "\n") {
			fmt.Fprintln(b, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
		}
		word = b.String()
	}

	// Surya and version.
	metaRowGap := 1
	maxVersionWidth := brandWidth - lipgloss.Width(charm) - metaRowGap
	version = ansi.Truncate(version, maxVersionWidth, "…")
	gap := max(0, brandWidth-lipgloss.Width(charm)-lipgloss.Width(version))
	metaRow := fg(o.BrandColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)
	word = metaRow + "\n" + word // do not TrimSpace to preserve leading spaces

	// Narrow version.
	diagRune := diag
	if dec.UseSimpleDiag {
		diagRune = "/"
	}
	if compact {
		// Sidebar mode: prefer stacked OPEN / PILOT if width allows clean render.
		openForms := []letterform{letterO, letterP, letterE, letterN}
		pilotForms := []letterform{letterP, letterI, letterL, letterO, letterT}
		openWord := renderWord(1, -1, openForms...)
		pilotWord := renderWord(1, -1, pilotForms...)
		openW := lipgloss.Width(openWord)
		pilotW := lipgloss.Width(pilotWord)
		stackW := max(openW, pilotW)
		stackOK := true
		// Require sufficient width to display stacked variant without truncation.
		if o.Width > 0 && stackW > o.Width-2 {
			stackOK = false
		}
		if stackOK {
			if dec.ApplyGradient {
				bOpen := new(strings.Builder)
				for r := range strings.SplitSeq(openWord, "\n") {
					fmt.Fprintln(bOpen, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
				}
				openWord = bOpen.String()
				bPilot := new(strings.Builder)
				for r := range strings.SplitSeq(pilotWord, "\n") {
					fmt.Fprintln(bPilot, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
				}
				pilotWord = bPilot.String()
			} else {
				plain := lipgloss.NewStyle().Foreground(o.TitleColorA)
				openWord = plain.Render(openWord)
				pilotWord = plain.Render(pilotWord)
			}
			// Meta row sized to stack width.
			maxVersionWidth = stackW - lipgloss.Width(charm) - metaRowGap
			if maxVersionWidth < 0 {
				maxVersionWidth = 0
			}
			version = ansi.Truncate(version, maxVersionWidth, "…")
			gap = max(0, stackW-lipgloss.Width(charm)-lipgloss.Width(version))
			metaRow = fg(o.BrandColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)
			content := metaRow + "\n" + openWord + "\n" + pilotWord
			// Pad lines to stackW to ensure stable width
			content = padLines(content, stackW)
			// Frame with two diag lines (top & bottom) sized to stackW.
			diagRune := diag
			if dec.UseSimpleDiag {
				diagRune = "/"
			}
			frame := fg(o.FieldColor, strings.Repeat(diagRune, stackW))
			block := strings.Join([]string{frame, content, frame, ""}, "\n")
			if o.Width > 0 {
				block = clampLines(block, o.Width)
			}
			return block
		}
		// Fallback to previous OP abbreviation logic.
		abbrLetters := []letterform{letterO, letterP}
		abbr := renderWord(1, -1, abbrLetters...)
		abbrWidth := lipgloss.Width(abbr)
		if dec.ApplyGradient {
			b2 := new(strings.Builder)
			for r := range strings.SplitSeq(abbr, "\n") {
				fmt.Fprintln(b2, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
			}
			abbr = b2.String()
		} else {
			abbr = lipgloss.NewStyle().Foreground(o.TitleColorA).Render(abbr)
		}
		if o.Width > 0 && abbrWidth > o.Width-2 {
			abbr = lipgloss.NewStyle().Foreground(o.TitleColorA).Render("OP")
			abbrWidth = lipgloss.Width(abbr)
		}
		maxVersionWidth = abbrWidth - lipgloss.Width(charm) - metaRowGap
		if maxVersionWidth < 0 {
			maxVersionWidth = 0
		}
		version = ansi.Truncate(version, maxVersionWidth, "…")
		gap = max(0, abbrWidth-lipgloss.Width(charm)-lipgloss.Width(version))
		metaRow = fg(o.BrandColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)
		mini := metaRow + "\n" + abbr
		mini = padLines(mini, abbrWidth)
		fieldWidth := abbrWidth
		if fieldWidth < 0 {
			fieldWidth = 0
		}
		diagRune := diag
		if dec.UseSimpleDiag {
			diagRune = "/"
		}
		field := fg(o.FieldColor, strings.Repeat(diagRune, fieldWidth))
		block := strings.Join([]string{field, mini, field, ""}, "\n")
		if o.Width > 0 {
			block = clampLines(block, o.Width)
		}
		return block
	}

	fieldHeight := lipgloss.Height(word)
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diagRune, dec.LeftFieldWidth))
	leftField := new(strings.Builder)
	for i := 0; i < fieldHeight; i++ {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightField := new(strings.Builder)
	for i := 0; i < fieldHeight; i++ {
		width := dec.RightFieldWidth
		if width < 0 {
			width = 0
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diagRune, width)), "\n")
	}

	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, word, hGap, rightField.String())
	// Ensure internal word section has consistent width before clamp
	if o.Width > 0 {
		logo = clampLines(logo, o.Width)
	}
	return logo
}

// regulateForAvailableWidth inspects the available width (o.Width) and the
// word width and forces the compact layout if the wide layout
// would not fit. This prevents generating oversized decorative fields that
// get truncated and corrupt ANSI sequences when the terminal is narrow
// (for example when a new chat opens in a narrow sidebar).
func regulateForAvailableWidth(compact bool, o Opts, brandWidth int) (bool, Opts) {
	if o.Width <= 0 {
		// no constraint provided; nothing to regulate
		return compact, o
	}

	const leftWidthLocal = 6
	const minRightWidth = 1 // allow at least one column for the right field
	const gaps = 2          // space between left/word and word/right

	// Minimum total width needed for the wide layout.
	minTotal := brandWidth + leftWidthLocal + gaps + minRightWidth

	// For chat mode with stacked layout, be more lenient about width requirements
	// since the stacked layout is more important for branding than decorative fields
	if o.Width < minTotal {
		// If we're very close to the requirement (within 6 chars), allow it but
		// reduce the right field to make it fit
		if o.Width >= brandWidth+leftWidthLocal+gaps {
			return compact, o // allow stacked layout with minimal right field
		}
		// Not enough room for even the basic layout — use compact rendering
		return true, o
	}

	// Enough room: leave as-is.
	return compact, o
}

// SmallRender renders a smaller version of the OpenPilot logo, suitable for
// smaller windows or sidebar usage.
func SmallRender(width int) string {
	t := styles.CurrentTheme()
	title := t.S().Base.Foreground(t.Secondary).Render("Surya™")
	title = fmt.Sprintf("%s %s", title, styles.ApplyBoldForegroundGrad("OpenPilot", t.Secondary, t.Primary))
	remainingWidth := width - lipgloss.Width(title) - 1 // 1 for the space after "OpenPilot"
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, t.S().Base.Foreground(t.Primary).Render(lines))
	}
	return title
}

// renderWord renders letterforms to fork a word. stretchIndex is the index of
// the letter to stretch, or -1 if no letter should be stretched.
func renderWord(spacing int, stretchIndex int, letterforms ...letterform) string {
	if spacing < 0 {
		spacing = 0
	}

	renderedLetterforms := make([]string, len(letterforms))

	// pick one letter randomly to stretch
	for i, letter := range letterforms {
		renderedLetterforms[i] = letter(i == stretchIndex)
	}

	if spacing > 0 {
		// Add spaces between the letters and render.
		renderedLetterforms = slice.Intersperse(renderedLetterforms, strings.Repeat(" ", spacing))
	}
	return strings.TrimSpace(
		lipgloss.JoinHorizontal(lipgloss.Top, renderedLetterforms...),
	)
}

// letterO renders an O-like rounded block.
func letterO(stretch bool) string {
	left := heredoc.Doc(`
		▄▀▀▀▄
		█   █
		▀▄▄▄▀
	`)
	return joinLetterform(left)
}

// letterP renders a P-like block.
func letterP(stretch bool) string {
	top := heredoc.Doc(`
		█▀▀▀█
		█▄▄▄█
		█
	`)
	return joinLetterform(top)
}

// letterE renders an E-like block.
func letterE(stretch bool) string {
	parts := heredoc.Doc(`
		█▀▀▀
		█▀▀
		█▄▄▄
	`)
	return joinLetterform(parts)
}

// letterN renders an N-like block.
func letterN(stretch bool) string {
	left := heredoc.Doc(`
		█▌ █
		██▌█
		█ ▐█
	`)
	return joinLetterform(left)
}

// letterI renders a thin I block.
func letterI(stretch bool) string {
	i := heredoc.Doc(`
		▀█▀
		 █
		▄█▄
		
	`)
	return joinLetterform(i)
}

// letterL renders an L-like block.
func letterL(stretch bool) string {
	l := heredoc.Doc(`
		█
		█
		█▄▄▄
	`)
	return joinLetterform(l)
}

// letterT renders a T-like block.
func letterT(stretch bool) string {
	t := heredoc.Doc(`
		▀▀█▀▀
		  █
		  █
	`)
	return joinLetterform(t)
}

func joinLetterform(letters ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, letters...)
}

// RenderChat renders a stacked logo specifically used when a chat with a
// model is started: the word is split into two lines with "OPEN" above
// and "PILOT" below to emphasize branding in the chat wide layout.
// If there isn't enough width it falls back to the compact rendering.
func RenderChat(version string, o Opts) string {
	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	openLetters := []letterform{letterO, letterP, letterE, letterN}
	pilotLetters := []letterform{letterP, letterI, letterL, letterO, letterT}

	// Base renders without stretch for measurement.
	openBase := renderWord(1, -1, openLetters...)
	pilotBase := renderWord(1, -1, pilotLetters...)
	openW := lipgloss.Width(openBase)
	pilotW := lipgloss.Width(pilotBase)
	brandWidth := max(openW, pilotW)

	// Decide layout (treat combined height similar to wide layout).
	dec := decideLayout(o.Width, brandWidth, 6, 15, len(openLetters)+len(pilotLetters))
	if dec.Compact {
		return Render(version, true, o)
	}

	// Derive deterministic stretch indices separately (split space of letter count).
	stretchOpen := -1
	stretchPilot := -1
	if dec.StretchIndex >= 0 {
		if dec.StretchIndex < len(openLetters) {
			stretchOpen = dec.StretchIndex
		} else {
			stretchPilot = dec.StretchIndex - len(openLetters)
			if stretchPilot >= len(pilotLetters) {
				stretchPilot = -1
			}
		}
	}

	open := renderWord(1, stretchOpen, openLetters...)
	pilot := renderWord(1, stretchPilot, pilotLetters...)

	// Apply gradients if enabled.
	if dec.ApplyGradient {
		b := new(strings.Builder)
		for r := range strings.SplitSeq(open, "\n") {
			fmt.Fprintln(b, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
		}
		open = b.String()
		b = new(strings.Builder)
		for r := range strings.SplitSeq(pilot, "\n") {
			fmt.Fprintln(b, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
		}
		pilot = b.String()
	}

	const charm = " Surya™"
	metaRowGap := 1
	maxVersionWidth := brandWidth - lipgloss.Width(charm) - metaRowGap
	version = ansi.Truncate(version, maxVersionWidth, "…")
	gap := max(0, brandWidth-lipgloss.Width(charm)-lipgloss.Width(version))
	metaRow := fg(o.BrandColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)

	fieldHeight := lipgloss.Height(open) + lipgloss.Height(pilot) + 2
	// Left field.
	diagRune := diag
	if dec.UseSimpleDiag {
		diagRune = "/"
	}
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diagRune, dec.LeftFieldWidth))
	leftField := new(strings.Builder)
	for i := 0; i < fieldHeight; i++ {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightField := new(strings.Builder)
	for i := 0; i < fieldHeight; i++ {
		width := dec.RightFieldWidth
		if width < 0 {
			width = 0
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diagRune, width)), "\n")
	}

	content := metaRow + "\n" + open + "\n\n" + pilot
	content = padLines(content, brandWidth)
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, content, hGap, rightField.String())

	if o.Width > 0 {
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

// padLines right-pads each visual line (ANSI-aware) to target width using spaces.
func padLines(s string, width int) string {
	if width <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		w := lipgloss.Width(line)
		if w < width {
			lines[i] = line + strings.Repeat(" ", width-w)
		}
	}
	return strings.Join(lines, "\n")
}

// clampLines truncates each line to width (ANSI-aware) without adding ellipsis.
func clampLines(s string, width int) string {
	if width <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, width, "")
	}
	return strings.Join(lines, "\n")
}

// layoutDecision represents the decisions made for logo layout based on
// available width and other factors.
type layoutDecision struct {
	Compact         bool
	StretchIndex    int  // -1 if no stretch
	ApplyGradient   bool // disable on very tight widths
	UseSimpleDiag   bool // fallback to '/'
	RightFieldWidth int
	LeftFieldWidth  int
}

// decideLayout centralizes width-based decisions to keep Render deterministic.
func decideLayout(availableWidth int, wordWidth int, baseLeft int, defaultRight int, letterCount int) layoutDecision {
	d := layoutDecision{
		Compact:       false,
		StretchIndex:  -1,
		ApplyGradient: true,
		UseSimpleDiag: false,
		RightFieldWidth: func() int {
			if availableWidth > 0 {
				return max(1, availableWidth-wordWidth-baseLeft-2)
			}
			return defaultRight
		}(),
		LeftFieldWidth: baseLeft,
	}
	if availableWidth <= 0 { // unconstrained: ok to stretch deterministically by width hash
		if letterCount > 0 {
			idx := (wordWidth*31 + letterCount*17) % letterCount
			d.StretchIndex = idx
		}
		return d
	}
	minTotal := wordWidth + baseLeft + 2 + 1
	if availableWidth < minTotal {
		d.Compact = true
		d.RightFieldWidth = 0
	}
	if availableWidth < wordWidth+baseLeft+6 {
		d.ApplyGradient = false
	}
	if availableWidth < 50 {
		d.UseSimpleDiag = true
	}
	// Disable stretching for constrained width (stability); only allow if abundant slack >= +12.
	if !d.Compact && availableWidth >= minTotal+12 && letterCount > 0 {
		idx := (availableWidth + wordWidth*7 + letterCount*13) % letterCount
		d.StretchIndex = idx
	}
	if d.RightFieldWidth < 0 {
		d.RightFieldWidth = 0
	}
	return d
}
