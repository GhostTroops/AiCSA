# bleve mapping UI

A reusable, web-based editor and viewer UI for bleve IndexMapping
JSON based on angular JS, angular-bootstrap and the angular-ui-control.

## Demo

Build the sample webapp...

    go build ./cmd/sample

Run the sample webapp...

    ./sample

Browse to this URL...

    http://localhost:9090/sample.html

## Screenshot

![screenshot](https://raw.githubusercontent.com/blevesearch/bleve-mapping-ui/master/docs/screenshot.png)

## License

Apache License Version 2.0

## For bleve mapping UI developers

### Code generation

There's static bindata resources, which can be regenerated using...

    go generate

### Unit tests

There are some "poor man's" unit tests, which you can run by visiting...

    http://localhost:9090/mapping_test.html

