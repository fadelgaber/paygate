# Transfer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ID** | **string** | ID to uniquely identify this transfer. If omitted, one will be generated | [optional] 
**TransferType** | **string** | Type of transaction being actioned against the receiving institution. Expected values are pull (debits) or push (credits). Only one period used to signify decimal value will be included. | [optional] 
**Amount** | **string** | Amount of money. USD - United States. | 
**Originator** | **string** | ID of the Originator account initiating the transfer. | 
**OriginatorDepository** | **string** | ID of the Originating Depository used with this transfer. | [optional] 
**Receiver** | **string** | ID of the Receiver account the transfer was sent to. | 
**ReceiverDepository** | **string** | ID of the Receiving Depository used with this transfer. | [optional] 
**Description** | **string** | Brief description of the transaction, that may appear on the receiving entity’s financial statement | 
**StandardEntryClassCode** | **string** | Standard Entry Class code will be generated based on Receiver type for CCD and PPD | [optional] 
**Status** | [**TransferStatus**](TransferStatus.md) |  | [optional] 
**SameDay** | **bool** | When set to true this indicates the transfer should be processed the same day if possible. | [optional] [default to false]
**ReturnCode** | [**ReturnCode**](ReturnCode.md) |  | [optional] 
**Created** | [**time.Time**](time.Time.md) |  | [optional] 
**CCDDetail** | [**CcdDetail**](CCDDetail.md) |  | [optional] 
**IATDetail** | [**IatDetail**](IATDetail.md) |  | [optional] 
**PPDDetail** | [**PpdDetail**](PPDDetail.md) |  | [optional] 
**TELDetail** | [**TelDetail**](TELDetail.md) |  | [optional] 
**WEBDetail** | [**WebDetail**](WEBDetail.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


