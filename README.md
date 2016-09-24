[![wercker status](https://app.wercker.com/status/c1941900b79b7ec8d1a01c7f023ce11b/s/master "wercker status")](https://app.wercker.com/project/bykey/c1941900b79b7ec8d1a01c7f023ce11b)

# apnsapi

Simple apns api (http/2) client for golang.

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

### APNs Provider Certificates

```
$ go run _example/main.go --p12 /path/to/file.p12 --topic $topic --token $token
```

### Provider Authentication Tokens

```
$ go run _example/main.go --key /path/to/file.p8 --kid $keyIdentifier --teamId $teamID --topic $topic --token $token
```

## SEE ALSO

* https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/Introduction.html
* https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/APNsProviderAPI.html
* https://developer.apple.com/videos/play/wwdc2016/724/

## LICENSE

MIT
