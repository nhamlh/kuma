package generate_test

import (
	"bytes"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	kumactl_resources "github.com/kumahq/kuma/app/kumactl/pkg/resources"

	"github.com/kumahq/kuma/app/kumactl/cmd"
	kumactl_cmd "github.com/kumahq/kuma/app/kumactl/pkg/cmd"
	"github.com/kumahq/kuma/app/kumactl/pkg/tokens"
	config_proto "github.com/kumahq/kuma/pkg/config/app/kumactl/v1alpha1"
)

type staticZoneIngressTokenGenerator struct {
	err error
}

var _ tokens.ZoneIngressTokenClient = &staticZoneIngressTokenGenerator{}

func (s *staticZoneIngressTokenGenerator) Generate(zone string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return fmt.Sprintf("token-for-%s", zone), nil
}

var _ = Describe("kumactl generate zone-ingress-token", func() {

	var rootCmd *cobra.Command
	var buf *bytes.Buffer
	var generator *staticZoneIngressTokenGenerator
	var ctx *kumactl_cmd.RootContext

	BeforeEach(func() {
		generator = &staticZoneIngressTokenGenerator{}
		ctx = &kumactl_cmd.RootContext{
			Runtime: kumactl_cmd.RootRuntime{
				NewZoneIngressTokenClient: func(*config_proto.ControlPlaneCoordinates_ApiServer) (tokens.ZoneIngressTokenClient, error) {
					return generator, nil
				},
				NewAPIServerClient: kumactl_resources.NewAPIServerClient,
			},
		}

		rootCmd = cmd.NewRootCmd(ctx)

		buf = &bytes.Buffer{}
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
	})

	type testCase struct {
		args   []string
		result string
	}
	DescribeTable("should generate token",
		func(given testCase) {
			// when
			rootCmd.SetArgs(given.args)
			err := rootCmd.Execute()

			// then
			Expect(err).ToNot(HaveOccurred())

			// and
			Expect(buf.String()).To(Equal(given.result))
		},
		Entry("for zone", testCase{
			args:   []string{"generate", "zone-ingress-token", "--zone=my-zone"},
			result: "token-for-my-zone",
		}),
		Entry("for empty zone", testCase{
			args:   []string{"generate", "zone-ingress-token"},
			result: "token-for-",
		}),
	)

	It("should write error when generating token fails", func() {
		// setup
		generator.err = errors.New("could not connect to API")

		// when
		rootCmd.SetArgs([]string{"generate", "zone-ingress-token", "--zone=example"})
		err := rootCmd.Execute()

		// then
		Expect(err).To(HaveOccurred())

		// and
		Expect(buf.String()).To(Equal("Error: failed to generate a zone ingress token: could not connect to API\n"))
	})

})
