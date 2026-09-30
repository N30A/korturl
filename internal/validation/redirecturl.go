package validation

import "net/url"

func RedirectURL(redirectURL string) bool {
	u, err := url.ParseRequestURI(redirectURL)
	if err != nil {
		return false
	}

	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
