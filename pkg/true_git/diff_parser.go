
	"github.com/flant/werf/pkg/path_matcher"
func makeDiffParser(out io.Writer, pathMatcher path_matcher.PathMatcher) *diffParser {
		PathMatcher: pathMatcher,
	PathMatcher path_matcher.PathMatcher
			return p.handleIndexDiffLine(line)
			if !p.PathMatcher.MatchPath(path) {
			newPath := p.PathMatcher.TrimFileBaseFilepath(path)
			if !p.PathMatcher.MatchPath(path) {
			newPath := p.PathMatcher.TrimFileBaseFilepath(path)
// TODO: remove index line from resulting patch completely in v1.2
func (p *diffParser) handleIndexDiffLine(line string) error {
	p.state = modifyFileDiff

	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		// unexpected format
		return p.writeOutLine(line)
	}

	prefix, hashes, suffix := parts[0], parts[1], parts[2]

	hashesParts := strings.SplitN(hashes, "..", 2)
	if len(hashesParts) != 2 {
		// unexpected format
		return p.writeOutLine(line)
	}

	stripHashFunc := func(h string) string {
		if len(h) < 8 {
			return h
		}
		return h[:8]
	}

	var leftHashes []string
	for _, h := range strings.Split(hashesParts[0], ",") {
		leftHashes = append(leftHashes, stripHashFunc(h))
	}

	var rightHashes []string
	for _, h := range strings.Split(hashesParts[1], ",") {
		rightHashes = append(rightHashes, stripHashFunc(h))
	}

	newLine := fmt.Sprintf("%s %s..%s %s", prefix, strings.Join(leftHashes, ","), strings.Join(rightHashes, ","), suffix)

	return p.writeOutLine(newLine)
}

	newPath := p.PathMatcher.TrimFileBaseFilepath(path)
	newPath := p.PathMatcher.TrimFileBaseFilepath(path)
	newPath := p.PathMatcher.TrimFileBaseFilepath(path)
	newPath := p.PathMatcher.TrimFileBaseFilepath(path)