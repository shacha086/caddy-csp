Installation:
``xcaddy build --with github.com/shacha086/caddy-csp``

Example use:
```Caddyfile
:8080 {
    route {
        csp {
            add default-src http://example.com :https
            add img-src :https
            remove script-src "'unsafe-inline'"
            add !all "'self'"
            set script-src "'none'"
        }
        respond "Hello, CSP!"
    }
}
```