CREATE TABLE team_members (
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (owner_id, member_id)
);
