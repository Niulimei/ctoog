CREATE TABLE cc_repo (
    pvob text default '',
    component text default '',
    stream text default '',
    primary key (pvob, component)
);
