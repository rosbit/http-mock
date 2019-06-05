# http-mock

1. `http-mock` is intended to provide a mock HTTP server,  the repsonses of which are read from a YAML config file.
1. `http-mock` supports 4 types of body
    - body
       - The value of `body` as response body content will be output to HTTP client without any change.
    - file-body
       - The value of `file-body` is a file name, the content of which will be output to the client.
    - tmpl-body
       - The value of `tmpl-body` is looked as the content of a template, so template actions will be evaluated before outputing the result.
    - redirect-body
       - the value of `redirect-body` is a URL string, which will be redirected to with a 302 response status.

## Configuration

- http-mock makes use of YAML as its configuratoin.

- There's a sample configuration file [config.sample.yaml](config.sample.yaml)

 ```yaml
  port: 8080
  
  # remove static-home if no static files to show
  static-home:
    root: ./static
    alias: /static
  
  # remove any CORS header if needed
  cors:
    allow-origin: "*"
    allow-headers: "Content-Type, Authorization, Accept, X-Requested-With"
    allow-methods: "GET, PUT, POST, DELETE, OPTIONS"
    expose-headers: "X-Total-Count, X-Limit-Count, Link"
    allow-credentials: "*"
  
  default-content-type: "application/json"
  
  actions:
    -
      uri: /
      method: "GET, POST"
      response:
        status: 200
        headers:
          content-type: "text/plain"
        cookies:
          cookie-name: cookie-value
        tmpl-body: |-
          Welcome to use
          http-mock
          {{ .r.Method }}
    -
      uri: /ping
      method: "GET, POST"
      response:
        body: |-
          {
             "code": 200,
             "msg": "OK"
          }
    -
      uri: /api/users/:id
      response:
        tmpl-body: '{"code":200,"msg":"OK", "id": "{{pathVar .r "id"}}"}'
    -
      uri: /showfile
      response:
        content-type: "text/plain"
        file-body: a-file.txt
    -
      uri: /redirect
      response:
        redirect-body: /
  
 ```

## Usage
 1. with go1.11.x or above

 2. change to any directory, run the commands

    ```bash
    $ git clone https://github.com/rosbit/http-mock
    $ cd http-mock
    $ go build
    ```

 3. if everything is ok, there will be an executable `http-mock` under the directory

 4. for binary output under Linux/macos, click [releases](https://github.com/rosbit/http-mock/releases) to download it.

 5. run `http-mock` with an env variable `CONF_FILE` pointed to the configuration file

     ```bash
     $ CONF_FILE=./config.sample.yaml ./http-mock
     I am listening at :8080...
     ```

 6. open another terminal window, run `curl` to see the result

     ```bash
     $ curl http://localhost:8080/
     Welcome to use
     http-mock
     GET
     
     $ curl http://localhost:8080/showfile
     this is the content of a-file
     
     $ curl http://localhost:8080/api/users/1
     {"code":200,"msg":"OK", "id": "1"}
     
     $ curl -X POST http://localhost:8080/ping
     {
        "code": 200,
        "msg": "OK"
     }
     
     $ curl http://localhost:8080/static/a-file.txt  # static file
     this is the content of a-file
     
     $ curl -D - http://localhost:8080/redirect
     HTTP/1.1 301 Moved Permanently
     Access-Control-Allow-Credentials: *
     Access-Control-Allow-Headers: Content-Type, Authorization, Accept, X-Requested-With
     Access-Control-Allow-Methods: GET, PUT, POST, DELETE, OPTIONS
     Access-Control-Allow-Origin: *
     Access-Control-Expose-Headers: X-Total-Count, X-Limit-Count, Link
     Content-Type: text/html; charset=utf-8
     Location: /
     Date: Wed, 05 Jun 2019 06:43:58 GMT
     Content-Length: 36
     
     <a href="/">Moved Permanently</a>.
     
     ```

## Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.

