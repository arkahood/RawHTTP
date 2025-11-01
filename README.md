# RawHTTP
HTTP 1.1 server from scratch

## HTTP Message Structure

According to RFC 7230, HTTP messages follow this structure:

```
start-line CRLF
*( field-line CRLF )
*( field-line CRLF )
...
CRLF
[ message-body ]
```

Where:
- **start-line**: Request line (method, URI, version) or status line
- **field-line**: HTTP headers (key-value pairs) (The RFC uses the term)
- **CRLF**: Carriage return + line feed (`\r\n`)
- **message-body**: Optional request/response body