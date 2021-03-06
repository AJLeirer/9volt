package monitor

import (
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/9corp/9volt/util"
)

var _ = Describe("dns_monitor", func() {
	var (
		monitor *DnsMonitor
		config  *RootMonitorConfig
	)

	BeforeEach(func() {
		config = &RootMonitorConfig{
			Config: &MonitorConfig{
				Host:          "beowulf",
				DnsTarget:     "grendel",
				Expect:        "IN A",
				DnsRecordType: "A",
				Interval:      util.CustomDuration(3 * time.Second),
			},
			Log: log.New(),
		}

		monitor = NewDnsMonitor(config)
	})

	Context("NewDnsMonitor", func() {
		It("returns a properly configured instance with defaults", func() {
			Expect(monitor.Timeout).To(Equal(DEFAULT_DNS_TIMEOUT))
			Expect(monitor.RecordType).To(Equal(DEFAULT_DNS_RECORD_TYPE))
			Expect(monitor.Client).NotTo(BeNil())
			Expect(monitor.MonitorFunc).NotTo(BeNil())
			Expect(monitor.Expect).NotTo(BeNil())
		})

		It("sources some instance configs from the main config", func() {
			Expect(monitor.Expect).NotTo(BeNil())
		})
	})

	Context("Validate", func() {
		It("returns nil with correct settings", func() {
			Expect(monitor.Validate()).To(BeNil())
		})

		Context("with bad settings", func() {
			It("bad record type", func() {
				config.Config.DnsRecordType = "foobar"
				err := monitor.Validate()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("Unknown record type"))
			})

			It("bad timeouts", func() {
				config.Config.Interval = 0
				err := monitor.Validate()

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("cannot equal or exceed"))
			})

			It("bad DNS target host", func() {
				config.Config.DnsTarget = ""
				err := monitor.Validate()

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("No DNS target"))
			})

			It("bad regex", func() {
				config.Config.Expect = "["
				err := monitor.Validate()

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("Unable to compile"))
			})
		})
	})
})
