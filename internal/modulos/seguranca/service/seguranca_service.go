package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"campus_connect_api/internal/modelos"
)

type ConteudoToken struct {
	UsuarioID string `json:"user_id"`
	Email     string `json:"email"`
	Perfil    string `json:"role"`
	ExpiraEm  int64  `json:"exp"`
}

const duracaoTokenPadrao = 12 * time.Hour

var ErrCredenciaisInvalidas = errors.New("credenciais invalidas")

type SegurancaService struct {
	repositorio SegurancaRepository
}

func NovoSegurancaService(repositorio SegurancaRepository) *SegurancaService {
	return &SegurancaService{repositorio: repositorio}
}

func segredoAssinatura() string {
	if arquivo := strings.TrimSpace(os.Getenv("API_SECRET_FILE")); arquivo != "" {
		dados, err := os.ReadFile(filepath.Clean(arquivo))
		if err == nil {
			if s := strings.TrimSpace(string(dados)); s != "" {
				return s
			}
		}
	}
	if s := strings.TrimSpace(os.Getenv("API_SECRET")); s != "" {
		return s
	}
	return "campusconnect-dev-secret"
}

func GerarToken(usuarioID, email, perfil string, duracao time.Duration) (string, error) {
	if usuarioID == "" || email == "" || perfil == "" {
		return "", errors.New("dados obrigatorios ausentes para token")
	}
	payload := ConteudoToken{
		UsuarioID: usuarioID,
		Email:     email,
		Perfil:    perfil,
		ExpiraEm:  time.Now().Add(duracao).Unix(),
	}
	bruto, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(bruto)
	assinatura := assinar(payloadB64, segredoAssinatura())
	return payloadB64 + "." + assinatura, nil
}

func (servico *SegurancaService) RealizarLogin(contexto context.Context, email, senha string) (modelos.RespostaLogin, error) {
	usuario, err := servico.repositorio.Autenticar(contexto, email, senha)
	if err != nil {
		if errors.Is(err, ErrAutenticacaoInvalida) {
			return modelos.RespostaLogin{}, ErrCredenciaisInvalidas
		}
		return modelos.RespostaLogin{}, err
	}
	token, err := GerarToken(usuario.ID, usuario.Email, usuario.PerfilCodigo, duracaoTokenPadrao)
	if err != nil {
		return modelos.RespostaLogin{}, err
	}
	return modelos.RespostaLogin{
		Token:     token,
		TipoToken: "Bearer",
		Login:     usuario.Email,
		ExpiraEm:  int64(duracaoTokenPadrao.Seconds()),
		Perfil:    usuario.PerfilCodigo,
		UsuarioID: usuario.ID,
	}, nil
}

func ValidarToken(token string) (ConteudoToken, error) {
	partes := strings.Split(token, ".")
	if len(partes) != 2 {
		return ConteudoToken{}, errors.New("token invalido")
	}
	payloadB64, assinaturaRecebida := partes[0], partes[1]
	assinaturaEsperada := assinar(payloadB64, segredoAssinatura())
	if !hmac.Equal([]byte(assinaturaEsperada), []byte(assinaturaRecebida)) {
		return ConteudoToken{}, errors.New("assinatura invalida")
	}
	bruto, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return ConteudoToken{}, errors.New("payload invalido")
	}
	var payload ConteudoToken
	if err := json.Unmarshal(bruto, &payload); err != nil {
		return ConteudoToken{}, errors.New("payload invalido")
	}
	if time.Now().Unix() >= payload.ExpiraEm {
		return ConteudoToken{}, errors.New("token expirado")
	}
	return payload, nil
}

func assinar(payload, segredo string) string {
	mac := hmac.New(sha256.New, []byte(segredo))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
