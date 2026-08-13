package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

const CAN_ECU_ENGINE_ID = 0x0C4

func main() {
	socket, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		log.Fatalf("Erro ao abrir socket CAN: %v", err)
	}
	defer unix.Close(socket)

	ifi, err := net.InterfaceByName("vcan0")
	if err != nil {
		log.Fatalf("vcan0 não encontrada: %v", err)
	}

	addr := &unix.SockaddrCAN{Ifindex: ifi.Index}
	if err := unix.Bind(socket, addr); err != nil {
		log.Fatalf("Erro ao vincular socket: %v", err)
	}

	fmt.Println("Simulador de ECU iniciado na vcan0! Transmitindo telemetria...")

	var rpm uint16 = 800
	var speed uint8 = 0
	var temp uint8 = 75
	accelerating := true

	for {
		if accelerating {
			rpm += 150
			speed += 2
			if temp < 92 {
				temp++
			}
			if rpm >= 6500 {
				accelerating = false
			}
		} else {
			rpm -= 250
			if speed > 0 {
				speed--
			}
			if rpm <= 1500 {
				accelerating = true
			}
		}

		frameBuf := make([]byte, 16)

		binary.LittleEndian.PutUint32(frameBuf[0:4], CAN_ECU_ENGINE_ID)
		frameBuf[4] = 4

		binary.BigEndian.PutUint16(frameBuf[8:10], rpm)
		frameBuf[10] = speed
		frameBuf[11] = temp

		_, err := unix.Write(socket, frameBuf)
		if err != nil {
			log.Printf("Erro ao enviar frame CAN: %v", err)
		}

		time.Sleep(50 * time.Millisecond)
	}
}