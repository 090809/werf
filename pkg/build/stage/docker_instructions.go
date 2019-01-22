package stage

import (
	"github.com/flant/werf/pkg/config"
	"github.com/flant/werf/pkg/image"
	"github.com/flant/werf/pkg/util"
)

func GenerateDockerInstructionsStage(imageConfig *config.Image, baseStageOptions *NewBaseStageOptions) *DockerInstructionsStage {
	if imageConfig.Docker != nil {
		return newDockerInstructionsStage(imageConfig.Docker, baseStageOptions)
	}

	return nil
}

func newDockerInstructionsStage(instructions *config.Docker, baseStageOptions *NewBaseStageOptions) *DockerInstructionsStage {
	s := &DockerInstructionsStage{}
	s.instructions = instructions
	s.BaseStage = newBaseStage(DockerInstructions, baseStageOptions)
	return s
}

type DockerInstructionsStage struct {
	*BaseStage

	instructions *config.Docker
}

func (s *DockerInstructionsStage) GetDependencies(_ Conveyor, _ image.ImageInterface) (string, error) {
	var args []string

	args = append(args, s.instructions.Volume...)
	args = append(args, s.instructions.Expose...)

	for k, v := range s.instructions.Env {
		args = append(args, k, v)
	}

	for k, v := range s.instructions.Label {
		args = append(args, k, v)
	}

	args = append(args, s.instructions.Cmd...)
	args = append(args, s.instructions.Onbuild...)
	args = append(args, s.instructions.Entrypoint...)
	args = append(args, s.instructions.Workdir)
	args = append(args, s.instructions.User)
	args = append(args, s.instructions.StopSignal)
	args = append(args, s.instructions.HealthCheck)

	return util.Sha256Hash(args...), nil
}

func (s *DockerInstructionsStage) PrepareImage(c Conveyor, prevBuiltImage, image image.ImageInterface) error {
	if err := s.BaseStage.PrepareImage(c, prevBuiltImage, image); err != nil {
		return err
	}

	imageCommitChangeOptions := image.Container().CommitChangeOptions()
	imageCommitChangeOptions.AddVolume(s.instructions.Volume...)
	imageCommitChangeOptions.AddExpose(s.instructions.Expose...)
	imageCommitChangeOptions.AddEnv(s.instructions.Env)
	imageCommitChangeOptions.AddCmd(s.instructions.Cmd...)
	imageCommitChangeOptions.AddOnbuild(s.instructions.Onbuild...)
	imageCommitChangeOptions.AddEntrypoint(s.instructions.Entrypoint...)
	imageCommitChangeOptions.AddUser(s.instructions.User)
	imageCommitChangeOptions.AddWorkdir(s.instructions.Workdir)
	imageCommitChangeOptions.AddStopSignal(s.instructions.StopSignal)
	imageCommitChangeOptions.AddHealthCheck(s.instructions.HealthCheck)

	return nil
}
