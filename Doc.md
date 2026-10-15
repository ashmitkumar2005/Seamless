# Seamless project learning notes

This document collects the explanations from our first networking exercises. It is a learning record; it does not replace your source code. You write and maintain the project files yourself.

## What the first TCP server taught us

The first version was a raw TCP server. It listened on the laptop, accepted one connection, read newline-delimited input, printed it, and could send a reply. The packages came from Go's standard library:

- `fmt` prints formatted text.
- `net` creates TCP listeners and connections.
- `bufio` scans the connection a line at a time.
- `os` exposes operating-system streams such as standard error (`os.Stderr`).

An import group uses parentheses so each package path can be listed on its own line:

```go
import (
	"bufio"
	"fmt"
	"net"
	"os"
)
```

### Starting the listener

```go
listener, err := net.Listen("tcp", ":9000")
```

- `net.Listen` starts listening for network connections.
- `"tcp"` selects TCP, a reliable, ordered byte stream.
- `":9000"` selects port `9000` on all local network interfaces. That lets other devices connect through Wi-Fi as well as local programs.
- `:=` declares new variables and assigns the function's results.
- `listener` is the listening socket; `err` is the error value, if startup failed.

```go
if err != nil {
```

This begins an error check. In Go, `nil` means there is no error. The opening brace begins the block that runs only when an error exists.

```go
fmt.Fprintf(os.Stderr, "Failed to start TCP server: %v\n", err)
```

- `fmt.Fprintf` writes formatted output to a destination supplied as its first argument.
- `os.Stderr` is the process's standard error stream, a conventional place for error messages.
- `%v` is replaced by the value of `err`.
- `\n` adds a newline.

```go
return
```

This exits `main`. If the listener could not start, there is no usable listener to continue with.

```go
}
```

This closes the startup error-check block. If there was no error, execution continues after it.

```go
defer listener.Close()
```

`defer` schedules `listener.Close()` to run when `main` returns. Closing releases the listening socket.

```go
fmt.Println("Listening on TCP port 9000...")
```

`fmt.Println` writes its argument to the terminal and appends a newline. This status message confirms the listener started; it does not itself accept a connection.

### Accepting and reading a connection

```go
connection, err := listener.Accept()
```

`Accept` waits for a client to connect. When one connects, it returns a `connection` for exchanging data and an `err` value. This first server accepts one client at a time.

```go
if err != nil {
	fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
	return
}
```

The `if` checks whether accepting failed. If it did, the server prints the reason to standard error and returns because it has no client connection to use.

```go
defer connection.Close()
```

This schedules the accepted client connection to close when `main` returns.

```go
scanner := bufio.NewScanner(connection)
```

`bufio.NewScanner` creates a scanner that reads from the connection. By default, it splits input into lines. `:=` declares `scanner` and stores the returned scanner in it.

```go
for scanner.Scan() {
	fmt.Println(scanner.Text())
}
```

- `for` repeats the block.
- `scanner.Scan()` reads the next line and returns `true` when a line is available. The loop ends when input ends or scanning fails.
- `scanner.Text()` returns the current line.
- `fmt.Println` prints that line on the laptop.

To reply to each line from inside the loop, we added:

```go
fmt.Fprintln(connection, "Message received")
```

`fmt.Fprintln` writes to its first argument and adds a newline. Here the destination is the TCP `connection`, so the text travels back to the client. This simple version does not check for a write error.

After the loop:

```go
if err := scanner.Err(); err != nil {
	fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
}
```

`scanner.Err()` reports an error that stopped scanning. In `if err := ...`, the new `err` variable is scoped to this `if` statement and its block. If it is not `nil`, the code prints it.

### What the browser test showed

A browser is an HTTP client. Opening the TCP server's address made the browser send an HTTP request such as `GET / HTTP/1.1`, followed by headers. The raw TCP server printed those bytes because it does not understand HTTP. It also did not send a valid HTTP response, so the browser could wait or show an error. This is why we moved to Go's `net/http` package for browser testing.

## The browser-based HTTP version

For the browser version, the imports became:

```go
import (
	"fmt"
	"net/http"
	"os"
)
```

- `fmt` prints messages and writes response text.
- `net/http` parses HTTP requests and serves HTTP responses.
- `os` provides `os.Stderr` for startup errors.
- `net` and `bufio` are no longer needed once the raw TCP listener and scanner are removed.

### Registering routes and starting HTTP

```go
http.HandleFunc("/", handleHome)
```

`http.HandleFunc` registers a function for a URL path. `"/"` is the root page, and `handleHome` is the function that serves it.

```go
http.HandleFunc("/send", handleSend)
```

This registers a second handler for the `/send` path, where the HTML form submits its message.

```go
err := http.ListenAndServe(":9000", nil)
```

- `http.ListenAndServe` opens an HTTP server on port `9000` and waits while it handles requests.
- `nil` means to use Go's default HTTP request router, where the two `HandleFunc` routes were registered.
- `err :=` declares `err` and stores the result. If the server stops, the result explains why.

```go
if err != nil {
	fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
}
```

This checks whether the HTTP server stopped with an error, then prints that error to standard error. `Fprintf` is the formatted version of `Fprint`; `%v` only acts as a placeholder with a formatted-print function.

### Serving the page

```go
func handleHome(w http.ResponseWriter, r *http.Request) {
```

This declares the root-page handler function. `w` is the response writer used to send data back to the browser. `r` contains details of the incoming HTTP request. The opening brace starts the function body.

```go
fmt.Fprint(w, "<h1>Seamless Messenger</h1>")
```

`fmt.Fprint` writes its text to `w`. The browser recognizes `<h1>` and `</h1>` as HTML tags and displays a heading.

