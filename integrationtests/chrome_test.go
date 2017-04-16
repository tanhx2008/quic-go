package integrationtests

import (
	"context"
	"fmt"
	"strconv"

	"github.com/knq/chromedp"

	"github.com/lucas-clemente/quic-go/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Chrome tests", func() {
	It("does not work with mismatching versions", func() {
		versionForUs := protocol.SupportedVersions[0]
		versionForChrome := protocol.SupportedVersions[len(protocol.SupportedVersions)-1]

		// If both are equal, this test doesn't make any sense.
		if versionForChrome == versionForUs {
			return
		}

		supportedVersionsBefore := protocol.SupportedVersions
		protocol.SupportedVersions = []protocol.VersionNumber{versionForUs}

		ctx, chrome := chromeForVersion(versionForChrome)

		defer func() {
			protocol.SupportedVersions = supportedVersionsBefore
		}()

		navigate(ctx, chrome, "https://quic.clemente.io/hello")
		Consistently(func() string {
			res := ""
			chrome.Run(ctx, chromedp.Tasks{chromedp.Text("body", &res)})
			return res
		}).ShouldNot(ContainSubstring("Hello, World!"))
	})

	for i := range protocol.SupportedVersions {
		version := protocol.SupportedVersions[i]

		Context(fmt.Sprintf("with quic version %d", version), func() {
			var (
				ctx                     context.Context
				chrome                  *chromedp.CDP
				supportedVersionsBefore []protocol.VersionNumber
			)

			BeforeEach(func() {
				// fmt.Printf("MAP quic.clemente.io:443 localhost:%s", port)
				// time.Sleep(time.Hour)
				supportedVersionsBefore = protocol.SupportedVersions
				protocol.SupportedVersions = []protocol.VersionNumber{version}
				ctx, chrome = chromeForVersion(version)
			})

			AfterEach(func() {
				protocol.SupportedVersions = supportedVersionsBefore
			})

			It("loads a simple hello world page using quic", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/hello")
				waitForText(ctx, chrome, "Hello, World!")
			})

			It("downloads a small file", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/downloadtest?num=1&len="+strconv.Itoa(dataLen))
				waitForText(ctx, chrome, "dltest ok")
			})

			It("downloads a large file", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/downloadtest?num=1&len="+strconv.Itoa(dataLongLen))
				waitForText(ctx, chrome, "dltest ok")
			})

			It("loads a large number of files", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/downloadtest?num=4&len=100")
				waitForText(ctx, chrome, "dltest ok")
			})

			It("uploads a small file", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/uploadtest?num=1&len="+strconv.Itoa(dataLen))
				Eventually(func() int32 { return nFilesUploaded }).Should(BeEquivalentTo(1))
			})

			It("uploads a large file", func() {
				navigate(ctx, chrome, "https://quic.clemente.io/uploadtest?num=1&len="+strconv.Itoa(dataLongLen))
				Eventually(func() int32 { return nFilesUploaded }, 30).Should(BeEquivalentTo(1))
			})

			It("uploads many small files", func() {
				num := protocol.MaxStreamsPerConnection + 20
				navigate(ctx, chrome, "https://quic.clemente.io/uploadtest?len="+strconv.Itoa(dataLen)+"&num="+strconv.Itoa(num))
				Eventually(func() int32 { return nFilesUploaded }, 30).Should(BeEquivalentTo(num))
			})
		})
	}
})
