	"path/filepath"
			p.state = modifyFileDiff
		if strings.HasPrefix(line, "index ") {
			return p.handleIndexDiffLine(line)
		}
		if strings.HasPrefix(line, "index ") {
			return p.handleIndexDiffLine(line)
		}
		if strings.HasPrefix(line, "index ") {
			return p.handleIndexDiffLine(line)
		}
			newPath := p.trimFileBaseFilepath(path)
			newPath := p.trimFileBaseFilepath(path)
	var prefix, hashes, suffix string
	if len(parts) == 3 {
		prefix, hashes, suffix = parts[0], parts[1], parts[2]
	} else if len(parts) == 2 {
		prefix, hashes = parts[0], parts[1]
	} else {
		// TODO: remove index line from resulting patch completely in v1.2
	var newLine string

	if suffix == "" {
		newLine = fmt.Sprintf("%s %s..%s", prefix, strings.Join(leftHashes, ","), strings.Join(rightHashes, ","))
	} else {
		newLine = fmt.Sprintf("%s %s..%s %s", prefix, strings.Join(leftHashes, ","), strings.Join(rightHashes, ","), suffix)
	}
	newPath := p.trimFileBaseFilepath(path)
	newPath := p.trimFileBaseFilepath(path)
	if strings.HasSuffix(line, " (commits not present)") {
		return fmt.Errorf("cannot handle \"commits not present\" in git diff line %q, check specified submodule commits are correct", line)
	}
	newPath := p.trimFileBaseFilepath(path)
func (p *diffParser) trimFileBaseFilepath(path string) string {
	return filepath.ToSlash(p.PathMatcher.TrimFileBaseFilepath(filepath.FromSlash(path)))
}

	newPath := p.trimFileBaseFilepath(path)