package v3

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	envoy_accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
)

var _ = Describe("defaultHandler", func() {
	Describe("error path", func() {
		type testCase struct {
			msg         *envoy_accesslog.StreamAccessLogsMessage
			expectedErr string
		}

		DescribeTable("should fail if configuration is not valid",
			func(given testCase) {
				// when
				_, err := defaultHandler(nil, given.msg)
				// then
				Expect(err).To(HaveOccurred())
				// and
				Expect(err.Error()).To(Equal(given.expectedErr))
			},
			Entry("empty `identifier.log_name` field", testCase{
				expectedErr: `log name "" has invalid format: expected 2 components separated by ';', got 1`,
			}),
			Entry("invalid access log format string", testCase{
				msg: &envoy_accesslog.StreamAccessLogsMessage{
					Identifier: &envoy_accesslog.StreamAccessLogsMessage_Identifier{
						LogName: ";%bytes_sent%",
					},
				},
				expectedErr: `format string is not valid: expected a command operator to start at position 1, instead got: "%bytes_sent%"`,
			}),
		)
	})
})
