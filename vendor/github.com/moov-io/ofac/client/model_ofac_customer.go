/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// OFAC Customer and metadata
type OfacCustomer struct {
	// OFAC Customer ID
	Id        string    `json:"id,omitempty"`
	SDN       Sdn       `json:"SDN,omitempty"`
	Addresses []Address `json:"Addresses,omitempty"`
	Alts      []Alt     `json:"Alts,omitempty"`
}
