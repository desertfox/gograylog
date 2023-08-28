# GoGraylog  v1.6.3

This package is a thin http client designed to facilitate simple graylog searches

## Install

``` bash
$ go get https://github.com/desertfox/gograylog@v1.6.3
```

## Usage

``` go

    import (
        "github.com/desertfox/gograylog"
    )

	client :=  gograylog.Client{
		Host:  graylogHostUrl,
		Session: &Session{},
		HttpClient: &http.Client{},
	}

    #Session creation
    err := client.Login("username", "password")
    ...

    #Session required request
    byteJSON, err := client.Streams()
    ...

    query := gograylog.Query{
        StreamID:    "Stream id",
        QueryString: "error OR warn",
        Frequency:   "3600",
        Fields:      "message,timestamp",
        Limit:       10000,
	}

    #Session required request
    byteCSV, err := client.Search(query)
    ...

```

## Change log

Please see [CHANGELOG](CHANGELOG.md) for more information on what has changed recently.

## Testing

``` bash
$ go test ./...
```

## Credits

- [DesertFox][link-author]