
package history

import (
	"fmt"
	"sync"
	"time"
	
	"networkmanager-tui/logger"
)

type Action struct {
	Timestamp time.Time
	UserID    string
	Action    string
	Details   string
}

var (
	actions []Action
	mutex   sync.RWMutex
)

func AddAction(userID, actionType, details string) {
	mutex.Lock()
	defer mutex.Unlock()

	action := Action{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    actionType,
		Details:   details,
	}
	actions = append(actions, action)

	// Log a ação também
	logAction(action)
}

func GetHistory() []Action {
	mutex.RLock()
	defer mutex.RUnlock()

	// Retorna uma cópia para evitar modificações externas
	history := make([]Action, len(actions))
	copy(history, actions)
	return history
}

func logAction(action Action) {
	logEntry := fmt.Sprintf("[%s] User %s: %s - %s",
		action.Timestamp.Format("2006-01-02 15:04:05"),
		action.UserID,
		action.Action,
		action.Details)
	
	logger.LogInfo(logEntry)
}
