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

	http.HandleFunc("/validate", validate)
	http.HandleFunc("/mutate", mutate)
	log.Printf("Starting webhook server on port %d", port)

	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key", nil)
	if err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	admissionReview, err := decodeAdmissionReview(r)
	if err != nil {
		log.Printf("Error decoding admission review request: %v", err)
		http.Error(w, "Error decoding admission review request", http.StatusBadRequest)
		return
	}

	AdmissionResponse := validationAdmissionReview(admissionReview)

	responseBody, err := json.Marshal(AdmissionResponse)
	if err != nil {
		log.Printf("Error marshalling admission response: %v", err)
		http.Error(w, "Error marshalling admission response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func mutate(w http.ResponseWriter, r *http.Request) {
	admissionReview, err := decodeAdmissionReview(r)
	if err != nil {
		log.Printf("Error decoding admission review request: %v", err)
		http.Error(w, "Error decoding admission review request", http.StatusBadRequest)
		return
	}

	// mutate handler
	AdmissionResponse, err := mutateAdmissionReview(admissionReview)
	if err != nil {
		log.Printf("Error mutating admission review: %v", err)
		http.Error(w, "Error mutating admission review", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(AdmissionResponse)
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

func validationAdmissionReview(review *admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
	AdmissionResponse := &admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
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
		Response: AdmissionResponse,
	}
}

func mutateAdmissionReview(review *admissionv1.AdmissionReview) (*admissionv1.AdmissionReview, error) {
	pod := corev1.Pod{}
	err := json.Unmarshal(review.Request.Object.Raw, &pod)
	if err != nil {
		return nil, err
	}

	initContainer := corev1.Container{
		Name:  "busybox-by-mutatehandler",
		Image: "busybox",
	}
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, initContainer)

	containersBytes, err := json.Marshal(&pod.Spec.InitContainers)
	if err != nil {
		return nil, err
	}

	patch := []JSONPatchEntry{
		{
			OP:    "add",
			Path:  "/spec/initContainers",
			Value: containersBytes,
		},
	}

	fmt.Printf("%v", patch)

	// Marshal the patch into JSON bytes
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		return nil, err
	}

	patchType := admissionv1.PatchTypeJSONPatch

	// Create the admission review response with the patch
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:       review.Request.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
	}, nil
}

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}
