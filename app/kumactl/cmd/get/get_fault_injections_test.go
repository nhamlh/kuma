package get_test

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
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	kumactl_resources "github.com/kumahq/kuma/app/kumactl/pkg/resources"

	"github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/app/kumactl/cmd"
	kumactl_cmd "github.com/kumahq/kuma/app/kumactl/pkg/cmd"
	config_proto "github.com/kumahq/kuma/pkg/config/app/kumactl/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	memory_resources "github.com/kumahq/kuma/pkg/plugins/resources/memory"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
)

var _ = Describe("kumactl get fault-injections", func() {

	faultInjectionResources := []*mesh.FaultInjectionResource{
		{
			Spec: &v1alpha1.FaultInjection{
				Sources: []*v1alpha1.Selector{
					{
						Match: map[string]string{
							"service": "frontend",
							"version": "0.1",
						},
					},
				},
				Destinations: []*v1alpha1.Selector{
					{
						Match: map[string]string{
							"service": "backend",
						},
					},
				},
				Conf: &v1alpha1.FaultInjection_Conf{
					Delay: &v1alpha1.FaultInjection_Conf_Delay{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						Value:      &durationpb.Duration{Seconds: 5},
					},
					Abort: &v1alpha1.FaultInjection_Conf_Abort{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						HttpStatus: &wrapperspb.UInt32Value{Value: 500},
					},
					ResponseBandwidth: &v1alpha1.FaultInjection_Conf_ResponseBandwidth{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						Limit:      &wrapperspb.StringValue{Value: "50 mbps"},
					},
				},
			},
			Meta: &test_model.ResourceMeta{
				Mesh: "default",
				Name: "fi1",
			},
		},
		{
			Spec: &v1alpha1.FaultInjection{
				Sources: []*v1alpha1.Selector{
					{
						Match: map[string]string{
							"service": "web",
							"version": "0.1",
						},
					},
				},
				Destinations: []*v1alpha1.Selector{
					{
						Match: map[string]string{
							"service": "redis",
						},
					},
				},
				Conf: &v1alpha1.FaultInjection_Conf{
					Delay: &v1alpha1.FaultInjection_Conf_Delay{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						Value:      &durationpb.Duration{Seconds: 5},
					},
					Abort: &v1alpha1.FaultInjection_Conf_Abort{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						HttpStatus: &wrapperspb.UInt32Value{Value: 500},
					},
					ResponseBandwidth: &v1alpha1.FaultInjection_Conf_ResponseBandwidth{
						Percentage: &wrapperspb.DoubleValue{Value: 50},
						Limit:      &wrapperspb.StringValue{Value: "50 mbps"},
					},
				},
			},
			Meta: &test_model.ResourceMeta{
				Mesh: "default",
				Name: "fi2",
			},
		},
	}

	Describe("GetFaultInjectionCmd", func() {

		var rootCtx *kumactl_cmd.RootContext
		var rootCmd *cobra.Command
		var buf *bytes.Buffer
		var store core_store.ResourceStore
		rootTime, _ := time.Parse(time.RFC3339, "2008-04-27T16:05:36.995Z")
		BeforeEach(func() {
			// setup
			rootCtx = &kumactl_cmd.RootContext{
				Runtime: kumactl_cmd.RootRuntime{
					Now: func() time.Time { return rootTime },
					NewResourceStore: func(*config_proto.ControlPlaneCoordinates_ApiServer) (core_store.ResourceStore, error) {
						return store, nil
					},
					NewAPIServerClient: kumactl_resources.NewAPIServerClient,
				},
			}

			store = core_store.NewPaginationStore(memory_resources.NewStore())

			for _, ds := range faultInjectionResources {
				err := store.Create(context.Background(), ds, core_store.CreateBy(core_model.MetaToResourceKey(ds.GetMeta())))
				Expect(err).ToNot(HaveOccurred())
			}

			rootCmd = cmd.NewRootCmd(rootCtx)
			buf = &bytes.Buffer{}
			rootCmd.SetOut(buf)
		})

		type testCase struct {
			outputFormat string
			goldenFile   string
			pagination   string
			matcher      func(interface{}) gomega_types.GomegaMatcher
		}

		DescribeTable("kumactl get fault-injections -o table|json|yaml",
			func(given testCase) {
				// given
				rootCmd.SetArgs(append([]string{
					"--config-file", filepath.Join("..", "testdata", "sample-kumactl.config.yaml"),
					"get", "fault-injections"}, given.outputFormat, given.pagination))

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
				goldenFile:   "get-fault-injections.golden.txt",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support Table output explicitly", testCase{
				outputFormat: "-otable",
				goldenFile:   "get-fault-injections.golden.txt",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support pagination", testCase{
				outputFormat: "-otable",
				goldenFile:   "get-fault-injections.pagination.golden.txt",
				pagination:   "--size=1",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support JSON output", testCase{
				outputFormat: "-ojson",
				goldenFile:   "get-fault-injections.golden.json",
				matcher:      MatchJSON,
			}),
			Entry("should support YAML output", testCase{
				outputFormat: "-oyaml",
				goldenFile:   "get-fault-injections.golden.yaml",
				matcher:      MatchYAML,
			}),
		)
	})
})
