package model

type FinancialCharges struct {
    Id                  int     `db:"id" json:"id"`
    Tariff              float64 `db:"tariff" json:"tariff"`
    WarehouseCommission float64 `db:"warehouse_commission" json:"warehouse_commission"`
    BankInstallment     float64 `db:"bank_installment" json:"bank_installment"`
}
