CREATE TABLE financial_charges (
    id SERIAL PRIMARY KEY,
    tariff FLOAT NOT NULL,
    warehouse_commission FLOAT NOT NULL,
    bank_installment FLOAT NOT NULL
);
