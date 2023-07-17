CREATE OR REPLACE PROCEDURE transfer (
		sender VARCHAR,
		receiver VARCHAR,
		amount DECIMAL,
        errorDescription out VARCHAR
)
language plpgsql
as $$
DECLARE
    d_receiver_account_type integer;
    d_receiver_account_balance decimal;
    d_sender_account_type integer;
    d_sender_account_balance decimal;
BEGIN
	SELECT a.type, a.balance INTO d_receiver_account_type, d_receiver_account_balance
	                         FROM accounts AS a WHERE a.id = receiver FOR UPDATE;

    SELECT  b.type, b.balance INTO d_sender_account_type, d_sender_account_balance
	                          FROM accounts AS b WHERE b.id = sender FOR UPDATE;

    if d_receiver_account_type = 1 AND d_sender_account_type = 2 THEN
        if d_sender_account_balance < amount OR d_receiver_account_balance < amount THEN
            errorDescription = 'Unsuccessful operation'
            ROLLBACK;
            return;
        END IF;


        UPDATE accounts SET balance = balance - amount WHERE id = sender;

        UPDATE accounts SET balance = balance - amount WHERE id = receiver;

        COMMIT;
        errorDescription = 'ok'
        return;
    END IF;

    if d_receiver_account_type = 2 AND d_sender_account_type = 1 THEN

        UPDATE accounts SET balance = balance + amount WHERE id = sender;

        UPDATE accounts SET balance = balance + amount WHERE id = receiver;
        COMMIT;
        errorDescription = 'ok'
        return;
    END IF;

    if d_receiver_account_type = 1 AND d_sender_account_type = 1 THEN
        if d_receiver_account_balance < amount THEN
            errorDescription = 'Not found any debt'
            ROLLBACK;
            return;
        END IF;

        UPDATE accounts SET balance = balance + amount WHERE id = sender;

        UPDATE accounts SET balance = balance - amount WHERE id = receiver;

        COMMIT;
        errorDescription = 'ok'
        return;
    END IF;

    if d_receiver_account_type = 2 AND d_sender_account_type = 2 THEN
        if d_sender_account_balance < amount THEN
            errorDescription = 'Insufficient balance'
            ROLLBACK;
            return;
        END IF;

        UPDATE accounts SET balance = balance - amount WHERE id = sender;

        UPDATE accounts SET balance = balance + amount WHERE id = receiver;

        COMMIT;
        errorDescription = 'ok'
        return;
    END IF;
END $$;
