# JWKS Server

A simple service that queries a JWKS uri endpoint and filters JWK entries that have no alg property.

Usage:

```
Usage of jwkssvr:
  -dry-run
    	print the settings and exit
  -issuer string
    	issuer, when a valid URL this will be used to discover the jwksUri [ISSUER]
    	If discovery fails it uses -jwks-uri to retrieve the JWKS.
  -jwks-uri string
    	remote jwks uri [JWKS_URI]
    	This flag is ignored when the -issuer is a url and can be discovered.
  -log-format string
    	log format [ text | json ] [LOG_FORMAT]
  -log-level string
    	log level [info | warn | error | debug], defaults to info [LOG_LEVEL]
  -port string
    	port number to listen on, listens on all interfaces, defaults to 8080 [PORT]
  -version
    	print the version and exit
```