
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

### 6. Códigos de Retorno

- 0: Sucesso
- 1: Erro genérico
- 2: Argumento inválido
- 3: Timeout
- 4: Interface não existe
- 5: Conexão falhou
- 6: Permissão negada
- 7: Não conectado
- 8: Conexão não encontrada
