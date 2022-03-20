package main

type Record struct {
	ID           int    `json:"id"`   // auto increment
	Name         string `json:"name"` // User name from int
	Morning      string `json:"morning"`
	PreBreakfast string `json:"pre_breakfast"`
	Breakfast    string `json:"breakfast"`
	Noon         string `json:"noon"`
	Lunch        string `json:"lunch"`
	Evening      string `json:"evening"`
	Dinner       string `json:"dinner"`
	PostDinner   string `json:"post_dinner"`
	Night        string `json:"night"`
}

type LabelTime struct {
	Label string `json:"label"`
	Time  string `json:"time"`
}

type DietResponse struct {
	Data   []Record             `json:"data"`
	Header map[string]LabelTime `json:"header"`
}

var defaultHeader = map[string]LabelTime{
	"morning":       {Label: "Morning", Time: "7:30 AM"},
	"pre_breakfast": {Label: "Pre Breakfast", Time: "8:40 AM"},
	"breakfast":     {Label: "Breakfast", Time: "9:00 AM"},
	"noon":          {Label: "Noon", Time: "12:00 PM"},
	"lunch":         {Label: "Lunch", Time: "2:00 PM"},
	"evening":       {Label: "Evening", Time: "5:00 PM"},
	"dinner":        {Label: "Dinner", Time: "7:30 PM"},
	"post_dinner":   {Label: "PostDinner", Time: "9:00 PM"},
	"night":         {Label: "Night", Time: "10:00 PM"},
}
