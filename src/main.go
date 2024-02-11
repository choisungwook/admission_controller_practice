package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 443, "Port to listen on")
	flag.Parse()

	http.HandleFunc("/validate", handleValidate)
	log.Printf("Starting webhook server on port %d", port)

	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key", nil)
	if err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received admission review request: %v", r.Body)

	decoder := json.NewDecoder(r.Body)
	var admissionReview admissionv1.AdmissionReview
	if err := decoder.Decode(&admissionReview); err != nil {
		log.Printf("Error decoding admission review request: %v", err)
		http.Error(w, "Error decoding admission review request", http.StatusBadRequest)
		return
	}

	// log.Printf("Received admission review request: %s", admissionReview)

	admissionResponse := handleAdmissionReview(admissionReview)

	responseBody, err := json.Marshal(admissionResponse)
	if err != nil {
		log.Printf("Error marshalling admission response: %v", err)
		http.Error(w, "Error marshalling admission response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func handleAdmissionReview(review admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
		Result: &metav1.Status{
			Code:    200,
			Message: "Success",
		},
	}

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
	}
}
