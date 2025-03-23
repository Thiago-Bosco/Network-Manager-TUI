package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"networkmanager-tui/history"
	"networkmanager-tui/i18n"
	"networkmanager-tui/logger"
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
	// Parse command line flags
	devMode := flag.Bool("dev", false, "Enable development mode")
	flag.Parse()

	// Inicializa o sistema de logs
	if err := logger.Init(); err != nil {
		fmt.Printf("Erro ao inicializar logs: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Verifica privilégios root (skip in dev mode)
	if !*devMode && os.Geteuid() != 0 {
		fmt.Println(i18n.T("error_root_required"))
		fmt.Println("Por favor, execute com sudo ou como usuário root.")
		os.Exit(1)
	}
	
	// Verificação de permissões
	if _, err := os.Stat("/etc/network"); os.IsPermission(err) {
		fmt.Println("Erro: Sem permissão para acessar configurações de rede")
		os.Exit(1)
	}

	// Cria uma nova aplicação tview
	app := tview.NewApplication()

	// Aplica o tema personalizado
	setTheme(app)

	// Adiciona handler global para a tecla Esc retornar ao menu principal
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Se pressionar Esc, volta para o menu principal
		if event.Key() == tcell.KeyEscape {
			menu.StartMenu(app)
			return nil
		}
		return event
	})

	// Registra início da aplicação no histórico
	history.AddAction("system", "app_start", "Aplicação iniciada", "", "system")
	
	// Inicia o menu principal
	menu.StartMenu(app)

	// Inicia animação de título
	animateTitle(app, "Network Manager TUI")

	// Inicia a aplicação
	if err := app.Run(); err != nil {
		panic(err)
	}
}