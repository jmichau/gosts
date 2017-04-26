# hsts provide middleware for support Strict-Transport-Security header.
## Install:
`go get github.com/mj420/go-hsts`

## Usage:
this is part of my application that explain how to use this middleware
* main.go
```
...
// Load reads the configuration file
func Load(configFile string) (*Info, error) {
	conf := new(Info)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}

func main() {
    	// Load the configuration file
    	conf, err := config.Load("config.yml")
    	if err != nil {
    		log.Fatalln(err)
    	}

        // Start the HTTP and HTTPS listeners
        err = server.Run(route.Load(), config.Server)
        if err != nil {
            log.Fatalln(err)
        }
}
...
```

* server.go
```
...
// Run starts the HTTP and/or HTTPS listener and wait for shutdown signal.
// Finally gracefully stop listeners
func (i *Info) Run(handlers http.Handler) error {
	var (
		err               error
		srvHTTP, srvHTTPS http.Server
	)

	// Start server
	if i.UseHTTPS {
		go func() {
			srvHTTPS, err = startHTTPS(handlers, i)
			if err != nil {
				log.Fatalln(err)
			}
		}()

		// Redirect HTTP to HTTPS
		go func() {
			srvHTTP, err = startHTTP(http.HandlerFunc(redirectToHTTPS), i)
			if err != nil {
				log.Fatalln(err)
			}
		}()
	} else {
		go func() {
			srvHTTP, err = startHTTP(handlers, i)
			if err != nil {
				log.Fatalln(err)
			}
		}()
	}

	// Graceful shutdown server
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Graceful stopping HTTP listeners:", sig)
		done <- true
	}()

	<-done

	if i.UseHTTPS {
		err = srvHTTPS.Shutdown(nil)
		if err != nil {
			return err
		}
	}
	err = srvHTTP.Shutdown(nil)
	if err != nil {
		return err
	}
	return nil
}
...

```

* config.yml:
```
...

Server:
  Hostname: localhost

  # This property is required even if HTTPS is used, cause of redirect from HTTP to HTTPS.
  # 80 is default HTTP port, so don't change it on production.
  HTTPPort: 80

  # If this field is set to true app will be available via HTTPS
  # and HTTP listener only is used to redirect to HTTPS.
  UseHTTPS: true

  # -----------------------------------------------------------------------------
  # Lines above are only processed if UseHTTPS is set to true, so if UseHTTPS
  # property is set to false then lines above does not have to be valid.
  # -----------------------------------------------------------------------------

  # 443 is default HTTPS port, so don't change it on production.
  HTTPSPort: 443
  CertFile: tls/server.crt
  KeyFile: tls/server.key

  # HTTP Strict Transport Security (HSTS) header. RFC 6797.
  HSTS:
    # MaxAge sets the duration (number of seconds) that the HSTS is valid for.
    # This value can't be less than zero. This property is required
    # cause of Expires that use MaxAge as a fallback if duration between now and Expires is less than zero.
    MaxAge: 86400

    # Expires sets the date and time after which the header will not be valid.
    # This property is automatically converted to MaxAge every request (duration between now and Expires).
    # If this property is set the MaxAge property is not respected,
    # but if duration between now and Expires is less than zero then MaxAge will be respected automatically.
    # This is a good strategy if you do not know whether you will renew your SSL certificate in the future.
    # Use RFC 3339 format.
    Expires: 2017-04-18T00:02:00+02:00

    # IncludeSubDomains specifying that this HSTS Policy also applies to any hosts
    # whose domain names are subdomains of the Known HSTS Host's domain name.
    IncludeSubDomains: true

    # SendPreloadDirective sets whether the preload directive should be set. The directive allows browsers to
    # confirm that the site should be added to a preload list. (see https://hstspreload.appspot.com/)
    SendPreloadDirective: true
...
```

## TODO:
- [ ] add json config support
- [ ] rewrite readme guide how to use this package (I know that is maybe not simple to understand but see also source code that can better explain what this package do)
