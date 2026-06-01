INSERT INTO account_number_sequences (key, next_number, prefix)
VALUES ('default', 1, '408178100')
ON CONFLICT (key) DO NOTHING;

INSERT INTO clients (id, last_name, first_name, middle_name, phone)
VALUES ('00000000-0000-0000-0000-000000000001', 'Касса', '', '', '00000000000')
ON CONFLICT DO NOTHING;

INSERT INTO accounts (id, client_id, account_number, account_type, status, balance)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    '20202810000000000001',
    'А', 'OPEN', 0
)
ON CONFLICT DO NOTHING;

COMMIT;
