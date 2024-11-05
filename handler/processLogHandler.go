package handler

import (
	"ConcurrentLogProcessor/processors"
	"encoding/json"
	"net/http"
)

func ProcessLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported for this endpoint", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("upload-file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	keywords := r.Form["keywords"]
	if len(keywords) == 0 {
		http.Error(w, "Please provide keywords", http.StatusBadRequest)
		return
	}

	outputs, err := processors.ProcessLogFile(file, keywords)

	if err != nil {
		http.Error(w, "Error processing the file", http.StatusBadRequest)
		return
	}

	var result []map[string]interface{}

	for _, output := range outputs {
		result = append(result, map[string]interface{}{output.Key: output.Value})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}

}
