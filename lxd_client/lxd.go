package main


// Client for interacting with and LXD server deployed on AWS

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"io/ioutil"
)

func main() {

	fmt.Println("HELLO")
	cert_path := "./lxd/lxd.crt"

	//var cert *x509.Certificate

	certByt, err := ioutil.ReadFile(cert_path)
	if err != nil {
		panic(err)
	}

	key , err := ioutil.ReadFile("./lxd/lxd.key")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(key))

	fmt.Println(string(certByt))

	//fmt.Println(cert.)
	//Create the args
	args := &lxd.ConnectionArgs{
		TLSClientCert: string(certByt),
		TLSClientKey: string(key),
		SkipGetServer: true,
		InsecureSkipVerify: true,
	}
	//
	// Connect to the remote server


	c, err := lxd.ConnectLXD(
		"https://18.220.157.48:8443",
		args)
	if err != nil {
		panic(err)
	}
	//
	c.RequireAuthenticated(false)

	// Get the name of the first container
	conts , err := c.GetContainers()
	fmt.Println(conts[0].Name)
}
