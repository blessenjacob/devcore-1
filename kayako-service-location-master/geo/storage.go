package geo

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"
)

// Names of the databases to be used
const (
	DBCity       = "GeoIP2-City"
	DBISP        = "GeoIP2-ISP"
	DBConnection = "GeoIP2-Connection-Type"
)

// connections are the connections to different databases
var connections map[string]*maxminddb.Reader

// dataRetriever retrieves data from database
type dataRetriever func(net.IP, interface{}) error

// dataRetrievers is a map from database type to lookup method
var dataRetrievers = map[string]dataRetriever{
	DBCity:       voidLookup,
	DBISP:        voidLookup,
	DBConnection: voidLookup,
}

// Connect connects to maxmind db
func Connect(filenames ...string) error {
	connections = make(map[string]*maxminddb.Reader)

	for _, n := range filenames {
		log.Infof("Connecting to database %s...", n)
		c, err := maxminddb.Open(n)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Info("Connection successful, verifying database integrity...")
		if err = c.Verify(); err != nil {
			log.Errorf("Database %s is corrupted, skipping connection", n)
			continue
		}
		log.Info("Database integrity verified")

		connections[c.Metadata.DatabaseType] = c
		dataRetrievers[c.Metadata.DatabaseType] = c.Lookup
		log.Infof("%s database is ready", c.Metadata.DatabaseType)
	}

	if len(connections) < 1 {
		return errors.New("No database connection can be established, cannot continue without a database")
	}

	return nil
}

// Disconnect closes the connection
func Disconnect() {
	for n, c := range connections {
		err := c.Close()
		if err != nil {
			log.Errorf("Failed to close the connection %s", n)
		} else {
			log.Infof("%s connection closed", n)
		}
	}
}

// lookupError lists the errors that happens during lookup
type lookupError struct {
	City error
	ISP  error
	Conn error
}

// Check is there is any error
func (l lookupError) hasErr() bool {
	return l.City != nil || l.ISP != nil || l.Conn != nil
}

// Get the error message
func (l lookupError) Err() error {
	m := ""
	if l.City != nil {
		m += "City: " + l.City.Error() + "\n"
	}
	if l.ISP != nil {
		m += "ISP: " + l.ISP.Error() + "\n"
	}
	if l.Conn != nil {
		m += "Connection: " + l.Conn.Error() + "\n"
	}

	if m != "" {
		return errors.New(strings.TrimRight(m, "\n"))
	}

	return nil
}

var wg = sync.WaitGroup{}

// LookupCity finds the city information for an IP address
func LookupCity(ip net.IP) (*City, error) {
	var (
		c   = City{}
		err = lookupError{}
	)

	wg.Add(3)

	go asyncLookup(dataRetrievers[DBCity], ip, &c, err.City)
	go asyncLookup(dataRetrievers[DBISP], ip, &c.ISP, err.ISP)
	go asyncLookup(dataRetrievers[DBConnection], ip, &c.ConnType, err.Conn)

	wg.Wait()

	return &c, err.Err()
}

// voidLookup is used when no database is available for lookup
func voidLookup(net.IP, interface{}) error {
	return nil
}

// asyncLookup performs the lookup asynchronously and
// returns the errors on error channel
func asyncLookup(fn dataRetriever, ip net.IP, v interface{}, err error) {
	err = fn(ip, v)
	wg.Done()
}
