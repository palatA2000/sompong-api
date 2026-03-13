package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type lineProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
}

// LIFFAuth verifies LINE access token from Authorization: Bearer <token>
func LIFFAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) < 8 || token[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid authorization header",
			})
		}
		accessToken := token[7:]

		profile, err := verifyLINEToken(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid LINE access token",
			})
		}

		c.Locals("lineUserID", profile.UserID)
		c.Locals("lineDisplayName", profile.DisplayName)
		c.Locals("linePictureURL", profile.PictureURL)

		return c.Next()
	}
}

func verifyLINEToken(accessToken string) (*lineProfile, error) {
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fiber.ErrUnauthorized
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile lineProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
