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

func NewCollector(browserBin, userData string) *Collector {
	ws := launcher.NewUserMode().
		Bin(browserBin).
		UserDataDir(userData).
		MustLaunch()

	browser := rod.New().ControlURL(ws).MustConnect().NoDefaultDevice()

	base, _ := url.Parse("https://store.epicgames.com")

	return &Collector{
		browser:    browser,
		baseUrl:    base,
		userData:   userData,
		browserBin: browserBin,
	}
}

type CollectorBuilder struct {
	autoFarmer Collector
}

func (afb *CollectorBuilder) WithBrowser(browser *rod.Browser) *CollectorBuilder {
	afb.autoFarmer.browser = browser
	return afb
}
func (afb *CollectorBuilder) WithBaseUrl(baseUrl *url.URL) *CollectorBuilder {
	afb.autoFarmer.baseUrl = baseUrl
	return afb
}
func (afb *CollectorBuilder) WithUserData(userData string) *CollectorBuilder {
	afb.autoFarmer.userData = userData
	return afb
}
func (afb *CollectorBuilder) WithBrowserBin(browserBin string) *CollectorBuilder {
	afb.autoFarmer.browserBin = browserBin
	return afb
}
func (afb *CollectorBuilder) Build() Collector {
	return afb.autoFarmer
}

type CollectorBuilderDirector struct {
	autoFarmerBuilder *CollectorBuilder
}

func (ad *CollectorBuilderDirector) NewChromeFarmer() Collector {
	ws := launcher.NewUserMode().
		Bin(chromeBin).
		UserDataDir(chromeProfile).
		MustLaunch()

	browser := rod.New().ControlURL(ws).MustConnect().NoDefaultDevice()

	base, _ := url.Parse("https://store.epicgames.com")

	return ad.autoFarmerBuilder.
		WithBaseUrl(base).
		WithBrowser(browser).
		WithBrowserBin(chromeBin).
		WithUserData(chromeProfile).
		Build()
}

func NewAutoFarmDirector() *CollectorBuilderDirector {
	return &CollectorBuilderDirector{autoFarmerBuilder: &CollectorBuilder{}}
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

	PurchaseContainer.MustElement(".payment-btn.payment-order-confirm__btn.payment-btn--primary").MustClick()
}
