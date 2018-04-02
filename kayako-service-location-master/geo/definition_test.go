package geo

import "testing"

func TestLocaleNames(t *testing.T) {
	tests := []struct {
		names LocaleNames
		key   string
		value string
	}{
		{
			names: LocaleNames{"Key": "Value"},
			key:   "Key",
			value: "Value",
		},
		{
			names: LocaleNames{"Key": "Value"},
			key:   "NonExisting",
			value: "",
		},
	}

	for _, test := range tests {
		v := test.names.Get(test.key, "")
		if v != test.value {
			t.Errorf("Expected value for key %s to be %s, found %s", test.key, test.value, v)
		}
	}
}

func TestResourceData(t *testing.T) {
	d := GeoIP{
		IP:                  "12.34.56.78",
		Latitude:            12.34,
		Longitude:           12.34,
		AreaCode:            "12345",
		PostalCode:          "12345",
		MetroCode:           1,
		City:                "Gurgaon",
		Region:              "DLF",
		Country:             "India",
		ISP:                 "Airtel",
		Organization:        "Kayako",
		NetSpeed:            "Broadband",
		CountryCode:         "IN",
		TimeZone:            "IST",
		LocaleCode:          "HI",
		ContinentCode:       "ASIA",
		ContinentName:       "Asia",
		RegionOneName:       "India",
		RegionTwoCode:       "IN",
		RegionTwoName:       "India",
		IsAnonymousProxy:    false,
		IsSatelliteProvider: false,
	}

	c, err := d.XML()
	if err != nil {
		t.Errorf("Failed to encode resource data to XML. %s", err.Error())
	}

	ex := "<GeoIP><GeoIPResult><IP>12.34.56.78</IP><City><Latitude>12.34</Latitude><Longitude>12.34</Longitude><AreaCode>12345</AreaCode><PostalCode>12345</PostalCode><MetroCode>1</MetroCode><Name>Gurgaon</Name><Region>DLF</Region><Country>India</Country></City><ISP>Airtel</ISP><Organization>Kayako</Organization><NetSpeed>Broadband</NetSpeed></GeoIPResult></GeoIP>"

	if string(c[:]) != ex {
		t.Errorf("Expected XML content to be \n%s\n, have \n%s\n", ex, c)
	}
}
