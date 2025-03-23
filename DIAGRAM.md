
# Diagrama do Network Manager TUI

```mermaid
graph TD
    %% Componente Principal
    A[main.go] --> B[Menu Principal]
    A --> C[Sistema de Idiomas]
    A --> D[Sistema de Logs]
    
    %% Menu e Navegação
    B --> E[Gerenciamento de Rede]
    B --> F[Status do Sistema]
    B --> G[Teste de Conexão]
    B --> H[Informações]
    B --> I[Ajuda]
    B --> J[Reiniciar]
    B --> K[Desligar]
    
    %% Gerenciamento de Rede
    E --> L[Configuração IPv4]
    E --> M[Configuração IPv6]
    E --> N[DHCP]
    E --> O[DNS]
    E --> P[Interfaces]
    
    %% Status e Monitoramento
    F --> Q[Status da Rede]
    F --> R[Conexões Ativas]
    F --> S[Interfaces de Rede]
    
    %% Segurança
    T[Sistema de Segurança] --> U[Verificação Root]
    T --> V[Validação de Entrada]
    T --> W[Modo Desenvolvimento]
    T --> X[Proteção Comandos]
    
    %% Logs e História
    D --> Y[Logs do Sistema]
    D --> Z[Histórico de Ações]
    D --> AA[Limpeza Automática]
    
    %% Interface do Usuário
    AB[Interface TUI] --> AC[Temas]
    AB --> AD[Navegação]
    AB --> AE[Feedback Visual]
    AB --> AF[Atalhos]

    %% Subgráfico de Componentes
    subgraph Componentes_Principais
        A
        B
        C
        D
        T
    end

    %% Subgráfico de Interface
    subgraph Interface_Usuario
        AB
        AC
        AD
        AE
        AF
    end

    %% Subgráfico de Funcionalidades
    subgraph Funcionalidades
        E
        F
        G
        H
    end
```

## Descrição da Arquitetura

### 1. Componentes Principais
- **main.go**: Ponto de entrada da aplicação
- **Menu Principal**: Interface central de navegação
- **Sistema de Idiomas**: Suporte a múltiplos idiomas
- **Sistema de Logs**: Registro de atividades
- **Sistema de Segurança**: Proteção e validação

### 2. Gerenciamento de Rede
- Configuração completa de rede IPv4/IPv6
- Suporte a DHCP e DNS
- Gerenciamento de interfaces de rede
- Monitoramento de conexões ativas

### 3. Segurança
- Verificação de privilégios root
- Validação de entrada de dados
- Modo de desenvolvimento seguro
- Proteção contra comandos perigosos
- Histórico de alterações

### 4. Monitoramento
- Status da rede em tempo real
- Informações do sistema
- Teste de conectividade
- Monitoramento de recursos

### 5. Interface do Usuário
- Temas personalizados
- Navegação intuitiva
- Feedback visual
- Atalhos de teclado
- Suporte a múltiplos idiomas

### 6. Sistema de Logs
- Registro detalhado de ações
- Histórico de modificações
- Limpeza automática após 90 dias
- Exportação de logs

### 7. Recursos Adicionais
- Reinicialização segura
- Desligamento controlado
- Sistema de ajuda integrado
- Backup de configurações
