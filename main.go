package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"io"
	"client/modules"
)
// ===== STRUCT ORDER =====
type Order struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Service string `json:"service"`
	Message string `json:"message"`
}

func main() {

    // ===== API =====
    http.HandleFunc("/api/order", orderHandler)
    http.HandleFunc("/api/check-order", checkOrderHandler)

    // ===== STATIC FOLDERS =====
    http.Handle("/order/",
        http.StripPrefix("/order",
            http.FileServer(http.Dir("./order")),
        ),
    )

    http.Handle("/pesanan/",
        http.StripPrefix("/pesanan",
            http.FileServer(http.Dir("./pesanan")),
        ),
    )

    http.Handle("/portofolio/",
        http.StripPrefix("/portofolio",
            http.FileServer(http.Dir("./portofolio")),
        ),
    )

    http.Handle("/services/",
        http.StripPrefix("/services",
            http.FileServer(http.Dir("./services")),
        ),
    )

    http.HandleFunc("/filedone/files/", fileProxyHandler)

    // ? ROOT HARUS PALING BAWAH
    http.Handle("/", http.FileServer(http.Dir("./landing")))

    log.Println("Server jalan di :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	var req Order
	json.NewDecoder(r.Body).Decode(&req)

	o := modules.Order{
		ID:       modules.GenerateOrderID(),
		Passcode: modules.GeneratePasscode(),
		Name:     req.Name,
		Contact:  req.Contact,
		Service:  req.Service,
		Message:  req.Message,
		Status:   "baru",
	}

	go notifyAdmin(o)

	json.NewEncoder(w).Encode(map[string]string{
		"id":       o.ID,
		"passcode": o.Passcode,
	})
}

func notifyAdmin(o modules.Order) {
	data, err := json.Marshal(o)
	if err != nil {
		log.Println("CLIENT: gagal marshal:", err)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:9090/api/admin/order",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Println("CLIENT: gagal buat request:", err)
		return
	}

	// HEADER WAJIB
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ADMIN-KEY", "S1kr3tB4ng3t") // sementara hardcode

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("CLIENT: gagal kirim ke admin:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("CLIENT: admin response =", resp.Status)
}

func checkOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("? HIT /api/check-order", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Println("? METHOD SALAH:", r.Method)
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		ID       string `json:"id"`
		Passcode string `json:"passcode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("? JSON ERROR:", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	log.Println("?? DATA DARI CLIENT:", req.ID, req.Passcode)

	body, _ := json.Marshal(req)

	request, err := http.NewRequest(
		"POST",
		"http://localhost:9090/api/admin/check-order",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("? GAGAL BUAT REQUEST KE ADMIN:", err)
		http.Error(w, "Internal error", 500)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-ADMIN-KEY", "S1kr3tB4ng3t")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("? ADMIN OFFLINE:", err)
		http.Error(w, "Admin offline", 500)
		return
	}
	defer resp.Body.Close()

	log.Println("? ADMIN BALIK:", resp.Status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func fileProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := "http://localhost:9090" + r.URL.Path

	log.Println("PROXY FILE:", targetURL)

	req, err := http.NewRequest(r.Method, targetURL, nil)
	if err != nil {
		http.Error(w, "Proxy error", 500)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "File server offline", 502)
		return
	}
	defer resp.Body.Close()

	// copy headers (content-type, content-disposition, dll)
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

