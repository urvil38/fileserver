package main

import "fmt"

func yellow(s string) string {
	return fmt.Sprintf("\x1b[1;33m%v\x1b[0m", s)
}

func red(s string) string {
	return fmt.Sprintf("\x1b[1;31m%v\x1b[0m", s)
}

func yellowUnderline(s string) string {
	return fmt.Sprintf("\x1b[4;33m%v\x1b[0m", s)
}

func redUnderline(s string) string {
	return fmt.Sprintf("\x1b[4;31m%v\x1b[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[1;32m%v\x1b[0m", s)
}
