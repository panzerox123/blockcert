# BlockCert
**Blockchain based digital certificates for PDFs and other documents.**

## Build
```
go build blockcert.go
```

## Usage
### Generating a Key Pair
```
blockcert keygen <filename>
```

Generates a keypair. Store this safely. The Private key will be used to generate your certificates, users can use the public key to verify the certificate