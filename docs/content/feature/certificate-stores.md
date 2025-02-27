---
title: "Certificate Stores"
since: "1.2"
---

Support for dynamic certificate stores which allow you to store certificates in
a central place and update them at runtime or generate them on the fly without
restart. You can store certificates in files, directories, on HTTP
servers in [Consul](https://consul.io/) or in
[Vault](https://vaultproject.io/).
You can use [Vault](https://vaultproject.io/) to generate certificates on the fly.

<!--more-->

Starting with version 1.2 fabio has support for dynamic certificate stores
which allow you to store certificates in a central place and update them at
runtime without restarting fabio. As of Go <= 1.7 only TLS certificates can be
changed at runtime. For updating client auth certificates [golang issue
16066](https://github.com/golang/go/issues/16066) is open.

Certificate stores are configured with the `proxy.cs` option.
You can configure one or multiple stores.

Each certificate source is configured with a list of
key/value options.

    cs=<name>;type=<type>;opt=arg;opt[=arg];...

Each source must have a **unique name** which is then
referenced in a listener configuration.

    proxy.cs = cs=mycerts;type=...
    proxy.addr = 1.2.3.4:9999;cs=mycerts;...

All certificates must be provided in **PEM format**

The following types of certificate sources are available:

 * [`file`](#file): legacy store for a single TLS and a set of client auth certificates
 * [`path`](#path): load certificates from a directory (e.g. managed by puppet/chef/ansible/...)
 * [`http`](#http): load certificates from an HTTP server
 * [`consul`](#consul) : load certificates from [Consul](https://consul.io/) KV store
 * [`vault`](#vault) : load certificates from [Vault](https://vaultproject.io/)

All certificate stores offer a set of [common options](#common-options). If you want to use
client certificate authentication with an Amazon API gateway check the `caupgcn` option there.

At the end you also find a list of [examples](#examples).

### File

The file certificate source supports one certificate which is loaded at
startup and is cached until the service exits.

The `cert` option contains the path to the certificate file. The `key`
option contains the path to the private key file. If the certificate file
contains both the certificate and the private key the `key` option can be
omitted. The `clientca` option contains the path to one or more client
authentication certificates.

##### Example

    cs=<name>;type=file;cert=p/a-cert.pem;key=p/a-key.pem;clientca=p/clientAuth.pem

### Path

The path certificate source loads certificates from a directory in
alphabetical order and refreshes them periodically.

The `cert` option provides the path to the TLS certificates and the
`clientca` option provides the path to the certificates for client
authentication.

TLS certificates are stored either in one or two files:

    www.example.com.pem or www.example.com-{cert,key}.pem

TLS certificates are loaded in alphabetical order and the first certificate
is the default for clients which do not support SNI.

The `refresh` option can be set to specify the refresh interval for the TLS
certificates. Client authentication certificates cannot be refreshed since
Go does not provide a mechanism for that yet.

The default refresh interval is 3 seconds and cannot be lower than 1 second
to prevent busy loops. To load the certificates only once and disable
automatic refreshing set `refresh` to zero.

##### Example

    cs=<name>;type=path;cert=path/to/certs;clientca=path/to/clientcas;refresh=3s

### HTTP

The http certificate source loads certificates from an HTTP/HTTPS server.

The `cert` option provides a URL to a text file which contains all files
that should be loaded from this directory. The filenames follow the same
rules as for the path source. The text file can be generated with:

    ls -1 *.pem > list

The `clientca` option provides a URL for the client authentication
certificates analogous to the `cert` option.

Authentication credentials can be provided in the URL as request parameter,
as basic authentication parameters or through a header.

The `refresh` option can be set to specify the refresh interval for the TLS
certificates. Client authentication certificates cannot be refreshed since
Go does not provide a mechanism for that yet.

The default refresh interval is 3 seconds and cannot be lower than 1 second
to prevent busy loops. To load the certificates only once and disable
automatic refreshing set `refresh` to zero.

##### Example

    cs=<name>;type=http;cert=https://host.com/path/to/cert/list&token=123
    cs=<name>;type=http;cert=https://user:pass@host.com/path/to/cert/list
    cs=<name>;type=http;cert=https://host.com/path/to/cert/list;hdr=Authorization: Bearer 1234

### Consul

The consul certificate source loads certificates from [Consul](https://consul.io/).

The `cert` option provides a KV store URL where the the TLS certificates are
stored.

The `clientca` option provides a URL to a path in the KV store where the the
client authentication certificates are stored.

The filenames follow the same rules as for the [`path`](#path) source.

The TLS certificates are updated automatically whenever the KV store
changes. The client authentication certificates cannot be updated
automatically since Go does not provide a mechanism for that yet.
(See [golang issue 16066](https://github.com/golang/go/issues/16066))

##### Example

    cs=<name>;type=consul;cert=http://localhost:8500/v1/kv/path/to/cert&token=123

### Vault

The Vault certificate store uses HashiCorp Vault as the certificate
store.

The `cert` option provides the path to the TLS certificates and the
`clientca` option provides the path to the certificates for client
authentication.

The `refresh` option can be set to specify the refresh interval for the TLS
certificates. Client authentication certificates cannot be refreshed since
Go does not provide a mechanism for that yet.

Certificate has to be stored in value. It means you have to write your cert into a *cert* and *key* fields of secret, that has to be your domain name.
Example:
```
vault write secret/fabio/certs/www.domain.com cert=@cert.pem key=@key.pem
```

The path to vault must be provided in the `VAULT_ADDR` environment
variable. The token must be provided in the `VAULT_TOKEN` environment
variable.

**fabio versions <= 1.2.1** require a token with root and/or sudo privileges to create an orphan
token for itself. This required fabio to have more privileges than it needs
and it also prevented revoking the fabio token if the parent token was revoked.
Therefore, supplying a token with root and/or sudo privileges is now deprecated
and will be removed in a later release.

**fabio versions > 1.2.1** will no longer attempt to create a token itself and
instead solely rely on the provided token. The provided token can be an orphan
and should be renewable for the duration fabio is expected to run. It is
recommended not to set the `explicit_max_ttl` unless fabio is restarted
before that time expires.

fabio needs the following policies set on the path where the
certificates are stored, for example:

      # For Vault < 0.7
      path "secret/fabio/cert" {
        capabilities = ["list"]
      }

      # For Vault >= 0.7; note the trailing slash
      path "secret/fabio/cert/" {
        capabilities = ["list"]
      }

      path "secret/fabio/cert/*" {
        capabilities = ["read"]
      }
      
      # To check for kv version in newer vault instances
      path "sys/internal/ui/mounts/secret/fabio/certs" {
        capabilities = ["read"]
      }
      
      # Allow our own token to be renewed.
      path "auth/token/renew-self" {
        capabilities = ["update"]
      }

##### Example

    cs=<name>;type=vault;cert=secret/fabio/certs

### Common options

All certificate stores support the following options:

 * `caupgcn` : Upgrade a self-signed client auth certificate with this common-name
            to a CA certificate. Typically used for self-singed certificates
            for the Amazon AWS API Gateway certificates which do not have the
            CA flag set which makes them unsuitable for client certificate
            authentication in Go. For the AWS API Gateway set this value
            to 'ApiGateway' to allow client certificate authentication.
            This replaces the deprecated parameter 'aws.apigw.cert.cn'
            which was introduced in version 1.1.5.

### Examples

     # file based certificate source
     proxy.cs = cs=some-name;type=file;cert=p/a-cert.pem;key=p/a-key.pem

     # path based certificate source
     proxy.cs = cs=some-name;type=path;cert=path/to/certs

     # HTTP certificate source
     proxy.cs = cs=some-name;type=http;cert=https://user:pass@host:port/path/to/certs

     # Consul certificate source
     proxy.cs = cs=some-name;type=consul;cert=https://host:port/v1/kv/path/to/certs?token=abc123

     # Vault certificate source
     proxy.cs = cs=some-name;type=vault;cert=secret/fabio/certs

     # Multiple certificate sources
     proxy.cs = cs=srcA;type=path;cert=path/to/certs,\
                cs=srcB;type=http;cert=https://user:pass@host:port/path/to/certs

     # path based certificate source for AWS Api Gateway
     proxy.cs = cs=some-name;type=path;cert=path/to/certs;clientca=path/to/clientcas;caupgcn=ApiGateway
