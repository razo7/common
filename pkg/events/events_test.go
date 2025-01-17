package events

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

var r *record.FakeRecorder

var _ = Describe("Emit custom formatted Event", func() {
	BeforeEach(func() {
		r = record.NewFakeRecorder(4)
	})

	Context("Emit event via custom API", func() {
		DescribeTable("Custom APIs",
			func(obj runtime.Object, eventType string, eventReason string, message string, expected string, args ...interface{}) {
				if eventType == "Normal" {
					if len(args) != 0 {
						NormalEventf(r, obj, eventReason, message, args...)
					} else {
						NormalEvent(r, obj, eventReason, message)
					}
				} else {
					if len(args) != 0 {
						WarningEventf(r, obj, eventReason, message, args...)
					} else {
						WarningEvent(r, obj, eventReason, message)
					}
				}
				for {
					select {
					case got := <-r.Events:
						Expect(got).To(Equal(expected))
						break
					case <-time.After(1 * time.Second):
						Fail("Timeout waiting for event")
					}
					break

				}
			},
			Entry("Emit normal event",
				nil,
				"Normal", "thisReason", "normal event message",
				"Normal thisReason [remediation] normal event message"),
			Entry("Emit normal event with args",
				nil,
				"Normal", "thisReason", "normal event message with some arguments: %s%d, %s%d",
				"Normal thisReason [remediation] normal event message with some arguments: somevalue1, somevalue2",
				"somevalue", 1, "somevalue", 2),
			Entry("Emit warning event",
				nil,
				"Warning", "thisReason", "warning event message",
				"Warning thisReason [remediation] warning event message"),
			Entry("Emit warning event with args",
				nil,
				"Warning", "thisReason", "warning event message with some arguments: %s%d, %s%d",
				"Warning thisReason [remediation] warning event message with some arguments: somevalue1, somevalue2",
				"somevalue", 1, "somevalue", 2),
		)
	})
})
