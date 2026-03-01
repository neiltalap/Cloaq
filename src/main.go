// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.
package main

import (
	"cloaq/src/tun"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Использование: sudo ./cloaq run <ip_сервера:порт>")
		return
	}

	switch os.Args[1] {
	case "run":
		runCommand()
	case "help":
		helpCommand()
	case "settings":
		settingsCommand()
	default:
		log.Println("Неизвестная команда:", os.Args[1])
	}
}

func runCommand() {
	fmt.Println("Запуск Cloaq...")

	id, err := GenerateIdentity()
	if err != nil {
		log.Fatalf("Ошибка Identity: %v", err)
	}
	fmt.Printf("Ваш Public Key: %s\n", id.String())

	dev, err := tun.InitDevice()
	if err != nil {
		log.Fatal("Ошибка TUN:", err)
	}
	fmt.Println("Интерфейс готов:", dev.Name())

	if err := dev.Start(); err != nil {
		log.Fatal("Ошибка старта TUN:", err)
	}

	tr, err := NewTransport(":9000")
	if err != nil {
		log.Fatal("Ошибка транспорта:", err)
	}

	var targetPeers []net.UDPAddr
	if len(os.Args) > 2 {
		addr, err := net.ResolveUDPAddr("udp", os.Args[2])
		if err == nil {
			targetPeers = append(targetPeers, *addr)
			fmt.Printf("Добавлен пир для рассылки: %s\n", addr.String())
		}
	}

	incoming := make(chan []byte, 1024)
	go tr.Listen(incoming)

	go func() {
		for data := range incoming {
			if len(data) < 1 {
				continue
			}

			if data[0] == 0x01 {
				ipPacket := data[1:]

				err := tun.WritePacket(dev, ipPacket)
				if err != nil {
					log.Println("Ошибка записи в TUN:", err)
				} else {
					log.Printf("[NET -> TUN] Получен и расшифрован пакет (%d байт)", len(ipPacket))
				}
			}
		}
	}()

	fmt.Println("Cloaq работает. Ожидание трафика...")
	if err := ReadLoop(dev, tr, targetPeers); err != nil {
		log.Fatal("Ошибка ReadLoop:", err)
	}
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
