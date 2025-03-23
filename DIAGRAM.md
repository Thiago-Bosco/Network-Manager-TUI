
# Diagrama do Network Manager TUI

```mermaid
graph TD
    A[main.go] --> B[menu/menu.go]
    A --> C[i18n/i18n.go]
    
    B --> D[network/network.go]
    B --> E[sysinfo/sysinfo.go]
    
    D --> F[Configuração de Rede]
    D --> G[Status da Rede]
    D --> H[Teste de Ping]
    
    E --> I[Info do Sistema]
    E --> J[Hardware]
    E --> K[Recursos]

    subgraph Interface
        B --> L[Menu Principal]
        L --> M[Configurar Rede]
        L --> N[Status da Rede]
        L --> O[Teste de Ping]
        L --> P[Info do Sistema]
        L --> Q[Ajuda]
        L --> R[Reiniciar]
        L --> S[Desligar]
        L --> T[Idioma]
    end

    subgraph Componentes
        D --> U[IPv4]
        D --> V[IPv6]
        D --> W[DHCP]
        D --> X[DNS]
    end
```
