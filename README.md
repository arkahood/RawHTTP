# RawHTTP
A from-scratch HTTP/1.1 server implementation in Go, built directly on TCP sockets.

## Project Overview

RawHTTP is an HTTP/1.1 server implementation that demonstrates low-level network programming concepts by building a complete HTTP server from first principles. Unlike frameworks that abstract away protocol details, this implementation provides direct insight into:

- Raw TCP socket handling
- HTTP message parsing and generation
- State machine-based request processing
- Buffer management and streaming

explore how web servers work at the protocol level.

## HTTP Protocol Deep Dive

### HTTP Message Structure

According to RFC 7230, HTTP messages follow this precise structure:

```
start-line CRLF
*( field-line CRLF )
CRLF
[ message-body ]
```

Where:
- **start-line**: Request line (method, URI, version) or status line
- **field-line**: HTTP headers (key-value pairs)
- **CRLF**: Carriage return + line feed (`\r\n`) - critical for protocol compliance
- **message-body**: Optional request/response body

### Request Message Anatomy

```
GET /api/users HTTP/1.1\r\n          ← Request Line
Host: example.com\r\n                ← Header Field
User-Agent: RawHTTP/1.0\r\n          ← Header Field
Accept: application/json\r\n          ← Header Field
\r\n                                  ← Header/Body Separator
[optional message body]               ← Body (for POST, PUT, etc.)
```

#### Request Line Components

1. **HTTP Method**: Defines the action to be performed
   - `GET`: Retrieve data
   - `POST`: Submit data
   - `PUT`: Update/create resource
   - `DELETE`: Remove resource

2. **Request Target**: Identifies the resource
   - **origin-form**: `/path/to/resource?query=value`
   - **absolute-form**: `http://example.com/path`
   - **authority-form**: `example.com:80` (CONNECT only)
   - **asterisk-form**: `*` (OPTIONS only)

3. **HTTP Version**: Currently `HTTP/1.1`

### Response Message Anatomy

```
HTTP/1.1 200 OK\r\n                  ← Status Line
Content-Type: text/html\r\n          ← Header Field
Content-Length: 1234\r\n             ← Header Field
Connection: close\r\n                ← Header Field
\r\n                                  ← Header/Body Separator
<html>...</html>                      ← Body
```


### Header Field Theory

#### Field Name Constraints
- **tchar**: `!#$%&'*+-.^_`|~` plus ALPHA and DIGIT
- Case-insensitive (normalized to lowercase in our implementation)
- No whitespace between field-name and colon

#### Field Value Processing
- Leading/trailing whitespace is removed
- Multiple values can be comma-separated

#### Critical Headers
- **Host**: Required in HTTP/1.1, enables virtual hosting
- **Content-Length**: **Byte count of message body**

## Running the Server

### Quick Start
```bash
# Clone the repository
git clone https://github.com/arkahood/RawHTTP.git
cd RawHTTP

# Run the HTTP server
go run ./cmd/httpserver/main.go

# Server starts on port 8080
# Visit http://localhost:8080 in your browser
```

### Testing with curl
```bash
# Basic GET request
curl -v http://localhost:8080/
```

### Development Commands
```bash
# Run tests
go test ./...

# Build binary
go build -o httpserver cmd/httpserver/main.go

# View test coverage
go test -cover ./...
```