ALTER TABLE procedures DROP COLUMN drg_definition_tokens;

DROP TRIGGEr process_procedures ON procedures;

CREATE OR REPLACE FUNCTION process_procedures()
RETURNS trigger AS                                                            
$BODY$                                                                                            
BEGIN                                                                
NEW.id = split_part(NEW.drg_definition, ' - ', 1);return new;
end;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER process_procedures
BEFORE INSERT
ON procedures 
for each ROW
EXECUTE FUNCTION process_procedures();

CREATE INDEX ON procedures USING GIN (drg_definition gin_trgm_ops);