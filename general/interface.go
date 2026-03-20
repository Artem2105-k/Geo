package general

type GeoProvider interface {
	AddressSearch(input string) ([]*Address, error)
	GeoCode(lat, lng string) ([]*Address, error)
}
type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}
