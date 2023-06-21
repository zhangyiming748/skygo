package mongo_model

type FirmWareRtsCrt struct {
	PublicKey   []PublicKeyInfo   `bson:"public_key" json:"public_key"`
	Certificate []CertificateInfo `bson:"certificate" json:"certificate"`
	PrivateKey  []PrivateKeyInfo  `bson:"private_key" json:"private_key"`
}

type PublicKeyInfo struct {
	PublicInfo InfoInfo `bson:"info" json:"info"`
	Path       string   `bson:"path" json:"path"`
	Content    string   `bson:"content" json:"content"`
}

type CertificateInfo struct {
	Info    CertInfo
	Path    string
	Content string
}

type PrivateKeyInfo struct {
	Info    string `bson:"info" json:"info"`
	Path    string `bson:"path" json:"path"`
	Content string `bson:"content" json:"content"`
}

type InfoInfo struct {
	Json      JsonInfo `bson:"json" json:"json"`
	Printable string   `bson:"printable" json:"printable"`
}

type JsonInfo struct {
	Length    int
	Modulus   string
	Exponent  string
	Algorithm string
}

type CertInfo struct {
	//Json CertJson `bson:"json" json:"json"`
	Printable string `bson:"printable" json:"printable`
}

//type CertJson struct {
//	Issuer string
//	Subject string
//	Version string
//	Validity
//		Not After
//		Not Before
//	Extensions string
//	Public Key
//		Length int
//		Modulus string
//	Algorithm ID string
//	Serial number string
//	Certificate Signature
//		Algorithm
//		Signature
//}
