package tab

const SuspectURLRegex = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`

type PredictableInputValue struct {
	Keyword []string
	Value   string
}

var (
	PredictableInputValues = map[string]PredictableInputValue{
		"username": {
			Keyword: []string{"username", "user", "name", "account", "login", "mail", "email", "yonghuming", "yonghu"},
			Value:   "test123",
		},
		"password": {
			Keyword: []string{"password", "pwd", "passwd", "pass", "mima"},
			Value:   "Test2021.",
		},
		"tel": {
			Keyword: []string{"tel", "phone", "shouji", "shoujihaoma", "mobile", "phonenumber"},
			Value:   "15812345678",
		},
		"captcha": {
			Keyword: []string{"code", "verifycode", "verify", "captcha", "yanzhengma"},
			Value:   "123q",
		},
		"email": {
			Keyword: []string{"mail", "email"},
			Value:   "test@test.com",
		},
		"qq": {
			Keyword: []string{"qq", "weixin", "tencent", "wechat"},
			Value:   "123456789",
		},
		"idCard": {
			Keyword: []string{"idcard", "card", "shenfenzheng"},
			Value:   "632323190605264827",
		},
		"url": {
			Keyword: []string{"url", "link", "href", "site", "blog", "web", "website"},
			Value:   "https://github.com",
		},
		"number": {
			Keyword: []string{"age", "count", "num", "number"},
			Value:   "18",
		},
	}
)