```go
fmt.Fprint(w, "<form action=\"/send\" method=\"POST\"><input name=\"message\"><button type=\"submit\">Send</button></form>")
```

This writes a simple HTML form:

- `<form>` begins the form.
- `action="/send"` tells the browser which path receives the submission.
- `method="POST"` tells the browser to send the form fields in an HTTP POST request.
- `<input name="message">` creates a text input. Its `name` is the key the server uses to retrieve its value.
- `<button type="submit">` creates a button that submits the form.
- The backslashes (`\"`) let quote characters appear inside a Go double-quoted string.

The closing brace ends `handleHome`.

### Receiving a submitted message

```go
func handleSend(w http.ResponseWriter, r *http.Request) {
```

This defines the handler registered for `/send`. It receives the browser request in `r` and writes the HTTP response through `w`.

```go
if err := r.ParseForm(); err != nil {
```

`r.ParseForm()` reads and parses submitted form fields. The short declaration creates an `err` value scoped to this `if`. `err != nil` checks whether parsing failed.

```go
http.Error(w, "Could not read form", http.StatusBadRequest)
```

`http.Error` sends an error response. `http.StatusBadRequest` is HTTP status code 400, meaning the request could not be understood.

```go
return
```

This stops `handleSend` after the error response so the handler does not try to use invalid form data.

```go
}
```

This closes the parsing error check. If parsing succeeded, execution continues afterward.

```go
message := r.FormValue("message")
```

`FormValue` looks up the form field named `message` and returns its submitted text. `:=` declares the variable and stores the result.

```go
fmt.Println("Received:", message)
```

This prints a label and the message in the laptop terminal. `Println` separates the values with a space and ends the output with a newline.

```go
fmt.Fprint(w, "Message received: ", message)
```

This writes a confirmation and the submitted message in the HTTP response. The browser displays that response after submitting the form.

The closing brace ends `handleSend`.

## Network troubleshooting notes

### Localhost and device addresses

`127.0.0.1` (often called `localhost`) always means “this same device.” On the phone it points to the phone, not the laptop. To connect from the phone, use the laptop's current Wi-Fi IPv4 address and port, for example `http://10.148.245.37:9000`.

The laptop's address was checked with:

```sh
ip -4 -br addr
```

This displays IPv4 addresses in a compact form. The active Wi-Fi interface was `wlp1s0`; on the phone hotspot it had `10.148.245.37/24`.

The phone's Wi-Fi address and route were checked with:

```sh
adb shell ip -4 addr show wlan0
adb shell ip route
```

The phone had an address on the same hotspot subnet. A ping test succeeded:

```sh
adb shell ping -c 3 10.148.245.37
```

This confirmed basic IP reachability, but it did not prove that TCP port `9000` was allowed.

### Checking the listener and firewall

```sh
ss -ltnp | grep 9000
```

This checks for a listening TCP socket on port `9000`. The output `*:9000` means the server is listening on all local IPv4/IPv6 interfaces, rather than only loopback.

```sh
ss -tnp | grep ':9000'
```

This checks current TCP connections containing port `9000` in their displayed address. No output means it showed no matching active TCP connection; it does not show listening sockets because `-l` was not requested.

`sudo nft list ruleset` showed UFW-managed firewall rules. The incoming chain had a drop policy, and the user input rules did not include TCP port `9000`. That explained why ping worked while the TCP connection timed out.

The rule added for the laptop's hotspot subnet was:

```sh
sudo ufw allow in on wlp1s0 from 10.148.245.0/24 to any port 9000 proto tcp
```

- `ufw allow in` allows incoming traffic.
- `on wlp1s0` applies to the Wi-Fi interface.
- `from 10.148.245.0/24` limits the source to that local subnet.
- `to any port 9000 proto tcp` selects TCP traffic addressed to port `9000`.

After this rule, the phone reached the listener and the request appeared in the server terminal.

### Sending a test line from the phone

The phone's built-in `toybox nc` command was used through ADB to send a plain line to the TCP server:

```sh
adb shell 'echo hello | toybox nc -w 2 10.148.245.37 9000'
```

- `adb shell` runs a command on the connected Android phone.
- `echo hello` writes `hello` followed by a newline.
- `|` sends that output to the next command as input.
- `toybox nc` is Android's netcat utility, a simple TCP client.
- `-w 2` sets a two-second timeout.
- `10.148.245.37 9000` are the laptop address and server port.

With the raw TCP server running, `hello` appeared in its terminal. After a reply line was added to the server, the phone's netcat command also received the reply. This verified two-way TCP communication.

## Go errors encountered during the HTTP transition

- `undefined: net` and `undefined: bufio` meant old raw-TCP code still referred to packages that had been removed from the import block. Removing the old TCP listener/scanner statements was necessary when switching to `net/http`.
- `http.Handlefunc` failed because Go identifiers are case-sensitive. The standard-library name is `http.HandleFunc`, with a capital `F`.
- `no new variables on left side of :=` meant `err` had already been declared in the same scope. Remove the leftover TCP block (which declared `err`) or use assignment only if `err` already exists and is in scope. Short declarations (`:=`) need at least one new variable on their left side.
- `http.StatusBarRequest` was a misspelling. The standard constant is `http.StatusBadRequest`.
- Use `fmt.Fprintf` when a format string contains placeholders such as `%v`. `fmt.Fprint` writes its arguments as-is and does not substitute format verbs.

## Current learning milestone

The browser can load the Go HTTP server's page, submit a message, and receive a confirmation. The server prints the message in the laptop terminal. This is a basic request/response interaction; the page currently submits and navigates to the response. A later improvement is to keep the page in place and display replies dynamically, then add conversation history or persistent real-time messaging.
