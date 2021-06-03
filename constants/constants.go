package constants

import "time"

const UserImages = "/static/images/user_images/"

const LinkPath = "/static/images/site_images/social_icons/SVG/"

var LinksImagesPairs = map[string]string{
	"vk.com":        "vk.svg",
	"facebook.com":  "facebook.svg",
	"habr.com":      "habr.svg",
	"github.com":    "github.svg",
	"instagram.com": "instagram.svg",
	"linkedin.com":  "linkedin.svg",
	"ok.ru":         "ok.svg",
	"pinterest.com": "pinterest.svg",
	"reddit.com":    "reddit.svg",
	"snapchat.com":  "snapchat.svg",
	"tumblr.com":    "tumblr.svg",
	"twitter.com":   "twitter.svg",
	"youtube.com":   "youtube.svg",
	"t.me":          "telegram.svg",
	"other":         "other.svg",
}

const JWTExpireTime = time.Second * 10
const RefreshTokenExpireTime = time.Second * 20
