package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenPass(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s $password [$cost]\n", args[0])
		os.Exit(1)
	}
	pass := args[1]
	cost, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse %s: %v\n", args[2], err)
		os.Exit(1)
	}

	t0 := time.Now()
	out, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't hash password: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
	fmt.Fprintf(os.Stderr, "%0.2fs elapsed\n", time.Since(t0).Seconds())
}
