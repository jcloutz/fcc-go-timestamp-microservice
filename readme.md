# FCC Go Timestamp Microservice

[ ![Codeship Status for jcloutz/fcc-go-timestamp-microservice](https://app.codeship.com/projects/ad59a520-bb33-0134-8fb7-56b7c9b632c8/status?branch=master)](https://app.codeship.com/projects/195610)

This app represents a first attempt at a microservice with Go.

* Requirements: https://www.freecodecamp.com/challenges/timestamp-microservice
* Example App: https://timestamp-ms.herokuapp.com/

### Usage

```
// pass a valid unix timestamp
http://example.com/1488153600
```

```
// pass a valid date string (must be this exact format)
http://example.com/February%2027,%202017
```

#### Return Value
```json
{
  "unix": 1488153600,
  "natural": "February 27, 2017"
}
```
