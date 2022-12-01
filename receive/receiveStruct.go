package receive

type Data struct {
	CollectionName string `bson:"collectionName" json:"collectionName"`
	Message        string `bson:"message" json:"message"`
	Date           string `bson:"date" json:"date"`
}
