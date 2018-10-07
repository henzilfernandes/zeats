package types

type Product struct {
	Id string `json:"productId,omitempty"`
	Name string `json:"name,omitempty"`
	ImageClosed string `json:"image_closed,omitempty"`
	ImageOpen string `json:"image_open,omitempty"`
	Description string `json:"description,omitempty"`
	Story string `json:"story,omitempty"`
	SourcingValues []string `json:"sourcing_values,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	AllergyInfo string `json:"allergy_info,omitempty"`
	DietaryCertifications string `json:"dietary_certifications,omitempty"`
}

