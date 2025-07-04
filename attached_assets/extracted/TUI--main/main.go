package main

import (
	"datadike-tui/menu"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Define o tema personalizado (escuro e moderno)
func setTheme(app *tview.Application) {
	// Cores RGB??///??
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
	// Cria uma nova aplicação tview
	app := tview.NewApplication()

	// Aplica o tema personalizado
	setTheme(app)

	// Inicia o menu principal
	menu.Start(app)

	// Inicia animação de título
	animateTitle(app, "DataDike TUI")

	// Inicia a aplicação
	if err := app.Run(); err != nil {
		panic(err)
	}
}
