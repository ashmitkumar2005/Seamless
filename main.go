package main

import (
	"fmt"
	"os"
	"net/http"
)

func main(){
	// listener, err := net.Listen("tcp",":9000")
	// if err!=nil{
	// 	fmt.Fprintf(os.Stderr,"Failed to start TCP server: %v\n",err)
	// 	return
	// }
	// defer listener.Close()

	// fmt.Println("Listneing to TCP port 9000...")
	
	// connection, err := listener.Accept()
	
	// if err!=nil{
	// 	fmt.Fprintf(os.Stderr,"Failed to accept connection %v\n",err)
	// 	return
	
	// }
	// defer connection.Close()
	// scanner := bufio.NewScanner(connection)
	
	// for scanner.Scan(){
	// 	fmt.Println(scanner.Text())
	// 	fmt.Fprintf(connection,"Message recieved")
	// }
	// if err := scanner.Err(); err!=nil{
	// 	fmt.Fprintf(os.Stderr,"Error reading message: %v\n",err)
	// }

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/send", handleSend)
	
	err := http.ListenAndServe(":9000", nil)
	if err!=nil{
		fmt.Fprint(os.Stderr, "HTTP server error: %v\n", err)
	}

}

func handleHome(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "<h1>Seamless Messenger</h1>")
	fmt.Fprint(w,"<form action=\"/send\" method=\"POST\"><input name=\"message\"><button type=\"submit\">Send</button></form>")
}

func handleSend(w http.ResponseWriter, r *http.Request){
	if err:= r.ParseForm(); err != nil {
		http.Error(w, "Could not read form", http.StatusBadRequest)
		return
	}
	message:= r.FormValue("message")
	fmt.Println("Recieved:", message)
	fmt.Fprint(w, "Message recieved: ", message)
}
