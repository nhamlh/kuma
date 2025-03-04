package inspect_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	gomega_types "github.com/onsi/gomega/types"
	"github.com/spf13/cobra"

	kumactl_resources "github.com/kumahq/kuma/app/kumactl/pkg/resources"

	"github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/app/kumactl/cmd"
	kumactl_cmd "github.com/kumahq/kuma/app/kumactl/pkg/cmd"
	"github.com/kumahq/kuma/app/kumactl/pkg/resources"
	config_proto "github.com/kumahq/kuma/pkg/config/app/kumactl/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/test/resources/model"
)

type testServiceOverviewClient struct {
	total     uint32
	overviews []*core_mesh.ServiceOverviewResource
}

func (c *testServiceOverviewClient) List(_ context.Context, mesh string) (*core_mesh.ServiceOverviewResourceList, error) {
	return &core_mesh.ServiceOverviewResourceList{
		Items: c.overviews,
		Pagination: core_model.Pagination{
			Total: c.total,
		},
	}, nil
}

var _ resources.ServiceOverviewClient = &testServiceOverviewClient{}

var _ = Describe("kumactl inspect services", func() {

	var rootCtx *kumactl_cmd.RootContext
	var rootCmd *cobra.Command
	var buf *bytes.Buffer
	rootTime, _ := time.Parse(time.RFC3339, "2008-04-27T16:05:36.995Z")

	serviceOverviewResources := []*core_mesh.ServiceOverviewResource{
		{
			Meta: &model.ResourceMeta{Mesh: "mesh-1", Name: "backend"},
			Spec: &v1alpha1.ServiceInsight_Service{
				Status: v1alpha1.ServiceInsight_Service_partially_degraded,
				Dataplanes: &v1alpha1.ServiceInsight_Service_DataplaneStat{
					Online: 5,
					Total:  10,
				},
			},
		},
		{
			Meta: &model.ResourceMeta{Mesh: "mesh-1", Name: "web"},
			Spec: &v1alpha1.ServiceInsight_Service{
				Status: v1alpha1.ServiceInsight_Service_online,
				Dataplanes: &v1alpha1.ServiceInsight_Service_DataplaneStat{
					Online: 20,
					Total:  20,
				},
			},
		},
		{
			Meta: &model.ResourceMeta{Mesh: "mesh-1", Name: "orders"},
			Spec: &v1alpha1.ServiceInsight_Service{
				Status: v1alpha1.ServiceInsight_Service_offline,
				Dataplanes: &v1alpha1.ServiceInsight_Service_DataplaneStat{
					Online: 0,
					Total:  5,
				},
			},
		},
	}

	BeforeEach(func() {
		rootCtx = &kumactl_cmd.RootContext{
			Runtime: kumactl_cmd.RootRuntime{
				Now: func() time.Time { return rootTime },
				NewServiceOverviewClient: func(server *config_proto.ControlPlaneCoordinates_ApiServer) (resources.ServiceOverviewClient, error) {
					return &testServiceOverviewClient{
						total:     uint32(len(serviceOverviewResources)),
						overviews: serviceOverviewResources,
					}, nil
				},
				NewAPIServerClient: kumactl_resources.NewAPIServerClient,
			},
		}

		rootCmd = cmd.NewRootCmd(rootCtx)
		buf = &bytes.Buffer{}
		rootCmd.SetOut(buf)
	})

	type testCase struct {
		outputFormat string
		goldenFile   string
		matcher      func(interface{}) gomega_types.GomegaMatcher
	}

	DescribeTable("kumactl inspect meshes -o table|json|yaml",
		func(given testCase) {
			// given
			rootCmd.SetArgs(append([]string{
				"--config-file", filepath.Join("..", "testdata", "sample-kumactl.config.yaml"),
				"inspect", "services"}, given.outputFormat))

			// when
			err := rootCmd.Execute()
			// then
			Expect(err).ToNot(HaveOccurred())

			// when
			expected, err := ioutil.ReadFile(filepath.Join("testdata", given.goldenFile))
			// then
			Expect(err).ToNot(HaveOccurred())
			// and
			Expect(buf.String()).To(given.matcher(expected))
		},
		Entry("should support Table output by default", testCase{
			outputFormat: "",
			goldenFile:   "inspect-services.golden.txt",
			matcher: func(expected interface{}) gomega_types.GomegaMatcher {
				return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
			},
		}),
		Entry("should support Table output explicitly", testCase{
			outputFormat: "-otable",
			goldenFile:   "inspect-services.golden.txt",
			matcher: func(expected interface{}) gomega_types.GomegaMatcher {
				return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
			},
		}),
		Entry("should support JSON output", testCase{
			outputFormat: "-ojson",
			goldenFile:   "inspect-services.golden.json",
			matcher:      MatchJSON,
		}),
		Entry("should support YAML output", testCase{
			outputFormat: "-oyaml",
			goldenFile:   "inspect-services.golden.yaml",
			matcher:      MatchYAML,
		}),
	)
})
