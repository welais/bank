CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS clients (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    last_name   VARCHAR(100) NOT NULL,
    first_name  VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    phone       VARCHAR(20)  NOT NULL UNIQUE,
    created_at  TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_clients_lastname
    ON clients USING BTREE (last_name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_clients_phone
    ON clients USING BTREE (phone);

CREATE TABLE IF NOT EXISTS account_number_sequences (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    key         VARCHAR(50) NOT NULL UNIQUE,
    next_number BIGINT      NOT NULL DEFAULT 1,
    prefix      VARCHAR(9)  NOT NULL DEFAULT '408178100',
    created_at  TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_account_number_sequences_key
    ON account_number_sequences USING BTREE (key);

CREATE TABLE IF NOT EXISTS accounts (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id      UUID          NOT NULL,
    account_number VARCHAR(20)   NOT NULL UNIQUE,
    account_type   VARCHAR(10)   NOT NULL DEFAULT 'А'
                                 CONSTRAINT chk_account_type CHECK (account_type IN ('А', 'П')),
    status         VARCHAR(10)   NOT NULL DEFAULT 'OPEN'
                                 CONSTRAINT chk_status CHECK (status IN ('OPEN', 'CLOSE')),
    balance        DECIMAL(20,2) NOT NULL DEFAULT 0,
    opened_at      DATE          NOT NULL DEFAULT now(),
    closed_at      DATE,
    created_at     TIMESTAMP     NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP     NOT NULL DEFAULT now(),
    CONSTRAINT fk_accounts_clients FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_account_number
    ON accounts USING BTREE (account_number);
CREATE INDEX IF NOT EXISTS idx_accounts_client_id
    ON accounts USING BTREE (client_id);

CREATE TABLE IF NOT EXISTS journal_entries (
    id                UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_number      VARCHAR(30)   NOT NULL UNIQUE,
    entry_date        DATE          NOT NULL,
    debit_account_id  UUID          NOT NULL,
    credit_account_id UUID          NOT NULL,
    amount            DECIMAL(20,2) NOT NULL,
    payment_purpose   VARCHAR(500),
    created_at        TIMESTAMP     NOT NULL DEFAULT now(),
    CONSTRAINT fk_transactions_debit_account  FOREIGN KEY (debit_account_id)  REFERENCES accounts(id),
    CONSTRAINT fk_transactions_credit_account FOREIGN KEY (credit_account_id) REFERENCES accounts(id),
    CONSTRAINT chk_debit_not_equal_credit     CHECK (debit_account_id <> credit_account_id),
    CONSTRAINT chk_amount_positive            CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_journal_entry_date
    ON journal_entries USING BTREE (entry_date);
CREATE INDEX IF NOT EXISTS idx_journal_entry_debit_account
    ON journal_entries USING BTREE (debit_account_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_credit_account
    ON journal_entries USING BTREE (credit_account_id);

CREATE TABLE IF NOT EXISTS account_statements (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id       UUID          NOT NULL,
    entry_id         UUID          NOT NULL,
    statement_date   DATE          NOT NULL,
    side             VARCHAR(6)    NOT NULL
                                   CONSTRAINT chk_side_values CHECK (side IN ('Debit', 'Credit')),
    incoming_balance DECIMAL(20,2) NOT NULL DEFAULT 0,
    amount           DECIMAL(20,2) NOT NULL
                                   CONSTRAINT chk_stmt_amount_positive CHECK (amount > 0),
    outgoing_balance DECIMAL(20,2) NOT NULL DEFAULT 0,
    created_at       TIMESTAMP     NOT NULL DEFAULT now(),
    CONSTRAINT fk_account_statements_account FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT fk_account_statements_entry   FOREIGN KEY (entry_id)   REFERENCES journal_entries(id),

    CONSTRAINT chk_balance_calculation CHECK (
        (side = 'Debit'  AND outgoing_balance = incoming_balance - amount) OR
        (side = 'Credit' AND outgoing_balance = incoming_balance + amount)
    )
);

CREATE INDEX IF NOT EXISTS idx_statements_account_id_date
    ON account_statements USING BTREE (account_id, statement_date, entry_id);
CREATE INDEX IF NOT EXISTS idx_statements_account_id
    ON account_statements USING BTREE (account_id);
