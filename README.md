distrobox enter
sudo dnf install -y can-utils
go mod init go-can-sim
go get golang.org/x/sys/unix
🏎️ Simulação de Telemetria Automotiva (CAN Bus) com Go e Linux

Introdução
--------
Este projeto simula e analisa em tempo real uma rede CAN (Controller Area Network) de um veículo. Demonstra como usar chamadas de sistema no Linux (`golang.org/x/sys/unix`) para ler e escrever dados via SocketCAN, usando Go. A proposta também mostra um fluxo de trabalho com containers (Distrobox) em distribuições imutáveis como Fedora Silverblue / Bluefin.

Pré-requisitos
-------------
- Linux com suporte a SocketCAN (módulos do kernel CAN/vcan).
- Go 1.18+ instalado no host ou no container de desenvolvimento.
- `distrobox` (opcional) se quiser isolar as ferramentas de inspeção.

1) Ativar a interface CAN virtual (vcan0) no Host
-------------------------------------------------
Execute no terminal do sistema host para carregar o módulo do kernel e criar a interface virtual de testes:

```bash
# 1. Carrega o módulo vcan no kernel
sudo modprobe vcan

# 2. Cria a interface vcan0
sudo ip link add dev vcan0 type vcan

# 3. Sobe a interface
sudo ip link set up vcan0

# Verifica
ip link show vcan0
```

2) Configurar o container Distrobox para ferramentas de inspeção (opcional)
--------------------------------------------------------------------------
Se você usar Distrobox (ou outro container que compartilhe o kernel e rede do host), instale a suíte `can-utils` para inspecionar o barramento:

```bash
 # Entra no container Distrobox
 distrobox enter

 # Em Fedora-based containers
 sudo dnf install -y can-utils

 # Em Debian/Ubuntu-based containers
 sudo apt update && sudo apt install -y can-utils
```

Observação: como o Distrobox compartilha o kernel e a pilha de rede com o host (`--net=host`), ferramentas como `candump` enxergarão a `vcan0` criada no host.

Configuração do módulo Go
-------------------------
No diretório do projeto, inicialize o módulo Go e adicione a dependência para chamadas de sistema Linux:

```bash
go mod init go-can-sim
go get golang.org/x/sys/unix
```

Código-fonte
------------
O projeto traz dois binários de exemplo:

1) `simulador.go` — ECU / Gerador de telemetria

🏎️ Simulação de Telemetria Automotiva (CAN Bus) — Instalação e funcionamento

Introdução
--------
Este documento agora contém apenas as instruções para configurar o ambiente, instalar os pacotes necessários e verificar o funcionamento básico do barramento CAN virtual (`vcan0`).

Pré-requisitos
-------------
- Sistema Linux com permissões de administrador (sudo).
- Pacotes: `iproute2` (para `ip`), `can-utils` (para inspecionar CAN). Opcional: `distrobox` para isolamento de ferramentas.

1) Ativar a interface CAN virtual (`vcan0`) no host
--------------------------------------------------
No terminal do host, execute:

```bash
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
ip link show vcan0
```

2) Instalar `can-utils` para inspeção do barramento
---------------------------------------------------
Instale `can-utils` no host ou em um container que compartilhe o kernel e a rede (ex.: Distrobox).

Em sistemas Fedora / RHEL:

```bash
sudo dnf install -y can-utils
```

Em Debian / Ubuntu:

```bash
sudo apt update
sudo apt install -y can-utils
```

Se estiver usando um container (Distrobox) para as ferramentas, entre no container e instale `can-utils` lá. Observação: containers que compartilham o kernel e a pilha de rede (`--net=host`) enxergam a `vcan0` do host.

3) Verificar o tráfego CAN (inspeção)
------------------------------------
Com a `vcan0` ativa e `can-utils` instalado, use `candump` para observar o tráfego no barramento:

```bash
candump vcan0
```

Isso exibirá pacotes CAN em tempo real (hexadecimal) que trafegarem pela `vcan0`.

Notas finais
-----------
- Esta versão do README removeu exemplos e trechos de código. Se quiser, eu posso adicionar instruções para compilar/rodar binários específicos, fornecer scripts de inicialização ou incluir notas sobre permissões de usuário para acessar sockets CAN.

Caso queira que eu gere um script de setup automático, confirme e eu crio o `setup.sh` com os comandos acima.