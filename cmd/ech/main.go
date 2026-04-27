package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hiericho/architect"
)

func main() {
	domain := "crypto.cloudflare.com"
	url := fmt.Sprintf("https://%s/cdn-cgi/trace", domain)
	
	fmt.Printf("1. Resolviendo HTTPS/SVCB record para %s vía DoH...\n", domain)
	echConfig, err := architect.FetchECHConfig(domain)
	if err != nil {
		log.Printf("Aviso ECH: No se encontró configuración o falló la resolución: %v\n", err)
		log.Println("Continuando sin ECH (Fallback)...")
	} else {
		fmt.Printf("2. ECH Config encontrado! [%d bytes]\n", len(echConfig))
	}

	// Clonamos el perfil de Safari y le inyectamos la configuración obtenida
	profile := architect.Safari17iOS
	profile.ECHConfig = echConfig

	// Inicializamos el ProfileManager
	manager := architect.NewProfileManager()
	manager.Register(profile)

	tr := &architect.Transport{
		Manager:            manager,
		DefaultProfile:     profile,
		// Crypto.cloudflare.com uses valid certs, but we skip verify just in case of environment date issues
		InsecureSkipVerify: true, 
	}
	client := &http.Client{Transport: tr}

	fmt.Println("3. Iniciando Handshake TLS 1.3 con ECH...")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", profile.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Extraer el campo 'sni' del trace de Cloudflare para confirmar que fue enmascarado/encriptado.
	// Cuando ECH funciona, el SNI expuesto en el trace suele ser el "Outer SNI" (ej. cloudflare-ech.com) 
	// o directamente la conexión se marca como validada internamente.
	fmt.Println("\n--- Cloudflare Trace Result ---")
	lines := strings.Split(bodyStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "sni=") || strings.HasPrefix(line, "tls=") {
			fmt.Println(line)
		}
	}
}
