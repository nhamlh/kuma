package generator_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/kumahq/kuma/pkg/dns/vips"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	mesh_core "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	model "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/dns/resolver"
	. "github.com/kumahq/kuma/pkg/test/matchers"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_common "github.com/kumahq/kuma/pkg/xds/envoy"
	"github.com/kumahq/kuma/pkg/xds/generator"
)

var _ = Describe("DNSGenerator", func() {

	type testCase struct {
		dataplaneFile string
		expected      string
	}

	DescribeTable("Generate Envoy xDS resources",
		func(given testCase) {
			// setup
			gen := &generator.DNSGenerator{}

			dnsResolver := resolver.NewDNSResolver("mesh")
			dnsResolver.SetVIPs(vips.List{
				vips.NewServiceEntry("backend_test-ns_svc_8080"): "240.0.0.0",
				vips.NewServiceEntry("httpbin"):                  "240.0.0.1",
			})
			ctx := xds_context.Context{
				ConnectionInfo: xds_context.ConnectionInfo{
					Authority: "kuma-system:5677",
				},
				ControlPlane: &xds_context.ControlPlaneContext{
					SdsTlsCert:  []byte("12345"),
					DNSResolver: dnsResolver,
				},
				Mesh: xds_context.MeshContext{
					Resource: &mesh_core.MeshResource{
						Meta: &test_model.ResourceMeta{
							Name: "default",
						},
						Spec: &mesh_proto.Mesh{
							Mtls: &mesh_proto.Mesh_Mtls{
								EnabledBackend: "builtin",
								Backends: []*mesh_proto.CertificateAuthorityBackend{
									{
										Name: "builtin",
										Type: "builtin",
									},
								},
							},
						},
					},
				},
			}

			dataplane := mesh_proto.Dataplane{}
			dpBytes, err := ioutil.ReadFile(filepath.Join("testdata", "dns", given.dataplaneFile))
			Expect(err).ToNot(HaveOccurred())
			Expect(util_proto.FromYAML(dpBytes, &dataplane)).To(Succeed())
			proxy := &model.Proxy{
				Id: *model.BuildProxyId("", "side-car"),
				Dataplane: &mesh_core.DataplaneResource{
					Meta: &test_model.ResourceMeta{
						Version: "1",
					},
					Spec: &dataplane,
				},
				APIVersion: envoy_common.APIV3,
				Routing: model.Routing{
					OutboundTargets: map[model.ServiceName][]model.Endpoint{
						"httpbin": {
							{
								Target: "httpbin.org",
							},
						},
					},
				},
				Metadata: &model.DataplaneMetadata{
					DNSPort:      53001,
					EmptyDNSPort: 53002,
				},
			}

			// when
			rs, err := gen.Generate(ctx, proxy)

			// then
			Expect(err).ToNot(HaveOccurred())

			// and output matches golden files
			resp, err := rs.List().ToDeltaDiscoveryResponse()
			Expect(err).ToNot(HaveOccurred())
			actual, err := util_proto.ToYAML(resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(MatchGoldenYAML(filepath.Join("testdata", "dns", given.expected)))
		},
		Entry("01. DNS enabled", testCase{
			dataplaneFile: "1-dataplane.input.yaml",
			expected:      "1-envoy-config.golden.yaml",
		}),
		Entry("02. DNS disabled", testCase{
			dataplaneFile: "2-dataplane.input.yaml",
			expected:      "2-envoy-config.golden.yaml",
		}),
	)
})
