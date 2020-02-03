ALTER TABLE procedures DROP COLUMN average_total_payments;
ALTER TABLE procedures DROP COLUMN average_medicare_payments;
ALTER TABLE procedures DROP COLUMN average_covered_charges;

ALTER TABLE provider_procedures ADD COLUMN average_total_payments numeric;