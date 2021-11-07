package wanted

type Response struct {
	Data []Data
}

type Data struct {
	ID       int
	Address  Address
	Company  Company
	LogoImg  LogoImage
	Position string
}

type Address struct {
	Country  string
	Location string
}

type Company struct {
	IndustryName string
	Name         string
}

type LogoImage struct {
	Origin string
}
