package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetSyncChanges(w http.ResponseWriter, r *http.Request) {
	// В реальной системе здесь бы проверялся LastTimestamp от клиента
	// и возвращались только новые записи из sync_log
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready_for_sync"})
}

func (h *Handler) PushSyncChanges(w http.ResponseWriter, r *http.Request) {
	// Прием изменений от клиента и применение их к БД сервера
	w.WriteHeader(http.StatusAccepted)
}
