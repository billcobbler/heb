package heb

type Timeslot struct {
	ID               string `json:"id"`
	Date             string `json:"date"`
	AllowsAlcohol    bool   `json:"allowAlcohol"`
	StoreID          string `json:"storeId"`
	FullfillmentType string `json:"fulfillmentType"`
	Capacity         int    `json:"capacity"`
	DayOfWeek        int    `json:"dayOfWeek"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
}

type TimeslotResponse struct {
	Store Store          `json:"pickupStore"`
	Items []TimeslotItem `json:"items"`
}

type TimeslotItem struct {
	Date     string   `json:"date"`
	Timeslot Timeslot `json:"timeslot"`
}
