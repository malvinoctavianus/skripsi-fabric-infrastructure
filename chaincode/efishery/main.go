package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the smart contract structure
type SmartContract struct {
	contractapi.Contract
}

// Invoice defines the structure for an invoice
type Invoice struct {
	ID      string `json:"id"`
	Tanggal string `json:"tanggal"`
	Nominal string `json:"nominal"`
	Pemasok string `json:"pemasok"`
	Status  string `json:"status"` // Uploaded, Verified, Approved
}

// UploadInvoice - Called by Purchasing
func (s *SmartContract) UploadInvoice(ctx contractapi.TransactionContextInterface, id string, tanggal string, nominal string, pemasok string) error {
	invoice := Invoice{
		ID:      id,
		Tanggal: tanggal,
		Nominal: nominal,
		Pemasok: pemasok,
		Status:  "Uploaded",
	}

	invoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, invoiceJSON)
}

// VerifyInvoice - Called by Finance
func (s *SmartContract) VerifyInvoice(ctx contractapi.TransactionContextInterface, id string) error {
	invoiceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("Gagal membaca dari blockchain: %v", err)
	}
	if invoiceJSON == nil {
		return fmt.Errorf("Invoice %s tidak ditemukan", id)
	}

	var invoice Invoice
	err = json.Unmarshal(invoiceJSON, &invoice)
	if err != nil {
		return err
	}

	if invoice.Status != "Uploaded" {
		return fmt.Errorf("Invoice %s tidak dalam status Uploaded. Status saat ini: %s", id, invoice.Status)
	}

	invoice.Status = "Verified"
	updatedInvoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedInvoiceJSON)
}

// ApproveInvoice - Called by Manager
func (s *SmartContract) ApproveInvoice(ctx contractapi.TransactionContextInterface, id string) error {
	invoiceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("Gagal membaca dari blockchain: %v", err)
	}
	if invoiceJSON == nil {
		return fmt.Errorf("Invoice %s tidak ditemukan", id)
	}

	var invoice Invoice
	err = json.Unmarshal(invoiceJSON, &invoice)
	if err != nil {
		return err
	}

	if invoice.Status != "Verified" {
		return fmt.Errorf("Invoice %s belum diverifikasi oleh Finance. Status saat ini: %s", id, invoice.Status)
	}

	invoice.Status = "Approved"
	updatedInvoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedInvoiceJSON)
}

// QueryInvoice - Called by Investor (Read-Only)
func (s *SmartContract) QueryInvoice(ctx contractapi.TransactionContextInterface, id string) (*Invoice, error) {
	invoiceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca dari blockchain: %v", err)
	}
	if invoiceJSON == nil {
		return nil, fmt.Errorf("Invoice %s tidak ditemukan", id)
	}

	var invoice Invoice
	err = json.Unmarshal(invoiceJSON, &invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating eFishery chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting eFishery chaincode: %v", err)
	}
}
