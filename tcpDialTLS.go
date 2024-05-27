package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	var address string

	fmt.Println("---------------------------------------------------")
	fmt.Println("| Masukkan Web Address (contoh: Binusmaya.com) |")
	fmt.Println("---------------------------------------------------")
	fmt.Print("> ")
	fmt.Scanln(&address)

	address = strings.TrimPrefix(address, "https://")
	address = strings.TrimPrefix(address, "http://")

	if !strings.Contains(address, ":") {
		address = address + ":443"
	}

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		fmt.Printf("Gagal terhubung: %v\n", err)
		return
	}
	defer conn.Close()

	config := &tls.Config{
		InsecureSkipVerify: true, 
	}

	tlsConn := tls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		fmt.Printf("Gagal menyelesaikan TLS handshake: %v\n", err)
		return
	}
	defer tlsConn.Close()

	state := tlsConn.ConnectionState()

	tlsVersion := tlsVersionString(state.Version)
	fmt.Printf("Versi TLS: %s\n", tlsVersion)

	cipherSuite := tls.CipherSuiteName(state.CipherSuite)
	fmt.Printf("CipherSuite: %s\n", cipherSuite)

	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		if len(cert.Issuer.Organization) > 0 {
			fmt.Printf("Organisasi Penerbit: %s\n", cert.Issuer.Organization[0])
		} else {
			fmt.Println("Organisasi Penerbit: Tidak tersedia")
		}
	} else {
		fmt.Println("Sertifikat Tidak Ditemukan")
	}
}

func tlsVersionString(version uint16) string {
	versionMap := map[uint16]string{
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
	}
	if str, ok := versionMap[version]; ok {
		return str
	}
	return "Unknown"
}
