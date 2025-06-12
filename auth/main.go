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

type Authorizator struct {
	baseUrl *url.URL
	browser *rod.Browser
}

type AuthorizatorBuilder struct {
	auth Authorizator
}

func NewAuthorizatorBuilder() *AuthorizatorBuilder {
	return &AuthorizatorBuilder{}
}

func (a *AuthorizatorBuilder) WithBaseUrl(baseUrl *url.URL) *AuthorizatorBuilder {
	a.auth.baseUrl = baseUrl
	return a
}

func (a *AuthorizatorBuilder) WithBrowser(browser *rod.Browser) *AuthorizatorBuilder {
	a.auth.browser = browser
	return a
}

func (a *AuthorizatorBuilder) Build() Authorizator {
	return a.auth
}

type AuthorizatorDirector struct {
	authBuilder *AuthorizatorBuilder
}

func NewAuthDirector() *AuthorizatorDirector {
	authBuilder := NewAuthorizatorBuilder()
	return &AuthorizatorDirector{authBuilder: authBuilder}
}

func (ad *AuthorizatorDirector) ChromeBrowser() Authorizator {
	baseUrl, err := url.Parse("https://store.epicgames.com/")
	if err != nil {
		fmt.Println(err)
	}
	websocket := launcher.NewUserMode().Bin(chromePath).UserDataDir(chromeProfile).MustLaunch()

	browser := rod.New().ControlURL(websocket).MustConnect().NoDefaultDevice()

	return ad.authBuilder.WithBaseUrl(baseUrl).WithBrowser(browser).Build()
}

func (a Authorizator) Login() {
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
