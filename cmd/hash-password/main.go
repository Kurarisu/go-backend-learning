package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password asli
	password := "secret123"

	// Hash password pakai bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("Password asli :", password)
	fmt.Println("Password hash :", string(hashed))

	// Contoh verifikasi password
	err = bcrypt.CompareHashAndPassword(hashed, []byte("secret123"))
	if err != nil {
		fmt.Println("❌ Password salah")
	} else {
		fmt.Println("✅ Password benar")
	}
}
