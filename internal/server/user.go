package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func addUserPageHandler(s fiber.Router) {
	s.Get(
		"/", func(c *fiber.Ctx) error {
			headers := c.GetReqHeaders()
			type headerData struct {
				Header string
				Value  string
			}
			var hd []headerData
			for h, vs := range headers {
				hd = append(
					hd, headerData{
						Header: h,
						Value:  strings.Join(vs, ", "),
					},
				)
			}
			return render(
				c, "user", map[string]interface{}{
					"headers":  hd,
					"username": c.Get("X-Forwarded-User"),
				},
			)
		},
	)
}
