package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/unix"
)

const CAN_ECU_ENGINE_ID = 0x0C4

func main() {
	// 1. Abrir Socket AF_CAN do Linux nativamente via pacote unix
	socket, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		log.Fatalf("Erro ao criar socket CAN: %v", err)
	}
	defer unix.Close(socket)

	// 2. Obter o índice da interface vcan0
	ifi, err := net.InterfaceByName("vcan0")
	if err != nil {
		log.Fatalf("Interface vcan0 não encontrada. Verifique se ativou o módulo virtual: %v", err)
	}

	// 3. Bind do Socket na vcan0 usando SockaddrCAN do unix
	addr := &unix.SockaddrCAN{Ifindex: ifi.Index}
	if err := unix.Bind(socket, addr); err != nil {
		log.Fatalf("Erro ao dar bind no socket: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, unix.SIGTERM)

	fmt.Print("\033[H\033[2J")

	frameBuf := make([]byte, 16)

	go func() {
		for {
			n, err := unix.Read(socket, frameBuf)
			if err != nil || n < 16 {
				continue
			}

			canID := binary.LittleEndian.Uint32(frameBuf[0:4]) & 0x1FFFFFFF
			dataLen := frameBuf[4]
			data := frameBuf[8:16]

			if canID == CAN_ECU_ENGINE_ID && dataLen >= 4 {
				rpm := binary.BigEndian.Uint16(data[0:2])
				speed := uint8(data[2])
				temp := uint8(data[3])

				renderDashboard(rpm, speed, temp)
			}
		}
	}()

	<-sigChan
	fmt.Println("\n\nEncerrando Analisador CAN...")
}

func renderDashboard(rpm uint16, speed uint8, temp uint8) {
	fmt.Print("\033[H")

	fmt.Println("==================================================")
	fmt.Println("       AUTOMOTIVE CAN BUS - ANALYZER (GO)         ")
	fmt.Println("==================================================")
	fmt.Printf(" Interface: vcan0  | Frame ID: 0x%03X               \n", CAN_ECU_ENGINE_ID)
	fmt.Println("--------------------------------------------------")

	rpmBars := int(rpm) / 250
	if rpmBars > 32 {
		rpmBars = 32
	}
	rpmBarStr := ""
	for i := 0; i < 32; i++ {
		if i < rpmBars {
			if i > 24 {
				rpmBarStr += "\033[31m|\033[0m"
			} else {
				rpmBarStr += "\033[32m|\033[0m"
			}
		} else {
			rpmBarStr += "."
		}
	}

	fmt.Printf(" [RPM]  %4d  [%s]\n", rpm, rpmBarStr)
	fmt.Printf(" [VEL]  %3d km/h\n", speed)
	fmt.Printf(" [TEMP] %3d °C\n", temp)
	fmt.Println("--------------------------------------------------")
	fmt.Printf(" Timestamp: %s\n", time.Now().Format("15:04:05.000"))
	fmt.Println(" Pressione Ctrl+C para sair.                       ")
}