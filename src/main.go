package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
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
	admissionReview, err := decodeAdmissionReview(r)
	if err != nil {
		log.Printf("Error decoding admission review request: %v", err)
		http.Error(w, "Error decoding admission review request", http.StatusBadRequest)
		return
	}

	// for debug
	pod, err := decodePodSpec(admissionReview)
	if err != nil {
		log.Printf("Error decoding pod spec: %v", err)
		http.Error(w, "Error decoding pod spec", http.StatusBadRequest)
		return
	}
	log.Printf("Unmarshal Pod: %v", pod)

	// validation handler
	admissionResponse := validationAdmissionReview(admissionReview)

	responseBody, err := json.Marshal(admissionResponse)
	if err != nil {
		log.Printf("Error marshalling admission response: %v", err)
		http.Error(w, "Error marshalling admission response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func decodeAdmissionReview(r *http.Request) (*admissionv1.AdmissionReview, error) {
	var admissionReview admissionv1.AdmissionReview

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	err := json.Unmarshal(body, &admissionReview)
	if err != nil {
		return nil, err
	}

	return &admissionReview, nil
}

func decodePodSpec(admissionReview *admissionv1.AdmissionReview) (*corev1.Pod, error) {
	pod := corev1.Pod{}
	err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
	if err != nil {
		return nil, err
	}

	return &pod, nil
}

func validationAdmissionReview(review *admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
	admissionResponse := &admissionv1.AdmissionResponse{
		UID: review.Request.UID,
		Allowed: true,
		Result: &metav1.Status{
			Code:    http.StatusOK,
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
