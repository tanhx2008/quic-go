package integrationtests

import (
	"fmt"

	"github.com/knq/chromedp"
	"github.com/knq/chromedp/runner"
	"github.com/lucas-clemente/quic-go/protocol"
	"golang.org/x/net/context"

	. "github.com/onsi/gomega"
)

type chromeInstance struct {
	ctx    context.Context
	chrome *chromedp.CDP
	cancel context.CancelFunc
}

var chromes = map[protocol.VersionNumber]chromeInstance{}

func initChromes() {
	for _, v := range protocol.SupportedVersions {
		ctx, cancel := context.WithCancel(context.Background())
		cdp, err := chromedp.New(ctx, chromedp.WithRunnerOptions(
			runner.Port(10000+int(v)),
			runner.Flag("enable-quic", true),
			runner.Flag("no-proxy-server", true),
			runner.Flag("origin-to-force-quic-on", "quic.clemente.io:443"),
			runner.Flag("host-resolver-rules",
				fmt.Sprintf("MAP quic.clemente.io:443 localhost:%s", port)),
			runner.Flag("quic-version", fmt.Sprintf("QUIC_VERSION_%d", v)),
		))
		Expect(err).NotTo(HaveOccurred())
		chromes[v] = chromeInstance{ctx: ctx, chrome: cdp, cancel: cancel}
	}
}

func killChromes() {
	for _, c := range chromes {
		defer c.cancel()
		err := c.chrome.Shutdown(c.ctx)
		Expect(err).NotTo(HaveOccurred())
		err = c.chrome.Wait()
		Expect(err).NotTo(HaveOccurred())
	}
}

func chromeForVersion(v protocol.VersionNumber) (context.Context, *chromedp.CDP) {
	return chromes[v].ctx, chromes[v].chrome
}

func navigate(ctx context.Context, chrome *chromedp.CDP, url string) {
	err := chrome.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
	})
	Expect(err).NotTo(HaveOccurred())
}

func waitForText(ctx context.Context, chrome *chromedp.CDP, text string) {
	Eventually(func() string {
		var res string
		err := chrome.Run(ctx, chromedp.Tasks{chromedp.Text("body", &res)})
		Expect(err).NotTo(HaveOccurred())
		return res
	}, 20).Should(ContainSubstring(text))
}
