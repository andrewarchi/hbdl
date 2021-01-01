package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/andrewarchi/hbdl/hb"
)

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		fmt.Fprintf(os.Stderr, "usage: %s command\n", os.Args[0])
		os.Exit(2)
	}
	cmd := os.Args[1]
	var cookieFile string
	flag.StringVar(&cookieFile, "cookies", "cookies.json", "location to store login cookies")
	if err := flag.CommandLine.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	c, err := hb.NewClient()
	try(err)
	switch cmd {
	case "login":
		try(login(c, cookieFile))
	case "help":

	}
}

func login(c *hb.Client, cookieFile string) error {
	fmt.Print("Username: ")
	username, err := readLine()
	if err != nil {
		return err
	}
	fmt.Print("Password: ")
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	password := string(p)
	fmt.Println()

	loginErr := c.Login(username, password)
	if loginErr == hb.ErrGuardRequired {
		fmt.Println("Enter the code sent to your email address to verify your account.")
		fmt.Print("Code: ")
		guard, err := readLine()
		if err != nil {
			return err
		}
		loginErr = c.LoginGuard(username, password, guard)
	} else if loginErr == hb.Err2FARequired {
		fmt.Println("Enter the code from your two-factor authenticator to verify your account.")
		fmt.Print("Code: ")
		code, err := readLine()
		if err != nil {
			return err
		}
		loginErr = c.Login2FA(username, password, code)
	}
	if loginErr != nil {
		return loginErr
	}
	return c.SaveCookies(cookieFile)
}

func readLine() (string, error) {
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return line, err
	}
	if len(line) >= 2 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], nil
	}
	return line[:len(line)-1], nil
}

func try(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
