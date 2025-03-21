package main

import (
        "flag"
        "fmt"
        "os"
        "time"

        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
        
        "networkmanager-tui/i18n"
        "networkmanager-tui/menu"
)

// Define o tema personalizado (escuro e moderno)
func setTheme(app *tview.Application) {
        // Cores RGB
        tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(30, 30, 30) // Fundo escuro
        tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(30, 30, 30)  // Fundo escuro
        tview.Styles.PrimaryTextColor = tcell.NewRGBColor(255, 255, 255)      // Texto branco
        tview.Styles.SecondaryTextColor = tcell.NewRGBColor(200, 200, 200)    // Texto cinza claro
        tview.Styles.TertiaryTextColor = tcell.NewRGBColor(180, 180, 180)     // Texto cinza médio
        tview.Styles.BorderColor = tcell.NewRGBColor(100, 100, 100)           // Bordas cinza escuro
        tview.Styles.TitleColor = tcell.NewRGBColor(0, 255, 0)                // Títulos verde
}

func animateTitle(app *tview.Application, title string) {
        go func() {
                for {
                        for i := 0; i < 5; i++ {
                                app.QueueUpdateDraw(func() {
                                        // Adiciona efeito de animação mudando a cor do título
                                        tview.Styles.TitleColor = tcell.NewRGBColor(
                                                int32(255-i*50), int32(100+i*20), int32(100-i*10))
                                })
                                // Pausa para dar efeito de animação
                                time.Sleep(time.Millisecond * 500)
                        }
                }
        }()
}

func main() {
        // Flag para desenvolvimento (desativa a verificação de root)
        devMode := flag.Bool("dev", false, "Ativa o modo de desenvolvimento (desativa verificação de privilégios)")
        flag.Parse()

        // Check root privileges unless in development mode
        if !*devMode && os.Geteuid() != 0 {
                fmt.Println(i18n.T("error_root_required"))
                fmt.Println("Por favor, execute com sudo ou como usuário root.")
                fmt.Println("Ou use a flag -dev para testes em desenvolvimento: go run main.go -dev")
                os.Exit(1)
        }

        // Aviso se estiver em modo de desenvolvimento
        if *devMode {
                fmt.Println("ATENÇÃO: Modo de desenvolvimento ativado. A verificação de privilégios root está desativada.")
                fmt.Println("Algumas funcionalidades que modificam o sistema operacional não funcionarão corretamente.")
        }

        // Cria uma nova aplicação tview
        app := tview.NewApplication()

        // Aplica o tema personalizado
        setTheme(app)

        // Inicia o menu principal
        menu.StartMenu(app)

        // Inicia animação de título
        animateTitle(app, "Network Manager TUI")

        // Inicia a aplicação
        if err := app.Run(); err != nil {
                panic(err)
        }
}
