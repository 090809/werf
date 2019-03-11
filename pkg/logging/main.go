package logging

import (
	"fmt"

	"github.com/flant/werf/pkg/logger"
)

var (
	imageNameFormat    = "⛵ image %s"
	artifactNameFormat = "🛸 artifact %s"
)

func Init() error {
	return logger.Init()
}

func EnableLogColor() {
	logger.EnableLogColor()
}

func DisableLogColor() {
	logger.DisableLogColor()
}

func SetTerminalWidth(value int) {
	logger.SetTerminalWidth(value)
}

func DisablePrettyLog() {
	imageNameFormat = "image %s"
	artifactNameFormat = "artifact %s"

	logger.DisablePrettyLog()
}

func ImageLogName(name string, isArtifact bool) string {
	if !isArtifact {
		if name == "" {
			name = "~"
		}
	}

	return name
}

func ImageLogProcessName(name string, isArtifact bool) string {
	logName := ImageLogName(name, isArtifact)
	if !isArtifact {
		return fmt.Sprintf(imageNameFormat, logName)
	} else {
		return fmt.Sprintf(artifactNameFormat, logName)
	}
}
