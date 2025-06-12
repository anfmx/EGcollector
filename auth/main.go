package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const (
	chromePath    = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	chromeProfile = `./ChromeProfile`
)

type Authenticator struct {
	baseUrl *url.URL
	browser *rod.Browser
}

type AuthenticatorBuilder struct {
	auth Authenticator
}

func NewAuthenticatorBuilder() *AuthenticatorBuilder {
	return &AuthenticatorBuilder{}
}

func (a *AuthenticatorBuilder) WithBaseUrl(baseUrl *url.URL) *AuthenticatorBuilder {
	a.auth.baseUrl = baseUrl
	return a
}

func (a *AuthenticatorBuilder) WithBrowser(browser *rod.Browser) *AuthenticatorBuilder {
	a.auth.browser = browser
	return a
}

func (a *AuthenticatorBuilder) Build() Authenticator {
	return a.auth
}

type AuthenticatorDirector struct {
	authBuilder *AuthenticatorBuilder
}

func NewAuthDirector() *AuthenticatorDirector {
	authBuilder := NewAuthenticatorBuilder()
	return &AuthenticatorDirector{authBuilder: authBuilder}
}

func (ad *AuthenticatorDirector) ChromeBrowser() Authenticator {
	baseUrl, err := url.Parse("https://store.epicgames.com/")
	if err != nil {
		fmt.Println(err)
	}
	websocket := launcher.NewUserMode().Bin(chromePath).UserDataDir(chromeProfile).MustLaunch()

	browser := rod.New().ControlURL(websocket).MustConnect().NoDefaultDevice()

	return ad.authBuilder.WithBaseUrl(baseUrl).WithBrowser(browser).Build()
}

func (a Authenticator) Login() {
	page := a.browser.MustPage(a.baseUrl.String())
	page.MustSetViewport(1920, 1080, 1, false)
	page.MustWaitLoad()

	loginLink := *(page.
		Timeout(time.Second * 3).
		MustElement("egs-navigation").
		MustShadowRoot().
		MustElement(".dropdown__button.secondary-cta").
		MustAttribute("href"))

	parsedLoginLink, err := url.Parse(loginLink)
	if err != nil {
		return
	}
	fullUrl := a.baseUrl.ResolveReference(parsedLoginLink).String()
	a.browser.MustPage(fullUrl)
}

func main() {
	director := NewAuthDirector()
	authenticator := director.ChromeBrowser()
	authenticator.Login()
}
