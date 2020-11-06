package werf_chart

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"

	"github.com/werf/logboek"
	"github.com/werf/werf/pkg/image"
)

type ServiceValuesOptions struct {
	Env    string
	IsStub bool
}

func GetServiceValues(ctx context.Context, projectName string, repo, namespace string, imageInfoGetters []*image.InfoGetter, opts ServiceValuesOptions) (map[string]interface{}, error) {
	werfInfo := map[string]interface{}{
		"name":      projectName,
		"repo":      repo,
		"namespace": namespace,
		"is_stub":   opts.IsStub,
	}

	if opts.IsStub {
		werfInfo["stub_image"] = fmt.Sprintf("%s:TAG", repo)
	}

	if opts.Env != "" {
		werfInfo["env"] = opts.Env
	}

	for _, imageInfoGetter := range imageInfoGetters {
		if imageInfoGetter.IsNameless() {
			werfInfo["is_nameless_image"] = true
			werfInfo["image"] = imageInfoGetter.GetName()
		} else {
			if werfInfo["image"] == nil {
				werfInfo["image"] = map[string]interface{}{}
			}
			werfInfo["image"].(map[string]interface{})[imageInfoGetter.GetWerfImageName()] = imageInfoGetter.GetName()
		}
	}

	res := map[string]interface{}{"werf": werfInfo}

	data, err := yaml.Marshal(res)
	logboek.Context(ctx).Debug().LogF("GetServiceValues result (err=%s):\n%s\n", err, data)

	return res, nil
}
