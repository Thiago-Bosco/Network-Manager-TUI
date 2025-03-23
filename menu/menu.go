
package menu

import (
        "flag"
        "fmt"
        "os/exec"
        "time"

        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"

        "networkmanager-tui/i18n"
        "networkmanager-tui/network"
        "networkmanager-tui/sysinfo"
)

// Cores padrão para UI
var (
        corBorda        = tcell.ColorDeepSkyBlue    
        corFundo        = tcell.ColorBlack          
        corTextoPrimario = tcell.ColorWhite         
        corSecundaria   = tcell.ColorAqua           
        corDestaque     = tcell.ColorTurquoise      
        corTitulo       = tcell.ColorTurquoise      
        corFundoBotao   = tcell.ColorTurquoise      
        corTextoBotao   = tcell.ColorBlack          
        corSucesso      = tcell.ColorPaleGreen      
        corErro         = tcell.ColorSalmon         
        corFundoCampo   = tcell.ColorMidnightBlue   
        corCabecalho    = tcell.ColorDodgerBlue     
        corInfo         = tcell.ColorLightSkyBlue   
)

// IniciarMenu apresenta o menu principal do sistema
func IniciarMenu(app *tview.Application) {
        menuPrincipal := criarMenuPrincipal(app)
        app.SetRoot(menuPrincipal, true)
}

[Continua com o resto do arquivo...]
