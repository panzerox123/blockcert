# BlockCert : Internship Project
**Blockchain based digital certificates for PDFs and other documents.**

## Build
```
$ go build blockcert.go
```

## Usage
### Generating a Key Pair
```
$ blockcert keygen <filename>
```

Generates a keypair. Store this safely. The Private key will be used to generate your certificates, users can use the public key to verify the certificate

### Starting a blockcert shell

```
$ blockcert shell
```

#### Connecting to a node directly
```
$ blockcert shell <peer_addr>
```

#### Blockcert shell commands
- `addcert <filename> <private_key>`

    Create a certificate for a file
- `verifyallcert`

    Verify the blockchain
- `checkcert <filename> <public_key>`

    Check the certificate for a particular file
- `showallcerts`

    Display blockchain info
- `peers`

    List all peers
