package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the Smart Contract structure
type SmartContract struct {
	contractapi.Contract
}

// Invoice defines the data structure for an invoice
type Invoice struct {
	ID       string `json:"id"`
	Supplier string `json:"supplier"`
	Amount   int    `json:"amount"`
	Status   string `json:"status"` // Status: "Uploaded", "Verified", "Approved"
	Uploader string `json:"uploader"`
}

// UploadInvoice (Purchasing): Menambahkan invoice baru ke dalam blockchain
func (s *SmartContract) UploadInvoice(ctx contractapi.TransactionContextInterface, id string, supplier string, amount int, uploader string) error {
	exists, err := s.InvoiceExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("invoice dengan ID %s sudah ada", id)
	}

	invoice := Invoice{
		ID:       id,
		Supplier: supplier,
		Amount:   amount,
		Status:   "Uploaded",
		Uploader: uploader,
	}

	invoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, invoiceJSON)
}

// VerifyInvoice (Finance): Mengubah status invoice menjadi "Verified"
func (s *SmartContract) VerifyInvoice(ctx contractapi.TransactionContextInterface, id string) error {
	invoice, err := s.QueryInvoice(ctx, id)
	if err != nil {
		return err
	}

	if invoice.Status != "Uploaded" {
		return fmt.Errorf("invoice %s tidak bisa diverifikasi, status saat ini: %s", id, invoice.Status)
	}

	invoice.Status = "Verified"
	invoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, invoiceJSON)
}

// ApproveInvoice (Manager): Mengubah status invoice menjadi "Approved"
func (s *SmartContract) ApproveInvoice(ctx contractapi.TransactionContextInterface, id string) error {
	invoice, err := s.QueryInvoice(ctx, id)
	if err != nil {
		return err
	}

	if invoice.Status != "Verified" {
		return fmt.Errorf("invoice %s tidak bisa di-approve, status saat ini: %s", id, invoice.Status)
	}

	invoice.Status = "Approved"
	invoiceJSON, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, invoiceJSON)
}

// QueryInvoice (Investor/Global): Membaca data invoice dari ledger
func (s *SmartContract) QueryInvoice(ctx contractapi.TransactionContextInterface, id string) (*Invoice, error) {
	invoiceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca data dari world state: %v", err)
	}
	if invoiceJSON == nil {
		return nil, fmt.Errorf("invoice %s tidak ditemukan", id)
	}

	var invoice Invoice
	err = json.Unmarshal(invoiceJSON, &invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// InvoiceExists: Fungsi bantuan untuk mengecek ketersediaan data
func (s *SmartContract) InvoiceExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	invoiceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("gagal membaca data dari world state: %v", err)
	}
	return invoiceJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create eFishery chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting eFishery chaincode: %s", err.Error())
	}
}
