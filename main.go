package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func factor(n int64) []int64 {
	r := []int64{}
	for i := int64(2); i <= n; i++ {
		for n%i == 0 {
			n /= i
			r = append(r, i)
		}
	}
	return r
}

func divisor(n int64) []int64 {
	r := []int64{}
	for i := int64(1); i <= n/i; i++ {
		if n%i == 0 {
			r = append(r, i)
			if n/i != i {
				r = append(r, n/i)
			}
		}
	}
	sort.Sort(sortBy(r))
	return r
}

func prime(n int64) []int64 {
	r := []int64{}
	p := make([]bool, n+1)
	for i := int64(2); i <= n; i++ {
		if p[i] {
			continue
		}
		r = append(r, i)
		if i > n/i {
			continue
		}
		for j := i * i; j <= n; j += i {
			p[j] = true
		}
	}
	return r
}

func gcd(a, b int64) int64 {
	if a < b {
		a, b = b, a
	}
	if a%b == 0 {
		return b
	}
	return gcd(b, a%b)
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

type sortBy []int64

func (a sortBy) Len() int           { return len(a) }
func (a sortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBy) Less(i, j int) bool { return a[i] < a[j] }

func parseInt(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	return n
}

func printInts(a []int64) {
	if numPerLine {
		for _, x := range a {
			fmt.Println(x)
		}
	} else {
		fmt.Println(strings.Trim(fmt.Sprint(a), "[]"))
	}
}

var (
	numPerLine = !isatty.IsTerminal(os.Stdout.Fd())
)

func main() {
	rootCmd := &cobra.Command{Use: "num"}
	rootCmd.PersistentFlags().BoolVarP(&numPerLine, "each-per-line", "1", numPerLine, "print each number per line")
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "factor [number]",
			Short: "Print prime factors",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				n := parseInt(args[0])
				printInts(factor(n))
			},
		},
		&cobra.Command{
			Use:   "divisor [number]",
			Short: "Print divisors",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				n := parseInt(args[0])
				printInts(divisor(n))
			},
		},
		&cobra.Command{
			Use:   "prime [number]",
			Short: "Print prime numbers less or equal to N",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				n := parseInt(args[0])
				printInts(prime(n))
			},
		},
		&cobra.Command{
			Use:   "gcd [number...]",
			Short: "Print GCD of numbers",
			Args:  cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				a := parseInt(args[0])
				for i := 1; i < len(args); i++ {
					b := parseInt(args[i])
					a = gcd(a, b)
					if a == 1 {
						break
					}
				}
				fmt.Println(a)
			},
		},
		&cobra.Command{
			Use:   "lcm [number...]",
			Short: "Print LCM of numbers",
			Args:  cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				a := int64(1)
				for i := 0; i < len(args); i++ {
					b := parseInt(args[i])
					a = lcm(a, b)
				}
				fmt.Println(a)
			},
		},
	)
	rootCmd.Execute()
}
