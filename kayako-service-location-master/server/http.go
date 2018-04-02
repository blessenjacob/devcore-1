package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/braintree/manners"
	"github.com/julienschmidt/httprouter"
	"github.com/kayako/service-location/geo"
	log "github.com/sirupsen/logrus"
)

// StatusValidationFailed denotes a validation failure
const StatusValidationFailed = 422

var errValidation = errors.New("Invalid IP address provided")
var errInternal = errors.New("An internal error occurred, please try again after some time")

// supported api version
var contentTypeGeo = "application/vnd.geo.v1+json"
var contentTypeLegacy = "application/xml"

// Warning message that will be send to a consumer who tries
// to use the legacy or unversioned API
var deprecationWarn = `"Unversioned API has been deprecated, use 'application/vnd.geo.v1+json'" "` + time.Now().String() + `"`

// SigShutdown receives the signal to shut down the HTTP server
var SigShutdown chan os.Signal

type httpErr string

func (e httpErr) JSON() []byte {
	msg := struct {
		Message string `json:"error"`
	}{
		Message: string(e),
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Failed to marshal error to JSON. %s", err.Error())
	}

	return b
}

// Mount the HTTP server on given network interface
func Mount(iface string) error {
	ip, port, err := net.SplitHostPort(iface)
	if err != nil {
		return err
	}

	router := httprouter.New()

	// Backward compatible endpoint for Novo framework
	router.GET("/", handleRequest)
	router.GET("/lookup", handleRequest)

	router.GET("/_version", showVersion)

	// Backward compatible endpoint for Swift Framework
	router.POST("/Location/Search/IP", handleRequest)
	router.POST("/index.php", handleRequest)

	srv := manners.NewServer()
	srv.Addr = ip + ":" + port
	srv.Handler = router

	go trapShutdown(srv)

	var crash error
	for {
		log.Infof("Mounting HTTP server on network interface %s", srv.Addr)
		crash = srv.ListenAndServe()
		if crash == nil {
			break
		}

		switch crash.(type) {
		case *net.OpError, *net.AddrError:
			return fmt.Errorf("HTTP server cannot proceed with malfunctioned configuration (%s)", crash.Error())
		}

		log.Errorf("HTTP server crashed with error - %s", crash.Error())
		log.Info("Respawning the HTTP server")
	}

	log.Info("HTTP server has shut down")
	return nil
}

// Handle the HTTP request
func handleRequest(res http.ResponseWriter, req *http.Request, param httprouter.Params) {
	// Check if the requested content type is versioned, otherwise
	// a data format compatible with Novo framework will be returned
	v := req.Header.Get("content-type")

	// Swift makes a POST call with 'ipaddress' parameter, this method
	// will be deprecated in coming release
	if req.Method == "POST" {
		serveXML(res, req, param)
	} else {
		serveJSON(res, req, param, v == contentTypeGeo)
	}
}

// serve JSON content to JSON client or Novo
func serveJSON(res http.ResponseWriter, req *http.Request, param httprouter.Params, versioned bool) {
	requested := req.URL.Query().Get("ip")
	city, err := lookup(requested)
	if err != nil {
		writeErr(res, err)
		return
	}

	cnt := []byte{}
	if versioned {
		cnt, err = city.JSON()
	} else {
		cnt, err = city.ToNovoResource("200", requested).JSON()
		res.Header().Add("Warning", "299 geo.kayako.com "+deprecationWarn)
	}

	if err != nil {
		writeErr(res, err)
		return
	}

	writeBody(res, http.StatusOK, cnt)
}

// Lookup the Geo information about an IP address and
// return city
func lookup(requested string) (*geo.City, error) {
	ip := net.ParseIP(requested)
	if ip == nil {
		return nil, errValidation
	}

	city, err := geo.LookupCity(ip)
	if err != nil {
		return nil, errInternal
	}

	return city, nil
}

// serveXMl serves the content to Swift framework
func serveXML(res http.ResponseWriter, req *http.Request, param httprouter.Params) {
	requested := req.FormValue("ipaddress")

	city, err := lookup(requested)
	if err != nil {
		writeErr(res, err)
		return
	}

	b, err := city.ToNovoResource("200", requested).Data.XML()
	if err != nil {
		writeErr(res, err)
		return
	}

	res.Header().Add("content-type", contentTypeLegacy)
	writeBody(res, http.StatusOK, b)
}

// writeErr writes a string to response stream
func writeErr(r http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err == errValidation {
		status = StatusValidationFailed
	}

	writeBody(r, status, httpErr(err.Error()).JSON())
}

// writeBody writes a slice of bytes to response stream
func writeBody(r http.ResponseWriter, status int, body []byte) {
	if r.Header().Get("content-type") == "" {
		r.Header().Add("content-type", contentTypeGeo)
	}

	r.WriteHeader(status)
	r.Write(body)
}

// showVersion returns the application version
func showVersion(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	vinfo := map[string]string{
		"version":    geo.AppVersion,
		"branch":     geo.BuildBranch,
		"build time": geo.BuildTime,
	}

	v, err := json.Marshal(vinfo)
	if err != nil {
		log.Errorf("failed to marshal version information to JSON, error: %s", err.Error())
		return
	}

	_, err = res.Write(v)
	if err != nil {
		log.Errorf("failed to write version information, error: %s", err.Error())
	}
}

// trapShutdown traps the shutdown signal and gracefully
// shuts down the HTTP server after draining the request pool
func trapShutdown(srv *manners.GracefulServer) {
	SigShutdown = make(chan os.Signal)
	signal.Notify(SigShutdown, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL)

	<-SigShutdown

	log.Info("Shutting down the HTTP server...")
	srv.Close()
}
