package ci_env_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"

	"github.com/flant/werf/pkg/testing/utils"
	"github.com/flant/werf/pkg/werf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ci-env", func() {
	BeforeEach(func() {
		Ω(werf.Init("", "")).Should(Succeed())
	})

	ciSystems := []string{
		"gitlab",
		"github",
	}

	for i := range ciSystems {
		ciSystem := ciSystems[i]

		Context(ciSystem, func() {
			It("should print only script path", func() {
				output := utils.SucceedCommandOutputString(
					"",
					werfBinPath,
					utils.WerfBinArgs("ci-env", ciSystem, "--as-file")...,
				)

				expectedPathGlob := filepath.Join(
					werf.GetServiceDir(),
					"tmp",
					"ci_env",
					"source_*_*",
				)

				resultPath := strings.TrimSuffix(output, "\n")
				matched, err := doublestar.PathMatch(expectedPathGlob, resultPath)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(matched).Should(BeTrue(), output)
				Ω(resultPath).Should(BeARegularFile())
			})

			It("should print only shell script", func() {
				output := utils.SucceedCommandOutputString(
					"",
					werfBinPath,
					utils.WerfBinArgs("ci-env", ciSystem)...,
				)

				useAsFileOutput := utils.SucceedCommandOutputString(
					"",
					werfBinPath,
					utils.WerfBinArgs("ci-env", ciSystem, "--as-file")...,
				)

				scriptPath := strings.TrimSpace(useAsFileOutput)
				scriptData, err := ioutil.ReadFile(scriptPath)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(len(string(scriptData))).Should(Equal(len(output)))
			})
		})
	}
})
