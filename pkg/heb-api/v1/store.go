package heb

type Store struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PostalCode string `json:"postalCode"`
	State      string `json:"state"`
}

type LocatorResponse struct {
	Stores []LocatorStore `json:"stores"`
}

type LocatorStore struct {
	Distance            float32 `json:"distance"`
	Store               Store   `json:"store"`
	SupportsMedTimeslot bool    `json:"supportsMedTimeslot"`
}
