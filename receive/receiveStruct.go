package receive

type Data struct {
	Message string `bson:"message" json:"message"`
	Date    string `bson:"date" json:"date"`
}
