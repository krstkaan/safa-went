CREATE TABLE IF NOT EXISTS print_requests (
    id BIGSERIAL PRIMARY KEY,
    requested_at TIMESTAMP WITH TIME ZONE NOT NULL,
    color_copies INT NOT NULL DEFAULT 0,
    bw_copies INT NOT NULL DEFAULT 0,
    description TEXT,
    requester_id BIGINT NOT NULL REFERENCES requesters(id),
    approver_id BIGINT NOT NULL REFERENCES approvers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_print_requests_requester_id ON print_requests (requester_id);
CREATE INDEX IF NOT EXISTS idx_print_requests_approver_id ON print_requests (approver_id);
CREATE INDEX IF NOT EXISTS idx_print_requests_requested_at ON print_requests (requested_at);
CREATE INDEX IF NOT EXISTS idx_print_requests_deleted_at ON print_requests (deleted_at);
