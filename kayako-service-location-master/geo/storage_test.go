// +build bench

package geo

import (
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
)

var ipList = []net.IP{
	net.ParseIP("122.162.147.193"),
	net.ParseIP("122.162.147.150"),
	net.ParseIP("122.162.147.100"),
	net.ParseIP("122.162.147.200"),
	net.ParseIP("122.162.147.120"),
	net.ParseIP("122.162.147.130"),
}

func TestMain(m *testing.M) {
	dbf := os.Getenv("GEOCITY_DB")
	err := Connect(dbf)
	if err != nil {
		log.Fatalf("Cannot benchmark without a valid database connection, %s", err.Error())
	}

	os.Exit(m.Run())
}

func BenchmarkCityLookup(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := rand.Perm(len(ipList))
			ip := ipList[s[0]]
			_, err := LookupCity(ip)

			if err != nil {
				b.Error(err.Error())
			}
		}
	})
}
