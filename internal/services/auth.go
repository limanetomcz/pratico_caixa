package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	EmpresaID *int   `json:"empresa_id"`
}

// ValidateJWT valida o token JWT com o serviço de autenticação
func ValidateJWT(token string) (*User, error) {
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://auth_nginx:80"
	}

	// Remover barra final se houver
	if len(authURL) > 0 && authURL[len(authURL)-1] == '/' {
		authURL = authURL[:len(authURL)-1]
	}

	url := fmt.Sprintf("%s/api/v1/me", authURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Ler o body completo para debug e decodificação
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token inválido (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			User User `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w, body: %s", err, string(bodyBytes))
	}

	if !response.Success {
		return nil, fmt.Errorf("resposta não indica sucesso")
	}

	return &response.Data.User, nil
}

