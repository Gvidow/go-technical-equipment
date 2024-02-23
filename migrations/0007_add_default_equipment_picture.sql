ALTER TABLE equipment ALTER COLUMN picture SET DEFAULT('http://localhost:9000/equipment/default.png');
UPDATE equipment SET picture = DEFAULT WHERE picture IS NULL OR picture = '';
ALTER TABLE equipment ALTER COLUMN picture SET NOT NULL;
