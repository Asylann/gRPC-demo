CREATE TABLE IF NOT EXISTS public.etag_versions (
    id serial NOT NULL,
    userid integer NOT NULL,
    version integer NOT NULL,
    CONSTRAINT etag_versions_pkey PRIMARY KEY (id)
);