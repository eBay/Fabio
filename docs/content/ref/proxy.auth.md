---
title: "proxy.auth"
---

`proxy.auth` configures one or more authorization schemes.

Each authorization scheme is configured with a list of
key/value options. Each scheme must have a unique
name which can then be referred to in a routing rule.

    name=<name>;type=<type>;opt=arg;opt[=arg];...

The following types of authorization schemes are available:

#### Basic

The basic authorization scheme leverages [Http Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication) and reads a [htpasswd](https://httpd.apache.org/docs/2.4/misc/password_encryptions.html) file at startup and credentials are cached until the service exits.

The `file` option contains the path to the htpasswd file. The `realm` parameter is optional (default is to use the `name`). The `refresh` option can set the htpasswd file refresh interval. Minimal refresh interval is `1s` to void busy loop. By default refresh is disabled i.e. set to zero.
Note: removing the htpasswd file will cause all requests to fail with HTTP status code 401 (Unauthorized) until the file is restored.

    name=<name>;type=basic;file=<file>;realm=<realm>;refresh=<interval>

Supported htpasswd formats are detailed [here](https://github.com/tg123/go-htpasswd)

#### External

This authorization scheme sends the incoming http request without body to the specified endpoint. The url path is appended to the endpoint value, thus the endpoint value should not end in a `/`. When the endpoint returns with an http 200 status code, the request is regarded as authorized and forwarded. If the endpoint returns with an http 302 status code, the redirection is send back to the client. For any other status code, the request will be regarded as unauthorized.

This scheme supports the `append-auth-headers` and `set-auth-headers` options. These configure headers from the authorization endpoint to copy to the upstream request. The `set-auth-headers` option replaces any header from the original client request, and the `append-auth-headers` option appends the header instead. Multiple headers can be specified separated by commas.

    name=<name>;type=external;endpoint=http://oathkeeper-api:4456/decisions;append-auth-headers=x-user,authorization

#### Examples

    # single basic auth scheme
    name=mybasicauth;type=basic;file=p/creds.file;

    # single basic auth scheme with refresh interval set to 30 seconds
    name=mybasicauth;type=basic;file=p/creds.htpasswd;refresh=30s

    # basic auth with multiple schemes
    proxy.auth = name=mybasicauth;type=basic;file=p/creds.htpasswd;refresh=30s
                 name=myotherauth;type=basic;file=p/other-creds.htpasswd;realm=myrealm

    # single ory oathkeeper decision api
    name=<name>;type=external;endpoint=http://oathkeeper-api:4456/decisions

The default is

    proxy.auth =
