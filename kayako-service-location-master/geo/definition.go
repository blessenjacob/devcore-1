package geo

import (
	"encoding/json"
	"encoding/xml"
)

// LangEN is name for language en
const LangEN = "en"

// VCS version information about the app
var (
	AppVersion  string
	BuildBranch string
	BuildTime   string
)

// LocaleNames are names in supported languages
type LocaleNames map[string]string

// Get a name by locale or return fallback
func (n LocaleNames) Get(locale, fallback string) string {
	if v, ok := n[locale]; ok {
		return v
	}

	return fallback
}

// SubDivision in a city
type SubDivision struct {
	GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
	IsoCode   string      `maxminddb:"iso_code" json:"iso"`
	Names     LocaleNames `maxminddb:"names" json:"names"`
}

// SubDivisions is a collection of SubDivision
type SubDivisions []SubDivision

// GetI returns the subdivision at index i.
// Index starts from 0
func (d SubDivisions) GetI(i int) SubDivision {
	if len(d)-1 < i {
		return SubDivision{}
	}

	return d[i]
}

// City contains the geo information about a city
type City struct {
	City struct {
		GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
		Names     LocaleNames `maxminddb:"names" json:"names"`
	} `maxminddb:"city" json:"city"`

	Continent struct {
		GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
		Code      string      `maxminddb:"code" json:"code"`
		Names     LocaleNames `maxminddb:"names" json:"names"`
	} `maxminddb:"continent" json:"continent"`

	Country struct {
		GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
		IsoCode   string      `maxminddb:"iso_code" json:"iso"`
		Names     LocaleNames `maxminddb:"names" json:"names"`
	} `maxminddb:"country" json:"country"`

	Location struct {
		Latitude  float64 `maxminddb:"latitude" json:"latitude"`
		Longitude float64 `maxminddb:"longitude" json:"longitude"`
		MetroCode uint    `maxminddb:"metro_code" json:"metro_code"`
		TimeZone  string  `maxminddb:"time_zone" json:"timezone"`
	} `maxminddb:"location" json:"location"`

	Postal struct {
		Code string `maxminddb:"code" json:"code"`
	} `maxminddb:"postal" json:"postal"`

	RegisteredCountry struct {
		GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
		IsoCode   string      `maxminddb:"iso_code" json:"iso"`
		Names     LocaleNames `maxminddb:"names" json:"names"`
	} `maxminddb:"registered_country" json:"registered_country"`

	RepresentedCountry struct {
		GeoNameID uint        `maxminddb:"geoname_id" json:"-"`
		IsoCode   string      `maxminddb:"iso_code" json:"iso"`
		Type      string      `maxminddb:"type" json:"type"`
		Names     LocaleNames `maxminddb:"names" json:"names"`
	} `maxminddb:"represented_country" json:"represented_country"`

	Subdivisions SubDivisions `maxminddb:"subdivisions" json:"subdivisions"`

	ISP ISP `json:"isp"`

	ConnType ConnectionType `json:"connection_type"`

	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy" json:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider" json:"is_satellite_provider"`
	} `maxminddb:"traits" json:"traits"`
}

// ISP structure corresponds to the data in the GeoIP2 ISP database.
type ISP struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number" json:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization" json:"autonomous_system_organization"`
	ISP                          string `maxminddb:"isp" json:"isp"`
	Organization                 string `maxminddb:"organization" json:"organization"`
}

// ConnectionType structure corresponds to the data in the GeoIP2
// Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type" json:"connection_type"`
}

// JSON converts City into JSON representation
func (c *City) JSON() ([]byte, error) {
	return json.Marshal(c)
}

// ToNovoResource converts the city data into Novo
// compatible data structure
func (c *City) ToNovoResource(status, ip string) *NovoResource {
	r := &NovoResource{}
	r.Status = status
	r.Data.IP = ip
	r.Data.Latitude = c.Location.Latitude
	r.Data.Longitude = c.Location.Longitude
	r.Data.AreaCode = c.Postal.Code
	r.Data.PostalCode = c.Postal.Code
	r.Data.MetroCode = c.Location.MetroCode
	r.Data.City = c.City.Names.Get(LangEN, "")
	r.Data.Region = c.Subdivisions.GetI(0).IsoCode
	r.Data.Country = c.Country.Names.Get(LangEN, "")
	r.Data.CountryCode = c.Country.IsoCode
	r.Data.ISP = c.ISP.ISP
	r.Data.Organization = c.ISP.Organization
	r.Data.NetSpeed = c.ConnType.ConnectionType
	r.Data.TimeZone = c.Location.TimeZone
	r.Data.LocaleCode = c.Country.IsoCode
	r.Data.ContinentCode = c.Continent.Code
	r.Data.ContinentName = c.Continent.Names.Get(LangEN, "")
	r.Data.RegionOneName = c.Subdivisions.GetI(0).Names.Get(LangEN, "")
	r.Data.RegionTwoName = c.Subdivisions.GetI(1).Names.Get(LangEN, "")
	r.Data.RegionTwoCode = c.Subdivisions.GetI(1).IsoCode
	r.Data.IsAnonymousProxy = c.Traits.IsAnonymousProxy
	r.Data.IsSatelliteProvider = c.Traits.IsSatelliteProvider

	return r
}

// NovoResource is the Novo framework compatible
// information data structure for backwards compatibility
type NovoResource struct {
	Status string `json:"status"`
	Data   GeoIP  `json:"data"`
}

// GeoIP contains the data in a Novo Resource
// object.
type GeoIP struct {
	IP                  string  `json:"ip"                    xml:"GeoIPResult>IP"`
	Latitude            float64 `json:"latitude"              xml:"GeoIPResult>City>Latitude"`
	Longitude           float64 `json:"longitude"             xml:"GeoIPResult>City>Longitude"`
	AreaCode            string  `json:"area_code"             xml:"GeoIPResult>City>AreaCode"`
	PostalCode          string  `json:"postal_code"           xml:"GeoIPResult>City>PostalCode"`
	MetroCode           uint    `json:"metro_code"            xml:"GeoIPResult>City>MetroCode"`
	City                string  `json:"city"                  xml:"GeoIPResult>City>Name"`
	Region              string  `json:"region"                xml:"GeoIPResult>City>Region"`
	Country             string  `json:"country"               xml:"GeoIPResult>City>Country"`
	ISP                 string  `json:"isp"                   xml:"GeoIPResult>ISP"`
	Organization        string  `json:"organization"          xml:"GeoIPResult>Organization"`
	NetSpeed            string  `json:"net_speed"             xml:"GeoIPResult>NetSpeed"`
	CountryCode         string  `json:"country_code"          xml:"-"`
	TimeZone            string  `json:"time_zone"             xml:"-"`
	LocaleCode          string  `json:"locale_code"           xml:"-"`
	ContinentCode       string  `json:"continent_code"        xml:"-"`
	ContinentName       string  `json:"continent_name"        xml:"-"`
	RegionOneName       string  `json:"region_1_name"         xml:"-"`
	RegionTwoCode       string  `json:"region_2_code"         xml:"-"`
	RegionTwoName       string  `json:"region_2_name"         xml:"-"`
	IsAnonymousProxy    bool    `json:"is_anonymous_proxy"    xml:"-"`
	IsSatelliteProvider bool    `json:"is_satellite_provider" xml:"-"`
}

// JSON converts NovoResponse into JSON representation
func (r *NovoResource) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// XML converts GeoIP into XML representation
func (d *GeoIP) XML() ([]byte, error) {
	return xml.Marshal(d)
}
