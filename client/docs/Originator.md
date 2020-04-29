# Originator

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ID** | **string** | Originator ID | [optional] 
**DefaultDepository** | **string** | The depository account to be used by default per transfer. ID must be a valid Originator Depository account | [optional] 
**Identification** | **string** | An identification number by which the receiver is known to the originator. | [optional] 
**CustomerID** | **string** | Optional ID when Originator data was created against Moov&#39;s Customers service | [optional] 
**BirthDate** | [**time.Time**](time.Time.md) | An optional value required for Know Your Customer (KYC) validation of this Originator. This field is not saved by PayGate.  | [optional] 
**Address** | [**Address**](Address.md) |  | [optional] 
**Metadata** | **string** | Populated into the Batch Header CompanyDiscretionaryData field | [optional] 
**Created** | [**time.Time**](time.Time.md) |  | [optional] 
**Updated** | [**time.Time**](time.Time.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


