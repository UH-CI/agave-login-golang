This is a GO implementation of the Ike Wai login API.

This is still under heavy development - please contact us if you plan to use this.

To use this you must configure the base URL for the Agave tenant you are accessing and then add the consumerKey and consumerSecret for the client you wish to use.  This can be found in auth.go.  You can also modify the port that the API will be served on- default is 8080.

A server.crt and server.key file is expected with the auth.go file to support SSL.

Once modified build the application:
```go build auth.go```

Then execute the application.  It will be availalbe on localhost:8080 (by default).

The API is accessed at localhost:8080/token and is expecting a Basic Auth username and password to login to the Agave tenant.