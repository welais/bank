CREATE OR REPLACE PROCEDURE add_client(
    p_last_name   VARCHAR,
    p_first_name  VARCHAR,
    p_middle_name VARCHAR,
    p_phone       VARCHAR,
    OUT p_id      UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO clients (last_name, first_name, middle_name, phone)
    VALUES (p_last_name, p_first_name, p_middle_name, p_phone)
    RETURNING id INTO p_id;
END;
$$;

CREATE OR REPLACE PROCEDURE update_client(
    p_id          UUID,
    p_last_name   VARCHAR,
    p_first_name  VARCHAR,
    p_middle_name VARCHAR,
    p_phone       VARCHAR
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE clients
    SET last_name   = p_last_name,
        first_name  = p_first_name,
        middle_name = p_middle_name,
        phone       = p_phone,
        updated_at  = now()
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Клиент с id % не найден', p_id;
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE open_account(
    p_client_id    UUID,
    p_account_type VARCHAR,
    OUT p_id       UUID,
    OUT p_number   VARCHAR
)
LANGUAGE plpgsql AS $$
DECLARE
    v_seq account_number_sequences%ROWTYPE;
BEGIN
    SELECT * INTO v_seq
    FROM account_number_sequences
    WHERE key = 'default'
    FOR UPDATE;

    p_number := v_seq.prefix || LPAD(v_seq.next_number::TEXT, 11, '0');

    UPDATE account_number_sequences
    SET next_number = next_number + 1
    WHERE key = 'default';


    INSERT INTO accounts (client_id, account_number, account_type, status, balance)
    VALUES (p_client_id, p_number, p_account_type, 'OPEN', 0)
    RETURNING id INTO p_id;
END;
$$;

CREATE OR REPLACE PROCEDURE close_account(
    p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE accounts
    SET status     = 'CLOSE',
        closed_at  = now(),
        updated_at = now()
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Счёт с id % не найден', p_id;
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE add_journal_entry(
    p_entry_number      VARCHAR,
    p_entry_date        DATE,
    p_debit_account_id  UUID,
    p_credit_account_id UUID,
    p_amount            DECIMAL,
    p_payment_purpose   VARCHAR,
    OUT p_id            UUID
)
LANGUAGE plpgsql AS $$
DECLARE
    v_debit_balance  DECIMAL(20,2);
    v_credit_balance DECIMAL(20,2);
BEGIN
    
    SELECT balance INTO v_debit_balance
    FROM accounts
    WHERE id = p_debit_account_id AND status = 'OPEN'
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Дебетовый счёт % не найден или закрыт', p_debit_account_id;
    END IF;

    SELECT balance INTO v_credit_balance
    FROM accounts
    WHERE id = p_credit_account_id AND status = 'OPEN'
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Кредитовый счёт % не найден или закрыт', p_credit_account_id;
    END IF;

  
    UPDATE accounts SET balance = v_debit_balance - p_amount, updated_at = now()
    WHERE id = p_debit_account_id;

    
    UPDATE accounts SET balance = v_credit_balance + p_amount, updated_at = now()
    WHERE id = p_credit_account_id;


    INSERT INTO journal_entries
        (entry_number, entry_date, debit_account_id, credit_account_id, amount, payment_purpose)
    VALUES
        (p_entry_number, p_entry_date, p_debit_account_id, p_credit_account_id, p_amount, p_payment_purpose)
    RETURNING id INTO p_id;


    INSERT INTO account_statements
        (account_id, entry_id, statement_date, side, incoming_balance, amount, outgoing_balance)
    VALUES
        (p_debit_account_id, p_id, p_entry_date, 'Debit',
         v_debit_balance, p_amount, v_debit_balance - p_amount);


    INSERT INTO account_statements
        (account_id, entry_id, statement_date, side, incoming_balance, amount, outgoing_balance)
    VALUES
        (p_credit_account_id, p_id, p_entry_date, 'Credit',
         v_credit_balance, p_amount, v_credit_balance + p_amount);
END;
$$;


CREATE OR REPLACE PROCEDURE delete_journal_entry(
    p_id UUID
)
LANGUAGE plpgsql AS $$
DECLARE
    v_entry            journal_entries%ROWTYPE;
    v_stmt             RECORD;
    v_running_bal      DECIMAL(20,2);

    v_debit_stmt_date  DATE;
    v_debit_stmt_ts    TIMESTAMP;
    v_credit_stmt_date DATE;
    v_credit_stmt_ts   TIMESTAMP;
BEGIN
    SELECT * INTO v_entry FROM journal_entries WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Проводка с id % не найдена', p_id;
    END IF;

 
    SELECT statement_date, created_at
    INTO v_debit_stmt_date, v_debit_stmt_ts
    FROM account_statements
    WHERE entry_id = p_id AND account_id = v_entry.debit_account_id
    LIMIT 1;

   
    SELECT statement_date, created_at
    INTO v_credit_stmt_date, v_credit_stmt_ts
    FROM account_statements
    WHERE entry_id = p_id AND account_id = v_entry.credit_account_id
    LIMIT 1;

  
    UPDATE accounts SET balance = balance + v_entry.amount, updated_at = now()
    WHERE id = v_entry.debit_account_id;

    UPDATE accounts SET balance = balance - v_entry.amount, updated_at = now()
    WHERE id = v_entry.credit_account_id;

    DELETE FROM account_statements WHERE entry_id = p_id;
    DELETE FROM journal_entries     WHERE id       = p_id;

 
    SELECT outgoing_balance INTO v_running_bal
    FROM account_statements
    WHERE account_id = v_entry.debit_account_id
      AND (statement_date, created_at) < (v_debit_stmt_date, v_debit_stmt_ts)
    ORDER BY statement_date DESC, created_at DESC
    LIMIT 1;


    IF v_running_bal IS NULL THEN
        v_running_bal := 0;
    END IF;

    FOR v_stmt IN
        SELECT id, side, amount
        FROM account_statements
        WHERE account_id = v_entry.debit_account_id
          AND (statement_date, created_at) >= (v_debit_stmt_date, v_debit_stmt_ts)
        ORDER BY statement_date, created_at
    LOOP
        UPDATE account_statements
        SET incoming_balance = v_running_bal,
            outgoing_balance = CASE v_stmt.side
                                   WHEN 'Debit'  THEN v_running_bal - v_stmt.amount
                                   WHEN 'Credit' THEN v_running_bal + v_stmt.amount
                               END
        WHERE id = v_stmt.id;

        v_running_bal := CASE v_stmt.side
                             WHEN 'Debit'  THEN v_running_bal - v_stmt.amount
                             WHEN 'Credit' THEN v_running_bal + v_stmt.amount
                         END;
    END LOOP;

    SELECT outgoing_balance INTO v_running_bal
    FROM account_statements
    WHERE account_id = v_entry.credit_account_id
      AND (statement_date, created_at) < (v_credit_stmt_date, v_credit_stmt_ts)
    ORDER BY statement_date DESC, created_at DESC
    LIMIT 1;

    IF v_running_bal IS NULL THEN
        v_running_bal := 0;
    END IF;

    FOR v_stmt IN
        SELECT id, side, amount
        FROM account_statements
        WHERE account_id = v_entry.credit_account_id
          AND (statement_date, created_at) >= (v_credit_stmt_date, v_credit_stmt_ts)
        ORDER BY statement_date, created_at
    LOOP
        UPDATE account_statements
        SET incoming_balance = v_running_bal,
            outgoing_balance = CASE v_stmt.side
                                   WHEN 'Debit'  THEN v_running_bal - v_stmt.amount
                                   WHEN 'Credit' THEN v_running_bal + v_stmt.amount
                               END
        WHERE id = v_stmt.id;

        v_running_bal := CASE v_stmt.side
                             WHEN 'Debit'  THEN v_running_bal - v_stmt.amount
                             WHEN 'Credit' THEN v_running_bal + v_stmt.amount
                         END;
    END LOOP;
END;
$$;
