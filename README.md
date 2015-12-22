# apnsapi

**require go1.6**

Simple apns api (http/2) client fo golang.

## USAGE

It provides only client, so should setup http/2 client, apns payload yourself.

```go
token := "..."
client := apnsapi.NewClient(apnsapi.ProductionServer, &http.Client{...})
header := &apnsapi.Header{ApnsTopic: "..."}
payload := `{ "aps" : { "alert" : "hi" } }`
if _, err = client.Do(token, header, []byte{payload}); err != nil {
    ...
}
```

```_example/main.go``` is sample implantation and used like this.

```
$ go run _example/main.go --p12 /path/to/file.p12 --topic $topic --token $token
```

## SEE ALSO

* https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/Introduction.html
* https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/APNsProviderAPI.html#//apple_ref/doc/uid/TP40008194-CH101-SW1

## LICENSE

MIT
