package midtrans

import (
	"strings"
)

func ResolveURL(order_str string) string {
	prefix := strings.SplitN(order_str, "-", 2)[0]

	for _, item := range WebhookConfig.URLs {
		if item.Code == prefix {
			return item.URL
		}
	}

	return ""
}
