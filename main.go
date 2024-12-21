package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Input struct {
	Array *string `json:"array"`
	Sign  *string `json:"sign"`
}

type Output struct {
	Result *string `json:"result"`
}

func CleanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

	var input Input

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if input.Array == nil || input.Sign == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("param is missing :("))
		return
	}

	sign := *input.Sign

	if len([]rune(sign)) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("некоректно введённый паметр sign: должен быть строкой длиной 1 символ"))
		return
	}

	cleanedString := RemoveCharacter(*input.Array, sign)

	output := Output{Result: &cleanedString}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func RemoveCharacter(input string, charToRemove string) string {
	if charToRemove == "" {
		return input
	}

	result := []rune{}
	for _, char := range input {
		if string(char) != charToRemove {
			result = append(result, char)
		}
	}

	return string(result)
}

func main() {
	http.HandleFunc("/clean", CleanHandler)
	fmt.Println("starting server on 127.0.0.1:8082...")
	err := http.ListenAndServe("127.0.0.1:8082", nil)
	if err != nil {
		panic(err)
	}
}
