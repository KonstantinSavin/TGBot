package telegram

import (
	"net/url"
	"projects/DAB/pkg/logging"

	"github.com/mymmrac/telego"
)

var logger logging.Logger

func IsURL(update telego.Update) bool {
	u, err := url.Parse(update.Message.Text)

	return err == nil && u.Host != ""
}
