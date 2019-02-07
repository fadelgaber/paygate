/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Physical address from OFAC list
type Address struct {
	EntityID                    string `json:"EntityID,omitempty"`
	AddressID                   string `json:"AddressID,omitempty"`
	Address                     string `json:"Address,omitempty"`
	CityStateProvincePostalCode string `json:"CityStateProvincePostalCode,omitempty"`
	Country                     string `json:"Country,omitempty"`
}
