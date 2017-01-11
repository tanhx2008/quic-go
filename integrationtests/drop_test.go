package integrationtests

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	_ "github.com/lucas-clemente/quic-clients" // download clients
	"github.com/lucas-clemente/quic-go/integrationtests/proxy"
	"github.com/lucas-clemente/quic-go/protocol"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Drop Proxy", func() {
	BeforeEach(func() {
		dataMan.GenerateData(dataLen)
	})

	clientPath := fmt.Sprintf(
		"%s/src/github.com/lucas-clemente/quic-clients/client-%s-debug",
		os.Getenv("GOPATH"),
		runtime.GOOS,
	)

	var dropproxy *proxy.UDPProxy
	proxyPort := 12345

	setupDropProxy := func(incomingPacketDropper, outgoingPacketDropper proxy.DropCallback) {
		iPort, _ := strconv.Atoi(port)
		var err error
		dropproxy, err = proxy.NewUDPProxy(proxyPort, "localhost", iPort, incomingPacketDropper, outgoingPacketDropper, 0, 0)
		Expect(err).ToNot(HaveOccurred())
	}

	AfterEach(func() {
		dropproxy.Stop()
		time.Sleep(time.Millisecond)
	})

	for i := range protocol.SupportedVersions {
		version := protocol.SupportedVersions[i]

		Context(fmt.Sprintf("with quic version %d", version), func() {
			Context("packet loss during the handshake", func() {
				It("dropping every 4th packet of all packets", func() {
					rand.Seed(time.Now().UTC().UnixNano())
					dropper := func(p protocol.PacketNumber) bool {
						return rand.Int31()%4 == 0
					}
					setupDropProxy(dropper, dropper)

					command := exec.Command(
						clientPath,
						"--quic-version="+strconv.Itoa(int(version)),
						"--host=127.0.0.1",
						"--port="+strconv.Itoa(proxyPort),
						"https://quic.clemente.io/hello",
					)
					session, err := Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					defer session.Kill()
					Eventually(session, 20).Should(Exit(0))
					Expect(session.Out).To(Say("Response:\nheaders: HTTP/1.1 200\nstatus: 200\n\nbody: Hello, World!\n"))
				}, 30)
			})

			Context("dropping every 4th packet after the crypto handshake", func() {
				dropper := func(p protocol.PacketNumber) bool {
					if p <= 10 { // don't interfere with the crypto handshake
						return false
					}
					return p%4 == 0
				}

				runDropTest := func(incomingPacketDropper, outgoingPacketDropper proxy.DropCallback, version protocol.VersionNumber) {
					setupDropProxy(incomingPacketDropper, outgoingPacketDropper)

					command := exec.Command(
						clientPath,
						"--quic-version="+strconv.Itoa(int(version)),
						"--host=127.0.0.1",
						"--port="+strconv.Itoa(proxyPort),
						"https://quic.clemente.io/data",
					)

					session, err := Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					defer session.Kill()
					Eventually(session, 20).Should(Exit(0))
					Expect(bytes.Contains(session.Out.Contents(), dataMan.GetData())).To(BeTrue())
				}

				It("gets a file when many outgoing packets are dropped", func() {
					runDropTest(nil, dropper, version)
				})

				It("gets a file when many incoming packets are dropped", func() {
					runDropTest(dropper, nil, version)
				})
			})
		})
	}
})
