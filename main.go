package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Customer struct {
	name    string
	balance float64
	owed    map[string]float64
}

type ATM struct {
	customers map[string]*Customer
	current   *Customer
}

func NewATM() *ATM {
	return &ATM{
		customers: make(map[string]*Customer),
	}
}

func (atm *ATM) login(name string) {
	if atm.current != nil {
		fmt.Printf("Error: Already logged in as %s\n", atm.current.name)
		return
	}
	customer, exists := atm.customers[name]
	if !exists {
		customer = &Customer{
			name: name,
			owed: make(map[string]float64),
		}
		atm.customers[name] = customer
	}
	atm.current = customer
	fmt.Printf("Hello, %s!\nYour balance is $%.2f\n", customer.name, customer.balance)

	if len(customer.owed) > 0 {
		for owedName, amount := range customer.owed {
			fmt.Printf("Owed $%.2f to %s\n", amount, owedName)
		}
	}

	amountsOwedByYou := make(map[string]float64)
	for _, target := range atm.customers {
		if amount, exists := target.owed[customer.name]; exists && amount > 0 && target.name != atm.current.name {
			amountsOwedByYou[target.name] = amount
		}
	}

	if len(amountsOwedByYou) > 0 {
		for owedName, amount := range amountsOwedByYou {
			fmt.Printf("Owed $%.2f from %s \n", amount, owedName)
			delete(amountsOwedByYou, owedName)
		}
	}
}

func (atm *ATM) deposit(amount float64) {
	if atm.current == nil {
		fmt.Println("Error: No customer logged in")
		return
	}

	atm.current.balance += amount

	for targetName, owedAmount := range atm.current.owed {
		if atm.current.balance >= owedAmount {
			atm.current.balance -= owedAmount
			delete(atm.current.owed, targetName)
			target, exists := atm.customers[targetName]
			if !exists {
				target = &Customer{
					name: targetName,
					owed: make(map[string]float64),
				}
				atm.customers[targetName] = target
			}
			target.balance += owedAmount
			fmt.Printf("Transferred $%.2f to %s \n", owedAmount, targetName)

		} else {
			atm.current.balance = 0
			atm.current.owed[targetName] -= amount
			target, exists := atm.customers[targetName]
			if !exists {
				target = &Customer{
					name: targetName,
					owed: make(map[string]float64),
				}
				atm.customers[targetName] = target
			}
			target.balance += amount
			fmt.Printf("Transferred $%.2f to %s \n", amount, targetName)
		}

	}

	fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)

	if len(atm.current.owed) > 0 {
		for targetName, owedAmount := range atm.current.owed {
			fmt.Printf("Owed $%.2f to %s\n", owedAmount, targetName)
		}
	}
}

func (atm *ATM) withdraw(amount float64) {
	if atm.current == nil {
		fmt.Println("Error: No customer logged in")
		return
	}
	if amount > atm.current.balance {
		fmt.Println("Error: Insufficient funds")
		return
	}
	atm.current.balance -= amount
	fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)
}

func (atm *ATM) transfer(targetName string, amount float64) {
	if atm.current == nil {
		fmt.Println("Error: No customer logged in")
		return
	}
	if targetName == atm.current.name {
		fmt.Println("Error: Cannot transfer to self")
		return
	}

	target, exists := atm.customers[targetName]
	if !exists {
		target = &Customer{
			name: targetName,
			owed: make(map[string]float64),
		}
		atm.customers[targetName] = target
	}

	if amount > atm.current.balance {
		diff := amount - atm.current.balance
		target.balance += atm.current.balance
		fmt.Printf("Transferred $%.2f to %s\n", atm.current.balance, targetName)
		atm.current.owed[targetName] += diff
		atm.current.balance = 0
		fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)
		fmt.Printf("Owed $%.2f to %s\n", atm.current.owed[targetName], targetName)
	} else {
		if owedAmount, exist := target.owed[atm.current.name]; exist {
			if owedAmount > amount {
				target.owed[atm.current.name] -= amount
				fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)
				fmt.Printf("Owed $%.2f from %s \n", target.owed[atm.current.name], target.name)
			} else {
				diff := amount - target.owed[atm.current.name]
				delete(target.owed, atm.current.name)
				atm.current.balance -= diff
				fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)
			}

		} else {
			atm.current.balance -= amount
			target.balance += amount
			fmt.Printf("Transferred $%.2f to %s\n", amount, targetName)
			fmt.Printf("Your balance is $%.2f\n\n", atm.current.balance)
			if _, exists := target.owed[atm.current.name]; exists {
				target.owed[atm.current.name] -= amount
				if target.owed[atm.current.name] <= 0 {
					delete(target.owed, atm.current.name)
				}
			}
		}

	}
}

func (atm *ATM) logout() {
	if atm.current == nil {
		fmt.Println("Error: No customer logged in")
		return
	}
	fmt.Printf("Goodbye, %s!\n", atm.current.name)
	atm.current = nil
}

func main() {
	atm := NewATM()

	for {

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("$ ")
		scanner.Scan()
		input := scanner.Text()

		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "login":
			if len(parts) != 2 {
				fmt.Println("Usage: login [name]")
				continue
			}
			name := parts[1]
			atm.login(name)

		case "deposit":
			if len(parts) != 2 {
				fmt.Println("Usage: deposit [amount]")
				continue
			}
			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			atm.deposit(amount)

		case "withdraw":
			if len(parts) != 2 {
				fmt.Println("Usage: withdraw [amount]")
				continue
			}
			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			atm.withdraw(amount)

		case "transfer":
			if len(parts) != 3 {
				fmt.Println("Usage: transfer [target] [amount]")
				continue
			}
			targetName := parts[1]
			amount, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				fmt.Println("Invalid amount")
				continue
			}
			atm.transfer(targetName, amount)

		case "logout":
			atm.logout()

		default:
			fmt.Println("Unknown command")
		}
	}
}
