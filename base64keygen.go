package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"os"
)

func main() {
	length := flag.Int("len", 32, "length of secret key")
	flag.Parse()

	secretKey := make([]byte, *length)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic(err)
	}

	w := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	defer w.Close()

	w.Write(secretKey)
}
