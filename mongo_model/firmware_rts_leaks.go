package mongo_model

type FirmWareRtsLeaks struct {
	Url         []UrlLeaks         `bson:"url" json:"url"`
	Email       []EmailLeaks       `bson:"email" json:"email"`
	IPV4Public  []IPV4PublicLeaks  `bson:"ipv4_public" json:"ipv4_public"`
	IPV4Private []IPV4PrivateLeaks `bson:"ipv4_private" json:"ipv4_private"`
}

type UrlLeaks struct {
	Info string `bson:"info" json:"info"`
	Path string `bson:"path" json:"path"`
}

type EmailLeaks struct {
	Info string `bson:"info" json:"info"`
	Path string `bson:"path" json:"path"`
}

type IPV4PublicLeaks struct {
	Info string `bson:"info" json:"info"`
	Path string `bson:"path" json:"path"`
}

type IPV4PrivateLeaks struct {
	Info string `bson:"info" json:"info"`
	Path string `bson:"path" json:"path"`
}
