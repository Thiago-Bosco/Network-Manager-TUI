
# Network Manager TUI - Documentação Técnica

## Comandos NMCLI por Função

### 1. GetNetworkConnections()
```bash
nmcli device status
```
- **Propósito**: Lista todas as interfaces de rede
- **Saída**: Nome da interface, tipo, estado e conexão
- **Exemplo**: eth0 ethernet connected Ethernet1

### 2. GetActiveConnection()
```bash
nmcli -t -f NAME connection show --active
```
- **Propósito**: Obtém conexão ativa atual
- **Flags**: 
  - `-t`: Formato tabulado
  - `-f NAME`: Filtra apenas nome

### 3. GetNetworkConnectionsInfo()
```bash
nmcli -t device status
nmcli -t device show [DEVICE]
```
- **Propósito**: Informações detalhadas da interface
- **Dados**: IP, DNS, Gateway, MAC

### 4. applyNetworkSettings()
#### IPv4 Manual
```bash
nmcli connection modify [INTERFACE] \
    ipv4.method manual \
    ipv4.addresses [IP/MASK] \
    ipv4.gateway [GATEWAY] \
    ipv4.dns [DNS]
```

#### IPv4 DHCP
```bash
nmcli connection modify [INTERFACE] ipv4.method auto
```

#### IPv6 Manual
```bash
nmcli connection modify [INTERFACE] \
    ipv6.method manual \
    ipv6.addresses [IPv6/PREFIX] \
    ipv6.gateway [GATEWAY] \
    ipv6.dns [DNS]
```

#### IPv6 DHCP
```bash
nmcli connection modify [INTERFACE] ipv6.method auto
```

#### IPv6 Desabilitado
```bash
# Limpa configurações
nmcli connection modify [INTERFACE] \
    ipv6.addresses "" \
    ipv6.gateway "" \
    ipv6.dns ""

# Desabilita IPv6
nmcli connection modify [INTERFACE] ipv6.method disabled
```

### 5. Parâmetros NMCLI

#### Métodos IPv4/IPv6
- `auto`: DHCP
- `manual`: Configuração estática
- `disabled`: Desabilitado (IPv6)

#### Propriedades
- `ipv4.addresses`: Endereço IP/máscara
- `ipv4.gateway`: Gateway padrão
- `ipv4.dns`: Servidores DNS
- `ipv6.addresses`: Endereço IPv6/prefixo
- `ipv6.gateway`: Gateway IPv6
- `ipv6.dns`: DNS IPv6

### 6. Ambiente Replit

#### 6.1 Arquivo .replit
O arquivo `.replit` é essencial para configurar o ambiente de desenvolvimento e produção:
- Define os módulos necessários (go-1.21)
- Configura workflows para diferentes modos de execução:
  - Development Mode: `go run main.go -dev`
  - Production Mode: `sudo go run main.go`

#### 6.2 Workflows
- **Development Mode**: Executa em modo desenvolvimento (sem sudo)
- **Production Mode**: Executa com privilégios sudo
- **Network Manager**: Workflow específico para gerenciamento de rede

#### 6.3 Deployment
- Build Command: Compilação do projeto
- Run Command: Execução do binário compilado
- Acesso via URL do Replit

#### 6.4 Estrutura no Replit
```
├── main.go          # Ponto de entrada
├── network/         # Gerenciamento de rede
├── menu/           # Interface do usuário
├── i18n/           # Internacionalização
└── logger/         # Sistema de logs
```

### 7. Códigos de Retorno

- 0: Sucesso
- 1: Erro genérico
- 2: Argumento inválido
- 3: Timeout
- 4: Interface não existe
- 5: Conexão falhou
- 6: Permissão negada
- 7: Não conectado
- 8: Conexão não encontrada
