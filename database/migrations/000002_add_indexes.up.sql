CREATE INDEX IF NOT EXISTS idx_user_role ON "user" ("role");

CREATE INDEX IF NOT EXISTS idx_pvz_city ON pvz (city);
CREATE INDEX IF NOT EXISTS idx_pvz_registration_date ON pvz ("registrationDate");

CREATE INDEX IF NOT EXISTS idx_reception_pvz_id ON reception ("pvzId");

CREATE INDEX IF NOT EXISTS idx_reception_datetime ON reception ("dateTime");
CREATE INDEX IF NOT EXISTS idx_reception_status ON reception (status);

CREATE INDEX IF NOT EXISTS idx_product_reception_id ON product ("receptionId");

CREATE INDEX IF NOT EXISTS idx_product_datetime ON product ("dateTime");
CREATE INDEX IF NOT EXISTS idx_product_type ON product (type);
