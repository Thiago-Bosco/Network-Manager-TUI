
# Network Manager TUI - Documentação

## 1. Visão Geral
O Network Manager TUI é uma interface de usuário baseada em terminal (TUI) para gerenciamento de redes em sistemas Linux. Desenvolvido em Go, oferece uma interface intuitiva e funcional para configuração e monitoramento de rede.

## 2. Arquitetura
- **Interface Principal**: Implementada com tview
- **Internacionalização**: Suporte para múltiplos idiomas
- **Modularização**: Componentes separados para cada funcionalidade

## 3. Funcionalidades
### 3.1 Configuração de Rede
- Configuração IPv4/IPv6
- Suporte DHCP e IP estático
- Configuração de DNS
- Gestão de interfaces de rede

### 3.2 Monitoramento
- Status da rede em tempo real
- Informações do sistema
- Teste de conectividade (ping)
- Monitoramento de recursos

### 3.3 Sistema
- Reinicialização do sistema
- Desligamento
- Troca de idioma
- Ajuda integrada

## 4. Segurança
- Verificação de privilégios root
- Validação de entrada
- Modo de desenvolvimento seguro
- Proteção contra comandos perigosos

## 5. Requisitos Técnicos
- Sistema Operacional: Linux
- Dependências: NetworkManager
- Linguagem: Go 1.21+
- Bibliotecas: tview, tcell

## 6. Interface do Usuário
- Menu intuitivo
- Navegação por teclado
- Feedback visual
- Suporte a cores e estilos

## 7. Desenvolvimento
### 7.1 Modos de Execução
- Modo Desenvolvimento: `go run main.go -dev`
- Modo Produção: `sudo go run main.go`

### 7.2 Estrutura do Projeto
```
├── main.go            # Ponto de entrada
├── menu/             # Interface principal
├── network/          # Gerenciamento de rede
├── sysinfo/          # Informações do sistema
└── i18n/             # Internacionalização
```

## 8. Melhorias Futuras
- [ ] Suporte a VPN
- [ ] Configuração de firewall
- [ ] Perfis de rede
- [ ] Backup de configurações
- [ ] Logs detalhados

## 9. Suporte e Manutenção
- Atualizações regulares
- Correções de bugs
- Melhorias de desempenho
- Documentação atualizada

## 10. Considerações de Uso
- Necessita privilégios root para operações do sistema
- Compatível com NetworkManager
- Interface adaptável a diferentes terminais
- Suporte a múltiplos idiomas
