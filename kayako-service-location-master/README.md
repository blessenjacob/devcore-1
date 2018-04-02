# Location
Location App for IP Intelligence Discovery

# Prerequisites
### Tools
Please ensure you have following tools installed on your machine before starting with installation:

 1. Golang >= v1.5.1
 1. GNU make
 1. Docker >= 1.8 (optional, only required if you want to build a docker container)
 1. curl
 1. tar

### GeoIP Databases
Location app uses MaxMind GeoIP databases under the hood, it can query from a GeoIP2 or GeoLite
[binary database][]. Databases are not required to build the app but the HTTP API will not work
without at least one database in place. It is recommended that you use following three databases but
only at least one of them is required.

 1. GeoIP2-City
 1. GeoIP2-ISP
 1. GeoIP2-Connection-Type

If you have a MaxMind premium subscription, you can automatically fetch the database files by
executing `make dbfiles MML=<Maxmind license key>`. If you want to use the GeoLite databases then
please run `make freedbfiles`.

[binary database]: http://maxmind.github.io/MaxMind-DB/

### Building
Building the location app is fairly straightforward, just run `make build` and a binary will be
produced at `build/location`. You can run `make install` to move the binary to
`/usr/local/bin/location` to make it available in your path.

### Building the container
If you want to build the docker container first make sure that you have the GeoIP databases in
place. The container will not fetch them for you automatically because there is no way to securely
put licensing information on the container.  
Once the databases are in place you can run `docker build -t vnd/geo:tag Dockerfile` to build a
container named `vnd/geo` and tag it `tag`.

### Usage
At the moment the app supports only one command

 1. `serve`

The serve command starts an HTTP web server to let a client lookup the GeoIP information via API.
Following flags are supported by `serve` command:

 - `--http-addr`: This specifies the network interface and port to be used by the HTTP server. It
   must be a string in host:port format where host can be omitted.
 - `--db`: Full path to MaxMind database file. This flag can be used multiple times to specify
   multiple databases.
 - `--log-level`: Minimum level of logs to write to destination stream. By default only logs with
   level `error` will be written. Set this to `debug` to increase the log verbosity.
