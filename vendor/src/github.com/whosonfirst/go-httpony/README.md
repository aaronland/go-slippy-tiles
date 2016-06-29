# go-httpony

Utility functions for HTTP ponies written in Go.

## Usage

### Crypto

```
package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-httpony/crypto"
)

func main() {

	var key = flag.String("key", "jwPsjM9rfZl73Pt0XURf0t9u8h5ZOpNT", "The key to encrypt and decrypt your text")

	flag.Parse()

	for _, text := range flag.Args() {

		c, err := crypto.NewCrypt(*key)

		if err != nil {
			panic(err)
		}

		enc, err := c.Encrypt(text)

		if err != nil {
			panic(err)
		}

		plain, err := c.Decrypt(enc)

		if err != nil {
			panic(err)
		}

		fmt.Println(text, enc, plain)
	}

}
```

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
