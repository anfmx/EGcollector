package autofarmer

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

type AutoFarmer struct {
	browser    *rod.Browser
	baseUrl    *url.URL
	userData   string
	browserBin string
}

func NewAutoFarmer(browserBin, userData string) *AutoFarmer {
	ws := launcher.NewUserMode().
		Bin(browserBin).
		UserDataDir(userData).
		MustLaunch()

	browser := rod.New().ControlURL(ws).MustConnect().NoDefaultDevice()

	base, _ := url.Parse("https://store.epicgames.com")

	return &AutoFarmer{
		browser:    browser,
		baseUrl:    base,
		userData:   userData,
		browserBin: browserBin,
	}
}

type AutoFarmerBuilder struct {
	autoFarmer AutoFarmer
}

func (afb *AutoFarmerBuilder) WithBrowser(browser *rod.Browser) *AutoFarmerBuilder {
	afb.autoFarmer.browser = browser
	return afb
}
func (afb *AutoFarmerBuilder) WithBaseUrl(baseUrl *url.URL) *AutoFarmerBuilder {
	afb.autoFarmer.baseUrl = baseUrl
	return afb
}
func (afb *AutoFarmerBuilder) WithUserData(userData string) *AutoFarmerBuilder {
	afb.autoFarmer.userData = userData
	return afb
}
func (afb *AutoFarmerBuilder) WithBrowserBin(browserBin string) *AutoFarmerBuilder {
	afb.autoFarmer.browserBin = browserBin
	return afb
}
func (afb *AutoFarmerBuilder) Build() AutoFarmer {
	return afb.autoFarmer
}

type AutoFarmerBuilderDirector struct {
	autoFarmerBuilder *AutoFarmerBuilder
}

func (ad *AutoFarmerBuilderDirector) NewChromeFarmer() AutoFarmer {
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

func NewAutoFarmDirector() *AutoFarmerBuilderDirector {
	return &AutoFarmerBuilderDirector{autoFarmerBuilder: &AutoFarmerBuilder{}}
}

func (a AutoFarmer) GetGames() []*rod.Element {
	page := a.browser.MustPage("https://store.epicgames.com/en-US/free-games")

	page.MustSetViewport(1920, 1080, 1, false)
	page.MustWaitLoad()

	container := page.MustElement(".css-2u323")
	return container.MustElements(".css-g3jcms")
}

func (a AutoFarmer) AddToCart(href string) {
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

func (a AutoFarmer) Checkout() {
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
