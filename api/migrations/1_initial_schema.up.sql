CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE FUNCTION public.tokenise_drg_definition() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
NEW.drg_definition_tokens = to_tsvector(NEW.drg_definition);
return new;
end;
$$;


SET default_tablespace = '';

SET default_with_oids = false;

CREATE TABLE public.inpatient_charge_data (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    provider_id integer NOT NULL,
    provider_name character varying NOT NULL,
    provider_street character varying NOT NULL,
    provider_city character varying NOT NULL,
    provider_state character(2) NOT NULL,
    provider_zip_code integer NOT NULL,
    hrr_description character varying,
    total_discharges integer DEFAULT 0,
    average_covered_charges numeric DEFAULT 0,
    average_total_payments numeric DEFAULT 0,
    average_medicare_payments numeric DEFAULT 0,
    drg_definition text NOT NULL,
    drg_definition_tokens tsvector
);

ALTER TABLE ONLY public.inpatient_charge_data
    ADD CONSTRAINT inpatient_charge_data_pkey PRIMARY KEY (id);

CREATE INDEX inpatient_charge_data_drg_definition_tokens_idx ON public.inpatient_charge_data USING gin (drg_definition_tokens);
CREATE TRIGGER tokenise_drg_def BEFORE INSERT ON public.inpatient_charge_data FOR EACH ROW EXECUTE PROCEDURE public.tokenise_drg_definition();