This is a GO implementation of the Ike Wai login API.

To use this you must configure the base URL for the Agave tenant you are accessing and then add the consumerKey and consumerSecret for the client you wish to use.  This can be found in auth.go.  You can also modify the port that the API will be served on- default is 8080.

Once modified build the application:
```go build auth.go```

Then execute the application.  It will be availalbe on localhost:8080 (by default).

The API is accessed at localhost:8080/token and is expecting a BasiC Auth username and password to login to the Agave tenant.