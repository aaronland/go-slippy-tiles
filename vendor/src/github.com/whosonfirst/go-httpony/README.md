# go-httpony

Utility functions for HTTP ponies written in Go.

## Usage

### TLS

```
import (
	"github.com/whosonfirst/go-httpony/tls"	
	"net/http"
)

// Ensures that httpony/certificates exists in your operating
// system's temporary directory and that its permissions are
// 0700. You do _not_ need to use this if you have your own
// root directory for certificates.

root, err := tls.EnsureTLSRoot()

if err != nil {
	panic(err)
}

// These are self-signed certificates so your browser _will_
// complain about them. All the usual caveats apply.

cert, key, err := tls.GenerateTLSCert(*host, root)
	
if err != nil {
	panic(err)
}

http.ListenAndServeTLS("localhost:443", cert, key, nil)
```

The details of setting up application specific HTTP handlers is left as an exercise to the reader.
