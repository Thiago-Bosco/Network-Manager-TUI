
package history

import (
	"fmt"
	"sync"
	"time"
	
	"networkmanager-tui/logger"
)

type Action struct {
	Timestamp   time.Time
	UserID      string
	Action      string
	Details     string
	Changes     string
	ModifiedBy  string
}

var (
	actions []Action
	mutex   sync.RWMutex
)

func AddAction(userID, actionType, details string, changes string, modifiedBy string) {
	mutex.Lock()
	defer mutex.Unlock()

	action := Action{
		Timestamp:   time.Now(),
		UserID:      userID,
		Action:      actionType,
		Details:     details,
		Changes:     changes,
		ModifiedBy:  modifiedBy,
	}
	actions = append(actions, action)

	// Log a ação com detalhes extras
	logAction(action)
}

func GetHistory() []Action {
	mutex.RLock()
	defer mutex.RUnlock()

	history := make([]Action, len(actions))
	copy(history, actions)
	return history
}

func logAction(action Action) {
	logEntry := fmt.Sprintf(`
=== Log de Ação ===
Data/Hora: %s
Usuário: %s
Ação: %s
Detalhes: %s
Modificado por: %s
Alterações: %s
==================`,
		action.Timestamp.Format("02/01/2006 15:04:05"),
		action.UserID,
		action.Action, 
		action.Details,
		action.ModifiedBy,
		action.Changes)
	
	logger.LogInfo(logEntry)
}
