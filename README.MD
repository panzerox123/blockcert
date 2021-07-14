# BlockCert : Internship Project
**Blockchain based digital certificates for PDFs and other documents.**

## Build
### Build using Go
```
$ go build blockcert.go
```
### Build using Docker
```
$ docker build docker/ -t <imagename>:<tagname>
```
#### Running the docker image
**From Dockerhub**
```
$ docker run -p <port>:8080 panzerox123/blockcert
```
**If you want to use your own image**
```
$ docker run -p <port>:8080 <imagename>:<tagname>
```

## Usage
### Generating a Key Pair
```
$ blockcert -k keygen -o <filename>
```

Generates a keypair. Store this safely. The Private key will be used to generate your certificates, users can use the public key to verify the certificate.
Default file is `certificate.key`.

### Starting blockcert (Start node along with a front-end server)
```
$ blockcert
```
Starts a blockcert node on your system along with a front-end server for the GUI. 
TODO: API Documentation

### Starting a blockcert debug shell(For testing functionality)

```
$ blockcert -debug
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
- `keygen <filename_optional>`
    
    Generate a keypair

- `discover`

    Run peer discovery function
