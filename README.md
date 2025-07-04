

# Network Manager TUI - Documentação

## 1. Visão Geral
O Network Manager TUI é uma interface de usuário baseada em terminal (TUI) para gerenciamento de redes em sistemas Linux. Desenvolvido em Go, oferece uma interface intuitiva e funcional para configuração e monitoramento de rede.

## 2. Arquitetura
- **Interface Principal**: Implementada com tview
- **Internacionalização**: Suporte para múltiplos idiomas
- **Modularização**: Componentes separados para cada funcionalidade
- **Segurança**: Verificação de privilégios root e validações

## 3. Funcionalidades e Comandos NMCLI

### 3.1 Configuração de Rede
#### IPv4 Manual
```bash
nmcli connection modify [INTERFACE] \
    ipv4.method manual \
    ipv4.addresses [IP/NETMASK] \
    ipv4.gateway [GATEWAY] \
    ipv4.dns [DNS1,DNS2]
```

#### IPv4 Automático (DHCP)
```bash
nmcli connection modify [INTERFACE] ipv4.method auto
```

#### IPv6 Manual
```bash
nmcli connection modify [INTERFACE] \
    ipv6.method manual \
    ipv6.addresses [IPv6/PREFIX] \
    ipv6.gateway [GATEWAY] \
    ipv6.dns [DNS1,DNS2]
```

#### IPv6 Automático
```bash
nmcli connection modify [INTERFACE] ipv6.method auto
```

#### IPv6 Desabilitado
```bash
nmcli connection modify [INTERFACE] \
    ipv6.addresses "" \
    ipv6.gateway "" \
    ipv6.dns "" \
    ipv6.method disabled
```

#### Aplicar Alterações
```bash
nmcli connection up [INTERFACE]
```

### 3.2 Monitoramento
#### Listar Conexões
```bash
nmcli device status
nmcli -t -f NAME connection show --active
```

#### Informações Detalhadas
```bash
nmcli -t device show [INTERFACE]
```

## 4. Modos de Execução
### 4.1 Desenvolvimento
```bash
go run main.go -dev
```
- Desativa verificação de privilégios root
- Permite testes sem sudo

### 4.2 Produção
```bash
sudo go run main.go
```
- Requer privilégios root
- Habilita todas as funcionalidades

## 5. Estrutura do Projeto
```
├── main.go            # Ponto de entrada
├── network/          # Gerenciamento de rede
├── i18n/             # Internacionalização
├── logger/           # Sistema de logs
└── menu/             # Interface principal
```

## 6. Códigos de Erro Comuns
- **exit status 1**: Erro genérico de execução
- **exit status 2**: Parâmetros inválidos
- **exit status 3**: Permissões insuficientes
- **exit status 4**: Interface não encontrada
- **exit status 5**: Configuração inválida

## 7. Validações
- Endereço IPv4: Formato xxx.xxx.xxx.xxx
- Máscara: Número CIDR (1-32) ou formato IPv4
- IPv6: Formato hexadecimal com ':'
- Prefixo IPv6: Número (1-128)

## 8. Dependências
- NetworkManager
- Go 1.21+
- tview
- tcell

## 9. Suporte e Manutenção
- Atualizações via Git
- Logs em /logs/
- Backup automático de configurações

## 10. Segurança
- Validação de entrada
- Verificação de privilégios
- Proteção contra comandos perigosos
- Logs de alterações
