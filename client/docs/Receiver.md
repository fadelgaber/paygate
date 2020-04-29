# Receiver

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ID** | **string** | Receiver ID | [optional] 
**Email** | **string** | The receivers email address | [optional] 
**DefaultDepository** | **string** | The depository account to be used by default per transfer. ID must be a valid Receiver Depository account | [optional] 
**Status** | **string** | Defines the status of the Receiver | [optional] 
**BirthDate** | [**time.Time**](time.Time.md) | An optional object required for Know Your Customer (KYC) validation of this Receiver. This field is not saved by PayGate.  | [optional] 
**Address** | [**Address**](Address.md) |  | [optional] 
**CustomerID** | **string** | Optional ID when Originator data was created against Moov&#39;s Customers service | [optional] 
**Metadata** | **string** | Populated into the Entry Detail IndividualName field | [optional] 
**Created** | [**time.Time**](time.Time.md) |  | [optional] 
**Updated** | [**time.Time**](time.Time.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


