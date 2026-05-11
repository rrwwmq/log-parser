CREATE SCHEMA IF NOT EXISTS logparser;

CREATE TABLE logparser.logs (
    id            SERIAL     PRIMARY KEY,
    file_name     TEXT          NOT NULL,
    status        TEXT          NOT NULL   DEFAULT 'processing',
    uploaded_at   TIMESTAMPTZ   NOT NULL   DEFAULT NOW(),
    node_count    INT           NOT NULL   DEFAULT 0,
    port_count    INT           NOT NULL   DEFAULT 0
);

CREATE TABLE logparser.nodes (
    id            SERIAL     PRIMARY KEY,
    log_id        INT           NOT NULL REFERENCES logparser.logs(id) ON DELETE CASCADE,
    node_guid     TEXT          NOT NULL,
    node_desc     TEXT          NOT NULL,
    node_type     INT           NOT NULL,
    num_ports     INT           NOT NULL
);

CREATE TABLE logparser.nodes_info (
    id                          SERIAL     PRIMARY KEY,
    node_id                     INT           NOT NULL REFERENCES logparser.nodes(id) ON DELETE CASCADE,
    serial_number               TEXT,
    part_number                 TEXT,
    revision                    TEXT,
    product_name                TEXT,
    endianness                  INT,
    enable_endianness_per_job   INT,
    reproducibility_disable     INT
);

CREATE TABLE logparser.ports (
    id            SERIAL     PRIMARY KEY,
    node_id       INT           NOT NULL REFERENCES logparser.nodes(id) ON DELETE CASCADE,
    port_guid     TEXT          NOT NULL,
    port_num      INT           NOT NULL,
    port_state    INT           NOT NULL,
    lid           INT           NOT NULL
);