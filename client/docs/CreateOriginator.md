# CreateOriginator

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultDepository** | **string** | The depository account to be used by default per transfer. ID must be a valid Originator Depository account | 
**Identification** | **string** | An identification number by which the receiver is known to the originator. | 
**BirthDate** | [**time.Time**](time.Time.md) | An optional value required for Know Your Customer (KYC) validation of this Originator. This field is not saved by PayGate.  | [optional] 
**Address** | [**Address**](Address.md) |  | [optional] 
**Metadata** | **string** | Populated into the Batch Header CompanyDiscretionaryData field | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


