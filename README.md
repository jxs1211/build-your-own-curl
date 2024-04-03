# build-your-own-curl

https://dev.to/ericbsantana/build-your-own-curl-in-go-2p71

# Build your own curl in Golang

[#curl](https://dev.to/t/curl)[#go](https://dev.to/t/go)

In this tutorial, we'll walk through the process of creating a simple command-line tool similar to `curl` using Go and [Cobra](https://github.com/spf13/cobra), a CLI library for Go.

> Disclaimer: Here we are going to use the `net` package instead of the `http` that is available natively in Go. The reason for using `net` is to get a bit into the basics of creating a HTTP request to a server from scrach. You could easily use `http` package to enjoy all that stuff that is a pain to make from scratch. For instance, `http` already handles HTTPS requests.

It is a challenge because handling TCP connections and HTTP requests can be complex, but we'll keep it simple and focus on the basics. But will be a good start to understand how `curl` works under the hood and how to build a simple HTTP client using Go.

## Prerequisites

Before we begin, make sure you have the following installed:

- Go (version 1.22)

## Setting Up Your Project

First, let's create a new Go module for our project:

```
mkdir build-your-own-curl
cd build-your-own-curl
go mod init build-your-own-curl
```

Next, let's install Cobra:

```
go get -u github.com/spf13/cobra/cobra
```

Now, let's initialize Cobra in our project:

```
cobra-cli init
```

This will create the structure for our CLI application with the following files:

```
├── cmd/
│   ├── root.go
├── go.mod
├── go.sum
├── main.go
├── LICENSE
```

The `cmd` directory contains the root command file, which is where we'll define our whole application. It is also possible to create subcommands in separate files within the `cmd` directory. For now, we'll keep it simple and define everything in the `root.go` file.

## Building the Root Command

Now that we have our project set up, let's define our CLI commands using Cobra.We'll create a simple command to make HTTP GET requests.

To achieve this, we need to parse the URL provided as an argument, extract the hostname, port, and path, and then make an HTTP GET request to the specified URL. Firstly, we should parse the URL and extract the necessary information to create our TCP connection.

```go
// cmd/root.go
package cmd

import (
    "net/url"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "build-your-own-curl",
    Short: "A brief description of your application",
    Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Args: cobra.ExactArgs(1),

    Run: func(cmd *cobra.Command, args []string) {
        u, err := url.Parse(args[0])

        if err != nil {
            panic(err)
        }

        host := u.Hostname()
        port := u.Port()
        path := u.Path

        println("Host:", host)
        println("Port:", port)
        println("Path:", path)
    },
}

// rest of the code
```

The `Args` field specifies the number of arguments the command expects. In this case, we expect exactly one argument, which is the URL we want to make a request to. You may want to add custom validators or other `Cobra` built-in validators if you want to expand this functionality. If you run the application without an argument, you should see an error message printed to the console.

```
go run main.go

Error: accepts 1 arg(s), received 0
Usage:
  build-your-own-curl [flags]

Flags:
  -h, --help     help for build-your-own-curl
  -t, --toggle   Help message for toggle

exit status 1
```

But if you run the application with an URL as an argument, you should see that hostname, port, and path are printed to the console.

```
go run main.go https://example.com/get

Host: example.com
Port:
Path: /get
```

As you probably have noticed, the port is not being extracted correctly. This is because the `url.Parse` function does not return the port if it is not specified in the URL.

Most of us do not specify the port when making a day-to-day HTTP request using `curl` or a browser. To make our UX better, let's set the default port to 80 if it is not specified in the URL. In this tutorial, I will not handle HTTPS requests, which is why we are going to use only port HTTP (80) for now.

```
// cmd/root.go
// ...
  Run: func(cmd *cobra.Command, args []string) {
    u, err := url.Parse(args[0])

    if err != nil {
      panic(err)
    }

    host := u.Hostname()
    port := u.Port()

    if port == "" {
      port = "80"
    }

    path := u.Path

    println("Host:", host)
    println("Port:", port)
    println("Path:", path)
  },
}
// ...
```

Run the application to see the default port bring printed to the console.

```
go run main.go https://example.com/get

Host: example.com
Port: 80
Path: /get
```

Now that we have the necessary information to create a TCP connection, let's make an HTTP GET request to the specified URL.

A basic HTTP GET request header consists of the following:

- Request line: `GET /path HTTP/1.0`
- Host header: `Host: hostname`

We will send this request to the server using `net.Dial` and read the response. We will then print the response to the console.

```
// cmd/root.go
package cmd

import (
    "fmt"
    "net"
    "net/url"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "build-your-own-curl",
    Short: "A brief description of your application",
    Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Args: cobra.ExactArgs(1),

    Run: func(cmd *cobra.Command, args []string) {
        u, err := url.Parse(args[0])

        if err != nil {
            panic(err)
        }

        host := u.Hostname()
        port := u.Port()
        path := u.Path

        if port == "" {
            port = "80"
        }

        conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

        if err != nil {
            panic(err)
        }

        defer conn.Close()

        fmt.Fprintf(conn, "GET %s HTTP/1.0\r\nHost: %s\r\n\r\n", path, host)

        buf := make([]byte, 1024)
        n, err := conn.Read(buf)

        if err != nil {
            panic(err)
        }

        fmt.Println(string(buf[:n]))
    },
}

// ...
```

In the code above, we create a TCP connection to the specified host and port.We then send an HTTP GET request to the server and read the response.Finally, we print the response to the console. If you run the application with a valid URL, you should see the HTTP response printed to the console.

```
go run main.go http://eu.httpbin.org/get
HTTP/1.1 200 OK
Date: Wed, 27 Mar 2024 22:35:13 GMT
Content-Type: application/json
Content-Length: 203
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Host": "eu.httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-66049f21-305d0735393fd4ae2bc554a0"
  },
  "url": "http://eu.httpbin.org/get"
}
```

And you have it! A command-line tool to make GET requests to any URL using Go and Cobra. Additional changes can be made to handle different HTTP methods, headers, and more.

A challenge for you:

- Add support for different HTTP methods (e.g., POST, PUT, DELETE).
- Add support for custom headers.
- Add support for HTTPS requests.

The first and second were made in my project [gurl](https://github.com/ericbsantana/gurl). You can check it out for more inspiration or even to contribute and improve it!

Another project that worth to mention here is [go-curling](https://github.com/cdwiegand/go-curling) from [@chriswiegand](https://dev.to/chriswiegand)

It uses `http` package instead of `net` which shows a more simple and fast way to implement a cURL client in Go without making requests from scratch using `net.Dial`.

## Conclusion

Congratulations! You've built a simple command-line tool similar to `curl` using Go and Cobra. Feel free to expand upon this project by adding more features like handling different HTTP methods, headers, and more.

I have made a project called [gurl](https://github.com/ericbsantana/gurl) that is a simple CLI tool that can make HTTP requests using Go. You can check it out for more inspiration.
