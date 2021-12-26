# Nap

Nap is a file-based framework for automating the execution of config-driven HTTP requests and scripts.

# Installation Options

## Using go get

```bash
$ go install nap@latest
```

## Building the Source

```bash
$ git clone https://nap.git
$ cd nap
$ go install nap
```

# Create a Request

```bash
$ nap create request my-request.yml
Created new request stub: my-request.yml
```

Creates a file like the following:

```yml
name: my-request
path: https://cat-fact.herokuapp.com/facts/
verb: GET
type: request
body: ""
headers: {}
prerequestscript: "" # not yet supported
postrequestscript: "" # not yet supported
```

# Run a Request

```bash
$ nap run my-request.yml
- my-request.yml: 200 OK
```

## Create a Configuration File

```bash
$ nap create config my-config.yml
Created new configuration file stub: my-config.yml
```

Creates an empty yml file.

## Variables

my-request.yml:
```yml
name: my-request
path: ${baseurl}/facts/
verb: GET
type: request
body: ""
headers: {}
prerequestscript: ""
postrequestscript: ""
```

my-config.yml:
```yml
baseurl: https://cat-fact.herokuapp.com
```

Executing:

```bash
$ nap run my-request.yml -c my-config.yml
- my-request.yml: 200 OK
```

## Verbose Mode

```bash
$ nap run my-request.yml -v
Target File Name: my-request.yml
Config File Name: 
Verbose Mode: true
- my-request.yml:

Running: my-request
Path: https://cat-fact.herokuapp.com/facts/
Verb: GET
Response Status: 200 OK (Content Length: 1859 bytes)
[
    {
        "status": {
            "verified": true,
            "sentCount": 1
        },
        "_id": "58e008800aac31001185ed07",
        "user": "58e007480aac31001185ecef",
        "text": "Wikipedia has a recording of a cat meowing, because why not?", 
        "__v": 0,
        "source": "user",
        "updatedAt": "2020-08-23T20:20:01.611Z",
        "type": "cat",
        "createdAt": "2018-03-06T21:20:03.505Z",
        "deleted": false,
        "used": false
    },
    ...
]
```