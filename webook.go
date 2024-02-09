package planetscale

type (
	ClerkPayload[T any] struct {
		Data   T      `json:"data"`
		Object string `json:"object"`
		Type   string `json:"type"`
	}

	ClerkUserPayload struct {
		EmailAddresses []struct {
			EmailAddress string `json:"email_address"`
			Id           string `json:"id"`
			Object       string `json:"object"`
			Verification struct {
				Status string `json:"status"`
			} `json:"verification"`
		} `json:"email_addresses"`
		FirstName       string `json:"first_name"`
		Id              string `json:"id"`
		LastName        string `json:"last_name"`
		Object          string `json:"object"`
		ProfileImageURL string `json:"profile_image_url"`
	}
)
