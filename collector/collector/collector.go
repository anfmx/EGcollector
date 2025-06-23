package collector

import (
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const (
	chromeBin     = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	chromeProfile = `./ChromeProfile`
)

type Collector struct {
	browser    *rod.Browser
	baseUrl    *url.URL
	userData   string
	browserBin string
}

type CollectorBuilder struct {
	collector Collector
}

func NewCollectorBuilder() *CollectorBuilder {
	return &CollectorBuilder{}
}

func (afb *CollectorBuilder) WithBrowser(browser *rod.Browser) *CollectorBuilder {
	afb.collector.browser = browser
	return afb
}
func (afb *CollectorBuilder) WithBaseUrl(baseUrl *url.URL) *CollectorBuilder {
	afb.collector.baseUrl = baseUrl
	return afb
}
func (afb *CollectorBuilder) WithUserData(userData string) *CollectorBuilder {
	afb.collector.userData = userData
	return afb
}
func (afb *CollectorBuilder) WithBrowserBin(browserBin string) *CollectorBuilder {
	afb.collector.browserBin = browserBin
	return afb
}
func (afb *CollectorBuilder) Build() Collector {
	return afb.collector
}

type CollectorBuilderDirector struct {
	collectorBuilder *CollectorBuilder
}

func (ad *CollectorBuilderDirector) NewChromeCollector() Collector {
	ws := launcher.NewUserMode().
		Bin(chromeBin).
		UserDataDir(chromeProfile).
		MustLaunch()

	browser := rod.New().ControlURL(ws).MustConnect().NoDefaultDevice()

	base, _ := url.Parse("https://store.epicgames.com")

	return ad.collectorBuilder.
		WithBaseUrl(base).
		WithBrowser(browser).
		WithBrowserBin(chromeBin).
		WithUserData(chromeProfile).
		Build()
}

func NewCollectorDirector() *CollectorBuilderDirector {
	collectorBuilder := NewCollectorBuilder()
	return &CollectorBuilderDirector{collectorBuilder: collectorBuilder}
}

func (a Collector) GetGames() []*rod.Element {
	page := a.browser.MustPage("https://store.epicgames.com/en-US/free-games")

	page.MustSetViewport(1920, 1080, 1, false)
	page.MustWaitLoad()

	container := page.MustElement(".css-2u323")

	return container.MustElements(".css-g3jcms")
}

func (a Collector) AddToCart(href string) {
	parsed, err := url.Parse(href)
	if err != nil {
		return
	}
	fullURL := a.baseUrl.ResolveReference(parsed).String()

	tab := a.browser.MustPage(fullURL)
	tab.MustWaitLoad()

	container, _ := tab.Timeout(time.Second * 2).Element(".css-1q94rgb")
	addToCartBtn, err := container.Timeout(time.Second * 2).Element("[data-testid=\"add-to-cart-cta-button\"]")
	if err != nil {
		return
	}

	addToCartBtn.MustClick()
}

func (a Collector) Checkout() {
	cart := a.browser.MustPage("https://store.epicgames.com/en-US/cart")

	cart.MustWaitLoad()

	checkoutContainer, err := cart.Timeout(time.Second * 2).Element(".css-1fkya4n")
	if err != nil {
		return
	}
	checkoutContainer.MustElement("button").MustClick()

	PurchaseContainer := cart.MustElement(".webPurchaseContainer").
		MustElement("iframe").
		MustFrame()

	time.Sleep(time.Second * 2)

	PurchaseContainer.MustElement(".payment-btn.payment-order-confirm__btn.payment-btn--primary").MustClick()
}
