CREATE OR REPLACE FUNCTION public.process_procedures() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
NEW.drg_definition_tokens = to_tsvector(NEW.drg_definition);
NEW.id = split_part(NEW.drg_definition, ' - ', 1);return new;
end;
$$;

CREATE OR REPLACE FUNCTION public.process_provider_address() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
NEW.address_line = NEW.name || ' ' || NEW.street|| ' ' || NEW.city || ' ' || NEW.state || ' ' || NEW.zip_code;
NEW.address_line_tokens = to_tsvector(NEW.address_line);return new;
end;
$$;

CREATE TABLE public.procedures (
    id character varying NOT NULL,
    average_total_payments numeric,
    average_covered_charges numeric,
    average_medicare_payments numeric,
    total_discharges integer,
    drg_definition text NOT NULL,
    drg_definition_tokens tsvector
);
    
CREATE TABLE public.provider_procedures (
    provider_id integer NOT NULL,
    procedure_id character varying NOT NULL
);

CREATE TABLE public.providers (
    id integer NOT NULL,
    name character varying NOT NULL,
    street character varying NOT NULL,
    city character varying NOT NULL,
    state character(2) NOT NULL,
    zip_code text NOT NULL,
    address_line text,
    address_line_tokens tsvector,
    hrr_description character varying NOT NULL
);

CREATE TABLE public.zip_code_lat_long (
    zip_code text NOT NULL,
    latitude numeric(6,3) NOT NULL,
    longitude numeric(6,3) NOT NULL
);

ALTER TABLE ONLY public.procedures
    ADD CONSTRAINT procedures_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.providers
    ADD CONSTRAINT providers_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.zip_code_lat_long
    ADD CONSTRAINT zip_code_lat_long_pkey PRIMARY KEY (zip_code);

CREATE INDEX procedures_average_total_payments_idx ON public.procedures USING btree (average_total_payments);
CREATE INDEX procedures_drg_definition_tokens_idx ON public.procedures USING gin (drg_definition_tokens);
CREATE INDEX provider_procedures_procedure_id_idx ON public.provider_procedures USING btree (procedure_id);
CREATE INDEX provider_procedures_provider_id_idx ON public.provider_procedures USING btree (provider_id);
CREATE INDEX providers_address_line_tokens_idx ON public.providers USING gin (address_line_tokens);
CREATE INDEX zip_code_lat_long_latitude_idx ON public.zip_code_lat_long USING btree (latitude);
CREATE INDEX zip_code_lat_long_longitude_idx ON public.zip_code_lat_long USING btree (longitude);
CREATE INDEX zip_code_lat_long_zip_code_idx ON public.zip_code_lat_long USING btree (zip_code);

CREATE TRIGGER process_procedures BEFORE INSERT ON public.procedures FOR EACH ROW EXECUTE FUNCTION public.process_procedures();
CREATE TRIGGER process_provider_address BEFORE INSERT ON public.providers FOR EACH ROW EXECUTE FUNCTION public.process_provider_address();

ALTER TABLE ONLY public.provider_procedures
    ADD CONSTRAINT provider_procedures_procedure_id_fkey FOREIGN KEY (procedure_id) REFERENCES public.procedures(id);

ALTER TABLE ONLY public.provider_procedures
    ADD CONSTRAINT provider_procedures_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES public.providers(id);

ALTER TABLE ONLY public.providers
    ADD CONSTRAINT providers_zip_code_fkey FOREIGN KEY (zip_code) REFERENCES public.zip_code_lat_long(zip_code);